package internal

import (
	"context"
	"fmt"

	"github.com/goh-chunlin/go-onedrive/onedrive"
	"golang.org/x/oauth2"
)

func Start(ctx context.Context) error {
	token, err := getStoredOrNewOAuthToken(ctx)
	if err != nil {
		return err
	}

	tokenSource := oauthConfig.TokenSource(ctx, token)

	tc := oauth2.NewClient(ctx, tokenSource)
	client := onedrive.NewClient(tc)

	response, err := client.Drives.List(ctx)
	if err != nil {
		return fmt.Errorf("list drives: %w", err)
	}

	for _, drive := range response.Drives {
		fmt.Println(*drive)
	}

	return nil
}
