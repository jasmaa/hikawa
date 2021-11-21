package gemini

import (
	"fmt"
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
			case 1:
				// 1X Input
				textChan <- "Input not supported yet by Hikawa"
			case 2:
				// 2X Success
				textChan <- resp.Body
				return
			case 3:
				// 3X Redirect
				url = resp.Header.Meta
			case 4:
				// 4X Temporary Failure
				textChan <- fmt.Sprintf(
					"Temporary Failure (%d %s): %s",
					resp.Header.Status,
					statusCodeToMessage(resp.Header.Status),
					resp.Header.Meta,
				)
			case 5:
				// 5X Permanent Failure
				textChan <- fmt.Sprintf(
					"Permanent Failure (%d %s): %s",
					resp.Header.Status,
					statusCodeToMessage(resp.Header.Status),
					resp.Header.Meta,
				)
			case 6:
				// 6X Client Certificate Required
				textChan <- fmt.Sprintf(
					"Client Certificate Required (%d %s): %s",
					resp.Header.Status,
					statusCodeToMessage(resp.Header.Status),
					resp.Header.Meta,
				)
			default:
				// Unrecognized status code
				textChan <- fmt.Sprintf("Unrecognized Code (%d): %s", resp.Header.Status, resp.Header.Meta)
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
