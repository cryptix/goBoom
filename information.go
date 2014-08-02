package goBoom

import (
	"errors"
	"net/url"
	"strings"

	"github.com/mitchellh/mapstructure"
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

	resp, err := i.c.api.Res("/1.0/items").Get(params)
	if err != nil {
		return 0, nil, err
	}

	arr, err := ProcessResponse(resp, err)
	if err != nil {
		return resp.Raw.StatusCode, nil, err
	}

	var infoResp []ItemInfo
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &infoResp}

	dec, err := mapstructure.NewDecoder(config)
	if err != nil {
		return 0, nil, errors.New("NewDecoder Error:" + err.Error())
	}

	err = dec.Decode(arr[1])
	if err != nil {
		return 0, nil, errors.New("Decode Error:" + err.Error())
	}

	// pretty.Println(arr[1])

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

	resp, err := i.c.api.Res("/1.0/du").Get(params)
	if err != nil {
		return 0, nil, err
	}

	arr, err := ProcessResponse(resp, err)
	if err != nil {
		return resp.Raw.StatusCode, nil, err
	}

	duResp := make(map[string]ItemSize)
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &duResp,
	}

	dec, err := mapstructure.NewDecoder(config)
	if err != nil {
		return 0, nil, errors.New("NewDecoder Error:" + err.Error())
	}

	err = dec.Decode(arr[1])
	if err != nil {
		return 0, nil, errors.New("Decode Error:" + err.Error())
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

	reqParams := make(url.Values, 2)
	reqParams.Set("token", i.c.User.session)
	reqParams.Set("item", item)

	req, err := i.c.NewRequest("GET", "ls", reqParams)
	if err != nil {
		return 0, nil, err
	}

	var lsResp LsInfo
	resp, err := i.c.DoJson(req, &lsResp)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, &lsResp, nil
}
