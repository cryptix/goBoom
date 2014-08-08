package goBoom

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFilesystemService(t *testing.T) {

	Convey("Given a newFilesystemService()", t, func() {
		setup()

		fs := newFilesystemService(client)
		fs.c.User.session = "testSession"

		Convey("Download() should send the request", func() {

			mux.HandleFunc("/1.0/dl", func(w http.ResponseWriter, r *http.Request) {
				So(r.Method, ShouldEqual, "GET")
				So(r.URL.Query().Get("token"), ShouldEqual, "testSession")
				So(r.URL.Query().Get("item"), ShouldEqual, "1234")
				fmt.Fprint(w, `[200, "testdl.host", "192388123-123-123123"]`)
			})

			code, resp, err := fs.Download("1234")
			So(err, ShouldBeNil)
			So(code, ShouldEqual, 200)
			So(resp.String(), ShouldEqual, "https://testdl.host/1.0/dlh?ticket=192388123-123-123123")

		})

		Reset(teardown)
	})

}
