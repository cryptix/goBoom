package goBoom

import (
	"fmt"
	"net/url"
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

type ItemInfo struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Root  string  `json:"root"`
	State string  `json:"state"`
	Type  string  `json:"type"`
	User  float64 `json:"user"`
}

func (i InformationService) Info(ids ...string) (int, []ItemInfo, error) {

	reqParams := make(url.Values, 2)
	reqParams.Set("token", i.c.User.session)
	reqParams.Set("items", strings.Join(ids, ","))

	req, err := i.c.NewRequest("GET", "info", reqParams)
	if err != nil {
		return 0, nil, err
	}

	var infoResp []ItemInfo
	apiResponseCode, resp, err := i.c.DoJson(req, &infoResp)
	if err != nil {
		return 0, nil, err
	}

	if resp.StatusCode != apiResponseCode {
		err = fmt.Errorf("resp.StatusCode[%d] != apiResponseCode[%d]", resp.StatusCode, apiResponseCode)
	}

	return resp.StatusCode, infoResp, nil
}

type ItemSize struct {
	Num  float64 `json:"num"`
	Size float64 `json:"size"`
}

func (i InformationService) Du() (int, map[string]ItemSize, error) {

	reqParams := make(url.Values, 2)
	reqParams.Set("token", i.c.User.session)

	req, err := i.c.NewRequest("GET", "du", reqParams)
	if err != nil {
		return 0, nil, err
	}

	duResp := make(map[string]ItemSize)
	apiResponseCode, resp, err := i.c.DoJson(req, &duResp)
	if err != nil {
		return 0, nil, err
	}

	if resp.StatusCode != apiResponseCode {
		err = fmt.Errorf("resp.StatusCode[%d] != apiResponseCode[%d]", resp.StatusCode, apiResponseCode)
	}

	return resp.StatusCode, duResp, nil
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
	apiResponseCode, resp, err := i.c.DoJson(req, &lsResp)
	if err != nil {
		return 0, nil, err
	}

	if resp.StatusCode != apiResponseCode {
		err = fmt.Errorf("resp.StatusCode[%d] != apiResponseCode[%d]", resp.StatusCode, apiResponseCode)
	}

	return resp.StatusCode, &lsResp, nil
}
