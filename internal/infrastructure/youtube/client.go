package youtube

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"youtube-autoposter/internal/domain"

	"google.golang.org/api/option"
	yt "google.golang.org/api/youtube/v3"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

func (c *Client) GetChannelsList(ctx context.Context) ([]domain.ChannelInfo, error) {
	service, err := yt.NewService(ctx, option.WithHTTPClient(c.httpClient))
	if err != nil {
		return nil, fmt.Errorf("gagal inisialisasi YouTube API client: %w", err)
	}

	call := service.Channels.List([]string{"snippet", "statistics"}).Mine(true)
	resp, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data channel YouTube: %w", err)
	}

	if len(resp.Items) == 0 {
		return nil, fmt.Errorf("tidak ditemukan channel YouTube yang terhubung dengan akun ini")
	}

	var channels []domain.ChannelInfo
	for _, ch := range resp.Items {
		var subsCount, videoCount uint64
		if ch.Statistics != nil {
			subsCount = ch.Statistics.SubscriberCount
			videoCount = ch.Statistics.VideoCount
		}
		title := "Channel Tanpa Nama"
		if ch.Snippet != nil {
			title = ch.Snippet.Title
		}
		channels = append(channels, domain.ChannelInfo{
			ID:              ch.Id,
			Title:           title,
			SubscriberCount: subsCount,
			VideoCount:      videoCount,
		})
	}
	return channels, nil
}

func (c *Client) GetChannelInfo(ctx context.Context) (*domain.ChannelInfo, error) {
	list, err := c.GetChannelsList(ctx)
	if err != nil {
		return nil, err
	}
	return &list[0], nil
}

func (c *Client) UploadVideo(ctx context.Context, video domain.Video) (*domain.UploadResult, error) {
	file, err := os.Open(video.FilePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFileNotFound, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("gagal membaca info file video: %w", err)
	}

	fmt.Printf("📂 File Video: %s (%.2f MB)\n", video.FilePath, float64(fileInfo.Size())/(1024*1024))
	fmt.Printf("📝 Title: %s\n", video.Title)
	fmt.Printf("🔒 Privacy Status: %s\n", video.Privacy)
	if len(video.Tags) > 0 {
		fmt.Printf("🏷️  Tags: %s\n", strings.Join(video.Tags, ", "))
	}

	service, err := yt.NewService(ctx, option.WithHTTPClient(c.httpClient))
	if err != nil {
		return nil, fmt.Errorf("gagal inisialisasi YouTube API client: %w", err)
	}

	uploadVideo := &yt.Video{
		Snippet: &yt.VideoSnippet{
			Title:       video.Title,
			Description: video.Description,
			Tags:        video.Tags,
			CategoryId:  video.CategoryID,
		},
		Status: &yt.VideoStatus{
			PrivacyStatus: string(video.Privacy),
		},
	}

	if video.PublishAt != nil {
		uploadVideo.Status.PublishAt = video.PublishAt.Format(time.RFC3339)
		if video.Privacy != domain.PrivacyPrivate {
			fmt.Println("⚠️  Catatan: Video scheduled publish diset otomatis ke privacyStatus 'private' sebelum waktu tayang.")
			uploadVideo.Status.PrivacyStatus = string(domain.PrivacyPrivate)
		}
	}

	call := service.Videos.Insert([]string{"snippet", "status"}, uploadVideo)
	call.Media(file)

	fmt.Println("\n🚀 Memulai proses upload video ke YouTube...")
	start := time.Now()

	res, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("gagal mengupload video ke YouTube: %w", err)
	}

	elapsed := time.Since(start)
	result := &domain.UploadResult{
		ID:          res.Id,
		WatchURL:    fmt.Sprintf("https://www.youtube.com/watch?v=%s", res.Id),
		Duration:    elapsed.Round(time.Second),
		DurationStr: elapsed.Round(time.Second).String(),
	}

	// Upload Thumbnail jika diberikan
	if video.ThumbnailPath != "" {
		fmt.Printf("\n🖼️  Memproses upload thumbnail: %s...\n", video.ThumbnailPath)
		thumbFile, err := os.Open(video.ThumbnailPath)
		if err != nil {
			fmt.Printf("⚠️  Gagal membuka file thumbnail (%s): %v\n", video.ThumbnailPath, err)
		} else {
			defer thumbFile.Close()
			_, err = service.Thumbnails.Set(res.Id).Media(thumbFile).Do()
			if err != nil {
				fmt.Printf("⚠️  Gagal mengupload thumbnail: %v\n", err)
			} else {
				fmt.Println("✨ Thumbnail kustom berhasil dipasang!")
				result.HasCustomThumb = true
			}
		}
	}

	return result, nil
}
