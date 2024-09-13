package filechunk

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
)

type Chunk struct {
	tmpPath string
}

func New(filepath string) *Chunk {
	return &Chunk{tmpPath: filepath}
}

func (c *Chunk) Write(index int64, hash string, chunk []byte) error {
	chunkPath := path.Join(c.tmpPath, fmt.Sprintf("%04d", index)+"_"+hash)
	if _, err := os.Stat(chunkPath); err == nil {
		return nil
	}
	f, err := os.OpenFile(chunkPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer func(f *os.File) { _ = f.Close() }(f)
	if _, err = f.Write(chunk); err != nil {
		return err
	}
	return nil
}

func (c *Chunk) Merge(filePath, fileHash string) (int64, error) {
	defer func() {
		_ = os.RemoveAll(path.Join(c.tmpPath))
		parentPath := path.Join(c.tmpPath, "../")
		if dir, _ := os.ReadDir(parentPath); len(dir) == 0 {
			_ = os.RemoveAll(parentPath)
		}
	}()
	files, err := os.ReadDir(c.tmpPath)
	if err != nil {
		return 0, err
	}
	completeFile, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(completeFile)
	var size int64
	for _, f := range files {
		if f.Name() == ".DS_Store" {
			continue
		}
		buf, err := os.ReadFile(path.Join(c.tmpPath, f.Name()))
		if err != nil {
			return 0, err
		}
		if _, err := completeFile.Write(buf); err != nil {
			return 0, err
		}
		size += int64(len(buf))
	}
	fileMd5, err := c.fileMD5(filePath)
	if err != nil {
		return 0, err
	}
	if fileMd5 != fileHash {
		return 0, fmt.Errorf("file md5 mismatch: %s not equal %s", fileMd5, fileHash)
	}
	return size, nil
}

func (c *Chunk) fileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(file)
	h := md5.New()
	if _, err = io.Copy(h, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
