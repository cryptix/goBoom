package goBoom

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"strings"
)

type InformationService struct {
	c *Client
}

func newInformationService(c *Client) *InformationService {
	i := &InformationService{}
	if c == nil {
		i.c = NewClient(nil)
	} else {
		i.c = c
	}

	return i
}

func (i InformationService) Info(ids ...string) (int, []ItemStat, error) {

	params := map[string]string{
		"token": i.c.User.session,
		"items": strings.Join(ids, ","),
	}

	resp, err := i.c.api.Res("info").Get(params)
	arr, err := ProcessResponse(resp, err)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	var infoResp []ItemStat
	if err = DecodeInto(&infoResp, arr[1]); err != nil {
		return resp.Raw.StatusCode, nil, err
	}

	return resp.Raw.StatusCode, infoResp, nil
}

type ItemSize struct {
	Num  int64 `json:"num"`
	Size int64 `json:"size"`
}

func (i InformationService) Du() (int, map[string]ItemSize, error) {

	params := map[string]string{
		"token": i.c.User.session,
	}

	resp, err := i.c.api.Res("du").Get(params)
	arr, err := ProcessResponse(resp, err)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	duResp := make(map[string]ItemSize)
	if err = DecodeInto(&duResp, arr[1]); err != nil {
		return resp.Raw.StatusCode, nil, err
	}

	return resp.Raw.StatusCode, duResp, nil
}

type LsInfo struct {
	Pwd   ItemStat
	Items []ItemStat
}

type ItemStat struct {
	Atime     string `mapstructure:"atime"`
	Ctime     string `mapstructure:"ctime"`
	Downloads int64  `mapstructure:"downloads"`
	ID        string `mapstructure:"id"`
	Mtime     string `mapstructure:"mtime"`
	Iname     string `mapstructure:"name"`
	Parent    string `mapstructure:"parent"`
	Root      string `mapstructure:"root"`
	State     string `mapstructure:"state"`
	Type      string `mapstructure:"type"`
	User      int64  `mapstructure:"user"`
	Isize     int64  `mapstructure:"size"`
	DDL       bool   `mapstructure:"ddl"`
	Mime      string `mapstructure:"mime"`
	Owner     bool   `mapstructure:"owner"`
}

func (i InformationService) Ls(item string) (int, *LsInfo, error) {

	params := map[string]string{
		"token": i.c.User.session,
		"item":  item,
	}

	resp, err := i.c.api.Res("ls").Get(params)
	arr, err := ProcessResponse(resp, err)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	var lsResp LsInfo
	if err = DecodeInto(&lsResp.Pwd, arr[1]); err != nil {
		return resp.Raw.StatusCode, nil, err
	}

	if err = DecodeInto(&lsResp.Items, arr[2]); err != nil {
		return resp.Raw.StatusCode, nil, err
	}

	return resp.Raw.StatusCode, &lsResp, nil
}

func (i ItemStat) IsDir() bool {
	return i.Type == "folder"
}

func (i ItemStat) ModTime() time.Time {
	const format = `2006-01-02 15:04:05.000000`
	var (
		t   = time.Now()
		err error
	)
	switch {
	case i.Mtime != "":
		t, err = time.Parse(format, i.Mtime)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ItemStat(%s).ModTime Parse(Mtime): %s\n", i.ID, err)
		}
	case i.Ctime != "":
		t, err = time.Parse(format, i.Mtime)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ItemStat(%s).ModTime Parse(Mtime): %s\n", i.ID, err)
		}
	case i.Atime != "":
		t, err = time.Parse(format, i.Atime)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ItemStat(%s).ModTime Parse(Mtime): %s\n", i.ID, err)
		}
	default:
		return t
	}

	return t
}

func (i ItemStat) Mode() os.FileMode {
	return os.ModePerm
}

func (i ItemStat) Name() string {
	return i.ID
}
func (i ItemStat) Size() int64 {
	return i.Isize
}

func (i ItemStat) Sys() interface{} {
	return nil
}
