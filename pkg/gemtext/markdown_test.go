package gemtext_test

import (
	"testing"

	"github.com/jasmaa/hikawa/pkg/gemtext"
	"github.com/stretchr/testify/assert"
)

// TestMarkdownLinks tests gemtext to markdown for links.
func TestMarkdownLinks(t *testing.T) {
	gemtextList := []string{
		`=>https://example.com A cool website`,
		`=>gopher://example.com      An even cooler gopherhole`,
		`=> gemini://example.com A supremely cool Gemini capsule`,
		`=>   sftp://example.com`,
	}
	targetMdList := []string{
		`[A cool website](https://example.com)`,
		`[An even cooler gopherhole](gopher://example.com)`,
		`[A supremely cool Gemini capsule](gemini://example.com)`,
		`[sftp://example.com](sftp://example.com)`,
	}
	for i := 0; i < len(gemtextList); i++ {
		gemtextContent := gemtextList[i]
		targetBbcodeContent := targetMdList[i]
		assert.Equal(t, targetBbcodeContent, gemtext.ConvertToMarkdown(gemtextContent))
	}
}

// TODO: add remaining tests when markdown mapped
