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

// Request requests with a url and returns a Response.
func Request(requestUrl string) (*Response, error) {
	u, err := url.ParseRequestURI(requestUrl)
	if err != nil {
		return nil, err
	}
	u.RawQuery = url.PathEscape(u.RawQuery)

	var host, port string
	hostport := strings.Split(u.Host, ":")
	if len(hostport) < 1 {
		return nil, errors.New("no hostname provided")
	}
	host = hostport[0]
	if len(hostport) > 1 {
		port = hostport[1]
	} else {
		port = "1965"
	}

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", host, port), conf)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(u.String() + "\r\n"))
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
