package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/internal/modules"
	"github.com/bndrmrtn/zxl/internal/runtimev2"
	"github.com/bndrmrtn/zxl/internal/version"
	"github.com/bndrmrtn/zxl/pkg/language"
	"github.com/fatih/color"
)

// Server is a language server
type Server struct {
	// ir is the language interpreter
	ir *language.Interpreter

	// root is the root string
	root string
	// colors is the color flag
	colors bool

	// useCaching is the caching flag
	useCaching bool
	// cache is the cache of the nodes
	cache map[string]*NodeCache
	// mu is the mutex for the cache
	mu sync.RWMutex
}

// NewServer creates a new server
func New(ir *language.Interpreter, root string, cache, colors bool) *Server {
	return &Server{
		ir:         ir,
		root:       filepath.Clean(root),
		colors:     colors,
		useCaching: cache,
		cache:      make(map[string]*NodeCache),
	}
}

// Serve starts the server
func (s *Server) Serve(addr string) error {
	blue := color.New(color.FgBlue, color.Bold)

	fmt.Printf("%s\n", blue.Sprint("Zx Web - ", version.Version))
	color.New(color.FgYellow).Printf("Server listening on %s\n", addr)
	color.New(color.FgRed).Printf("Press Ctrl+C to quit\n\n")

	return http.ListenAndServe(addr, s)
}

// ServeHTTP serves the http request
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Disable colors for http requests
	color.NoColor = true
	cached := false

	// Log the request
	start := time.Now()
	defer func() {
		color.NoColor = !s.colors
		logServerRequest(start, r.Method, r.URL.Path, cached)
	}()

	// Set the version header
	w.Header().Add("X-Zx-Version", version.Version)

	// Serve files if they exist
	path := filepath.Join(s.root, r.URL.Path[1:])
	if filepath.Ext(path) != ".zx" {
		if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
			http.ServeFile(w, r, path)
			return
		}
	}

	// Execute the cached nodes if they exist
	if s.useCaching {
		if nodes, ok := s.getCache(path); ok {
			cached = true
			s.executeNodes(nodes, w, r)
			return
		}
	}

	// Get the executable path
	zxPath, err := s.getExecutablePath(path)
	if err != nil {
		s.handleError(err, w, r)
		return
	}

	// Open the file
	file, err := os.Open(zxPath)
	if err != nil {
		s.handleError(err, w, r)
		return
	}
	defer file.Close()

	// Get the nodes
	nodes, err := s.ir.GetNodes(zxPath, file)
	if err != nil {
		s.handleError(err, w, r)
		return
	}

	// Cache the nodes for faster execution
	if s.useCaching {
		s.setCache(path, nodes)
	}

	// Execute the nodes
	s.executeNodes(nodes, w, r)
}

func (s *Server) executeNodes(nodes []*models.Node, w http.ResponseWriter, r *http.Request) {
	// Execute the nodes
	run := runtimev2.New()
	httpModule := modules.NewHttpServerModule(w, r)
	run.BindModule(httpModule)
	err := s.ir.ExecuteSourceFiles(run)
	if err != nil {
		s.handleError(err, w, r)
		return
	}

	// Execute the nodes
	if _, err := run.Execute(nodes); err != nil {
		s.handleError(err, w, r)
		return
	}

	// Write the response
	w.WriteHeader(httpModule.Code)
	w.Write(httpModule.Body.Bytes())
}
