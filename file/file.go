package file

import "bytes"

type File struct {
	name   string
	buffer *bytes.Buffer
}

func New(name string) *File {
	return &File{
		name:   name,
		buffer: &bytes.Buffer{},
	}
}

func (f *File) Write(chunk []byte) error {
	_, err := f.buffer.Write(chunk)
	return err
}

func (f *File) ReadAll() []byte {
	return f.buffer.Bytes()
}

func (f *File) Size() int64 {
	return int64(f.buffer.Len())
}
