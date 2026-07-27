package usecase

import (
	"context"

	"youtube-autoposter/internal/domain"
	"youtube-autoposter/internal/infrastructure/oauth"
	"youtube-autoposter/internal/infrastructure/youtube"
)

type GetChannelInfoUseCase struct {
	oauthProvider *oauth.GoogleOAuthProvider
}

func NewGetChannelInfoUseCase(oauthProvider *oauth.GoogleOAuthProvider) *GetChannelInfoUseCase {
	return &GetChannelInfoUseCase{
		oauthProvider: oauthProvider,
	}
}

func (uc *GetChannelInfoUseCase) Execute(ctx context.Context, secretFile, tokenFile string) (*domain.ChannelInfo, error) {
	httpClient, err := uc.oauthProvider.GetHTTPClient(ctx, secretFile, tokenFile)
	if err != nil {
		return nil, err
	}

	ytClient := youtube.NewClient(httpClient)
	return ytClient.GetChannelInfo(ctx)
}
