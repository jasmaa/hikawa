package gemini

import (
	"net/url"
	"path"
)

// NextUrl constructs next URL for link navigation
func NextUrl(currentUrl string, newUrl string) (string, error) {
	u, err := url.Parse(currentUrl)
	if err != nil {
		return "", err
	}
	newU, err := url.Parse(newUrl)
	if err != nil {
		return "", err
	}
	if len(newU.Scheme) == 0 {
		// Same host
		if path.IsAbs(newU.Path) {
			// Absolute url
			u.Path = newU.Path
		} else {
			// Relative url
			dir, _ := path.Split(u.Path)
			u.Path = path.Join(dir, newU.Path)
		}
		u.RawQuery = newU.RawQuery
		return u.String(), nil
	} else {
		// New host
		return newUrl, nil
	}
}
