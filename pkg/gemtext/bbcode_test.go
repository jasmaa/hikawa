package gemtext_test

import (
	"testing"

	"github.com/jasmaa/hikawa/pkg/gemtext"
	"github.com/stretchr/testify/assert"
)

// TestToBbcodeLinks tests gemtext to bbcode for links.
func TestToBbcodeLinks(t *testing.T) {
	gemtextList := []string{
		`=>https://example.com A cool website`,
		`=>gopher://example.com      An even cooler gopherhole`,
		`=> gemini://example.com A supremely cool Gemini capsule`,
		`=>   sftp://example.com`,
	}
	targetBbcodeList := []string{
		`[url=https://example.com]A cool website[/url]`,
		`[url=gopher://example.com]An even cooler gopherhole[/url]`,
		`[url=gemini://example.com]A supremely cool Gemini capsule[/url]`,
		`[url=sftp://example.com]sftp://example.com[/url]`,
	}
	for i := 0; i < len(gemtextList); i++ {
		gemtextContent := gemtextList[i]
		targetBbcodeContent := targetBbcodeList[i]
		assert.Equal(t, targetBbcodeContent, gemtext.ConvertToBbcode(gemtextContent))
	}
}

// TODO: add remaining tests when bbcode mapped
