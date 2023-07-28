package file

import (
	"os"
	"path"
)

type Storage struct {
	dir string
}

func NewStorage(dir string) (*Storage, error) {
	_, err := os.Stat(dir)
	if err != nil && !os.IsExist(err) {
		if err = os.MkdirAll(dir, 0777); err != nil {
			return nil, err
		}
	}
	return &Storage{
		dir: dir,
	}, nil
}

func (s *Storage) Store(file *File) error {
	return os.WriteFile(path.Join(s.dir, file.name), file.buffer.Bytes(), 0644)
}
