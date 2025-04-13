package cache

import (
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bndrmrtn/flare/internal/models"
)

// CacheDirectory is the directory where cache files are stored
const CacheDirectory = ".flcache"

// Cache is the cache information
type Cache struct {
	// Hash is the hash of the data
	Hash string
	// Nodes is the nodes
	Nodes []*models.Node
}

// Store stores the cache information to file
func Store(file string, data []byte, nodes []*models.Node) {
	fileName := strings.ReplaceAll(strings.ReplaceAll(file, "\\", "$"), "/", "$") + "bin"

	if info, err := os.Stat(CacheDirectory); err != nil || !info.IsDir() {
		// Create cache directory if it does not exist
		_ = os.MkdirAll(CacheDirectory, os.ModePerm)
	}

	hash := fmt.Sprintf("%x", md5.Sum(data))

	cacheData := Cache{
		Hash:  hash,
		Nodes: nodes,
	}

	f, err := os.Create(filepath.Join(CacheDirectory, fileName))
	if err != nil {
		fmt.Println("Error creating cache file:", err)
		return
	}
	defer f.Close()

	if err := gob.NewEncoder(f).Encode(cacheData); err != nil {
		fmt.Println("Error encoding cache data to file:", err)
	}
}

// Get gets the cache information from file
func Get(file string, data []byte) ([]*models.Node, bool) {
	fileName := strings.ReplaceAll(strings.ReplaceAll(file, "\\", "$"), "/", "$") + "bin"
	hash := fmt.Sprintf("%x", md5.Sum(data))

	f, err := os.Open(filepath.Join(CacheDirectory, fileName))
	if err != nil {
		return nil, false
	}
	defer f.Close()

	var cacheData Cache
	if err := gob.NewDecoder(f).Decode(&cacheData); err != nil {
		return nil, false
	}

	if cacheData.Hash != hash {
		return nil, false
	}

	return cacheData.Nodes, true
}
