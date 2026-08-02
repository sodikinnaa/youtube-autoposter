package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"youtube-autoposter/internal/cli"
	"youtube-autoposter/internal/domain"
	"youtube-autoposter/internal/infrastructure/oauth"
	"youtube-autoposter/internal/mcp"
	"youtube-autoposter/internal/usecase"
	"youtube-autoposter/internal/web"
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
		profileName   = flag.String("profile", "", "Pilih nama profile akun (contoh: 'akun_dua' -> token_akun_dua.json)")
		interactive   = flag.Bool("i", false, "Jalankan dalam mode interaktif (User POV)")
		jsonOutput    = flag.Bool("json", false, "Output hasil dalam format JSON (Machine / AI Agent POV)")

		// Web UI Dashboard & REST API Flags
		webMode    = flag.Bool("web", false, "Jalankan Web UI Dashboard & REST API Server (http://localhost:8080)")
		serverMode = flag.Bool("server", false, "Alias untuk -web: Jalankan Web UI Dashboard & REST API Server")
		port       = flag.Int("port", 8080, "Port HTTP untuk Web Server Dashboard (default: 8080)")

		// AI Agent Inspection & MCP Flags
		mcpMode      = flag.Bool("mcp", false, "Jalankan sebagai Model Context Protocol (MCP) Server via stdio")
		listProfiles = flag.Bool("list-profiles", false, "AI Agent: Tampilkan daftar profile token akun dalam JSON")
		listChannels = flag.Bool("list-channels", false, "AI Agent: Tampilkan daftar channel YouTube pada token dalam JSON")
		listVideos   = flag.Bool("list-videos", false, "AI Agent: Pindai dan tampilkan daftar file video dalam JSON")
	)

	flag.Parse()

	ctx := context.Background()

	// 0. Model Context Protocol (MCP) Server Mode
	if *mcpMode {
		mcp.StartMCPServer(ctx, *secretFile)
		return
	}

	// Handle Profile Name flag mapping
	if *profileName != "" {
		*tokenFile = oauth.GetTokenFileForProfile(*profileName)
	}

	// 1. AI Agent Inspection Flag: List Profiles
	if *listProfiles {
		profiles, err := oauth.ListProfiles()
		if err != nil {
			fmt.Printf(`{"status":"error","message":%q}`+"\n", err.Error())
			os.Exit(1)
		}
		json.NewEncoder(os.Stdout).Encode(profiles)
		return
	}

	// 2. AI Agent Inspection Flag: List Videos
	if *listVideos {
		foundVideos, err := cli.ScanVideoFiles(".")
		if err != nil {
			fmt.Printf(`{"status":"error","message":%q}`+"\n", err.Error())
			os.Exit(1)
		}
		var scanned []domain.ScannedVideo
		for _, v := range foundVideos {
			scanned = append(scanned, domain.ScannedVideo{
				Path:          v.Path,
				RelPath:       v.RelPath,
				SizeBytes:     v.SizeBytes,
				SizeFormatted: cli.FormatFileSize(v.SizeBytes),
			})
		}
		json.NewEncoder(os.Stdout).Encode(scanned)
		return
	}

	// Dependency Injection Setup (Clean Architecture)
	oauthProvider := oauth.NewGoogleOAuthProvider()
	getChannelInfoUseCase := usecase.NewGetChannelInfoUseCase(oauthProvider)
	uploadUseCase := usecase.NewUploadVideoUseCase(oauthProvider)

	// Web UI Dashboard & REST API Server Mode
	if *webMode || *serverMode {
		srv := web.NewServer(getChannelInfoUseCase, uploadUseCase, *secretFile, *tokenFile, *port)
		if err := srv.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error menjalankan Web Server: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// 3. AI Agent Inspection Flag: List Channels
	if *listChannels {
		channels, err := getChannelInfoUseCase.ExecuteList(ctx, *secretFile, *tokenFile)
		if err != nil {
			fmt.Printf(`{"status":"error","message":%q}`+"\n", err.Error())
			os.Exit(1)
		}
		json.NewEncoder(os.Stdout).Encode(channels)
		return
	}

	var input usecase.UploadVideoInput

	// Jika flag -file tidak diisi atau flag -i diaktifkan, masuk ke Mode Interaktif (User POV)
	if *videoPath == "" || *interactive {
		var ok bool
		input, ok = cli.RunInteractiveMode(ctx, *secretFile, *tokenFile, getChannelInfoUseCase)
		if !ok {
			os.Exit(0)
		}
	} else {
		// Mode Direct Flag (Scripting / Automation / AI Agent POV)
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

	result, err := uploadUseCase.Execute(ctx, input)
	if *jsonOutput {
		if err != nil {
			fmt.Printf(`{"status":"error","message":%q}`+"\n", err.Error())
			os.Exit(1)
		}
		json.NewEncoder(os.Stdout).Encode(result)
		return
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "\n❌ Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n=======================================================")
	fmt.Println("🎉 UPLOAD VIDEO BERHASIL DIPROSES!")
	fmt.Printf("📺 Video ID    : %s\n", result.ID)
	fmt.Printf("🔗 Link Watch  : %s\n", result.WatchURL)
	fmt.Printf("⏱️  Durasi Upload: %s\n", result.DurationStr)
	if result.HasCustomThumb {
		fmt.Println("🖼️  Thumbnail   : Custom Thumbnail Berhasil Dipasang")
	}
	fmt.Println("=======================================================")
}
