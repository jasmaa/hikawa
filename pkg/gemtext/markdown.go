package gemtext

import (
	"strings"
)

// ConvertToMarkdown converts gemtext to Dear ImGui markdown.
// https://github.com/juliettef/imgui_markdown
func ConvertToMarkdown(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	gemtextList := strings.Split(text, "\n")
	mdList := make([]string, 0)
	isPreformatMode := false
	for _, gemtext := range gemtextList {
		md := gemtext
		if !isPreformatMode && heading3Re.MatchString(gemtext) {
			// Heading 3
			// no-op
		} else if !isPreformatMode && heading2Re.MatchString(gemtext) {
			// Heading 2
			// no-op
		} else if !isPreformatMode && heading1Re.MatchString(gemtext) {
			// Heading 1
			// no-op
		} else if !isPreformatMode && linkRe.MatchString(gemtext) {
			// Link with name
			md = linkRe.ReplaceAllString(gemtext, "[${text}](${url})")
		} else if !isPreformatMode && linkRe2.MatchString(gemtext) {
			// Link without name
			md = linkRe2.ReplaceAllString(gemtext, "[${url}](${url})")
		} else if !isPreformatMode && listBulletRe.MatchString(gemtext) {
			// List
			md = listBulletRe.ReplaceAllString(gemtext, "  * ${text}")
		} else if !isPreformatMode && blockquoteRe.MatchString(gemtext) {
			// Blockquote
			// no-op
		} else if preformattedRe.MatchString(gemtext) {
			// Preformat
			// Wraps in preformat block in ```. This does not get renderered in the markdown widget.
			md = "```"
			isPreformatMode = !isPreformatMode
		}
		mdList = append(mdList, md)
	}
	return strings.Join(mdList, "\n")
}
