package storage

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestURLStorage(t *testing.T) {
	tempFile, err := os.CreateTemp("", "urlstorage_test_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	urlStorage, err := NewURLStorage(tempFile.Name())
	require.NoError(t, err)
	defer urlStorage.Close()

	url, ok := urlStorage.Get("test")
	require.False(t, ok, "Expected false")
	assert.Empty(t, url, "Expected empty string")

	urlStorage.Set("test", "https://example.com")
	url, ok = urlStorage.Get("test")
	require.True(t, ok, "Expected true")
	assert.Equal(t, "https://example.com", url, "Expected https://example.com")

	reloadedStorage, err := NewURLStorage(tempFile.Name())
	require.NoError(t, err)
	defer reloadedStorage.Close()

	url, ok = reloadedStorage.Get("test")
	require.True(t, ok, "Expected true")
	assert.Equal(t, "https://example.com", url, "Expected https://example.com")
}

func TestGetURLMap(t *testing.T) {
	tempFile, err := os.CreateTemp("", "urlstorage_test_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	urlStorage, err := NewURLStorage(tempFile.Name())
	require.NoError(t, err)
	defer urlStorage.Close()

	urlStorage.Set("test1", "https://example1.com")
	urlStorage.Set("test2", "https://example2.com")
	urlStorage.Set("test3", "https://example3.com")

	allURLs := urlStorage.GetURLMap()

	require.Len(t, allURLs, 3, "Expected length of 3")

	assert.Equal(t, "https://example1.com", allURLs["test1"])
	assert.Equal(t, "https://example2.com", allURLs["test2"])
	assert.Equal(t, "https://example3.com", allURLs["test3"])

	// Reload storage from file and verify data
	reloadedStorage, err := NewURLStorage(tempFile.Name())
	require.NoError(t, err)
	defer reloadedStorage.Close()

	allURLs = reloadedStorage.GetURLMap()
	require.Len(t, allURLs, 3, "Expected length of 3")

	assert.Equal(t, "https://example1.com", allURLs["test1"])
	assert.Equal(t, "https://example2.com", allURLs["test2"])
	assert.Equal(t, "https://example3.com", allURLs["test3"])
}

func TestURLStorageUUID(t *testing.T) {
	tempFile, err := os.CreateTemp("", "urlstorage_test_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	urlStorage, err := NewURLStorage(tempFile.Name())
	require.NoError(t, err)
	defer urlStorage.Close()

	assert.Equal(t, 1, urlStorage.UUID, "Expected initial UUID to be 1")

	urlStorage.Set("test1", "https://example1.com")
	assert.Equal(t, 2, urlStorage.UUID, "Expected UUID to be incremented to 2")

	urlStorage.Set("test2", "https://example2.com")
	assert.Equal(t, 3, urlStorage.UUID, "Expected UUID to be incremented to 3")

	// Reload storage from file and verify UUID
	reloadedStorage, err := NewURLStorage(tempFile.Name())
	require.NoError(t, err)
	defer reloadedStorage.Close()

	assert.Equal(t, 3, reloadedStorage.UUID, "Expected reloaded UUID to be 3")
}
