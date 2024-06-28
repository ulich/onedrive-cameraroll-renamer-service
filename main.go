package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ulich/onedrive-cameraroll-renamer-service/internal"
)

func main() {
	ctx := context.Background()

	err := internal.Start(ctx)
	if err != nil {
		slog.Error("error starting worker", "error", err)
		os.Exit(1)
	}
}
