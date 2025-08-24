package vbook

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

type LocalServer struct {
	server     *http.Server
	projectDir string
	port       int
}

func NewLocalServer(port int, projectDir string) *LocalServer {
	mux := http.NewServeMux()

	ls := &LocalServer{
		server: &http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: mux,
		},
		projectDir: projectDir,
		port:       port,
	}

	// Register handler
	mux.HandleFunc("/", ls.handleFileRequest)

	return ls
}

// Start starts the local HTTP server
func (ls *LocalServer) Start() error {
	return ls.server.ListenAndServe()
}

// Stop stops the local HTTP server
func (ls *LocalServer) Stop() error {
	return ls.server.Close()
}

// GetAddr returns the server address
func (ls *LocalServer) GetAddr() string {
	return ls.server.Addr
}

// GetPort returns the server port
func (ls *LocalServer) GetPort() int {
	return ls.port
}

// handleFileRequest handles file serving requests
func (ls *LocalServer) handleFileRequest(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	queryParams, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	// Get required parameters
	file := queryParams.Get("file")
	root := queryParams.Get("root")

	if file == "" || root == "" {
		http.Error(w, "Missing required query parameters: file and root", http.StatusBadRequest)
		return
	}

	// Construct file path
	// The root parameter should be relative to the project directory
	filePath := filepath.Join(ls.projectDir, root, file)

	// Security check: ensure the file is within the project directory
	absProjectDir, err := filepath.Abs(ls.projectDir)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// Check if the resolved path is within the project directory
	relPath, err := filepath.Rel(absProjectDir, absFilePath)
	if err != nil || filepath.IsAbs(relPath) || len(relPath) > 0 && relPath[0] == '.' {
		http.Error(w, "Access denied: file outside project directory", http.StatusForbidden)
		return
	}

	// Read the file
	data, err := os.ReadFile(absFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
		}
		return
	}

	// Encode file content as base64
	base64Data := base64.StdEncoding.EncodeToString(data)

	// Set response headers
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(base64Data)))

	// Write response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(base64Data))
}
