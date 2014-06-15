package goBoom

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the GitHub client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

// setup sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server.  Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// github client configured to use test server
	client = NewClient(nil)
	url, _ := url.Parse(server.URL + "/1.0/")
	client.baseURL = url
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func TestNewClient(t *testing.T) {
	var c *Client
	Convey("Given a new Client", t, func() {
		c = NewClient(nil)

		Convey("It should have the correct BaseURL", func() {
			So(c.baseURL.String(), ShouldEqual, defaultBaseURL)
		})

		Convey("It should have the correct UserAgent", func() {
			So(c.UserAgent, ShouldEqual, userAgent)
		})
	})
}

func TestNewRequest(t *testing.T) {
	var (
		c *Client
	)

	Convey("Given a new Client", t, func() {
		c = NewClient(nil)

		Convey("and a GET Request", func() {
			inURL, outURL := "foo", defaultBaseURL+"foo?param1=val1"

			params := url.Values{"param1": []string{"val1"}}
			req, err := c.NewRequest("GET", inURL, params)
			So(err, ShouldBeNil)

			Convey("It should be a GET request", func() {
				So(req.Method, ShouldEqual, "GET")
			})

			Convey("It should have its URL expanded with the parameters", func() {
				So(req.URL.String(), ShouldEqual, outURL)
			})

			Convey("It should have the default user-agent is attached to the request", func() {
				userAgent := req.Header.Get("User-Agent")
				So(c.UserAgent, ShouldEqual, userAgent)
			})
		})

		Convey("and a POST Request", func() {
			inURL, outURL := "foo", defaultBaseURL+"foo"

			params := url.Values{"param1": []string{"val1"}}
			req, err := c.NewReaderRequest("POST", inURL, strings.NewReader(params.Encode()), "")
			So(err, ShouldBeNil)

			Convey("It should be a POST request", func() {
				So(req.Method, ShouldEqual, "POST")
			})

			Convey("It should have its URL expanded", func() {
				So(req.URL.String(), ShouldEqual, outURL)
			})

			Convey("It should have the default user-agent is attached to the request", func() {
				userAgent := req.Header.Get("User-Agent")
				So(c.UserAgent, ShouldEqual, userAgent)
			})
		})
	})
}

func TestNewJsonRequest(t *testing.T) {
	var (
		c   *Client
		req *http.Request
	)

	type createPut struct {
		Name, Email string
	}

	Convey("Given a new Client", t, func() {
		c = NewClient(nil)

		Convey("and a valid Request", func() {
			inURL, outURL := "foo", defaultBaseURL+"foo"
			inBody, outBody := &createPut{Name: "l", Email: "hi@me.com"}, `{"Name":"l","Email":"hi@me.com"}`+"\n"
			req, _ = c.NewJsonRequest("PUT", inURL, inBody)

			Convey("It should have its URL expanded", func() {
				So(req.URL.String(), ShouldEqual, outURL)
			})

			Convey("It should encode the body in JSON", func() {
				body, _ := ioutil.ReadAll(req.Body)
				So(string(body), ShouldEqual, outBody)
			})

			Convey("It should have the default user-agent is attached to the request", func() {
				userAgent := req.Header.Get("User-Agent")
				So(c.UserAgent, ShouldEqual, userAgent)
			})

		})

		Convey("and an invalid Request", func() {
			type T struct {
				A map[int]interface{}
			}
			_, err := c.NewJsonRequest("GET", "/", &T{})

			Convey("It should return an error (beeing *json.UnsupportedTypeError)", func() {
				So(err, ShouldNotBeNil)
				So(err, ShouldHaveSameTypeAs, &json.UnsupportedTypeError{})
			})

		})

		Convey("and a bad Request URL", func() {
			_, err := c.NewJsonRequest("GET", ":", nil)
			Convey("It should return an error (beeing *url.Error{})", func() {
				So(err, ShouldNotBeNil)
				So(err, ShouldHaveSameTypeAs, &url.Error{})
			})
		})
	})
}

func TestDo(t *testing.T) {

	Convey("Given a clean test server", t, func() {
		setup()

		Convey("Do() should send the request", func() {

			type foo struct {
				A string
			}

			mux.HandleFunc("/1.0/", func(w http.ResponseWriter, r *http.Request) {
				So(r.Method, ShouldEqual, "GET")

				fmt.Fprint(w, `[200, {"A":"n"}]`)
			})

			req, err := client.NewJsonRequest("GET", "", nil)
			So(err, ShouldBeNil)

			var f foo
			_, err = client.DoJson(req, &f)
			So(err, ShouldBeNil)
			So(f, ShouldResemble, foo{"n"})
		})

		Convey("A plain request should get response", func() {

			want := `/1.0/servertime`

			mux.HandleFunc("/1.0/", func(w http.ResponseWriter, r *http.Request) {
				So(r.Method, ShouldEqual, "GET")
				fmt.Fprint(w, want)
			})

			req, err := client.NewRequest("GET", "", nil)
			So(err, ShouldBeNil)

			body, _, err := client.DoPlain(req)
			So(err, ShouldBeNil)
			So(string(body), ShouldEqual, want)
		})

		Convey("A bad plain request should return a http error", func() {

			mux.HandleFunc("/1.0/", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Bad Request", 400)
			})

			req, err := client.NewJsonRequest("GET", "", nil)
			So(err, ShouldBeNil)

			_, _, err = client.DoPlain(req)
			So(err, ShouldNotBeNil)
		})

		Reset(teardown)
	})

}
