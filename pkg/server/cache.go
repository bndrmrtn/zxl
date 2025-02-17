package server

import (
	"time"

	"github.com/bndrmrtn/zxl/internal/models"
)

// CacheDuration is the duration of the cache
const CacheDuration = time.Minute * 5

// NodeCache is the cache of the nodes
type NodeCache struct {
	Expiration time.Time
	Nodes      []*models.Node
}

// getCache gets the cache
func (s *Server) getCache(path string) ([]*models.Node, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cache, ok := s.cache[path]
	if !ok {
		return nil, false
	}

	if time.Now().After(cache.Expiration) {
		return nil, false
	}

	return cache.Nodes, true
}

// setCache sets the cache
func (s *Server) setCache(path string, nodes []*models.Node) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache[path] = &NodeCache{
		Expiration: time.Now().Add(CacheDuration),
		Nodes:      nodes,
	}
}
