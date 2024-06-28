package internal

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/goh-chunlin/go-onedrive/onedrive"
)

type fileProcessor struct {
	client         driveItemsClient
	targetFolderId string
}

type driveItemsClient interface {
	ListSpecial(ctx context.Context, special onedrive.DriveSpecialFolder) (*onedrive.OneDriveDriveItemsResponse, error)
	Move(ctx context.Context, parentFolderId string, itemId string, targetFolderId string) (*onedrive.MoveItemResponse, error)
	Rename(ctx context.Context, parentFolderId string, itemId string, newName string) (*onedrive.RenameItemResponse, error)
}

func (fp fileProcessor) processFiles(ctx context.Context) error {
	pictures, err := fp.client.ListSpecial(ctx, onedrive.CameraRoll)
	if err != nil {
		return fmt.Errorf("list pictures in camera roll folder: %w", err)
	}

	for _, picture := range pictures.DriveItems {
		if picture.Id != fp.targetFolderId {
			err := fp.processFile(ctx, picture)
			if err != nil {
				slog.Warn("failed to process file", "file", picture.Name, "error", err.Error())
			}
		}
	}

	return nil
}

func (fp fileProcessor) processFile(ctx context.Context, file *onedrive.DriveItem) error {
	newName, err := calcNewFilename(file.Name)
	if err != nil {
		return fmt.Errorf("calculate new filename for %s: %w", file.Name, err)
	}

	err = fp.mv(ctx, file, newName, newName, 1)
	if err != nil {
		return err
	}

	return nil
}

func (fp fileProcessor) mv(ctx context.Context, file *onedrive.DriveItem, firstNewName, newName string, attempt int) error {
	if newName != file.Name {
		slog.Info("renaming file", "source", file.Name, "destination", newName)
		_, err := fp.client.Rename(ctx, "", file.Id, newName)
		if err != nil {
			return fmt.Errorf("rename file %s to %s: %w", file.Name, newName, err)
		}
	}

	slog.Info("moving file to target folder", "old_name", file.Name, "new_name", newName)
	_, err := fp.client.Move(ctx, "", file.Id, fp.targetFolderId)
	if err != nil {
		if strings.HasPrefix(err.Error(), "nameAlreadyExists -") {
			fileExtension := filepath.Ext(firstNewName)
			baseFilename := strings.TrimSuffix(firstNewName, fileExtension)
			evenNewerFilename := fmt.Sprintf("%v_%v%s", baseFilename, attempt, fileExtension)

			slog.Info("file already exists in target folder, choosing different filename", "file", file.Name)
			return fp.mv(ctx, file, firstNewName, evenNewerFilename, attempt+1)
		}
		return fmt.Errorf("move file %s to target folder: %w", file.Name, err)
	}
	return nil
}
