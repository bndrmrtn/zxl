package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/flarelang/flare/internal/models"
	"github.com/flarelang/flare/internal/modules/servermodule"
	"github.com/flarelang/flare/internal/runtimev2"
	"github.com/flarelang/flare/internal/state"
	"github.com/flarelang/flare/internal/version"
	"github.com/flarelang/flare/pkg/language"
	"github.com/fatih/color"
	"github.com/flarelang/webrouter"
)

// Server is a language server
type Server struct {
	// ir is the language interpreter
	ir *language.Interpreter

	// root is the root string
	root string
	// isDir is the directory flag
	isDir bool
	// rootFile is the root file if isDir is false
	rootFile string
	// colors is the color flag
	colors bool
	// dev is the development flag
	dev bool
	// router is the router for the server
	wr *webrouter.Router

	// useCaching is the caching flag
	useCaching bool
	// cache is the cache of the nodes
	cache map[string]*NodeCache
	// serverStateProvider is the provider for the server state
	serverStateProvider *state.Provider

	// mu is the mutex for the cache
	mu sync.RWMutex
}

// NewServer creates a new server
func New(ir *language.Interpreter, root string, isDir, cache, colors, dev bool) *Server {
	var (
		rootFile string
		router   *webrouter.Router
	)

	if !isDir {
		rootFile = filepath.Base(root)
		root = filepath.Dir(root)
	} else {
		router = webrouter.New(filepath.Clean(root))
		router.Reload()
	}

	return &Server{
		ir:                  ir,
		root:                filepath.Clean(root),
		isDir:               isDir,
		rootFile:            rootFile,
		colors:              colors,
		useCaching:          cache,
		cache:               make(map[string]*NodeCache),
		serverStateProvider: state.Default(),
		dev:                 dev,
		wr:                  router,
	}
}

// Serve starts the server
func (s *Server) Serve(addr string) error {
	blue := color.New(color.FgBlue, color.Bold)

	fmt.Printf("%s\n", blue.Sprint("Flare Web - ", version.Version))
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
	w.Header().Add("X-Flare-Version", version.Version)

	if s.dev {
		if err := s.wr.Reload(); err != nil {
			s.handleError(err, w, r)
		}
	}

	route, ok := s.wr.Match(r.URL.Path)
	if !ok {
		s.handleError(errNotFound, w, r)
		return
	}

	if !route.IsExecutable {
		http.ServeFile(w, r, route.FilePath)
		return
	}

	ctx := context.WithValue(r.Context(), "__params__", route.Params)
	r = r.WithContext(ctx)

	s.serveRequest(w, r, route, &cached)
}

func (s *Server) serveRequest(w http.ResponseWriter, r *http.Request, route *webrouter.Route, cached *bool) {
	// Execute the cached nodes if they exist
	if s.useCaching {
		if nodes, ok := s.getCache(route.FilePath); ok {
			*cached = true
			s.executeNodes(nodes, w, r)
			return
		}
	}

	// Open the file
	file, err := os.Open(route.FilePath)
	if err != nil {
		s.handleError(err, w, r)
		return
	}
	defer file.Close()

	// Get the nodes
	nodes, err := s.ir.GetNodes(route.FilePath, file)
	if err != nil {
		s.handleError(err, w, r)
		return
	}

	// Cache the nodes for faster execution
	if s.useCaching {
		s.setCache(route.FilePath, nodes)
	}

	// Execute the nodes
	s.executeNodes(nodes, w, r)
}

func (s *Server) executeNodes(nodes []*models.Node, w http.ResponseWriter, r *http.Request) {
	// Execute the nodes
	run, err := runtimev2.New(s.serverStateProvider)
	if err != nil {
		s.handleError(err, w, r)
		return
	}

	httpModule := servermodule.New(w, r)
	run.BindModule(httpModule)

	// Execute the nodes
	if _, err := run.Execute(nodes); err != nil {
		s.handleError(err, w, r)
		return
	}

	if !httpModule.Written {
		// Write the response
		w.WriteHeader(httpModule.Code)
		_, _ = w.Write(httpModule.Body.Bytes())
	}
}
