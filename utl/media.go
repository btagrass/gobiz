package utl

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"time"

	"github.com/abema/go-mp4"
	"github.com/corona10/goimagehash"
	"github.com/dsoprea/go-exif/v3"
	exic "github.com/dsoprea/go-exif/v3/common"
	jpegs "github.com/dsoprea/go-jpeg-image-structure/v2"
	pngs "github.com/dsoprea/go-png-image-structure/v2"
	"github.com/spf13/cast"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

var (
	ImageFilters = []string{".jpg", ".png"}
	VideoFilters = []string{".mov", ".mp4"}
)

func HashImage(filePath string) (uint64, error) {
	if !HasSuffix(filePath, ImageFilters...) {
		return 0, fmt.Errorf("format not supported")
	}
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	image, _, err := image.Decode(file)
	if err != nil {
		return 0, err
	}
	hash, err := goimagehash.DifferenceHash(image)
	if err != nil {
		return 0, err
	}
	return hash.GetHash(), nil
}

func ReadMetadata(filePath string) (map[string]string, error) {
	data := make(map[string]string)
	if HasSuffix(filePath, ImageFilters...) {
		exifs, err := exif.SearchFileAndExtractExif(filePath)
		if err != nil {
			return nil, err
		}
		tags, _, err := exif.GetFlatExifDataUniversalSearch(exifs, &exif.ScanOptions{}, true)
		if err != nil {
			return nil, err
		}
		for _, t := range tags {
			if t.TagName == "GPSLatitude" || t.TagName == "GPSLongitude" {
				vals := [3]float64{}
				fs := Split(t.Formatted, '[', ']', ' ')
				for i, f := range fs {
					r, _ := Eval(f)
					vals[i] = cast.ToFloat64(r)
				}
				data[t.TagName] = cast.ToString(vals[0] + vals[1]/60.0 + vals[2]/3600.0)
			} else {
				data[t.TagName] = t.FormattedFirst
			}
		}
	} else if HasSuffix(filePath, VideoFilters...) {
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
			dateUtc := time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC)
			loc, _ := time.LoadLocation("Local")
			data["CreationTime"] = dateUtc.Add(time.Duration(mvhd.GetCreationTime()) * time.Second).In(loc).Format("2006:01:02 15:04:05")
			data["ModificationTime"] = dateUtc.Add(time.Duration(mvhd.GetModificationTime()) * time.Second).Format("2006:01:02 15:04:05")
			data["Timescale"] = cast.ToString(mvhd.Timescale)
			data["Duration"] = cast.ToString(mvhd.GetDuration())
			data["Rate"] = cast.ToString(mvhd.Rate)
			data["Volume"] = cast.ToString(mvhd.Volume)
		}
	} else {
		return nil, fmt.Errorf("format not supported")
	}
	return data, nil
}

func WriteMetadata(srcFilePath, dstFilePath string, data map[string]any) error {
	if HasSuffix(srcFilePath, ImageFilters[0]) {
		context, err := jpegs.NewJpegMediaParser().ParseFile(srcFilePath)
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
		ifdBuilder, _ := exifBuilder.ChildWithTagId(exic.IfdExifStandardIfdIdentity.TagId())
		if ifdBuilder == nil {
			mapping, err := exic.NewIfdMappingWithStandard()
			if err != nil {
				return err
			}
			index := exif.NewTagIndex()
			err = exif.LoadStandardTags(index)
			if err != nil {
				return err
			}
			ifdBuilder = exif.NewIfdBuilder(mapping, index, exic.IfdExifStandardIfdIdentity, exic.EncodeDefaultByteOrder)
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
		dstFile, err := os.Create(dstFilePath)
		if err != nil {
			return err
		}
		defer dstFile.Close()
		err = segments.Write(dstFile)
		if err != nil {
			return err
		}
	} else if HasSuffix(srcFilePath, ImageFilters[1]) {
		context, err := pngs.NewPngMediaParser().ParseFile(srcFilePath)
		if err != nil {
			return err
		}
		chunks := context.(*pngs.ChunkSlice)
		exifBuilder, err := chunks.ConstructExifBuilder()
		if err != nil {
			return err
		}
		ifdBuilder, _ := exifBuilder.ChildWithTagId(exic.IfdExifStandardIfdIdentity.TagId())
		if ifdBuilder == nil {
			mapping, err := exic.NewIfdMappingWithStandard()
			if err != nil {
				return err
			}
			index := exif.NewTagIndex()
			err = exif.LoadStandardTags(index)
			if err != nil {
				return err
			}
			ifdBuilder = exif.NewIfdBuilder(mapping, index, exic.IfdExifStandardIfdIdentity, exic.EncodeDefaultByteOrder)
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
		dstFile, err := os.Create(dstFilePath)
		if err != nil {
			return err
		}
		defer dstFile.Close()
		err = chunks.WriteTo(dstFile)
		if err != nil {
			return err
		}
	} else if HasSuffix(srcFilePath, VideoFilters...) {
		srcFile, err := os.Open(srcFilePath)
		if err != nil {
			return err
		}
		defer srcFile.Close()
		dstFile, err := os.Create(dstFilePath)
		if err != nil {
			return err
		}
		defer dstFile.Close()
		writer := mp4.NewWriter(dstFile)
		_, err = mp4.ReadBoxStructure(srcFile, func(handle *mp4.ReadHandle) (any, error) {
			if handle.BoxInfo.Type == mp4.BoxTypeMdat() || !handle.BoxInfo.IsSupportedType() {
				return nil, writer.CopyBox(srcFile, &handle.BoxInfo)
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
	} else {
		return fmt.Errorf("format not supported")
	}
	return nil
}

func TakeStream(dstNames map[string]map[string]any, srcName string, srcArgs ...map[string]any) error {
	srcArg := make(map[string]any)
	if len(srcArgs) > 0 {
		srcArg = srcArgs[0]
	}
	if HasPrefix(srcName, "rtsp://") {
		srcArg["rtsp_transport"] = "tcp"
	}
	stream := ffmpeg.Input(srcName, srcArg)
	var streams []*ffmpeg.Stream
	for k, v := range dstNames {
		err := MakeDir(filepath.Dir(k))
		if err != nil {
			return err
		}
		streams = append(streams, stream.Output(k, v))
	}
	return ffmpeg.MergeOutputs(streams...).OverWriteOutput().ErrorToStdOut().Run()
}
