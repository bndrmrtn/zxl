package modules

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bndrmrtn/zxl/internal/lang"
)

type IO struct{}

func NewIOModule() *IO {
	return &IO{}
}

func (*IO) Namespace() string {
	return "io"
}

func (*IO) Objects() map[string]lang.Object {
	return nil
}

func (*IO) Methods() map[string]lang.Method {
	return map[string]lang.Method{
		"open": lang.NewFunction([]string{"path"}, fnReadFile, nil),
	}
}

func fnReadFile(args []lang.Object) (lang.Object, error) {
	if args[0].Type() != lang.TString {
		return nil, fmt.Errorf("expected string, got %s", args[0].Type())
	}

	path := args[0]
	pathString := filepath.Clean(filepath.Join(filepath.Dir(path.Debug().File), path.Value().(string)))

	reader, err := os.Open(pathString)
	if err != nil {
		return nil, err
	}

	return lang.NewIOStream("file", reader), nil
}
