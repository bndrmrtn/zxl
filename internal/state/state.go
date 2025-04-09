package state

import (
	"sync"

	"github.com/bndrmrtn/zxl/lang"
)

// State represents a stateful object that can be used to store and retrieve data
// that is global and can be accessed by multiple threads.
type State interface {
	// Get retrieves the value associated with the given key.
	Get(key string) (lang.Object, bool)
	// Set sets the value associated with the given key.
	Set(key string, value lang.Object) bool
}

type DefaultState struct {
	m map[string]lang.Object

	mu sync.RWMutex
}

func NewDefaultState() *DefaultState {
	return &DefaultState{
		m: make(map[string]lang.Object),
	}
}

func (s *DefaultState) Get(key string) (lang.Object, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.m[key]
	return value, ok
}

func (s *DefaultState) Set(key string, value lang.Object) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.m[key] = value
	return true
}
