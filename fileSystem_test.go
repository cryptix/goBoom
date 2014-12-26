package goBoom

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilesystemService_DL(t *testing.T) {
	setup()
	defer teardown()

	fs := newFilesystemService(client)
	fs.c.User.session = "testSession"

	mux.HandleFunc("/1.0/dl", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Query().Get("token"), "testSession")
		assert.Equal(t, r.URL.Query().Get("item"), "1234")
		fmt.Fprint(w, `[200, "testdl.host", "192388123-123-123123"]`)
	})

	code, resp, err := fs.Download("1234")
	assert.Nil(t, err)
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, resp.String(), "https://testdl.host/1.0/dlh?ticket=192388123-123-123123")

}

func TestFilesystemService_UL_Server(t *testing.T) {
	setup()
	defer teardown()

	fs := newFilesystemService(client)

	mux.HandleFunc("/1.0/ul/server", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		fmt.Fprint(w, `[200, ["s7.oboom.com"]]`)
	})

	servers, err := fs.GetULServer()
	assert.Nil(t, err)
	assert.Len(t, servers, 1)
	assert.Equal(t, "s7.oboom.com", servers[0])

}

func TestFilesystemService_InterfaceFileSystem(t *testing.T) {
	setup()
	defer teardown()

	fs := newFilesystemService(client)
	fs.c.User.session = "testSession"

	mux.HandleFunc("/1.0/info", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "GET")
		assert.Equal(t, r.URL.Query().Get("token"), testSession)
		assert.Equal(t, r.URL.Query().Get("items"), "1")
		cpJson(t, w, "_tests/info.json")
	})

	isFS := func(http.FileSystem) {}
	isFS(fs)

	_, err := fs.Open("/")
	assert.Nil(t, err)
}
