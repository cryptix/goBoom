package goBoom

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInformationService(t *testing.T) {

	Convey("Given a newInformationService()", t, func() {
		setup()

		info := newInformationService(client)

		// fake login
		client.User.session = "testSession"

		Convey("Info(...) should send the request", func() {

			mux.HandleFunc("/1.0/info", func(w http.ResponseWriter, r *http.Request) {
				So(r.Method, ShouldEqual, "GET")
				So(r.URL.Query().Get("token"), ShouldEqual, "testSession")
				// So(r.ParseForm(), ShouldBeNil)
				// So(r.PostForm.Get("auth"), ShouldEqual, "test@mail.com")
				// So(r.PostForm.Get("pass"), ShouldEqual, "94406d8b3a3876308552d168e56a42f9")
				fmt.Fprint(w, `[200, ["1", "1C"]]`)
			})

			code, resp, err := info.Info()
			So(err, ShouldBeNil)
			So(code, ShouldEqual, 200)

			So(resp, ShouldHaveSameTypeAs, []string{})
			So(len(resp), ShouldEqual, 2)
			So(resp[0], ShouldEqual, "1")
			So(resp[1], ShouldEqual, "1C")
		})

		Reset(teardown)
	})

}
