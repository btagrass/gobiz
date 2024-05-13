package utl

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func CopyFile(inFilePath, outFilePath string) error {
	data, err := os.ReadFile(inFilePath)
	if err != nil {
		return err
	}
	return os.WriteFile(outFilePath, data, os.ModePerm)
}

func Exist(name string) bool {
	matches, err := filepath.Glob(name)
	if err != nil {
		return false
	}
	return len(matches) > 0
}

func Glob(patterns ...string) []string {
	var files []string
	for _, p := range patterns {
		matches, err := filepath.Glob(p)
		if err != nil {
			logrus.Error(err)
			continue
		}
		files = append(files, matches...)
	}
	return files
}

func MakeDir(names ...string) error {
	for _, name := range names {
		if !Exist(name) {
			err := os.MkdirAll(name, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Remove(names ...string) error {
	for _, name := range names {
		matches, err := filepath.Glob(name)
		if err != nil {
			return err
		}
		for _, m := range matches {
			err = os.Remove(m)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func RenameFile(inFilePath, outFilePath string) error {
	err := MakeDir(filepath.Dir(outFilePath))
	if err != nil {
		return err
	}
	err = os.Rename(inFilePath, outFilePath)
	return err
}

func UnzipFile(inFilePath string, outFilePath ...string) error {
	var outDir string
	if len(outFilePath) > 0 {
		outDir = outFilePath[0]
	}
	reader, err := zip.OpenReader(inFilePath)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, f := range reader.File {
		filePath := filepath.Join(outDir, f.Name)
		if f.FileInfo().IsDir() {
			err = MakeDir(filePath)
			if err != nil {
				return err
			}
			continue
		}
		outFile, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer outFile.Close()
		inFile, err := f.Open()
		if err != nil {
			return err
		}
		defer inFile.Close()
		_, err = io.Copy(outFile, inFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func ZipFile(inFilePath string, outFilePath ...string) error {
	var outDir string
	if len(outFilePath) > 0 {
		outDir = outFilePath[0]
	}
	filePath := filepath.Join(outDir, fmt.Sprintf("%s.zip", inFilePath))
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	writer := zip.NewWriter(outFile)
	defer writer.Close()
	filepath.Walk(inFilePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fileHeader, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		fileHeader.Name = path
		if !info.IsDir() {
			fileHeader.Method = zip.Deflate
		}
		outFile, err := writer.CreateHeader(fileHeader)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			inFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer inFile.Close()
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}
