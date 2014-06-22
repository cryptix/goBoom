package goBoom

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/cheggaaa/pb"
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
func (s *FilesystemService) Upload(fname string, input io.Reader, size int64) (resp *http.Response, newInfo *ItemInfo, err error) {

	pipeOut, pipeIn := io.Pipe()
	bar := pb.New(int(size)).SetUnits(pb.U_BYTES)
	bar.ShowSpeed = true

	writer := multipart.NewWriter(pipeIn)

	// do the request concurrently
	// var item ItemInfo
	done := make(chan error)
	go func() {

		// prepare request
		req, err := s.c.NewReaderRequest("POST", "ul", pipeOut, writer.FormDataContentType())
		if err != nil {
			done <- err
			return
		}
		// TODO calculate right amount of overhead
		req.ContentLength = size // filesize
		req.ContentLength += 227 // multipart header exclouding filename
		req.ContentLength += int64(len(fname))
		req.ContentLength -= 19

		// set  token
		q := req.URL.Query()
		q.Set("token", s.c.User.session)
		q.Set("parent", "1")
		req.URL.RawQuery = q.Encode()

		fmt.Println("Created Request")
		bar.Start()

		var fInfo ItemInfo
		// resp, err = s.c.c.Do(req)
		resp, err = s.c.DoJson(req, &fInfo)
		if err != nil {
			done <- err
			return
		}
		newInfo = &fInfo

		done <- nil
	}()

	part, err := writer.CreateFormFile("file", filepath.Base(fname))
	if err != nil {
		return
	}

	out := io.MultiWriter(part, bar)
	_, err = io.Copy(out, input)
	if err != nil {
		fmt.Println("pipe copy error")
		err = <-done
		return
	}

	err = writer.Close()
	if err != nil {
		return
	}

	err = pipeIn.Close()
	if err != nil {
		return
	}

	err = <-done
	bar.Finish()
	if err != nil {
		return
	}

	return
}
