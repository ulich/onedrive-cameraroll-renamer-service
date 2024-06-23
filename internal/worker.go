package internal

import (
	"context"
	"fmt"

	"github.com/goh-chunlin/go-onedrive/onedrive"
	"golang.org/x/oauth2"
)

func Start(ctx context.Context) error {
	token, err := getOauthToken(ctx)
	if err != nil {
		return fmt.Errorf("get oauth token: %w", err)
	}

	tc := oauth2.NewClient(ctx, oauthConfig.TokenSource(ctx, token))
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
