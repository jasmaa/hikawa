package gemini

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

// Client is a high-level client for Gemini to handle redirects and timeouts.
type Client struct {
	Timeout      time.Duration
	MaxRetries   int
	MaxRedirects int
}

// ClientResponse is a high-level client response.
type ClientResponse struct {
	Response  *Response
	Url       string
	MimeTypes map[string]bool
}

// MakeClient makes the default client
func MakeClient() Client {
	return Client{
		Timeout:      7 * time.Second,
		MaxRetries:   3,
		MaxRedirects: 5,
	}
}

// NavigatePage gets the new url and page content pointed at by `url`.
func (c *Client) NavigatePage(rawurl string) (*ClientResponse, error) {
	type result struct {
		Response *Response
		Err      error
	}
	resChan := make(chan result, 1)

	go func() {
		tries := 0
		redirects := 0
		for tries-1 < c.MaxRetries && redirects < c.MaxRedirects {
			resp, err := Request(rawurl)
			if err != nil {
				resChan <- result{Err: err}
				return
			}
			switch resp.Header.Status / 10 {
			case 1:
				// 1X Input
				resChan <- result{Response: resp}
				tries++
				return
			case 2:
				// 2X Success
				resChan <- result{Response: resp}
				tries++
				return
			case 3:
				// 3X Redirect
				u, _ := url.ParseRequestURI(resp.Header.Meta)
				if len(u.Scheme) == 0 {
					// Relative url
					prevUrl, _ := url.ParseRequestURI(rawurl)
					prevUrl.Path = resp.Header.Meta
					rawurl = prevUrl.String()
				} else {
					// Absolute url
					rawurl = resp.Header.Meta
				}
				redirects++
			case 4:
				// 4X Temporary Failure
				resChan <- result{Response: resp}
				tries++
			case 5:
				// 5X Permanent Failure
				resChan <- result{Response: resp}
				tries++
				return
			case 6:
				// 6X Client Certificate Required
				resChan <- result{Response: resp}
				tries++
				return
			default:
				// Unrecognized status code
				resChan <- result{
					Err: errors.New("unrecognized status code"),
				}
				tries++
				return
			}

			// Sleep before retrying
			time.Sleep(100 * time.Millisecond)
		}

		if tries >= c.MaxRetries {
			resChan <- result{
				Err: errors.New("exceeded maximum number of retries"),
			}
		} else if redirects > -c.MaxRedirects {
			resChan <- result{
				Err: errors.New("exceeded maximum number of redirects"),
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
			Url:       rawurl,
			MimeTypes: mimeTypes,
		}, nil
	case <-time.After(c.Timeout):
		return nil, errors.New("request timed out")
	}
}
