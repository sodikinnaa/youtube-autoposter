package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
		return nil, fmt.Errorf("gagal membaca client_secret file (%s): %w. Pastikan file credentials dari Google Cloud Console sudah ada", secretFile, err)
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

	server := &http.Server{Addr: ":8080"}
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			fmt.Fprintln(w, "<h1>Otentikasi Berhasil!</h1><p>Kamu bisa menutup tab ini dan kembali ke terminal.</p>")
			codeChan <- code
		} else {
			errMsg := r.URL.Query().Get("error")
			fmt.Fprintf(w, "<h1>Otentikasi Gagal</h1><p>%s</p>", errMsg)
			errChan <- fmt.Errorf("OAuth error: %s", errMsg)
		}
	})

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Ignore closed server
		}
	}()
	defer server.Shutdown(ctx)

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Println("\n=======================================================")
	fmt.Println("🔑 PERLU OTENTIKASI YOUTUBE OAUTH2")
	fmt.Println("Buka URL berikut di browser kamu untuk memberikan izin:")
	fmt.Printf("\n%s\n\n", authURL)
	fmt.Println("Sedang menunggu callback otomatis di http://localhost:8080/callback...")
	fmt.Println("=======================================================")

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
