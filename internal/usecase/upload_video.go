package usecase

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"youtube-autoposter/internal/domain"
	"youtube-autoposter/internal/infrastructure/oauth"
	"youtube-autoposter/internal/infrastructure/youtube"
)

type UploadVideoInput struct {
	FilePath      string
	ThumbnailPath string
	Title         string
	Description   string
	Tags          []string
	CategoryID    string
	PrivacyStatus string
	PublishAt     string
	SecretFile    string
	TokenFile     string
}

type UploadVideoUseCase struct {
	oauthProvider *oauth.GoogleOAuthProvider
}

func NewUploadVideoUseCase(oauthProvider *oauth.GoogleOAuthProvider) *UploadVideoUseCase {
	return &UploadVideoUseCase{
		oauthProvider: oauthProvider,
	}
}

func (uc *UploadVideoUseCase) Execute(ctx context.Context, input UploadVideoInput) (*domain.UploadResult, error) {
	if input.FilePath == "" {
		return nil, fmt.Errorf("file path video tidak boleh kosong")
	}

	title := input.Title
	if title == "" {
		baseName := filepath.Base(input.FilePath)
		title = strings.TrimSuffix(baseName, filepath.Ext(baseName))
	}

	privacy := domain.PrivacyStatus(input.PrivacyStatus)
	if privacy != domain.PrivacyPublic && privacy != domain.PrivacyPrivate && privacy != domain.PrivacyUnlisted {
		privacy = domain.PrivacyPrivate
	}

	var publishAtTime *time.Time
	if input.PublishAt != "" {
		t, err := time.Parse(time.RFC3339, input.PublishAt)
		if err != nil {
			return nil, domain.ErrInvalidPublishAt
		}
		publishAtTime = &t
	}

	video := domain.Video{
		FilePath:      input.FilePath,
		ThumbnailPath: input.ThumbnailPath,
		Title:         title,
		Description:   input.Description,
		Tags:          input.Tags,
		CategoryID:    input.CategoryID,
		Privacy:       privacy,
		PublishAt:     publishAtTime,
	}

	httpClient, err := uc.oauthProvider.GetHTTPClient(ctx, input.SecretFile, input.TokenFile)
	if err != nil {
		return nil, err
	}

	ytClient := youtube.NewClient(httpClient)
	result, err := ytClient.UploadVideo(ctx, video)
	if err != nil {
		return nil, err
	}

	return result, nil
}
