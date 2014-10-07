package goBoom

import "strings"

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

type ItemInfo struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Root  string  `json:"root"`
	State string  `json:"state"`
	Type  string  `json:"type"`
	User  float64 `json:"user"`
}

func (i InformationService) Info(ids ...string) (int, []ItemInfo, error) {

	params := map[string]string{
		"token": i.c.User.session,
		"items": strings.Join(ids, ","),
	}

	resp, err := i.c.api.Res("info").Get(params)
	arr, err := ProcessResponse(resp, err)
	if err != nil {
		return resp.Raw.StatusCode, nil, err
	}

	var infoResp []ItemInfo
	if err = DecodeInto(&infoResp, arr[1]); err != nil {
		return resp.Raw.StatusCode, nil, err
	}

	return resp.Raw.StatusCode, infoResp, nil
}

type ItemSize struct {
	Num  float64 `json:"num"`
	Size float64 `json:"size"`
}

func (i InformationService) Du() (int, map[string]ItemSize, error) {

	params := map[string]string{
		"token": i.c.User.session,
	}

	resp, err := i.c.api.Res("du").Get(params)
	arr, err := ProcessResponse(resp, err)
	if err != nil {
		return resp.Raw.StatusCode, nil, err
	}

	duResp := make(map[string]ItemSize)
	if err = DecodeInto(&duResp, arr[1]); err != nil {
		return resp.Raw.StatusCode, nil, err
	}

	return resp.Raw.StatusCode, duResp, nil
}

type LsInfo struct {
	Pwd   ItemInfo
	Items []ItemStat
}

type ItemStat struct {
	Atime     interface{} `json:"atime"`
	Ctime     string      `json:"ctime"`
	Downloads float64     `json:"downloads"`
	ID        string      `json:"id"`
	Mtime     string      `json:"mtime"`
	Name      string      `json:"name"`
	Parent    string      `json:"parent"`
	Root      string      `json:"root"`
	State     string      `json:"state"`
	Type      string      `json:"type"`
	User      float64     `json:"user"`
}

func (i InformationService) Ls(item string) (int, *LsInfo, error) {

	params := map[string]string{
		"token": i.c.User.session,
		"item":  item,
	}

	resp, err := i.c.api.Res("ls").Get(params)
	arr, err := ProcessResponse(resp, err)
	if err != nil {
		return resp.Raw.StatusCode, nil, err
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
