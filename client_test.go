package goBoom

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/cryptix/gocrayons"
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
	client.api = gocrayons.Api(server.URL + "/1.0/")
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func TestNewClient(t *testing.T) {
	var (
		c *Client
	)

	type createPut struct {
		Name, Email string
	}

	Convey("Given a new Client", t, func() {
		c = NewClient(nil)

		Convey("It should have the correct BaseURL", func() {
			So(c.baseURL.String(), ShouldEqual, defaultBaseURL)
		})

		Convey("It should have the correct userAgent", func() {
			So(c.userAgent, ShouldEqual, userAgent)
		})

		Convey("It should have a UserService", func() {
			So(c.User, ShouldHaveSameTypeAs, &UserService{})
			So(c.User, ShouldNotBeNil)
		})

		Convey("It should have a InformationService", func() {
			So(c.Info, ShouldHaveSameTypeAs, &InformationService{})
			So(c.Info, ShouldNotBeNil)
		})

		Convey("It should have a FilesystemService", func() {
			So(c.FS, ShouldHaveSameTypeAs, &FilesystemService{})
			So(c.FS, ShouldNotBeNil)
		})
	})

}
