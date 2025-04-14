package bundler

import (
	"encoding/gob"
	"os"
	"path/filepath"
)

// Idea - bundle flare app to a single flareb (Flare Bundle) file that can be imported with flare
// or can be run by executing the main function without any argument
// flare make:bundle --as=bundle.flb
// bundle.flb - flare run:bundle bundle.flb
// bundle.fb - flare run bundle.flb

type Bundler struct {
	root string
}

func New(root string) *Bundler {
	return &Bundler{
		root: filepath.Clean(root),
	}
}

func (b *Bundler) Bundle(out string) error {
	bundle, err := NewBundle(b.root)
	if err != nil {
		return err
	}

	file, err := os.Create(out + ".flb")
	if err != nil {
		return err
	}
	defer file.Close()

	return gob.NewEncoder(file).Encode(bundle)
}

func (b *Bundler) Unpack() (*Bundle, error) {
	file, err := os.Open(b.root)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var bundle Bundle
	if err := gob.NewDecoder(file).Decode(&bundle); err != nil {
		return nil, err
	}

	return &bundle, nil
}
