package goBoom

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUserService(t *testing.T) {

	Convey("Given a NewUser()", t, func() {
		setup()

		user := newUserService(client)

		Convey("Login() should send the request", func() {

			mux.HandleFunc("/1.0/login", func(w http.ResponseWriter, r *http.Request) {
				So(r.Method, ShouldEqual, "POST")
				So(r.ParseForm(), ShouldBeNil)
				So(r.PostForm.Get("auth"), ShouldEqual, "test@mail.com")
				So(r.PostForm.Get("pass"), ShouldEqual, "94406d8b3a3876308552d168e56a42f9")
				fmt.Fprint(w, `[200,{"cookie":"1000000000:efcb5ef3efec97aa50c33e1efb183e223633a3bf","user":{"id":"1000000000","name":"johndoe","email":"john@example.com","api_key":"d09272c7412aba77d1b06795bf9d8f701ee0171e","pro":"0000-00-00T00:00:00.000Z","webspace": 0.523,"traffic":{"current":0.532,"increase":0.532,"last":0.532,"max":0.532 },"balance":0,"settings":{"rewrite_behaviour":1,"ddl":0},"external_id":"EXTERNAL_OAUTH_PROVIDER_ID","ftp_username":"johndoe","partner":"3","partner_last":1392728258},"session":"cb597b3e-cfc4-4329-abe0-5dc2b64a8e9a"}]`)
			})

			code, resp, err := user.Login("test@mail.com", "1234")
			So(err, ShouldBeNil)

			So(code, ShouldEqual, 200)
			So(resp.Cookie, ShouldEqual, "1000000000:efcb5ef3efec97aa50c33e1efb183e223633a3bf")
			So(resp.User.Name, ShouldEqual, "johndoe")
			So(resp.Session, ShouldEqual, "cb597b3e-cfc4-4329-abe0-5dc2b64a8e9a")

			Convey("Should set the session to onto the user", func() {
				So(user.session, ShouldEqual, "cb597b3e-cfc4-4329-abe0-5dc2b64a8e9a")
			})
		})

		Reset(teardown)
	})

}
