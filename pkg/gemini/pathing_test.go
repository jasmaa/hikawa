package gemini_test

import (
	"testing"

	"github.com/jasmaa/hikawa/pkg/gemini"
	"github.com/stretchr/testify/assert"
)

// TestNextUrlNewHost tests link navigation to new host
func TestNextUrlNewHost(t *testing.T) {
	targetUrl, err := gemini.NextUrl("gemini://foo.com", "gemini://bar.com")
	if assert.Nil(t, err) {
		assert.Equal(t, "gemini://bar.com", targetUrl)
	}
}

// TestNextUrlSameHostAbsolute tests link navigation to same host, absolute url
func TestNextUrlSameHostAbsolute(t *testing.T) {
	targetUrl, err := gemini.NextUrl("gemini://foo.com/1/2/", "/bar")
	if assert.Nil(t, err) {
		assert.Equal(t, "gemini://foo.com/bar", targetUrl)
	}
}

// TestNextUrlSameHostRelative tests link navigation to same host, relative url
func TestNextUrlSameHostRelative(t *testing.T) {
	targetUrl, err := gemini.NextUrl("gemini://foo.com/1/2/", "bar")
	if assert.Nil(t, err) {
		assert.Equal(t, "gemini://foo.com/1/2/bar", targetUrl)
	}
}

// TestNextUrlFile tests link navigation when paths have files
func TestNextUrlFile(t *testing.T) {
	targetUrl, err := gemini.NextUrl("gemini://foo.com/1/2/content.gmi", "bar/otherContent.gmi")
	if assert.Nil(t, err) {
		assert.Equal(t, "gemini://foo.com/1/2/bar/otherContent.gmi", targetUrl)
	}
}
