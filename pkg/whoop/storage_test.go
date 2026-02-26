package whoop

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStorage_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	fs := NewFileStorageWithDir(dir, "test-passphrase")

	data := &TokenData{
		AccessToken:  "access123",
		RefreshToken: "refresh456",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		ExpiresAt:    1700000000,
		CreatedAt:    1699996400,
	}

	require.NoError(t, fs.Save(data))

	loaded, err := fs.Load()
	require.NoError(t, err)

	assert.Equal(t, data.AccessToken, loaded.AccessToken)
	assert.Equal(t, data.RefreshToken, loaded.RefreshToken)
	assert.Equal(t, data.ClientID, loaded.ClientID)
	assert.Equal(t, data.ClientSecret, loaded.ClientSecret)
	assert.InDelta(t, data.ExpiresAt, loaded.ExpiresAt, 0.01)
	assert.InDelta(t, data.CreatedAt, loaded.CreatedAt, 0.01)
}

func TestFileStorage_WrongPassphrase(t *testing.T) {
	dir := t.TempDir()
	fs1 := NewFileStorageWithDir(dir, "correct-passphrase")
	fs2 := NewFileStorageWithDir(dir, "wrong-passphrase")

	data := &TokenData{
		AccessToken: "access123",
		ClientID:    "client-id",
	}

	require.NoError(t, fs1.Save(data))

	_, err := fs2.Load()
	assert.Error(t, err)
}

func TestFileStorage_Delete(t *testing.T) {
	dir := t.TempDir()
	fs := NewFileStorageWithDir(dir, "test-passphrase")

	data := &TokenData{AccessToken: "access123", ClientID: "client-id"}
	require.NoError(t, fs.Save(data))

	require.NoError(t, fs.Delete())

	_, err := fs.Load()
	assert.Error(t, err)
}

func TestFileStorage_LoadNonexistent(t *testing.T) {
	dir := t.TempDir()
	fs := NewFileStorageWithDir(dir, "test-passphrase")

	_, err := fs.Load()
	assert.Error(t, err)
}

func TestFileStorage_DirectoryPermissions(t *testing.T) {
	dir := t.TempDir()
	fs := NewFileStorageWithDir(dir+"/nested/dir", "test-passphrase")

	data := &TokenData{AccessToken: "access123", ClientID: "client-id"}
	require.NoError(t, fs.Save(data))

	// Verify the nested directory was created
	info, err := os.Stat(dir + "/nested/dir")
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}
