package build

import (
	"fmt"
	"io"
	"io/fs"
	"path"
)

type Fs struct {
	Files []*File
}

// An openFile is a regular file open for reading.
type openFile struct {
	f      *File // the file itself
	offset int64 // current read offset
}

func (f *openFile) Close() error               { return nil }
func (f *openFile) Stat() (fs.FileInfo, error) { return f.f.Info, nil }

func (f *openFile) Read(b []byte) (int, error) {
	if f.offset >= f.f.Info.Size() {
		return 0, io.EOF
	}
	if f.offset < 0 {
		return 0, &fs.PathError{Op: "read", Path: f.f.Name, Err: fs.ErrInvalid}
	}
	n := copy(b, f.f.Data[f.offset:])
	f.offset += int64(n)
	return n, nil
}

func (f *openFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0:
		// offset += 0
	case 1:
		offset += f.offset
	case 2:
		offset += f.f.Info.Size()
	}
	if offset < 0 || offset > f.f.Info.Size() {
		return 0, &fs.PathError{Op: "seek", Path: f.f.Name, Err: fs.ErrInvalid}
	}
	f.offset = offset
	return offset, nil
}

func (f *Fs) Open(name string) (fs.File, error) {
	name = path.Clean("/" + name)
	for _, file := range f.Files {
		if file.Name == name {
			return &openFile{file, 0}, nil
		}
	}
	return nil, fmt.Errorf("file not found")
}
