package goBoom

import (
	"bytes"
	"io"
	"mime/multipart"
	"path/filepath"
)

type FilesystemService struct {
	c *Client
}

func newFilesystemService(c *Client) *FilesystemService {
	i := &FilesystemService{}
	if c == nil {
		i.c = NewClient(nil)
	} else {
		i.c = c
	}

	return i
}

// Upload
func (s *FilesystemService) Upload(fname string, input io.Reader) (int, *ItemInfo, error) {

	var bodyBuf bytes.Buffer
	writer := multipart.NewWriter(&bodyBuf)

	part, err := writer.CreateFormFile("file", filepath.Base(fname))
	if err != nil {
		return 0, nil, err
	}

	_, err = io.Copy(part, input)
	if err != nil {
		return 0, nil, err
	}

	err = writer.Close()
	if err != nil {
		return 0, nil, err
	}

	// prepare request
	req, err := s.c.NewReaderRequest("POST", "ul", &bodyBuf, writer.FormDataContentType())
	if err != nil {
		return 0, nil, err
	}

	// set  token
	q := req.URL.Query()
	q.Set("token", s.c.User.session)
	q.Set("parent", "1")
	req.URL.RawQuery = q.Encode()

	// do the request
	var item ItemInfo
	resp, err := s.c.DoJson(req, &item)
	if err != nil {
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, &item, nil
}
