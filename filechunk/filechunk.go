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

func New(tmpPath string) *Chunk {
	return &Chunk{tmpPath: tmpPath}
}

// Write 写入文件块
func (c *Chunk) Write(index int64, md5 string, chunk []byte) error {
	chunkPath := path.Join(c.tmpPath, fmt.Sprintf("%04d", index)+"_"+md5)
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

// Merge 合并文件块
func (c *Chunk) Merge(distPath, fileMD5 string) (int64, error) {
	defer c.removeAll()
	files, err := os.ReadDir(c.tmpPath)
	if err != nil {
		return 0, err
	}
	completeFile, err := os.Create(distPath)
	if err != nil {
		return 0, err
	}
	defer completeFile.Close()
	var size int64
	h := md5.New()
	buffer := make([]byte, 1024*1024) // 1M缓冲区
	for _, f := range files {
		if f.Name() == ".DS_Store" {
			continue
		}
		tmpFile, err := os.Open(path.Join(c.tmpPath, f.Name()))
		if err != nil {
			return 0, err
		}
		for {
			n, err := tmpFile.Read(buffer)
			if err != nil && err != io.EOF {
				_ = tmpFile.Close()
				return 0, err
			}
			if n == 0 {
				break
			}
			if _, err = completeFile.Write(buffer[:n]); err != nil {
				_ = tmpFile.Close()
				return 0, err
			}
			if _, err = h.Write(buffer[:n]); err != nil {
				_ = tmpFile.Close()
				return 0, err
			}
			size += int64(n)
		}
		_ = tmpFile.Close()
	}
	md5Hex := hex.EncodeToString(h.Sum(nil))
	if md5Hex != fileMD5 {
		return 0, fmt.Errorf("file md5 mismatch: %s not equal %s", md5Hex, fileMD5)
	}
	return size, nil
}

func (c *Chunk) removeAll() {
	_ = os.RemoveAll(path.Join(c.tmpPath))
	parentPath := path.Join(c.tmpPath, "../")
	if dir, _ := os.ReadDir(parentPath); len(dir) == 0 {
		_ = os.RemoveAll(parentPath)
	}
}
