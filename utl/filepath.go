package utl

import (
	"archive/zip"
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CopyFile(inFilePath, outFilePath string) error {
	err := MakeDir(filepath.Dir(outFilePath))
	if err != nil {
		return err
	}
	data, err := os.ReadFile(inFilePath)
	if err != nil {
		return err
	}
	err = os.WriteFile(outFilePath, data, os.ModePerm)
	if err != nil {
		return err
	}
	inFileInfo, err := os.Stat(inFilePath)
	if err != nil {
		return err
	}
	return os.Chtimes(outFilePath, time.Now(), inFileInfo.ModTime())
}

func Exist(name string) bool {
	matches, err := filepath.Glob(name)
	if err != nil {
		return false
	}
	return len(matches) > 0
}

func Glob(path string, patterns ...string) []string {
	var files []string
	filepath.WalkDir(path, func(path string, _ fs.DirEntry, _ error) error {
		for _, p := range patterns {
			matches, _ := filepath.Glob(fmt.Sprintf("%s/%s", path, p))
			files = append(files, matches...)
			matches, _ = filepath.Glob(fmt.Sprintf("%s/%s", path, strings.ToUpper(p)))
			files = append(files, matches...)
		}
		return nil
	})
	return files
}

func HashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
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
