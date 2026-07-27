package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"youtube-autoposter/internal/cli"
	"youtube-autoposter/internal/infrastructure/oauth"
	"youtube-autoposter/internal/usecase"
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
		interactive   = flag.Bool("i", false, "Jalankan dalam mode interaktif (User POV)")
	)

	flag.Parse()

	// Dependency Injection Setup (Clean Architecture)
	oauthProvider := oauth.NewGoogleOAuthProvider()
	uploadUseCase := usecase.NewUploadVideoUseCase(oauthProvider)

	var input usecase.UploadVideoInput

	// Jika flag -file tidak diisi atau flag -i diaktifkan, masuk ke Mode Interaktif (User POV)
	if *videoPath == "" || *interactive {
		var ok bool
		input, ok = cli.RunInteractiveMode()
		if !ok {
			os.Exit(0)
		}
	} else {
		// Mode Direct Flag (Scripting / Automation POV)
		var tags []string
		if *tagsStr != "" {
			for _, t := range strings.Split(*tagsStr, ",") {
				trimmed := strings.TrimSpace(t)
				if trimmed != "" {
					tags = append(tags, trimmed)
				}
			}
		}

		input = usecase.UploadVideoInput{
			FilePath:      *videoPath,
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
	}

	ctx := context.Background()
	result, err := uploadUseCase.Execute(ctx, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n❌ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n=======================================================")
	fmt.Println("🎉 UPLOAD VIDEO BERHASIL DIPROSES!")
	fmt.Printf("📺 Video ID    : %s\n", result.ID)
	fmt.Printf("🔗 Link Watch  : %s\n", result.WatchURL)
	fmt.Printf("⏱️  Durasi Upload: %s\n", result.Duration)
	if result.HasCustomThumb {
		fmt.Println("🖼️  Thumbnail   : Custom Thumbnail Berhasil Dipasang")
	}
	fmt.Println("=======================================================")
}
