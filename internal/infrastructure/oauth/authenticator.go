package oauth

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

type GoogleOAuthProvider struct{}

func NewGoogleOAuthProvider() *GoogleOAuthProvider {
	return &GoogleOAuthProvider{}
}

func (p *GoogleOAuthProvider) GetHTTPClient(ctx context.Context, secretFile, tokenFile string) (*http.Client, error) {
	b, err := os.ReadFile(secretFile)
	if err != nil {
		fmt.Printf("⚠️ File credentials '%s' tidak ditemukan. Membuka wizard penyiapan...\n", secretFile)
		if wizErr := RunCredentialsWizard(secretFile); wizErr != nil {
			return nil, fmt.Errorf("wizard setup credentials gagal: %w", wizErr)
		}
		// Retry membaca file setelah wizard selesai
		b, err = os.ReadFile(secretFile)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca client_secret file (%s) setelah wizard: %w", secretFile, err)
		}
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeUploadScope, youtube.YoutubeReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("gagal parse client secret JSON: %w", err)
	}

	if len(config.RedirectURL) == 0 {
		config.RedirectURL = "http://localhost:8080/callback"
	}

	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok, err = p.getTokenFromWeb(ctx, config)
		if err != nil {
			return nil, fmt.Errorf("gagal mendapatkan token OAuth: %w", err)
		}
		if err := saveToken(tokenFile, tok); err != nil {
			fmt.Printf("Peringatan: gagal menyimpan token ke %s: %v\n", tokenFile, err)
		}
	}

	return config.Client(ctx, tok), nil
}

func (p *GoogleOAuthProvider) getTokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	codeChan := make(chan string)
	errChan := make(chan error)

	mux := http.NewServeMux()
	handler := func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			fmt.Fprintln(w, "<div style='text-align:center;font-family:sans-serif;margin-top:50px;'><h1>✅ Otentikasi YouTube Berhasil!</h1><p>Kamu bisa menutup tab browser ini dan kembali ke terminal.</p></div>")
			select {
			case codeChan <- code:
			default:
			}
			return
		}

		errMsg := r.URL.Query().Get("error")
		if errMsg != "" {
			fmt.Fprintf(w, "<div style='text-align:center;font-family:sans-serif;margin-top:50px;'><h1>❌ Otentikasi Gagal</h1><p>%s</p></div>", errMsg)
			select {
			case errChan <- fmt.Errorf("OAuth error: %s", errMsg):
			default:
			}
			return
		}
	}

	mux.HandleFunc("/callback", handler)
	mux.HandleFunc("/", handler)

	// Coba bind listener di port 8080, 8089, atau 8989 secara fleksibel
	var listener net.Listener
	var bindErr error
	ports := []string{":8080", ":8089", ":8989", ":0"}

	for _, port := range ports {
		l, err := net.Listen("tcp", port)
		if err == nil {
			listener = l
			// Update RedirectURL sesuai port yang berhasil dibind
			addr := l.Addr().(*net.TCPAddr)
			config.RedirectURL = fmt.Sprintf("http://localhost:%d/callback", addr.Port)
			break
		}
		bindErr = err
	}

	server := &http.Server{Handler: mux}
	if listener != nil {
		go func() {
			if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
				// Ignore
			}
		}()
		defer server.Shutdown(ctx)
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Println("\n=======================================================")
	fmt.Println("🔑 PERLU OTENTIKASI YOUTUBE OAUTH2")
	fmt.Println("1. Buka URL berikut di browser kamu untuk memberikan izin:")
	fmt.Printf("\n%s\n\n", authURL)
	if listener != nil {
		fmt.Printf("2. Server callback lokal aktif di: %s\n", config.RedirectURL)
	} else {
		fmt.Printf("⚠️ Port lokal sibuk (%v). Gunakan paste manual di bawah ini.\n", bindErr)
	}
	fmt.Println("3. Jika browser menampilkan 404 / gagal redirect, SALIN URL LENGKAP dari browser dan PASTE di bawah ini:")
	fmt.Println("=======================================================")

	// Goroutine untuk membaca input manual dari terminal (Failproof fallback)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\n👉 PASTE KODE / URL BROWSER DI SINI: ")
		input, err := reader.ReadString('\n')
		if err == nil {
			input = strings.TrimSpace(input)
			if input != "" {
				code := extractCodeFromInput(input)
				if code != "" {
					select {
					case codeChan <- code:
					default:
					}
				}
			}
		}
	}()

	select {
	case code := <-codeChan:
		tok, err := config.Exchange(ctx, code)
		if err != nil {
			return nil, fmt.Errorf("gagal menukar code dengan token: %w", err)
		}
		return tok, nil
	case err := <-errChan:
		return nil, err
	}
}

// extractCodeFromInput extracts OAuth code whether full URL or raw code is pasted
func extractCodeFromInput(input string) string {
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		u, err := url.Parse(input)
		if err == nil {
			code := u.Query().Get("code")
			if code != "" {
				return code
			}
		}
	}
	return input
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Menyimpan token credential ke: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("gagal membuat file token %s: %w", path, err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}
