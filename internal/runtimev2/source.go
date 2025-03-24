package runtimev2

import (
	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/source"
	"go.uber.org/zap"
)

type Namespace struct {
	Name   string
	Nodes  []*models.Node
	Loaded bool
}

func (r *Runtime) LoadSourceNamespaces() error {
	zap.L().Info("loading source namespaces")

	namespaces, err := source.Get()
	if err != nil {
		return err
	}

	var m = make(map[string]*Namespace, len(namespaces))

	for _, nodes := range namespaces {
		if len(nodes) > 0 {
			ns := nodes[0].Content
			m[ns] = &Namespace{
				Name:   ns,
				Nodes:  nodes,
				Loaded: false, // not loaded yet
			}
		}
	}

	r.mu.Lock()
	r.sourceNamespaces = m
	r.mu.Unlock()
	return nil
}
