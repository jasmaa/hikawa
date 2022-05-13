package gemtext

import (
	"regexp"
)

var heading3Re = regexp.MustCompile(`^#{3}\s(?P<text>[^\n]+)$`)
var heading2Re = regexp.MustCompile(`^#{2}\s(?P<text>[^\n]+)$`)
var heading1Re = regexp.MustCompile(`^#\s(?P<text>[^\n]+)$`)
var linkRe = regexp.MustCompile(`^=>\s*(?P<url>\S+)\s+(?P<text>[^\n]+)$`)
var linkRe2 = regexp.MustCompile(`^=>\s*(?P<url>\S+)$`)
var listBulletRe = regexp.MustCompile(`^\*\s(?P<text>[^\n]+)$`)
var blockquoteRe = regexp.MustCompile(`^>(?P<text>[^\n]+)$`)
var preformattedRe = regexp.MustCompile("^```(?P<alt>[^\\n]*)$")
