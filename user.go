package goBoom

import (
	"code.google.com/p/go.crypto/pbkdf2"
	"crypto/sha1"
	"fmt"
	"net/url"
	"strings"
)

type UserService struct {
	c *Client

	session string
}

func newUserService(c *Client) *UserService {
	u := &UserService{}
	if c == nil {
		u.c = NewClient(nil)
	} else {
		u.c = c
	}

	return u
}

type loginResponse struct {
	Cookie  string `json:"cookie"`
	Session string `json:"session"`
	User    struct {
		ApiKey      string      `json:"api_key"`
		Balance     interface{} `json:"balance"`
		Email       string      `json:"email"`
		ExternalID  string      `json:"external_id"`
		FtpUsername string      `json:"ftp_username"`
		ID          string      `json:"id"`
		Name        string      `json:"name"`
		Partner     string      `json:"partner"`
		PartnerLast interface{} `json:"partner_last"`
		Pro         string      `json:"pro"`
		Settings    struct {
			Ddl              float64 `json:"ddl"`
			RewriteBehaviour float64 `json:"rewrite_behaviour"`
		} `json:"settings"`
		Traffic struct {
			Current  float64 `json:"current"`
			Increase float64 `json:"increase"`
			Last     float64 `json:"last"`
			Max      float64 `json:"max"`
		} `json:"traffic"`
		Webspace float64 `json:"webspace"`
	} `json:"user"`
}

func (u *UserService) Login(name, passw string) (int, *loginResponse, error) {

	derived := pbkdf2.Key([]byte(passw), []byte(reverse(passw)), 1000, 16, sha1.New)

	reqParams := url.Values{
		"auth": []string{name},
		"pass": []string{fmt.Sprintf("%x", derived)},
	}

	oldHost := u.c.baseURL.Host
	u.c.baseURL.Host = strings.Replace(u.c.baseURL.Host, "api.oboom.com", "www.oboom.com", 1)

	req, err := u.c.NewReaderRequest("POST", "login", strings.NewReader(reqParams.Encode()), "")
	if err != nil {
		return 0, nil, err
	}

	u.c.baseURL.Host = oldHost

	var liResp loginResponse
	liStatus, resp, err := u.c.DoJson(req, &liResp)
	if err != nil {
		return 0, nil, err
	}

	if resp.StatusCode != liStatus {
		return liStatus, nil, fmt.Errorf("resp.StatusCode[%d] != liStatus[%d]", resp.StatusCode, liStatus)
	}
	fmt.Printf("Status Code[%d]\n", liStatus)

	u.session = liResp.Session

	return resp.StatusCode, &liResp, nil
}
