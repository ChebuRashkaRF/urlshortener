package storage

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestURLStorage(t *testing.T) {
	urlStorage := NewURLStorage()

	url, ok := urlStorage.Get("test")
	require.False(t, ok, "Expected false")
	assert.Empty(t, url, "Expected empty string")

	urlStorage.Set("test", "https://example.com")
	url, ok = urlStorage.Get("test")
	require.True(t, ok, "Expected true")
	assert.Equal(t, "https://example.com", url, "Expected https://example.com")
}

func TestGetURLMap(t *testing.T) {
	urlStorage := NewURLStorage()

	urlStorage.Set("test1", "https://example1.com")
	urlStorage.Set("test2", "https://example2.com")
	urlStorage.Set("test3", "https://example3.com")

	allURLs := urlStorage.GetURLMap()

	require.Len(t, allURLs, 3, "Expected length of 3")

	assert.Equal(t, "https://example1.com", allURLs["test1"])
	assert.Equal(t, "https://example2.com", allURLs["test2"])
	assert.Equal(t, "https://example3.com", allURLs["test3"])
}
