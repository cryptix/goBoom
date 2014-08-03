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
func (s *FilesystemService) Upload(fname string, input io.Reader) (int, []ItemStat, error) {

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
	res := s.c.api.Res("/1.0/ul")
	res.Payload = &bodyBuf
	res.Headers.Set("Content-Type", writer.FormDataContentType())

	// set  token
	params := map[string]string{
		"token":  s.c.User.session,
		"parent": "1",
	}

	// do the request
	resp, err := res.FormPost(params)
	arr, err := ProcessResponse(resp, err)
	if err != nil {
		return resp.Raw.StatusCode, nil, err
	}

	var items []ItemStat
	if err = DecodeInto(&items, arr[1]); err != nil {
		return resp.Raw.StatusCode, nil, err
	}

	return resp.Raw.StatusCode, items, nil
}
