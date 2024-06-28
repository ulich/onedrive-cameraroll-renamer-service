package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/goh-chunlin/go-onedrive/onedrive"
	"golang.org/x/oauth2"
)

func Start(ctx context.Context) error {
	token, err := getStoredOrNewOAuthToken(ctx)
	if err != nil {
		return err
	}

	client, err := createClient(ctx, token)
	if err != nil {
		return err
	}

	fp := fileProcessor{
		client:         client.DriveItems,
		targetFolderId: os.Getenv("CAMERA_ROLL_TARGET_FOLDER_ID"),
	}
	return fp.processFiles(ctx)
}

func createClient(ctx context.Context, token *oauth2.Token) (*onedrive.Client, error) {
	tokenSource := oauthConfig.TokenSource(ctx, token)

	tc := oauth2.NewClient(ctx, tokenSource)
	client := onedrive.NewClient(tc)

	refreshedToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("get refreshed token: %w", err)
	}

	err = storeToken(refreshedToken)
	if err != nil {
		return nil, fmt.Errorf("store refreshed token: %w", err)
	}
	return client, nil
}
