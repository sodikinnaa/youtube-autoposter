package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var (
		videoPath     = flag.String("file", "", "Path ke file video (misal: ./video.mp4)")
		thumbnailPath = flag.String("thumbnail", "", "Path ke file gambar thumbnail (contoh: ./thumbnail.jpg)")
		title         = flag.String("title", "", "Judul video")
		description   = flag.String("description", "", "Deskripsi video")
		tagsStr       = flag.String("tags", "", "Tags video dipisahkan koma (contoh: coding,golang,tutorial)")
		categoryID    = flag.String("category", "22", "YouTube Category ID (22=People & Blogs, 28=Science & Tech, 27=Education)")
		privacyStatus = flag.String("privacy", "private", "Privacy Status: public, private, unlisted")
		publishAt     = flag.String("publish-at", "", "Jadwal tayang otomatis (format RFC3339, contoh: 2026-08-01T15:00:00Z)")
		secretFile    = flag.String("secret", "client_secret.json", "Path ke file OAuth client_secret.json")
		tokenFile     = flag.String("token", "token.json", "Path ke simpanan token.json")
	)

	flag.Parse()

	if *videoPath == "" {
		fmt.Println("🎬 YouTube Auto-Poster CLI (Golang)")
		fmt.Println("==================================================")
		fmt.Println("Penggunaan:")
		fmt.Println("  go run . -file video.mp4 -thumbnail thumb.png -title \"Judul Video\" -privacy private")
		fmt.Println("\nFlag lengkap:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *title == "" {
		// Gunakan nama file tanpa ekstensi sebagai judul default jika judul kosong
		baseName := filepath.Base(*videoPath)
		ext := filepath.Ext(baseName)
		*title = strings.TrimSuffix(baseName, ext)
	}

	var tags []string
	if *tagsStr != "" {
		for _, t := range strings.Split(*tagsStr, ",") {
			trimmed := strings.TrimSpace(t)
			if trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}

	opts := UploadOptions{
		VideoPath:     *videoPath,
		ThumbnailPath: *thumbnailPath,
		Title:         *title,
		Description:   *description,
		Tags:          tags,
		CategoryID:    *categoryID,
		PrivacyStatus: *privacyStatus,
		PublishAt:     *publishAt,
		SecretFile:    *secretFile,
		TokenFile:     *tokenFile,
	}

	ctx := context.Background()
	if err := UploadVideo(ctx, opts); err != nil {
		fmt.Fprintf(os.Stderr, "\n❌ Error: %v\n", err)
		os.Exit(1)
	}
}
