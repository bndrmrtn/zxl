package language

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bndrmrtn/zexlang/internal/errs"
	"github.com/bndrmrtn/zexlang/internal/runtime"
	"github.com/bndrmrtn/zexlang/internal/version"
	"github.com/fatih/color"
)

// Server is a language server
type Server struct {
	// ir is the language interpreter
	ir *Interpreter
	// root is the root string
	root string
	// colors is the color flag
	colors bool
}

// NewServer creates a new server
func NewServer(ir *Interpreter, root string, colors bool) *Server {
	return &Server{
		ir:     ir,
		root:   filepath.Clean(root),
		colors: colors,
	}
}

// Serve starts the server
func (s *Server) Serve(addr string) error {
	blue := color.New(color.FgBlue, color.Bold)

	fmt.Printf("%s\n", blue.Sprint("Zex Web - ", version.Version))
	color.New(color.FgYellow).Printf("Server listening on %s\n", addr)
	color.New(color.FgRed).Printf("Press Ctrl+C to quit\n")

	return http.ListenAndServe(addr, s)
}

// ServeHTTP serves the http request
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	color.NoColor = true
	defer func() {
		color.NoColor = !s.colors
	}()

	w.Header().Add("X-Zex-Version", version.Version)

	filePath := filepath.Join(s.root, filepath.Clean(r.URL.Path[1:])) + ".zx"

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		dirPath := filepath.Join(s.root, filepath.Clean(r.URL.Path[1:]))
		indexFilePath := filepath.Join(dirPath, "index.zx")

		_, err := os.Stat(indexFilePath)
		if os.IsNotExist(err) {
			http.FileServer(http.Dir(s.root)).ServeHTTP(w, r)
			return
		} else if err != nil {
			s.handleError(err, w)
			return
		}
		filePath = indexFilePath
	} else if err != nil {
		s.handleError(err, w)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		s.handleError(err, w)
		return
	}
	defer file.Close()

	nodes, err := s.ir.GetNodes(filePath, file)
	if err != nil {
		s.handleError(err, w)
		return
	}

	run := runtime.New(runtime.EntryPoint)

	httpModule := runtime.NewHttpModule(w, r)
	htmlModule := runtime.NewHtmlModule(httpModule)

	run.BindModule("http", httpModule)
	run.BindModule("html", htmlModule)

	if _, err = run.Execute(nodes); err != nil {
		s.handleError(err, w)
		return
	}

	w.WriteHeader(httpModule.StatusCode)
	w.Write(httpModule.Body.Bytes())
}

func (s *Server) handleError(err error, w http.ResponseWriter) {
	var de errs.DebugError
	if errors.As(err, &de) {
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(de.HttpError()))
		return
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
}
