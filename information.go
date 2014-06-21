package goBoom

import (
	"fmt"
	"net/url"
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

func (i InformationService) Info() (int, []string, error) {

	reqParams := url.Values{
		"token": []string{i.c.User.session},
	}

	req, err := i.c.NewRequest("GET", "info", reqParams)
	if err != nil {
		return 0, nil, err
	}

	var infoResp []string
	apiResponseCode, resp, err := i.c.DoJson(req, &infoResp)
	if err != nil {
		return 0, nil, err
	}

	if resp.StatusCode != apiResponseCode {
		err = fmt.Errorf("resp.StatusCode[%d] != apiResponseCode[%d]", resp.StatusCode, apiResponseCode)
	}

	return resp.StatusCode, infoResp, nil
}
