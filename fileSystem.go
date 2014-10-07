package goBoom

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/url"
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

// Upload pushes the input io.Reader to the service
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
	res := s.c.api.Res("ul")
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

// Download requests a download url for item
func (s *FilesystemService) Download(item string) (int, *url.URL, error) {
	if s.c.User == nil {
		return -1, nil, errors.New("non pro download not supported")
	}

	params := map[string]string{
		"token": s.c.User.session,
		"item":  item,
	}

	resp, err := s.c.api.Res("dl").Get(params)
	arr, err := ProcessResponse(resp, err)
	if err != nil {
		return resp.Raw.StatusCode, nil, err
	}

	var (
		u  url.URL
		ok bool
	)

	u.Host, ok = arr[1].(string)
	if !ok {
		return -1, nil, errors.New("arr[1] is not a string")
	}

	ticket, ok := arr[2].(string)
	if !ok {
		return -1, nil, errors.New("arr[2] is not a string")
	}

	u.Scheme = "https"
	u.Path = libraryVersion + "/dlh"

	qry := u.Query()
	qry.Set("ticket", ticket)
	u.RawQuery = qry.Encode()

	return resp.Raw.StatusCode, &u, nil
}
