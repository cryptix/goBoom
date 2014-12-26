package goBoom

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

type BoomFile struct {
	*bytes.Reader
	ID       string
	client   *Client
	readDirN int
	info     os.FileInfo
}

func NewBoomFile(c *Client, name string) (*BoomFile, error) {
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}

	if name == "" {
		return NewBoomFile(c, "1")
	}

	if name[len(name)-1] == '/' {
		name = name[:len(name)-2]
		return NewBoomFile(c, name)
	}

	_, info, err := c.Info.Info(name)
	if err != nil {
		return nil, err
	}

	if len(info) < 1 {
		return nil, errors.New("api: not found")
	}

	if info[0].Type == "folder" {
		return &BoomFile{ID: name, client: c, info: info[0]}, nil
	}

	if info[0].Size() > 1024*1024*5 {
		return nil, errors.New("can't inline file transfer over 5mb")
	}

	_, url, err := c.FS.Download(name)
	if err != nil {
		return nil, os.ErrNotExist
	}

	resp, err := c.c.Get(url.String())
	if err != nil {
		return nil, err
	}

	err = CheckResponse(resp)
	if err != nil {
		return nil, err
	}

	f := &BoomFile{
		ID:     name,
		client: c,
		info:   info[0],
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	f.Reader = bytes.NewReader(b)

	return f, resp.Body.Close()
}

func (b *BoomFile) Close() error {
	return nil
}

func (b *BoomFile) Readdir(n int) ([]os.FileInfo, error) {
	_, ls, err := b.client.Info.Ls(b.ID)
	if err != nil {
		return nil, err
	}

	if b.readDirN == len(ls.Items) {
		return nil, io.EOF
	}

	finfo := make([]os.FileInfo, len(ls.Items))
	for i := range ls.Items {
		finfo[i] = ls.Items[i]
	}
	b.readDirN += len(finfo)
	return finfo, nil
}

func (b *BoomFile) Stat() (os.FileInfo, error) {
	return b.info, nil
}
