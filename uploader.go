package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type UploadOptions struct {
	VideoPath     string
	ThumbnailPath string
	Title         string
	Description   string
	Tags          []string
	CategoryID    string
	PrivacyStatus string // public, private, unlisted
	PublishAt     string // RFC3339 format e.g. 2026-08-01T15:00:00Z
	SecretFile    string
	TokenFile     string
}

func UploadVideo(ctx context.Context, opts UploadOptions) error {
	file, err := os.Open(opts.VideoPath)
	if err != nil {
		return fmt.Errorf("gagal membuka file video (%s): %w", opts.VideoPath, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("gagal membaca info file: %w", err)
	}

	fmt.Printf("📂 File Video: %s (%.2f MB)\n", opts.VideoPath, float64(fileInfo.Size())/(1024*1024))
	fmt.Printf("📝 Title: %s\n", opts.Title)
	fmt.Printf("🔒 Privacy Status: %s\n", opts.PrivacyStatus)
	if len(opts.Tags) > 0 {
		fmt.Printf("🏷️  Tags: %s\n", strings.Join(opts.Tags, ", "))
	}

	client, err := getClient(ctx, opts.SecretFile, opts.TokenFile)
	if err != nil {
		return err
	}

	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("gagal membuat YouTube service client: %w", err)
	}

	uploadVideo := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       opts.Title,
			Description: opts.Description,
			Tags:        opts.Tags,
			CategoryId:  opts.CategoryID,
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus: opts.PrivacyStatus,
		},
	}

	if opts.PublishAt != "" {
		// Validasi format waktu publish jika diset
		if _, err := time.Parse(time.RFC3339, opts.PublishAt); err != nil {
			return fmt.Errorf("format publish-at tidak valid (%s). Gunakan format RFC3339, contoh: 2026-08-01T15:00:00Z", opts.PublishAt)
		}
		uploadVideo.Status.PublishAt = opts.PublishAt
		// Jika scheduled publish, status privacy harus private terlebih dahulu sampai waktu publish
		if opts.PrivacyStatus != "private" {
			fmt.Println("⚠️  Catatan: Video scheduled publish diset otomatis ke privacyStatus 'private' sebelum waktu tayang.")
			uploadVideo.Status.PrivacyStatus = "private"
		}
	}

	call := service.Videos.Insert([]string{"snippet", "status"}, uploadVideo)
	call.Media(file)

	fmt.Println("\n🚀 Memulai proses upload video ke YouTube...")
	start := time.Now()

	response, err := call.Do()
	if err != nil {
		return fmt.Errorf("gagal mengupload video ke YouTube: %w", err)
	}

	duration := time.Since(start)
	fmt.Println("\n=======================================================")
	fmt.Println("✅ UPLOAD VIDEO BERHASIL!")
	fmt.Printf("📺 Video ID : %s\n", response.Id)
	fmt.Printf("🔗 Link Watch: https://www.youtube.com/watch?v=%s\n", response.Id)
	fmt.Printf("⏱️  Waktu Upload: %s\n", duration.Round(time.Second))

	// Upload Thumbnail jika opsi thumbnail diberikan
	if opts.ThumbnailPath != "" {
		fmt.Printf("\n🖼️  Memproses upload thumbnail: %s...\n", opts.ThumbnailPath)
		thumbFile, err := os.Open(opts.ThumbnailPath)
		if err != nil {
			fmt.Printf("⚠️  Gagal membuka file thumbnail (%s): %v\n", opts.ThumbnailPath, err)
		} else {
			defer thumbFile.Close()
			_, err = service.Thumbnails.Set(response.Id).Media(thumbFile).Do()
			if err != nil {
				fmt.Printf("⚠️  Gagal mengupload thumbnail: %v\n", err)
			} else {
				fmt.Println("✨ Thumbnail kustom berhasil dipasang!")
			}
		}
	}

	fmt.Println("=======================================================")

	return nil
}
