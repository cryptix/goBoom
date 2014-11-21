package goBoom

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testSession = "testSession"

func TestInformationService_Info(t *testing.T) {
	setup()
	defer teardown()

	info := newInformationService(client)
	client.User.session = testSession

	mux.HandleFunc("/1.0/info", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Query().Get("token"), testSession)
		assert.Equal(t, r.URL.Query().Get("items"), "a,b,c")
		fmt.Fprint(w, `[200,[{"name":"trash","root":"1C","state":"online","user":298814,"type":"folder","id":"1C"},{"name":"public","root":"1","state":"online","user":298814,"type":"folder","id":"1"}]]`)
	})

	code, resp, err := info.Info("a", "b", "c")
	assert.Nil(t, err)
	assert.Equal(t, code, http.StatusOK)

	assert.IsType(t, resp, []ItemInfo{})
	assert.Len(t, resp, 2)
	assert.Equal(t, resp[0], ItemInfo{"1C", "trash", "1C", "online", "folder", 298814})
	assert.Equal(t, resp[1], ItemInfo{"1", "public", "1", "online", "folder", 298814})
}

func TestInformationService_Du(t *testing.T) {
	setup()
	defer teardown()

	info := newInformationService(client)
	client.User.session = testSession

	mux.HandleFunc("/1.0/du", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Query().Get("token"), testSession)
		fmt.Fprint(w, `[200,{"1":{"num":1,"size":2893557},"1C":{"num":0,"size":0},"total":{"num":1,"size":2893557}}]`)
	})

	code, resp, err := info.Du()
	assert.Nil(t, err)
	assert.Equal(t, code, http.StatusOK)

	dummyMap := make(map[string]ItemSize)
	assert.IsType(t, resp, dummyMap)
	assert.Len(t, resp, 3)

	assert.Equal(t, resp["total"], ItemSize{1, 2893557})
	assert.Equal(t, resp["1"], ItemSize{1, 2893557})
	assert.Equal(t, resp["1C"], ItemSize{0, 0})
}

func TestInformationService_Ls(t *testing.T) {
	setup()
	defer teardown()

	info := newInformationService(client)
	client.User.session = testSession

	mux.HandleFunc("/1.0/ls", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Query().Get("token"), testSession)
		assert.Equal(t, r.URL.Query().Get("item"), "pdfs")
		fmt.Fprint(w, `[200,{"name":"public","root":"1","state":"online","user":298814,"type":"folder","id":"1"},[{"name":"pdfs","parent":"1","type":"folder","downloads":0,"state":"online","user":298814,"mtime":"2014-06-21 23:23:46.615535","atime":null,"root":"1","id":"99QJ0C6Y","ctime":"2014-06-21 23:23:46.615535"},{"size":1,"name":"test1.txt","parent":"1","type":"file","downloads":0,"thumb_320":false,"state":"online","mime":"text/plain","user":298814,"mtime":"2014-06-22 00:24:15.259402","owner":true,"atime":null,"root":"1","id":"GE308U2K","ctime":"2014-06-22 00:24:10.954148"},{"size":1,"name":"test2.txt","parent":"1","type":"file","downloads":0,"thumb_320":false,"state":"online","mime":"text/plain","user":298814,"mtime":"2014-06-22 00:24:21.074755","owner":true,"atime":null,"root":"1","id":"I1EYJZTU","ctime":"2014-06-22 00:24:17.562414"}],1]`)
	})

	code, resp, err := info.Ls("pdfs")
	assert.Nil(t, err)
	assert.Equal(t, code, http.StatusOK)

	assert.Equal(t, resp.Pwd, ItemInfo{"1", "public", "1", "online", "folder", 298814})
	assert.Len(t, resp.Items, 3)
}
