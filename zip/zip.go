package zip

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexmullins/zip"
)

func Compresses(source, destination string, password ...string) error {
	if err := os.MkdirAll(filepath.Dir(destination), os.ModePerm); err != nil {
		return err
	}
	zipFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	return filepath.Walk(source, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		relativePath := strings.TrimPrefix(path, filepath.Dir(source)+"/")
		header.Name = relativePath
		header.Method = zip.Deflate
		header.Flags = 0x800
		if len(password) > 0 {
			header.SetPassword(password[0])
		}
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := io.Copy(writer, file); err != nil {
			return err
		}
		return nil
	})
}
