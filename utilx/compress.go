package utilx

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

// Untargz 解压 tar.gz
// @archive : the compress file's path , such as - "./test.tar.gz"
// @dest    : the file Untar dest, such as "./ttt", and you will get fold named "test" under the ttt
func Untargz(archive, dest string) error {
	if exist, _ := PathExists(archive); !exist {
		return fmt.Errorf("archive not found: %s", archive)
	}

	// 如果不存在，则新建dest
	if e, _ := PathExists(dest); !e {
		err := os.MkdirAll(dest, 0755)
		if err != nil {
			return err
		}
	}

	srcFile, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := path.Join(dest, hdr.Name) // dest + hdr.Name

		if hdr.FileInfo().IsDir() {
			// create path before create file in <create> func, continue here
			continue
		}

		file, err := create(filename)
		if err != nil {
			return err
		}
		if _, err = io.Copy(file, tr); err != nil {
			_ = file.Close()
			return err
		}
		_ = os.Chmod(filename, hdr.FileInfo().Mode())
		_ = file.Close()
	}
	return nil
}

func create(name string) (*os.File, error) {
	dir, _ := filepath.Split(name)
	// create dir before create file
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}
