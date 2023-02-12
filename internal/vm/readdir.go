package vm

import (
	"os"
	"strings"

	"github.com/jaronnie/gvm/internal/global"
)

type ReadDirVM struct{}

func NewReadDirVM() Interface {
	return &ReadDirVM{}
}

func (r *ReadDirVM) List() ([]string, error) {
	dir, err := os.ReadDir(global.GVM_CONFIG_DIR)
	if err != nil {
		return nil, err
	}

	var list []string

	for _, v := range dir {
		if v.IsDir() && strings.HasPrefix(v.Name(), "go") {
			list = append(list, v.Name())
		}
	}

	return list, nil
}
