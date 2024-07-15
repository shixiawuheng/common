package build

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func SaveDat(root string, outfile string) error {
	dir, err := buildDirTree(root)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(dir)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(outfile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	// 创建gzip解压器
	gw := gzip.NewWriter(file)
	if err != nil {
		return err
	}
	defer gw.Close()
	// 写入压缩后的数据
	_, err = gw.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func buildDirTree(root string) (fs.FS, error) {
	dir := Fs{
		Files: make([]*File, 0),
	}
	root = filepath.Clean(root)
	err := filepath.Walk(root, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			return nil
		}
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		file := File{
			Name: strings.ReplaceAll(strings.TrimPrefix(filePath, root), "\\", "/"),
			Data: data,
			Info: NewFileInfo(fileInfo.Name(), fileInfo.Size(), fileInfo.ModTime(), fileInfo.IsDir()),
		}
		dir.Files = append(dir.Files, &file)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &dir, nil
}

func LoadData(file string) (fs.FS, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 创建gzip解压器
	gr, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gr.Close()
	var packed Fs
	dec := gob.NewDecoder(gr)
	err = dec.Decode(&packed)
	if err != nil {
		return nil, err
	}
	return &packed, nil
}

func LoadDataIo(read io.Reader) (fs.FS, error) {
	// 创建gzip解压器
	gr, err := gzip.NewReader(read)
	if err != nil {
		return nil, err
	}
	defer gr.Close()
	var packed Fs
	dec := gob.NewDecoder(gr)
	err = dec.Decode(&packed)
	if err != nil {
		return nil, err
	}
	return &packed, nil
}

func ViewData(file string) ([]*File, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// 创建gzip解压器
	gr, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gr.Close()
	var packed Fs
	dec := gob.NewDecoder(gr)
	err = dec.Decode(&packed)
	if err != nil {
		return nil, err
	}
	return packed.Files, nil
}

func SaveZip(root, filep string) error {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		//relPath := encoder.String()
		relPath := strings.TrimPrefix(path, root+string(os.PathSeparator))
		if relPath == path {
			return nil
		}
		//fmt.Println("压缩文件：", relPath)
		//relPath, err = encoder.String(relPath)
		if err != nil {
			return err
		}
		if info.IsDir() {
			// 如果是目录，则在压缩文件中创建一个目录条目
			_, err = zipWriter.Create(relPath + "/")
			if err != nil {
				return err
			}

			return nil
		}
		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(zipEntry, file)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	err = zipWriter.Close()
	return ioutil.WriteFile(filep, buf.Bytes(), 0644)
}
