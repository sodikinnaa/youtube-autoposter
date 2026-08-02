package web

import (
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/http"
)

//go:embed static/*
var staticFS embed.FS

func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Static Web UI Dashboard
	subFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		return err
	}
	fileServer := http.FileServer(http.FS(subFS))
	noCacheFileServer := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		fileServer.ServeHTTP(w, r)
	})
	mux.Handle("/", noCacheFileServer)

	// CORS Middleware wrapper
	cors := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			h(w, r)
		}
	}

	mux.HandleFunc("/api/health", cors(s.handleHealth))
	mux.HandleFunc("/api/profiles", cors(s.handleProfiles))
	mux.HandleFunc("/api/channels", cors(s.handleChannels))
	mux.HandleFunc("/api/videos", cors(s.handleVideos))
	mux.HandleFunc("/api/upload", cors(s.handleUpload))

	// Auth Wizard Endpoints
	mux.HandleFunc("/api/auth/status", cors(s.handleAuthStatus))
	mux.HandleFunc("/api/auth/save-secret", cors(s.handleSaveSecret))
	mux.HandleFunc("/api/auth/url", cors(s.handleGetAuthURL))
	mux.HandleFunc("/api/auth/exchange", cors(s.handleExchangeCode))

	// Port Binding with Automatic Random/Free Port Fallback
	var listener net.Listener
	targetPort := s.port

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", targetPort))
	if err == nil {
		listener = l
	} else {
		fmt.Printf("⚠️ Port %d sedang digunakan (%v). Mencari port bebas/acak...\n", targetPort, err)
		// Try consecutive ports up to targetPort+50
		for i := 1; i <= 50; i++ {
			altPort := targetPort + i
			l, err := net.Listen("tcp", fmt.Sprintf(":%d", altPort))
			if err == nil {
				listener = l
				targetPort = altPort
				break
			}
		}
		// Final fallback: let OS pick random free port (:0)
		if listener == nil {
			l, err := net.Listen("tcp", ":0")
			if err != nil {
				return fmt.Errorf("gagal menemukan port bebas: %w", err)
			}
			listener = l
			targetPort = l.Addr().(*net.TCPAddr).Port
		}
	}

	s.port = targetPort

	fmt.Printf("\n🚀 Web UI & REST API Dashboard berjalan di: http://localhost:%d\n", s.port)
	fmt.Printf("🌐 Akses Dashboard Browser  : http://localhost:%d\n", s.port)
	fmt.Printf("📡 REST API Endpoints:\n")
	fmt.Printf("   - GET  http://localhost:%d/api/health\n", s.port)
	fmt.Printf("   - GET  http://localhost:%d/api/profiles\n", s.port)
	fmt.Printf("   - GET  http://localhost:%d/api/channels?profile=<name>\n", s.port)
	fmt.Printf("   - GET  http://localhost:%d/api/videos\n", s.port)
	fmt.Printf("   - POST http://localhost:%d/api/upload\n\n", s.port)

	return http.Serve(listener, mux)
}
