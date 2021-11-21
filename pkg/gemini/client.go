package gemini

import (
	"time"
)

// Client is a high-level client for Gemini.
type Client struct {
	Timeout time.Duration
}

// NavigatePage gets the new url and page content pointed at by `url`.
func (c *Client) NavigatePage(url string) (string, string) {
	textChan := make(chan string, 1)

	go func() {
		for {
			r, err := ParseRequest(url)
			if err != nil {
				textChan <- err.Error()
				return
			}
			resp, err := r.Send()
			if err != nil {
				textChan <- err.Error()
				return
			}
			switch resp.Header.Status / 10 {
			case 2:
				// 2X Success
				textChan <- resp.Body
				return
			case 3:
				// 3X Redirect
				url = resp.Header.Meta
			default:
				// Else, display meta
				textChan <- resp.Header.Meta
				return
			}
		}
	}()

	select {
	case res := <-textChan:
		return url, res
	case <-time.After(c.Timeout):
		return url, "request timed out"
	}
}
