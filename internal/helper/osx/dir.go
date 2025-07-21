package osx

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
)

func IsDir(file string) (bool, error) {
	fi, err := os.Stat(file)
	if err != nil {
		return false, err
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return true, err
	case mode.IsRegular():
		return false, err
	}

	return false, errors.New("unknow file type")
}

func IsFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func ReadDirNames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, err
}

func GetAllFiles(pathname string) ([]string, error) {
	s := make([]string, 0)
	rd, err := os.ReadDir(pathname)
	if err != nil {
		return s, err
	}
	for _, fi := range rd {
		if !fi.IsDir() {
			fullName := path.Join(pathname, fi.Name())
			s = append(s, fullName)
		}
	}

	return s, nil
}

func CreateDir(dirpath string) error {
	err := os.MkdirAll(dirpath, 0o766)
	return err
}

func CompressDir(srcDir, destDir string) error {
	zipFile, err := os.Create(destDir)
	if err != nil {
		return err
	}

	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 设置zip头
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// 设置文件名 zip archive ?
		header.Name, err = filepath.Rel(srcDir, path)

		// 如果文件是一个目录，添加它到zip archive
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Store
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}

			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
