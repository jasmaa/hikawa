package gemini_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/jasmaa/hikawa/pkg/gemini"
	"github.com/stretchr/testify/assert"
)

// TestParseRequestNoSubpath tests request parsing for url without subpath.
func TestParseRequestNoSubpath(t *testing.T) {
	r, err := gemini.ParseRequest("gemini://example.com")
	if assert.Nil(t, err) {
		assert.Equal(t, "example.com", r.Host)
		assert.Equal(t, "1965", r.Port)
		assert.Equal(t, "/", r.Path)
	}
}

// TestParseRequestSubpath tests request parsing for url with subpath.
func TestParseRequestSubpath(t *testing.T) {
	r, err := gemini.ParseRequest("gemini://example.com/foo/bar/page.gmi")
	if assert.Nil(t, err) {
		assert.Equal(t, "example.com", r.Host)
		assert.Equal(t, "1965", r.Port)
		assert.Equal(t, "/foo/bar/page.gmi", r.Path)
	}
}

// TestReadResponseSuccess tests reading successful response.
func TestReadResponseSuccess(t *testing.T) {
	status := gemini.STATUS_SUCCESS
	meta := "text/gemini"
	body := "this is some text\n" +
		"this is the second line\n" +
		"\r\n"
	rawresp := fmt.Sprintf("%d %s\r\n%s", status, meta, body)
	conn := bytes.NewReader([]byte(rawresp))
	resp, err := gemini.ReadResponse(conn)
	if assert.Nil(t, err) {
		assert.Equal(t, status, resp.Header.Status)
		assert.Equal(t, meta, resp.Header.Meta)
		assert.Equal(t, body, resp.Body)
	}
}

// TestReadResponseNotFound tests reading not found response.
func TestReadResponseNotFound(t *testing.T) {
	status := gemini.STATUS_NOT_FOUND
	meta := ""
	body := ""
	rawresp := fmt.Sprintf("%d %s\r\n%s", status, meta, body)
	conn := bytes.NewReader([]byte(rawresp))
	resp, err := gemini.ReadResponse(conn)
	if assert.Nil(t, err) {
		assert.Equal(t, status, resp.Header.Status)
	}
}

// TestReadResponseMetaGreaterThan1024 tests erroring when meta > 1024 bytes.
func TestReadResponseMetaGreaterThan1024(t *testing.T) {
	status := gemini.STATUS_SUCCESS
	meta := strings.Repeat("a", 1025)
	body := "this is some text\n" +
		"this is the second line\n" +
		"\r\n"
	rawresp := fmt.Sprintf("%d %s\r\n%s", status, meta, body)
	conn := bytes.NewReader([]byte(rawresp))
	_, err := gemini.ReadResponse(conn)
	assert.NotNil(t, err)
}

// TestReadResponseEmptyMeta tests response with empty meta
func TestReadResponseEmptyMeta(t *testing.T) {
	status := gemini.STATUS_SUCCESS
	rawresp := fmt.Sprintf("%d\r\n", status)
	conn := bytes.NewReader([]byte(rawresp))
	_, err := gemini.ReadResponse(conn)
	assert.NotNil(t, err)
}
