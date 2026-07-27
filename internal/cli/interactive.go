package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"youtube-autoposter/internal/domain"
	"youtube-autoposter/internal/infrastructure/oauth"
	"youtube-autoposter/internal/usecase"
)

// RunInteractiveMode launches a step-by-step user-friendly CLI prompt starting with credential & channel verification
func RunInteractiveMode(ctx context.Context, secretFile, tokenFile string, getChannelInfoUseCase *usecase.GetChannelInfoUseCase) (usecase.UploadVideoInput, bool) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n=======================================================")
	fmt.Println("🎬 YOUTUBE AUTO-POSTER - INTERACTIVE WIZARD")
	fmt.Println("=======================================================")

	// 0. MULTI-ACCOUNT PROFILE MANAGER
	profiles, _ := oauth.ListProfiles()
	if len(profiles) > 0 {
		fmt.Println("\n=======================================================")
		fmt.Println("👥 PILIH AKUN YOUTUBE (MULTI-ACCOUNT MANAGER)")
		fmt.Println("=======================================================")
		for i, p := range profiles {
			fmt.Printf("  [%d] %s (%s)\n", i+1, p.Name, p.TokenFile)
		}
		newAccIdx := len(profiles) + 1
		fmt.Printf("  [%d] + Tambah / Login Akun YouTube Baru\n", newAccIdx)
		fmt.Printf("Pilihan Akun (1-%d) [default: 1]: ", newAccIdx)

		pChoiceStr, _ := reader.ReadString('\n')
		pChoiceStr = strings.TrimSpace(pChoiceStr)

		var choiceIdx int
		if pChoiceStr == "" {
			choiceIdx = 1
		} else {
			fmt.Sscanf(pChoiceStr, "%d", &choiceIdx)
		}

		if choiceIdx == newAccIdx {
			fmt.Print("\nMasukkan Nama/Alias untuk Akun Baru ini (contoh: channel_dua): ")
			aliasInput, _ := reader.ReadString('\n')
			alias := strings.TrimSpace(aliasInput)
			if alias != "" {
				tokenFile = oauth.GetTokenFileForProfile(alias)
				fmt.Printf("✅ Menyiapkan profile token baru di '%s'...\n", tokenFile)
			}
		} else if choiceIdx >= 1 && choiceIdx <= len(profiles) {
			tokenFile = profiles[choiceIdx-1].TokenFile
		}
	}

	fmt.Println("\n🔍 Memeriksa kredensial OAuth & Mengambil Daftar Channel YouTube...")

	// 1. VERIFIKASI KREDENSIAL & LOAD DAFTAR CHANNEL YOUTUBE DARI KREDENSIAL INI
	channels, err := getChannelInfoUseCase.ExecuteList(ctx, secretFile, tokenFile)
	if err != nil {
		fmt.Printf("\n❌ Gagal terhubung ke YouTube API / OAuth Credentials Error:\n%v\n", err)
		return usecase.UploadVideoInput{}, false
	}

	var selectedChannel domain.ChannelInfo
	if len(channels) == 1 {
		selectedChannel = channels[0]
		fmt.Println("\n=======================================================")
		fmt.Println("📺 CHANNEL YOUTUBE TERHUBUNG SUKSES")
		fmt.Println("=======================================================")
		fmt.Printf("👤 Nama Channel : %s\n", selectedChannel.Title)
		fmt.Printf("🆔 Channel ID   : %s\n", selectedChannel.ID)
		fmt.Printf("👥 Subscribers  : %d\n", selectedChannel.SubscriberCount)
		fmt.Printf("📹 Total Video  : %d\n", selectedChannel.VideoCount)
		fmt.Println("=======================================================")
		fmt.Print("Lanjutkan upload video ke channel di atas? (Y/n): ")

		channelConfirm, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(channelConfirm)) == "n" {
			fmt.Println("\n🚫 Proses dibatalkan oleh pengguna.")
			return usecase.UploadVideoInput{}, false
		}
	} else {
		fmt.Println("\n=======================================================")
		fmt.Printf("📺 DITEMUKAN %d CHANNEL YOUTUBE PADA AKUN INI:\n", len(channels))
		fmt.Println("=======================================================")
		for i, ch := range channels {
			fmt.Printf("  [%d] %s (%d Subs | %d Video) - ID: %s\n", i+1, ch.Title, ch.SubscriberCount, ch.VideoCount, ch.ID)
		}
		fmt.Printf("Pilih Channel Target Upload (1-%d) [default: 1]: ", len(channels))

		chChoiceStr, _ := reader.ReadString('\n')
		chChoiceStr = strings.TrimSpace(chChoiceStr)

		var chIdx int
		if chChoiceStr == "" {
			chIdx = 1
		} else {
			fmt.Sscanf(chChoiceStr, "%d", &chIdx)
		}

		if chIdx < 1 || chIdx > len(channels) {
			chIdx = 1
		}
		selectedChannel = channels[chIdx-1]
		fmt.Printf("\n✅ Channel terpilih: %s (ID: %s)\n", selectedChannel.Title, selectedChannel.ID)
	}

	fmt.Println("\nSilakan isi detail video yang akan di-upload:")

	// 2. DETEKSI OTOMATIS & PILIH FILE VIDEO
	foundVideos, _ := ScanVideoFiles(".")
	var videoPath string

	if len(foundVideos) > 0 {
		fmt.Println("\n=======================================================")
		fmt.Printf("🎥 DETEKSI OTOMATIS: DITEMUKAN %d FILE VIDEO:\n", len(foundVideos))
		fmt.Println("=======================================================")
		for i, v := range foundVideos {
			fmt.Printf("  [%d] %s (%s)\n", i+1, v.RelPath, FormatFileSize(v.SizeBytes))
		}
		manualIdx := len(foundVideos) + 1
		fmt.Printf("  [%d] 📂 Ketik/Tulis Path File Video Manual\n", manualIdx)
		fmt.Printf("Pilih File Video (1-%d) [default: 1]: ", manualIdx)

		vChoiceStr, _ := reader.ReadString('\n')
		vChoiceStr = strings.TrimSpace(vChoiceStr)

		var vIdx int
		if vChoiceStr == "" {
			vIdx = 1
		} else {
			fmt.Sscanf(vChoiceStr, "%d", &vIdx)
		}

		if vIdx >= 1 && vIdx <= len(foundVideos) {
			videoPath = foundVideos[vIdx-1].RelPath
			fmt.Printf("✅ File video terpilih: %s\n", videoPath)
		}
	}

	if videoPath == "" {
		for {
			fmt.Print("\n🎥 Masukkan Path File Video (contoh: ./my_video.mp4): ")
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
	}

	// 3. PATH THUMBNAIL (OPSIONAL)
	fmt.Print("\n🖼️  Path Custom Thumbnail (opsional, contoh: ./thumb.jpg, tekan Enter jika tidak ada): ")
	thumbInput, _ := reader.ReadString('\n')
	thumbPath := strings.TrimSpace(thumbInput)
	if thumbPath != "" {
		if _, err := os.Stat(thumbPath); os.IsNotExist(err) {
			fmt.Printf("⚠️  Peringatan: File thumbnail '%s' tidak ditemukan. Upload akan dilanjutkan tanpa thumbnail.\n", thumbPath)
			thumbPath = ""
		}
	}

	// 4. JUDUL VIDEO
	defaultTitle := strings.TrimSuffix(filepath.Base(videoPath), filepath.Ext(videoPath))
	fmt.Printf("\n📝 Judul Video [default: %s]: ", defaultTitle)
	titleInput, _ := reader.ReadString('\n')
	title := strings.TrimSpace(titleInput)
	if title == "" {
		title = defaultTitle
	}

	// 5. DESKRIPSI VIDEO
	fmt.Print("\n📄 Deskripsi Video (opsional, tekan Enter jika kosong): ")
	descInput, _ := reader.ReadString('\n')
	description := strings.TrimSpace(descInput)

	// 6. TAGS VIDEO
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

	// 7. STATUS PRIVASI
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

	// 8. SCHEDULED PUBLISH
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

	// KONFIRMASI AKHIR SEBELUM UPLOAD
	fmt.Println("\n=======================================================")
	fmt.Println("📋 RINGKASAN KONFIGURASI UPLOAD")
	fmt.Println("=======================================================")
	fmt.Printf("📺 Channel Target: %s (%s)\n", selectedChannel.Title, selectedChannel.ID)
	fmt.Printf("🎥 File Video    : %s\n", videoPath)
	if thumbPath != "" {
		fmt.Printf("🖼️  Thumbnail     : %s\n", thumbPath)
	} else {
		fmt.Println("🖼️  Thumbnail     : (Tanpa Thumbnail)")
	}
	fmt.Printf("📝 Judul         : %s\n", title)
	if description != "" {
		fmt.Printf("📄 Deskripsi     : %s\n", description)
	}
	if len(tags) > 0 {
		fmt.Printf("🏷️  Tags          : %s\n", strings.Join(tags, ", "))
	}
	fmt.Printf("🔒 Privasi       : %s\n", privacyStatus)
	if publishAt != "" {
		fmt.Printf("⏱️  Jadwal        : %s\n", publishAt)
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
		SecretFile:    secretFile,
		TokenFile:     tokenFile,
	}, true
}
