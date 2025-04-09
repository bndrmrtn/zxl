package state

import "sync"

// Provider represents a state provider.
type Provider struct {
	States map[string]State

	mu         sync.RWMutex
	stateMaker func() State
}

// NewProvider creates a new state provider.
func NewProvider(stateMaker func() State) *Provider {
	return &Provider{
		States:     make(map[string]State),
		stateMaker: stateMaker,
	}
}

// State returns a state with the given name.
func (p *Provider) State(name string) State {
	p.mu.RLock()
	state, ok := p.States[name]
	p.mu.RUnlock()

	if ok {
		return state
	}

	s := p.stateMaker()

	p.mu.Lock()
	p.States[name] = s
	p.mu.Unlock()

	return s
}
