package gemini

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

// Request is a Gemini request.
type Request struct {
	Host string
	Port string
	Path string
}

// ResponseHeader is a Gemini response header.
type ResponseHeader struct {
	Status int
	Meta   string
}

// Response is a Gemini response.
type Response struct {
	Header ResponseHeader
	Body   string
}

// ParseRequest parses a URL into a Request.
func ParseRequest(rawurl string) (*Request, error) {
	var host, port, path string

	u, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return nil, err
	}

	// Validate scheme
	if u.Scheme != "gemini" {
		return nil, errors.New("scheme was not `gemini`")
	}

	// Set host and port
	hostport := strings.Split(u.Host, ":")
	host = hostport[0]
	if len(hostport) == 2 {
		port = hostport[1]
	} else {
		port = "1965"
	}

	// Set path
	path = u.Path
	if len(u.Path) == 0 {
		path = "/"
	}

	return &Request{
		Host: host,
		Port: port,
		Path: path,
	}, nil
}

// Send executes the Request and returns a Response.
func (r *Request) Send() (*Response, error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", r.Host, r.Port), conf)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(fmt.Sprintf("gemini://%s%s\r\n", r.Host, r.Path)))
	if err != nil {
		return nil, err
	}

	resp, err := ReadResponse(conn)
	return resp, err
}

// ReadResponse reads Response from connection reader.
func ReadResponse(conn io.Reader) (*Response, error) {
	// Read status
	buffer := make([]byte, 2)
	_, err := io.ReadFull(conn, buffer)
	if err != nil {
		return nil, err
	}
	status, err := strconv.Atoi(string(buffer[:2]))
	if err != nil {
		return nil, err
	}

	buffer = make([]byte, 1)
	_, err = io.ReadFull(conn, buffer)
	if err != nil {
		return nil, err
	}
	if buffer[0] != 0x20 {
		// TODO: Fix status reads for responses without meta
		return nil, errors.New("no meta found")
	}

	// Read meta, meta (1024) + CRLF (2)
	buffer = make([]byte, 1026)
	meta, beginIdx, endIdx, err := readMeta(conn, buffer)
	if err != nil {
		return nil, err
	}

	// Read body on 2X status
	rawbody := make([]byte, len(buffer[beginIdx:endIdx]))
	copy(rawbody, buffer[beginIdx:endIdx])
	n := 0
	if status/10 == 2 {
		for err != io.EOF {
			n, err = conn.Read(buffer)
			rawbody = append(rawbody, buffer[:n]...)
		}
	}

	return &Response{
		Header: ResponseHeader{
			Status: status,
			Meta:   meta,
		},
		Body: string(rawbody),
	}, nil
}

func readMeta(conn io.Reader, buffer []byte) (string, int, int, error) {
	n, _ := io.ReadFull(conn, buffer)
	var idx int
	for idx = 1; idx < n; idx++ {
		if buffer[idx-1] == 0x0D && buffer[idx] == 0x0A {
			meta := string(buffer[:idx-1])
			return meta, idx + 1, n, nil
		}
	}
	return "", -1, -1, errors.New("meta greater than 1024 bytes")
}
