package vm

import (
	"os"
	"strings"

	"github.com/jaronnie/gvm/internal/global"
)

type LocalVM struct{}

func NewLocalVM() Interface {
	return &LocalVM{}
}

func (r *LocalVM) List() ([]string, error) {
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
