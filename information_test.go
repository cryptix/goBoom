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

		Convey("Ls() should send the request", func() {
			mux.HandleFunc("/1.0/ls", func(w http.ResponseWriter, r *http.Request) {
				So(r.Method, ShouldEqual, "GET")
				So(r.URL.Query().Get("token"), ShouldEqual, testSession)
				So(r.URL.Query().Get("item"), ShouldEqual, "pdfs")
				fmt.Fprint(w, `[200,{"name":"public","root":"1","state":"online","user":298814,"type":"folder","id":"1"},[{"name":"pdfs","parent":"1","type":"folder","downloads":0,"state":"online","user":298814,"mtime":"2014-06-21 23:23:46.615535","atime":null,"root":"1","id":"99QJ0C6Y","ctime":"2014-06-21 23:23:46.615535"},{"size":1,"name":"test1.txt","parent":"1","type":"file","downloads":0,"thumb_320":false,"state":"online","mime":"text/plain","user":298814,"mtime":"2014-06-22 00:24:15.259402","owner":true,"atime":null,"root":"1","id":"GE308U2K","ctime":"2014-06-22 00:24:10.954148"},{"size":1,"name":"test2.txt","parent":"1","type":"file","downloads":0,"thumb_320":false,"state":"online","mime":"text/plain","user":298814,"mtime":"2014-06-22 00:24:21.074755","owner":true,"atime":null,"root":"1","id":"I1EYJZTU","ctime":"2014-06-22 00:24:17.562414"}],1]`)
			})

			code, resp, err := info.Ls("pdfs")
			So(err, ShouldBeNil)
			So(code, ShouldEqual, 200)

			So(resp.Pwd, ShouldResemble, ItemInfo{"1", "public", "1", "online", "folder", 298814})
			So(len(resp.Items), ShouldEqual, 3)
		})

		Reset(teardown)
	})

}
