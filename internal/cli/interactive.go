package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"youtube-autoposter/internal/usecase"
)

// RunInteractiveMode launches a step-by-step user-friendly CLI prompt
func RunInteractiveMode() (usecase.UploadVideoInput, bool) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n=======================================================")
	fmt.Println("🎬 YOUTUBE AUTO-POSTER - INTERACTIVE WIZARD")
	fmt.Println("=======================================================")
	fmt.Println("Silakan isi detail video di bawah ini (atau tekan Enter untuk default):")

	// 1. Path File Video
	var videoPath string
	for {
		fmt.Print("\n🎥 Path File Video (contoh: ./my_video.mp4): ")
		input, _ := reader.ReadString('\n')
		videoPath = strings.TrimSpace(input)
		if videoPath == "" {
			fmt.Println("❌ Path file video tidak boleh kosong!")
			continue
		}
		if _, err := os.Stat(videoPath); os.IsNotExist(err) {
			fmt.Printf("❌ File '%s' tidak ditemukan! Silakan masukkan path yang benar.\n", videoPath)
			continue
		}
		break
	}

	// 2. Path Thumbnail (Opsional)
	fmt.Print("\n🖼️  Path Custom Thumbnail (opsional, contoh: ./thumb.jpg, tekan Enter jika tidak ada): ")
	thumbInput, _ := reader.ReadString('\n')
	thumbPath := strings.TrimSpace(thumbInput)
	if thumbPath != "" {
		if _, err := os.Stat(thumbPath); os.IsNotExist(err) {
			fmt.Printf("⚠️  Peringatan: File thumbnail '%s' tidak ditemukan. Upload akan dilanjutkan tanpa thumbnail.\n", thumbPath)
			thumbPath = ""
		}
	}

	// 3. Judul Video
	defaultTitle := strings.TrimSuffix(filepath.Base(videoPath), filepath.Ext(videoPath))
	fmt.Printf("\n📝 Judul Video [default: %s]: ", defaultTitle)
	titleInput, _ := reader.ReadString('\n')
	title := strings.TrimSpace(titleInput)
	if title == "" {
		title = defaultTitle
	}

	// 4. Deskripsi Video
	fmt.Print("\n📄 Deskripsi Video (opsional, tekan Enter jika kosong): ")
	descInput, _ := reader.ReadString('\n')
	description := strings.TrimSpace(descInput)

	// 5. Tags Video
	fmt.Print("\n🏷️  Tags Video (pisahkan dengan koma, contoh: coding,golang,tutorial): ")
	tagsInput, _ := reader.ReadString('\n')
	var tags []string
	if strings.TrimSpace(tagsInput) != "" {
		for _, t := range strings.Split(tagsInput, ",") {
			trimmed := strings.TrimSpace(t)
			if trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}

	// 6. Status Privasi
	fmt.Println("\n🔒 Status Privasi Video:")
	fmt.Println("   [1] Private  (Hanya kamu yang bisa lihat - Default)")
	fmt.Println("   [2] Public   (Bisa ditonton siapa saja)")
	fmt.Println("   [3] Unlisted (Hanya yang punya link yang bisa lihat)")
	fmt.Print("Pilihan Status (1/2/3) [default: 1]: ")
	privInput, _ := reader.ReadString('\n')
	privChoice := strings.TrimSpace(privInput)

	privacyStatus := "private"
	switch privChoice {
	case "2":
		privacyStatus = "public"
	case "3":
		privacyStatus = "unlisted"
	}

	// 7. Scheduled Publish
	fmt.Print("\n⏱️  Jadwalkan Tayang Otomatis? (y/N): ")
	schedConfirm, _ := reader.ReadString('\n')
	var publishAt string
	if strings.ToLower(strings.TrimSpace(schedConfirm)) == "y" {
		fmt.Print("Masukkan tanggal & waktu (format: YYYY-MM-DD HH:MM, contoh: 2026-08-01 15:00): ")
		timeInput, _ := reader.ReadString('\n')
		timeStr := strings.TrimSpace(timeInput)
		parsedTime, err := time.Parse("2006-01-02 15:04", timeStr)
		if err != nil {
			fmt.Println("⚠️  Format waktu tidak valid. Jadwal Tayang Otomatis dibatalkan.")
		} else {
			publishAt = parsedTime.Format(time.RFC3339)
			fmt.Printf("✅ Tayang Otomatis dijadwalkan pada: %s\n", parsedTime.Format("02 Jan 2006 15:04 MST"))
		}
	}

	// Konfirmasi Sebelum Upload
	fmt.Println("\n=======================================================")
	fmt.Println("📋 RINGKASAN KONFIGURASI UPLOAD")
	fmt.Println("=======================================================")
	fmt.Printf("🎥 File Video : %s\n", videoPath)
	if thumbPath != "" {
		fmt.Printf("🖼️  Thumbnail  : %s\n", thumbPath)
	} else {
		fmt.Println("🖼️  Thumbnail  : (Tanpa Thumbnail)")
	}
	fmt.Printf("📝 Judul      : %s\n", title)
	if description != "" {
		fmt.Printf("📄 Deskripsi  : %s\n", description)
	}
	if len(tags) > 0 {
		fmt.Printf("🏷️  Tags       : %s\n", strings.Join(tags, ", "))
	}
	fmt.Printf("🔒 Privasi    : %s\n", privacyStatus)
	if publishAt != "" {
		fmt.Printf("⏱️  Jadwal     : %s\n", publishAt)
	}
	fmt.Println("=======================================================")
	fmt.Print("Apakah konfigurasi di atas sudah sesuai dan siap upload? (Y/n): ")

	confirm, _ := reader.ReadString('\n')
	confirmStr := strings.ToLower(strings.TrimSpace(confirm))
	if confirmStr == "n" {
		fmt.Println("\n🚫 Upload dibatalkan oleh pengguna.")
		return usecase.UploadVideoInput{}, false
	}

	return usecase.UploadVideoInput{
		FilePath:      videoPath,
		ThumbnailPath: thumbPath,
		Title:         title,
		Description:   description,
		Tags:          tags,
		CategoryID:    "22",
		PrivacyStatus: privacyStatus,
		PublishAt:     publishAt,
		SecretFile:    "client_secret.json",
		TokenFile:     "token.json",
	}, true
}
