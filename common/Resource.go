package common

import (
	"embed"
	"errors"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
)

type Resource struct {
	fs   embed.FS
	path string
}

func NewResource(fs embed.FS, path string) *Resource {
	return &Resource{
		fs:   fs,
		path: path,
	}
}

func (r *Resource) Open(name string) (fs.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	p1 := path.Clean("/" + name)
	p2 := path.Join(r.path, p1)
	f, e := r.fs.Open(p2)
	return f, e
}
