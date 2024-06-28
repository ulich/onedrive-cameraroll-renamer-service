package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/goh-chunlin/go-onedrive/onedrive"
	"github.com/stretchr/testify/mock"
)

const targetFolderId = "aTargetFolderId"

func TestProcessFiles(t *testing.T) {

	t.Run("moves files", func(t *testing.T) {
		client := newMockClient([]*onedrive.DriveItem{
			{Id: "id1", Name: "20190519_205940.jpg"},
			{Id: "id2", Name: "20190519_205941.jpg"},
		})

		client.stubMoveSuccess("id1", targetFolderId)
		client.stubMoveSuccess("id2", targetFolderId)

		fp := fileProcessor{client: client, targetFolderId: targetFolderId}

		err := fp.processFiles(context.Background())
		if err != nil {
			t.Errorf("wanted no error, got %v", err)
		}
	})

	t.Run("renames and moves files", func(t *testing.T) {
		client := newMockClient([]*onedrive.DriveItem{
			{Id: "id1", Name: "20190519_205940.JPG"}, // should be renamed
			{Id: "id2", Name: "20190519_205941.jpg"}, // should NOT be renamed
			{Id: "id3", Name: "20190519_205942.JPG"}, // should be renamed
		})

		client.stubRenameSuccess("id1", "20190519_205940.jpg")
		client.stubMoveSuccess("id1", targetFolderId)

		client.stubMoveSuccess("id2", targetFolderId)

		client.stubRenameSuccess("id3", "20190519_205942.jpg")
		client.stubMoveSuccess("id3", targetFolderId)

		fp := fileProcessor{client: client, targetFolderId: targetFolderId}

		err := fp.processFiles(context.Background())
		if err != nil {
			t.Errorf("wanted no error, got %v", err)
		}
	})

	t.Run("adds a suffix number if a file in the target folder already exists with the same name", func(t *testing.T) {
		client := newMockClient([]*onedrive.DriveItem{
			{Id: "id1", Name: "20190519_205940.jpg"},
		})

		client.stubMoveFailureAlreadyExists("id1", targetFolderId)
		client.stubRenameSuccess("id1", "20190519_205940_1.jpg")
		client.stubMoveFailureAlreadyExists("id1", targetFolderId)
		client.stubRenameSuccess("id1", "20190519_205940_2.jpg")
		client.stubMoveFailureAlreadyExists("id1", targetFolderId)
		client.stubRenameSuccess("id1", "20190519_205940_3.jpg")
		client.stubMoveSuccess("id1", targetFolderId)

		fp := fileProcessor{client: client, targetFolderId: targetFolderId}

		err := fp.processFiles(context.Background())
		if err != nil {
			t.Errorf("wanted no error, got %v", err)
		}
	})
}

func newMockClient(driveItems []*onedrive.DriveItem) *mockClient {
	client := mockClient{}
	client.On("ListSpecial", mock.Anything, onedrive.CameraRoll).Return(&onedrive.OneDriveDriveItemsResponse{
		DriveItems: driveItems,
	}, nil)
	return &client
}

type mockClient struct {
	mock.Mock
}

func (m *mockClient) ListSpecial(ctx context.Context, special onedrive.DriveSpecialFolder) (*onedrive.OneDriveDriveItemsResponse, error) {
	args := m.Called(ctx, special)
	return args.Get(0).(*onedrive.OneDriveDriveItemsResponse), args.Error(1)
}

func (m *mockClient) Move(ctx context.Context, parentFolderId string, itemId string, targetFolderId string) (*onedrive.MoveItemResponse, error) {
	args := m.Called(ctx, parentFolderId, itemId, targetFolderId)
	return args.Get(0).(*onedrive.MoveItemResponse), args.Error(1)
}

func (c *mockClient) stubMoveSuccess(fileId, targetFolderId string) *mock.Call {
	return c.On("Move", mock.Anything, "", fileId, targetFolderId).
		Return(&onedrive.MoveItemResponse{}, nil).
		Once()
}

func (c *mockClient) stubMoveFailureAlreadyExists(fileId, targetFolderId string) *mock.Call {
	return c.On("Move", mock.Anything, "", fileId, targetFolderId).
		Return(&onedrive.MoveItemResponse{}, fmt.Errorf("nameAlreadyExists - the file already exists in the target folder")).
		Once()
}

func (m *mockClient) Rename(ctx context.Context, parentFolderId string, itemId string, newName string) (*onedrive.RenameItemResponse, error) {
	args := m.Called(ctx, parentFolderId, itemId, newName)
	return args.Get(0).(*onedrive.RenameItemResponse), args.Error(1)
}

func (c *mockClient) stubRenameSuccess(fileId, newName string) *mock.Call {
	return c.On("Rename", mock.Anything, "", fileId, newName).
		Return(&onedrive.RenameItemResponse{}, nil).
		Once()
}
