package domain

import (
	"context"
	"errors"
	"time"
)

// PrivacyStatus represents YouTube video privacy settings
type PrivacyStatus string

const (
	PrivacyPublic   PrivacyStatus = "public"
	PrivacyPrivate  PrivacyStatus = "private"
	PrivacyUnlisted PrivacyStatus = "unlisted"
)

// Video represents a YouTube video entity
type Video struct {
	FilePath      string
	ThumbnailPath string
	Title         string
	Description   string
	Tags          []string
	CategoryID    string
	Privacy       PrivacyStatus
	PublishAt     *time.Time
}

// UploadResult represents the result of a successful upload
type UploadResult struct {
	ID        string
	WatchURL  string
	Duration  time.Duration
	HasCustomThumb bool
}

// Errors
var (
	ErrFileNotFound    = errors.New("file video tidak ditemukan")
	ErrInvalidPublishAt = errors.New("format jadwal publish-at tidak valid (gunakan RFC3339)")
)

// YouTubeService defines the interface for uploading videos & thumbnails to YouTube
type YouTubeService interface {
	UploadVideo(ctx context.Context, video Video) (*UploadResult, error)
}

// Authenticator defines the interface for OAuth2 authentication
type Authenticator interface {
	GetHTTPClient(ctx context.Context, secretFile, tokenFile string) (interface{}, error)
}
