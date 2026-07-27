package oauth

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type InstalledAppCredentials struct {
	Installed struct {
		ClientID                string   `json:"client_id"`
		ProjectID               string   `json:"project_id"`
		AuthURI                 string   `json:"auth_uri"`
		TokenURI                string   `json:"token_uri"`
		AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url"`
		ClientSecret            string   `json:"client_secret"`
		RedirectURIs            []string `json:"redirect_uris"`
	} `json:"installed"`
}

// RunCredentialsWizard guides the user through setting up client_secret.json interactively
func RunCredentialsWizard(secretFile string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n=======================================================")
	fmt.Println("🧙 WIZARD SETUP GOOGLE OAUTH CREDENTIALS")
	fmt.Println("=======================================================")
	fmt.Printf("File '%s' belum ditemukan.\n\n", secretFile)
	fmt.Println("Pilih cara penyiapan credentials:")
	fmt.Println("  [1] Input Client ID & Client Secret secara manual")
	fmt.Println("  [2] Masukkan path file JSON credentials yang sudah didownload")
	fmt.Println("  [3] Lihat Panduan Google Cloud Console & Keluar")
	fmt.Print("\nPilihan Kamu (1/2/3): ")

	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		return createSecretFromInputs(reader, secretFile)
	case "2":
		return copySecretFromFile(reader, secretFile)
	case "3":
		printGCPGuide()
		os.Exit(0)
		return nil
	default:
		fmt.Println("Pilihan tidak valid. Membatalkan setup wizard.")
		os.Exit(1)
		return nil
	}
}

func createSecretFromInputs(reader *bufio.Reader, targetPath string) error {
	fmt.Println("\n-------------------------------------------------------")
	fmt.Println("📝 Masukkan Client ID & Client Secret")
	fmt.Println("(Dapatkan dari Google Cloud Console > Credentials > OAuth 2.0 Client ID)")
	fmt.Println("-------------------------------------------------------")

	fmt.Print("Client ID     : ")
	clientID, _ := reader.ReadString('\n')
	clientID = strings.TrimSpace(clientID)

	fmt.Print("Client Secret : ")
	clientSecret, _ := reader.ReadString('\n')
	clientSecret = strings.TrimSpace(clientSecret)

	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("Client ID dan Client Secret tidak boleh kosong")
	}

	creds := InstalledAppCredentials{}
	creds.Installed.ClientID = clientID
	creds.Installed.ClientSecret = clientSecret
	creds.Installed.ProjectID = "youtube-autoposter"
	creds.Installed.AuthURI = "https://accounts.google.com/o/oauth2/auth"
	creds.Installed.TokenURI = "https://oauth2.googleapis.com/token"
	creds.Installed.AuthProviderX509CertURL = "https://www.googleapis.com/oauth2/v1/certs"
	creds.Installed.RedirectURIs = []string{"http://localhost:8080/callback"}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("gagal format JSON credentials: %w", err)
	}

	err = os.WriteFile(targetPath, data, 0600)
	if err != nil {
		return fmt.Errorf("gagal membuat file %s: %w", targetPath, err)
	}

	fmt.Printf("\n✅ Berhasil membuat file credentials di '%s'!\n\n", targetPath)
	return nil
}

func copySecretFromFile(reader *bufio.Reader, targetPath string) error {
	fmt.Print("\nMasukkan path lokasi file credentials JSON: ")
	sourcePath, _ := reader.ReadString('\n')
	sourcePath = strings.TrimSpace(sourcePath)

	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("gagal membaca file '%s': %w", sourcePath, err)
	}

	err = os.WriteFile(targetPath, data, 0600)
	if err != nil {
		return fmt.Errorf("gagal menyalin ke '%s': %w", targetPath, err)
	}

	fmt.Printf("\n✅ Berhasil menyalin credentials dari '%s' ke '%s'!\n\n", sourcePath, targetPath)
	return nil
}

func printGCPGuide() {
	fmt.Println("\n=======================================================")
	fmt.Println("📖 PANDUAN GOOGLE AUTH PLATFORM (TAMPILAN TERBARU)")
	fmt.Println("=======================================================")
	fmt.Println("1. Buka Google Cloud Console: https://console.cloud.google.com/")
	fmt.Println("2. Aktifkan API 'YouTube Data API v3' di menu API & Services > Library.")
	fmt.Println("3. Masuk ke menu 'Google Auth Platform':")
	fmt.Println("   a. Klik 'Branding' / 'Audience' (Configure App): Pilih External, isi Nama App & Email kamu.")
	fmt.Println("   b. Di menu 'Audience' > 'Test Users': Klik '+ Add Users', masukkan email Google/YouTube kamu.")
	fmt.Println("   c. Di menu 'Clients' (Credentials): Klik '+ Create Client', pilih Application Type: 'Desktop app'.")
	fmt.Println("4. Salin 'Client ID' dan 'Client Secret' yang muncul di layar, atau Download JSON-nya.")
	fmt.Println("5. Jalankan kembali ./youtube-autoposter untuk memasukkan Client ID & Secret tersebut.")
	fmt.Println("=======================================================")
}
