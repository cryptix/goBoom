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
				So(r.URL.Query().Get("items"), ShouldEqual, "a,b,c")
				fmt.Fprint(w, `[200,[{"name":"trash","root":"1C","state":"online","user":298814,"type":"folder","id":"1C"},{"name":"public","root":"1","state":"online","user":298814,"type":"folder","id":"1"}]]`)
			})

			code, resp, err := info.Info("a", "b", "c")
			So(err, ShouldBeNil)
			So(code, ShouldEqual, 200)

			So(resp, ShouldHaveSameTypeAs, []ItemInfo{})
			So(len(resp), ShouldEqual, 2)
			So(resp[0], ShouldResemble, ItemInfo{"1C", "trash", "1C", "online", "folder", 298814})
			So(resp[1], ShouldResemble, ItemInfo{"1", "public", "1", "online", "folder", 298814})
		})

		Reset(teardown)
	})

}
