package utl

import (
	"os"
	"path/filepath"
	"time"

	"github.com/abema/go-mp4"
	"github.com/dsoprea/go-exif/v3"
	exifc "github.com/dsoprea/go-exif/v3/common"
	jpegs "github.com/dsoprea/go-jpeg-image-structure/v2"
	pngs "github.com/dsoprea/go-png-image-structure/v2"
	"github.com/spf13/cast"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func ReadMetadata(filePath string) (map[string]any, error) {
	data := make(map[string]any)
	if HasSuffix(filePath, ".jpg", ".png") {
		exifs, err := exif.SearchFileAndExtractExif(filePath)
		if err != nil {
			return nil, err
		}
		tags, _, err := exif.GetFlatExifData(exifs, &exif.ScanOptions{})
		if err != nil {
			return nil, err
		}
		for _, t := range tags {
			if Contains(t.TagName, "DateTime") {
				dateTime, err := time.Parse("2006:01:02 15:04:05", cast.ToString(t.Value))
				if err != nil {
					return nil, err
				}
				t.Value = dateTime.Format(time.DateTime)
			}
			data[t.TagName] = t.Value
		}
	} else if HasSuffix(filePath, ".mov", ".mp4") {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		boxes, err := mp4.ExtractBoxWithPayload(file, nil, mp4.BoxPath{mp4.BoxTypeMoov(), mp4.BoxTypeMvhd()})
		if err != nil {
			return nil, err
		}
		if len(boxes) > 0 {
			mvhd := boxes[0].Payload.(*mp4.Mvhd)
			data["CreationTime"] = cast.ToTime(mvhd.GetCreationTime()).AddDate(-66, 0, 0).Format(time.DateTime)
			data["ModificationTime"] = cast.ToTime(mvhd.GetModificationTime()).AddDate(-66, 0, 0).Format(time.DateTime)
			data["Timescale"] = mvhd.Timescale
			data["Duration"] = mvhd.GetDuration()
			data["Rate"] = mvhd.Rate
			data["Volume"] = mvhd.Volume
		}
	}
	return data, nil
}

func WriteMetadata(inFilePath, outFilePath string, data map[string]any) error {
	if HasSuffix(inFilePath, ".jpg") {
		context, err := jpegs.NewJpegMediaParser().ParseFile(inFilePath)
		if err != nil {
			return err
		}
		segments := context.(*jpegs.SegmentList)
		exifBuilder, _ := segments.ConstructExifBuilder()
		if exifBuilder == nil {
			_, err = segments.DropExif()
			if err != nil {
				return err
			}
			exifBuilder, err = segments.ConstructExifBuilder()
			if err != nil {
				return err
			}
		}
		ifdBuilder, _ := exifBuilder.ChildWithTagId(exifc.IfdExifStandardIfdIdentity.TagId())
		if ifdBuilder == nil {
			mapping, err := exifc.NewIfdMappingWithStandard()
			if err != nil {
				return err
			}
			index := exif.NewTagIndex()
			err = exif.LoadStandardTags(index)
			if err != nil {
				return err
			}
			ifdBuilder = exif.NewIfdBuilder(mapping, index, exifc.IfdExifStandardIfdIdentity, exifc.EncodeDefaultByteOrder)
			err = exifBuilder.AddChildIb(ifdBuilder)
			if err != nil {
				return err
			}
		}
		for k, v := range data {
			if Contains(k, "DateTime") {
				dateTime, err := time.Parse(time.DateTime, cast.ToString(v))
				if err != nil {
					return err
				}
				v = dateTime.Format("2006:01:02 15:04:05")
			}
			err = ifdBuilder.SetStandardWithName(k, v)
			if err != nil {
				return err
			}
		}
		err = segments.SetExif(exifBuilder)
		if err != nil {
			return err
		}
		outFile, err := os.Create(outFilePath)
		if err != nil {
			return err
		}
		defer outFile.Close()
		err = segments.Write(outFile)
		if err != nil {
			return err
		}
	} else if HasSuffix(inFilePath, ".png") {
		context, err := pngs.NewPngMediaParser().ParseFile(inFilePath)
		if err != nil {
			return err
		}
		chunks := context.(*pngs.ChunkSlice)
		exifBuilder, err := chunks.ConstructExifBuilder()
		if err != nil {
			return err
		}
		ifdBuilder, _ := exifBuilder.ChildWithTagId(exifc.IfdExifStandardIfdIdentity.TagId())
		if ifdBuilder == nil {
			mapping, err := exifc.NewIfdMappingWithStandard()
			if err != nil {
				return err
			}
			index := exif.NewTagIndex()
			err = exif.LoadStandardTags(index)
			if err != nil {
				return err
			}
			ifdBuilder = exif.NewIfdBuilder(mapping, index, exifc.IfdExifStandardIfdIdentity, exifc.EncodeDefaultByteOrder)
			err = exifBuilder.AddChildIb(ifdBuilder)
			if err != nil {
				return err
			}
		}
		for k, v := range data {
			if Contains(k, "DateTime") {
				dateTime, err := time.Parse(time.DateTime, cast.ToString(v))
				if err != nil {
					return err
				}
				v = dateTime.Format("2006:01:02 15:04:05")
			}
			err = ifdBuilder.SetStandardWithName(k, v)
			if err != nil {
				return err
			}
		}
		err = chunks.SetExif(exifBuilder)
		if err != nil {
			return err
		}
		outFile, err := os.Create(outFilePath)
		if err != nil {
			return err
		}
		defer outFile.Close()
		err = chunks.WriteTo(outFile)
		if err != nil {
			return err
		}
	} else if HasSuffix(inFilePath, ".mov", ".mp4") {
		inFile, err := os.Open(inFilePath)
		if err != nil {
			return err
		}
		defer inFile.Close()
		outFile, err := os.Create(outFilePath)
		if err != nil {
			return err
		}
		defer outFile.Close()
		writer := mp4.NewWriter(outFile)
		_, err = mp4.ReadBoxStructure(inFile, func(handle *mp4.ReadHandle) (any, error) {
			if handle.BoxInfo.Type == mp4.BoxTypeMdat() || !handle.BoxInfo.IsSupportedType() {
				return nil, writer.CopyBox(inFile, &handle.BoxInfo)
			}
			_, err = writer.StartBox(&handle.BoxInfo)
			if err != nil {
				return nil, err
			}
			box, _, err := handle.ReadPayload()
			if err != nil {
				return nil, err
			}
			if handle.BoxInfo.Type == mp4.BoxTypeMvhd() {
				mvhd := box.(*mp4.Mvhd)
				for k, v := range data {
					if Contains(k, "Time") {
						v = cast.ToTime(v).AddDate(66, 0, 0).Unix()
					}
					if k == "CreationTime" {
						if mvhd.Version == 0 {
							mvhd.CreationTimeV0 = cast.ToUint32(v)
						} else if mvhd.Version == 1 {
							mvhd.CreationTimeV1 = cast.ToUint64(v)
						}
					} else if k == "ModificationTime" {
						if mvhd.Version == 0 {
							mvhd.ModificationTimeV0 = cast.ToUint32(v)
						} else if mvhd.Version == 1 {
							mvhd.ModificationTimeV1 = cast.ToUint64(v)
						}
					}
				}
			}
			_, err = mp4.Marshal(writer, box, handle.BoxInfo.Context)
			if err != nil {
				return nil, err
			}
			_, err = handle.Expand()
			if err != nil {
				return nil, err
			}
			_, err = writer.EndBox()
			return nil, err
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func TakeStream(outs map[string]map[string]any, in string, inArgs ...map[string]any) error {
	inArg := make(map[string]any)
	if len(inArgs) > 0 {
		inArg = inArgs[0]
	}
	if HasPrefix(in, "rtsp://") {
		inArg["rtsp_transport"] = "tcp"
	}
	stream := ffmpeg.Input(in, inArg)
	var streams []*ffmpeg.Stream
	for k, v := range outs {
		err := MakeDir(filepath.Dir(k))
		if err != nil {
			return err
		}
		streams = append(streams, stream.Output(k, v))
	}
	return ffmpeg.MergeOutputs(streams...).OverWriteOutput().ErrorToStdOut().Run()
}
