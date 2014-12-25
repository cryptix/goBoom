package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/cryptix/go/logging"
	"github.com/cryptix/goBoom"
)

var client *goBoom.Client

func init() {
	client = goBoom.NewClient(nil)

	code, _, err := client.User.Login("el.rey.de.wonns@gmail.com", "70e878c4")
	logging.CheckFatal(err)

	log.Println("Login Response: ", code)
}

func main() {
	var err error

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(http.FileServer(client.FS))

	port := os.Getenv("PORT")
	if port == "" {
		port = "0"
	}

	l, err := net.Listen("tcp", ":"+port)
	logging.CheckFatal(err)
	log.Printf("Serving at http://%s/", l.Addr())

	logging.CheckFatal(http.Serve(l, n))
}
