package utl

import (
	"archive/zip"
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

func CopyFile(srcFilePath, dstFilePath string) error {
	err := MakeDir(filepath.Dir(dstFilePath))
	if err != nil {
		return err
	}
	data, err := os.ReadFile(srcFilePath)
	if err != nil {
		return err
	}
	err = os.WriteFile(dstFilePath, data, os.ModePerm)
	if err != nil {
		return err
	}
	srcFileInfo, err := os.Stat(srcFilePath)
	if err != nil {
		return err
	}
	return os.Chtimes(dstFilePath, time.Now(), srcFileInfo.ModTime())
}

func Exist(name string) bool {
	matches, err := filepath.Glob(name)
	if err != nil {
		return false
	}
	return len(matches) > 0
}

func GlobDir(name string, patterns ...string) []string {
	var files []string
	filepath.WalkDir(name, func(path string, _ fs.DirEntry, _ error) error {
		for _, p := range patterns {
			matches, _ := filepath.Glob(fmt.Sprintf("%s/%s", path, p))
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

func MakeDir(dirs ...string) error {
	for _, d := range dirs {
		if Exist(d) {
			continue
		}
		err := os.MkdirAll(d, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func Remove(names ...string) error {
	for _, n := range names {
		matches, err := filepath.Glob(n)
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

func RenameFile(srcFilePath, dstFilePath string) error {
	err := MakeDir(filepath.Dir(dstFilePath))
	if err != nil {
		return err
	}
	err = os.Rename(srcFilePath, dstFilePath)
	return err
}

func UnzipFile(srcFilePath string, dstFilePath ...string) error {
	var dstDir string
	if len(dstFilePath) > 0 {
		dstDir = dstFilePath[0]
	}
	srcReader, err := zip.OpenReader(srcFilePath)
	if err != nil {
		return err
	}
	defer srcReader.Close()
	for _, f := range srcReader.File {
		filePath := filepath.Join(dstDir, f.Name)
		if f.FileInfo().IsDir() {
			err = MakeDir(filePath)
			if err != nil {
				return err
			}
			continue
		}
		dstFile, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer dstFile.Close()
		srcFile, err := f.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()
		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func ZipFile(srcFilePath string, dstFilePath ...string) error {
	var dstDir string
	if len(dstFilePath) > 0 {
		dstDir = dstFilePath[0]
	}
	filePath := filepath.Join(dstDir, fmt.Sprintf("%s.zip", srcFilePath))
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	writer := zip.NewWriter(outFile)
	defer writer.Close()
	filepath.Walk(srcFilePath, func(path string, info os.FileInfo, err error) error {
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
		dstFile, err := writer.CreateHeader(fileHeader)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			srcFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer srcFile.Close()
			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}
