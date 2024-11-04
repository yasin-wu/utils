package zip

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexmullins/zip"
)

func TestZip(t *testing.T) {
	err := Compresses("/Users/yasin/Downloads/20241104", "./20241104.zip")
	t.Fatal(err)
}

func Zip(zipFile string, files ...string) error {
	if err := os.MkdirAll(filepath.Dir(zipFile), os.ModePerm); err != nil {
		return err
	}
	archive, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer func(archive *os.File) {
		_ = archive.Close()
	}(archive)
	zipWriter := zip.NewWriter(archive)
	defer func(zipWriter *zip.Writer) {
		_ = zipWriter.Close()
	}(zipWriter)
	for _, file := range files {
		file = strings.TrimSuffix(file, string(os.PathSeparator))
		err = filepath.Walk(file, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			header.Method = zip.Deflate
			header.Name, err = filepath.Rel(filepath.Dir(file), path)
			if err != nil {
				return err
			}
			if info.IsDir() {
				header.Name += string(os.PathSeparator)
			}
			headerWriter, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(headerWriter, f)
			return err
		})
		if err != nil {
			return err
		}
	}
	return nil
}
