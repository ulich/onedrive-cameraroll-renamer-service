package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ulich/onedrive-cameraroll-renamer-service/internal"
)

func main() {
	ctx := context.Background()

	err := internal.Start(ctx)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
