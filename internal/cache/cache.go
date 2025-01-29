package cache

import (
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"os"

	"github.com/bndrmrtn/zexlang/internal/models"
)

// Store stores the cache information to file
func Store(data []byte, nodes []*models.Node) {
	if info, err := os.Stat(".zxcache/"); err != nil || !info.IsDir() {
		// Create cache directory if it does not exist
		_ = os.MkdirAll(".zxcache/", os.ModePerm)
	}

	hash := fmt.Sprintf("%x", md5.Sum(data))

	// Write cache information to file
	f, err := os.Create(".zxcache/" + hash + ".zxbin")
	if err != nil {
		return
	}
	defer f.Close()

	// Write cache information to file
	_ = gob.NewEncoder(f).Encode(nodes)
}

// Get gets the cache information from file
func Get(data []byte) ([]*models.Node, bool) {
	// Create cache directory if it does not exist
	_ = os.MkdirAll(".zxcache/", os.ModePerm)

	hash := fmt.Sprintf("%x", md5.Sum(data))

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
