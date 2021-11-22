package gemtext

import (
	"regexp"
	"strings"
)

var heading3Re = regexp.MustCompile(`^#{3}\s(?P<text>[^\n]+)$`)
var heading2Re = regexp.MustCompile(`^#{2}\s(?P<text>[^\n]+)$`)
var heading1Re = regexp.MustCompile(`^#\s(?P<text>[^\n]+)$`)
var linkRe = regexp.MustCompile(`^=>\s*(?P<url>\S+)\s+(?P<text>[^\n]+)$`)
var linkRe2 = regexp.MustCompile(`^=>\s*(?P<url>\S+)$`)
var listBulletRe = regexp.MustCompile(`^\*\s(?P<text>[^\n]+)$`)
var blockquoteRe = regexp.MustCompile(`^>(?P<text>[^\n]+)$`)
var preformattedRe = regexp.MustCompile("^```(?P<alt>[^\n]*)$")

// ConvertToBbcode converts gemtext to bbcode.
func ConvertToBbcode(text string) string {
	gemtextList := strings.Split(text, "\n")
	bbcodeList := make([]string, 0)
	isPreformatMode := false
	for _, gemtext := range gemtextList {
		bbcode := gemtext
		if !isPreformatMode && heading3Re.MatchString(gemtext) {
			// Heading 3
			bbcode = heading3Re.ReplaceAllString(gemtext, "[color=green]${text}[/color]")
		} else if !isPreformatMode && heading2Re.MatchString(gemtext) {
			// Heading 2
			bbcode = heading2Re.ReplaceAllString(gemtext, "[color=blue]${text}[/color]")
		} else if !isPreformatMode && heading1Re.MatchString(gemtext) {
			// Heading 1
			bbcode = heading1Re.ReplaceAllString(gemtext, "[color=red]${text}[/color]")
		} else if !isPreformatMode && linkRe.MatchString(gemtext) {
			// Link with name
			bbcode = linkRe.ReplaceAllString(gemtext, "[url=${url}]${text}[/url]")
		} else if !isPreformatMode && linkRe2.MatchString(gemtext) {
			// Link without name
			bbcode = linkRe2.ReplaceAllString(gemtext, "[url=${url}]${url}[/url]")
		} else if !isPreformatMode && listBulletRe.MatchString(gemtext) {
			// List
			bbcode = listBulletRe.ReplaceAllString(gemtext, "[indent]* ${text}[/indent]")
		} else if !isPreformatMode && blockquoteRe.MatchString(gemtext) {
			// Blockquote
			bbcode = blockquoteRe.ReplaceAllString(gemtext, "[center]${text}[/center]")
		} else if preformattedRe.MatchString(gemtext) {
			// Preformat
			if isPreformatMode {
				bbcode = "[code]"
			} else {
				bbcode = "[/code]"
			}
			isPreformatMode = !isPreformatMode
		}
		bbcodeList = append(bbcodeList, bbcode)
	}
	return strings.Join(bbcodeList, "\n")
}
