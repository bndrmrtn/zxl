package language

import (
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"io"
	"os"

	"github.com/bndrmrtn/zexlang/internal/models"
)

// storeCache stores the cache information to file
func (ir *Interpreter) storeCache(data io.Reader, nodes []*models.Node) {
	if !ir.cache {
		return
	}

	if info, err := os.Stat(".zxcache/"); err != nil || !info.IsDir() {
		// Create cache directory if it does not exist
		_ = os.MkdirAll(".zxcache/", os.ModePerm)
	}

	b, err := io.ReadAll(data)
	if err != nil {
		return
	}

	hash := fmt.Sprintf("%x", md5.Sum(b))

	// Write cache information to file
	f, err := os.Create(".zxcache/" + hash + ".zxbin")
	if err != nil {
		return
	}
	defer f.Close()

	// Write cache information to file
	_ = gob.NewEncoder(f).Encode(nodes)
}

// getCache gets the cache information from file
func (ir *Interpreter) getCache(data io.Reader) ([]*models.Node, bool) {
	if !ir.cache {
		return nil, false
	}

	// Create cache directory if it does not exist
	_ = os.MkdirAll(".zxcache/", os.ModePerm)

	b, err := io.ReadAll(data)
	if err != nil {
		return nil, false
	}

	hash := fmt.Sprintf("%x", md5.Sum(b))

	// Read cache information from file
	f, err := os.Open(".zxcache/" + hash + ".zxbin")
	if err != nil {
		return nil, false
	}
	defer f.Close()

	var nodes []*models.Node
	_ = gob.NewDecoder(f).Decode(&nodes)

	return nodes, true
}
