package build

import (
	"io"
	"io/fs"
)

type File struct {
	Name string
	Data []byte
	Info *FileInfo
}

func (f *File) Readdir(count int) ([]fs.FileInfo, error) {
	return nil, nil
}

func (f *File) Stat() (fs.FileInfo, error) {
	return f.Info, nil
}

func (f *File) Read(p []byte) (n int, err error) {
	return copy(p, f.Data), io.EOF
}

func (f *File) Close() error {
	return nil
}
