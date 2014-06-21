package goBoom

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const testSession = "testSession"

func TestInformationService(t *testing.T) {
	var info *InformationService

	Convey("Given a newInformationService()", t, func() {
		setup()

		info = newInformationService(client)

		// fake login
		client.User.session = testSession

		Convey("Info(...) should send the request", func() {
			mux.HandleFunc("/1.0/info", func(w http.ResponseWriter, r *http.Request) {
				So(r.Method, ShouldEqual, "GET")
				So(r.URL.Query().Get("token"), ShouldEqual, testSession)
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

		Convey("Du() should send the request", func() {
			mux.HandleFunc("/1.0/du", func(w http.ResponseWriter, r *http.Request) {
				So(r.Method, ShouldEqual, "GET")
				So(r.URL.Query().Get("token"), ShouldEqual, testSession)
				fmt.Fprint(w, `[200,{"1":{"num":1,"size":2893557},"1C":{"num":0,"size":0},"total":{"num":1,"size":2893557}}]`)
			})

			code, resp, err := info.Du()
			So(err, ShouldBeNil)
			So(code, ShouldEqual, 200)

			dummyMap := make(map[string]ItemSize)
			So(resp, ShouldHaveSameTypeAs, dummyMap)
			So(len(resp), ShouldEqual, 3)
			So(resp["total"], ShouldResemble, ItemSize{1, 2893557})
			So(resp["1"], ShouldResemble, ItemSize{1, 2893557})
			So(resp["1C"], ShouldResemble, ItemSize{0, 0})
		})

		Reset(teardown)
	})

}
