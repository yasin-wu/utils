package file

import "bytes"

type File struct {
	name     string
	path     string
	fileType string
	content  string
	size     int64
	buffer   *bytes.Buffer
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

func (f *File) Name() string {
	return f.name
}

func (f *File) Path() string {
	return f.path
}

func (f *File) Type() string {
	return f.fileType
}

func (f *File) Content() string {
	return f.content
}

func (f *File) Size() int64 {
	if f.size >= int64(f.buffer.Len()) {
		return f.size
	}
	return int64(f.buffer.Len())
}
