package build

import (
	"io/fs"
	"time"
)

type FileInfo struct {
	EName    string
	ESize    int64
	EModTime time.Time
	EIsDir   bool
}

func NewFileInfo(name string, size int64, modTime time.Time, isDir bool) *FileInfo {
	return &FileInfo{
		EName:    name,
		ESize:    size,
		EModTime: modTime,
		EIsDir:   isDir,
	}
}
func (f *FileInfo) Name() string {
	return f.EName
}

func (f *FileInfo) Size() int64 {
	return f.ESize
}

func (f *FileInfo) Mode() fs.FileMode {
	return 0444
}

func (f *FileInfo) ModTime() time.Time {
	return f.EModTime
}

func (f *FileInfo) IsDir() bool {
	return f.EIsDir
}

func (f *FileInfo) Sys() interface{} {
	return nil
}
