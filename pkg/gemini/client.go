package gemini

import (
	"errors"
	"strings"
	"time"
)

// Client is a high-level client for Gemini to handle redirects and timeouts.
type Client struct {
	Timeout time.Duration
}

// ClientResponse is a high-level client response.
type ClientResponse struct {
	Response  *Response
	Url       string
	MimeTypes map[string]bool
}

// NavigatePage gets the new url and page content pointed at by `url`.
func (c *Client) NavigatePage(url string) (*ClientResponse, error) {
	type result struct {
		Response *Response
		Err      error
	}
	resChan := make(chan result, 1)

	go func() {
		for {
			r, err := ParseRequest(url)
			if err != nil {
				resChan <- result{Err: err}
				return
			}
			resp, err := r.Send()
			if err != nil {
				resChan <- result{Err: err}
				return
			}
			switch resp.Header.Status / 10 {
			case 1:
				// 1X Input
				resChan <- result{
					Response: resp,
					Err:      errors.New("hikawa does not support input yet"),
				}
			case 2:
				// 2X Success
				resChan <- result{Response: resp}
				return
			case 3:
				// 3X Redirect
				url = resp.Header.Meta
			case 4:
				// 4X Temporary Failure
				resChan <- result{Response: resp}
			case 5:
				// 5X Permanent Failure
				resChan <- result{Response: resp}
			case 6:
				// 6X Client Certificate Required
				resChan <- result{Response: resp}
			default:
				// Unrecognized status code
				resChan <- result{
					Err: errors.New("unrecognized status code"),
				}
			}
		}
	}()

	select {
	case respRes := <-resChan:
		if respRes.Err != nil {
			return nil, respRes.Err
		}
		mimeTypes := make(map[string]bool)
		if respRes.Response.Header.Status == 20 {
			// Parse MIME types from meta
			for _, mime := range strings.Split(respRes.Response.Header.Meta, ";") {
				mimeTypes[mime] = true
			}
		}
		return &ClientResponse{
			Response:  respRes.Response,
			Url:       url,
			MimeTypes: mimeTypes,
		}, nil
	case <-time.After(c.Timeout):
		return nil, errors.New("request timed out")
	}
}
