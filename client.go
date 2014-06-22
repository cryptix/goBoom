package goBoom

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/kr/pretty"
)

const (
	libraryVersion = "0.1"
	defaultBaseURL = "https://api.oboom.com/1.0/"
	userAgent      = "goBoom/" + libraryVersion

	defaultAccept    = "application/json"
	defaultMediaType = "application/octet-stream"

	debug = true
)

var (
	ErrUnknwonFourResponseType = errors.New("Can only handle three responses for InformationService.Ls()")
)

type ErrStatusCodeMissmatch struct{ Http, Api int }

func (e ErrStatusCodeMissmatch) Error() string {
	return fmt.Sprintf("ErrStatusCodeMissmatch: http[%d] != api[%d]", e.Http, e.Api)
}

// A Client manages communication with the Pshdl Rest API.
type Client struct {
	// HTTP client used to communicate with the API.
	c *http.Client

	// Base URL for API requests.  baseURL should always be specified with a trailing slash.
	baseURL *url.URL

	// User agent used when communicating with the PSHDL REST API.
	userAgent string

	User *UserService
	Info *InformationService
}

// NewClient returns a new PSHDL REST API client.  If a nil httpClient is
// provided, http.DefaultClient will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		panic(err)
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	httpClient.Jar = jar
	client := &Client{c: httpClient, baseURL: baseURL, userAgent: userAgent}
	client.User = newUserService(client)
	client.Info = newInformationService(client)

	return client
}

// NewJsonRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the baseURL of the Client.
// Relative URLs should always be specified without a preceding slash.
func (c *Client) NewRequest(method, urlStr string, params url.Values) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(rel)
	if params != nil {
		u.RawQuery = params.Encode()
	}

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", defaultAccept)
	req.Header.Add("User-Agent", c.userAgent)
	return req, nil
}

// NewJsonRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the baseURL of the Client.
// Relative URLs should always be specified without a preceding slash.  If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewJsonRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(rel)
	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", defaultAccept)
	req.Header.Add("User-Agent", c.userAgent)
	return req, nil
}

// NewReaderRequest creates an API request. Uses a io.Reader and ctype instead of marshaling json.
func (c *Client) NewReaderRequest(method, urlStr string, body io.Reader, ctype string) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "text/plain")
	req.Header.Add("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", ctype)
	if ctype == "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return req, nil
}

// Do sends an API request and returns the API response.  The API response is
// decoded and stored in the value pointed to by v, or returned as an error if
// an API error has occurred.
func (c *Client) DoJson(req *http.Request, v interface{}) (statusCode int, resp *http.Response, err error) {
	if v == nil {
		err = fmt.Errorf("Dont DoJson() with nothing to unmarshal!")
		return
	}

	resp, err = c.c.Do(req)
	if err != nil {
		return
	}

	statusCode = resp.StatusCode
	err = CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return
	}

	var body io.Reader = resp.Body
	if debug == true {
		bodyBytes, err := ioutil.ReadAll(body)
		if err != nil {
			panic(err)
		}
		fmt.Printf("DEBUG[%s]\n", string(bodyBytes))

		body = bytes.NewReader(bodyBytes)
	}

	var apiResp []interface{}
	err = json.NewDecoder(body).Decode(&apiResp)
	if err != nil {
		err = fmt.Errorf("Json.Decode() failed:%s\n", err)
		return
	}
	resp.Body.Close()

	if len(apiResp) >= 1 {
		code, ok := apiResp[0].(float64)
		if !ok {
			return resp.StatusCode, resp, fmt.Errorf("first result was no float64")
		}

		if resp.StatusCode != int(code) {
			statusCode = 0
			err = ErrStatusCodeMissmatch{resp.StatusCode, int(code)}
			return
		}
		statusCode = int(code)
	}

	if len(apiResp) == 2 {
		err = jsonRemarshal(apiResp[1], &v)

	} else if len(apiResp) == 4 {
		// suspecting Ls() for pwd, []data until further occurance
		data, okTarget := v.(*LsInfo)
		pwd, okPwd := apiResp[1].(map[string]interface{})
		if !okTarget || !okPwd {
			err = ErrUnknwonFourResponseType
			return
		}

		err = jsonRemarshal(pwd, &(data.Pwd))
		if err != nil {
			return
		}
		err = jsonRemarshal(apiResp[2], &(data.Items))

	} else {
		fmt.Printf("DEBUG:%# v\n", len(apiResp), pretty.Formatter(apiResp))
		err = fmt.Errorf("Unknown amount of apiResponses: %d", len(apiResp))
		return

	}

	return
}

// DoPlain sends an API request and returns the API response as a slice of bytes.
func (c *Client) DoPlain(req *http.Request) ([]byte, *http.Response, error) {
	req.Header.Set("Accept", "text/plain")

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return nil, resp, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	return data, resp, err
}

/*
An ErrorResponse reports one or more errors caused by an API request.

PSHDL REST API docs: http://developer.github.com/v3/#client-errors
*/
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Message  interface{}
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message)
}

// CheckResponse checks the API response for errors, and returns them if
// present.  A response is considered an error if it has a status code outside
// the 200 range.  API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.  Any other
// response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}
