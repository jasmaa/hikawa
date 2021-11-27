package browsing_test

import (
	"testing"

	"github.com/jasmaa/hikawa/pkg/browsing"
	"github.com/stretchr/testify/assert"
)

// TestEmpty tests empty history.
func TestEmpty(t *testing.T) {
	h := browsing.NewHistory()
	_, err := h.GetCurrentUrl()
	assert.NotNil(t, err)
	assert.False(t, h.CanGoBack())
	assert.False(t, h.CanGoForward())
}

// TestPushOnce tests pushing one page.
func TestPush(t *testing.T) {
	h := browsing.NewHistory()
	nextUrl := "gemini://example.com"
	h.Push(nextUrl)
	currentUrl, err := h.GetCurrentUrl()
	if assert.Nil(t, err) {
		assert.Equal(t, nextUrl, currentUrl)
	}
	assert.False(t, h.CanGoBack())
	assert.False(t, h.CanGoForward())
}

// TestPushTwice tests pushing two pages.
func TestPushTwice(t *testing.T) {
	h := browsing.NewHistory()
	nextUrl1 := "gemini://example.com/foo"
	nextUrl2 := "gemini://example.com/bar"
	h.Push(nextUrl1)
	h.Push(nextUrl2)
	currentUrl, err := h.GetCurrentUrl()
	if assert.Nil(t, err) {
		assert.Equal(t, nextUrl2, currentUrl)
	}
	assert.True(t, h.CanGoBack())
	assert.False(t, h.CanGoForward())
}

// TestPushTwiceAndBackOnce tests pushing two pages and going back once.
func TestPushTwiceAndBackOnce(t *testing.T) {
	h := browsing.NewHistory()
	nextUrl1 := "gemini://example.com/foo"
	nextUrl2 := "gemini://example.com/bar"
	h.Push(nextUrl1)
	h.Push(nextUrl2)
	h.GoBack()
	currentUrl, err := h.GetCurrentUrl()
	if assert.Nil(t, err) {
		assert.Equal(t, nextUrl1, currentUrl)
	}
	assert.False(t, h.CanGoBack())
	assert.True(t, h.CanGoForward())
}

// TestPushTwiceAndBackOnceAndPushOnce tests pushing two pages, going back, and pushing a new page.
func TestPushTwiceAndBackOnceAndPushOnce(t *testing.T) {
	h := browsing.NewHistory()
	nextUrl1 := "gemini://example.com/foo"
	nextUrl2 := "gemini://example.com/bar"
	nextUrl3 := "gemini://example.com/thirdPage"
	h.Push(nextUrl1)
	h.Push(nextUrl2)
	h.GoBack()
	h.Push(nextUrl3)
	currentUrl, err := h.GetCurrentUrl()
	if assert.Nil(t, err) {
		assert.Equal(t, nextUrl3, currentUrl)
	}
	assert.True(t, h.CanGoBack())
	assert.False(t, h.CanGoForward())
}
