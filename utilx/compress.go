package utilx

import (
	"archive/tar"
	"archive/zip"
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
		err := os.MkdirAll(dest, 0o755)
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
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}

// Unzip extracts a zip archive to the destination directory
// @archive : the zip file's path, such as - "./test.zip"
// @dest    : the file unzip dest, such as "./ttt"
func Unzip(archive, dest string) error {
	if exist, _ := PathExists(archive); !exist {
		return fmt.Errorf("archive not found: %s", archive)
	}

	// create dest if not exists
	if e, _ := PathExists(dest); !e {
		err := os.MkdirAll(dest, 0o755)
		if err != nil {
			return err
		}
	}

	r, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip vulnerability
		if !filepath.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Create parent directory if needed
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
