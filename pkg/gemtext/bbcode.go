package gemtext

import (
	"strings"
)

// ConvertToBbcode converts gemtext to bbcode.
func ConvertToBbcode(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	gemtextList := strings.Split(text, "\n")
	bbcodeList := make([]string, 0)
	isPreformatMode := false
	for _, gemtext := range gemtextList {
		bbcode := gemtext
		if !isPreformatMode && heading3Re.MatchString(gemtext) {
			// Heading 3
			bbcode = heading3Re.ReplaceAllString(gemtext, "[font=res://assets/fonts/Ubuntu/UbuntuR_H3.tres]${text}[/font]")
		} else if !isPreformatMode && heading2Re.MatchString(gemtext) {
			// Heading 2
			bbcode = heading2Re.ReplaceAllString(gemtext, "[font=res://assets/fonts/Ubuntu/UbuntuR_H2.tres]${text}[/font]")
		} else if !isPreformatMode && heading1Re.MatchString(gemtext) {
			// Heading 1
			bbcode = heading1Re.ReplaceAllString(gemtext, "[font=res://assets/fonts/Ubuntu/UbuntuR_H1.tres]${text}[/font]")
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
			if !isPreformatMode {
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
