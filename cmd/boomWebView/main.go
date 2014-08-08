package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/cryptix/goBoom"
)

// tmplFiles template with list of files found in the workspace
var tmplFiles = template.Must(template.New("tmplFiles").Parse(`
<!doctype html>
<html>
<head>
	<title>List of Files for {{.Pwd.Name}}</title>
</head>
<body>
<h1>List:  {{.Pwd}} </h1>

<p> <a href="/">Root Folder</a></p>
<ul>
{{range .Items}}
	{{if eq .Type "folder"}}
		<li>Folder: <a href="/?item={{.ID}}">{{.Name}}</a></li>
	{{else}}
		<li>File: <a href="/get?item={{.ID}}">{{.Name}}</a></li>
	{{end}}

{{end}}
</ul>
</body>
</html>
`))

func listHandler(rw http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if item == "" {
		item = "1"
	}

	_, ls, err := client.Info.Ls(item)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		log.Println("listHandler Err:", err)
		return
	}

	check(tmplFiles.Execute(rw, ls))
}

func getHandler(rw http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if item == "" {
		http.Error(rw, "No item ID", http.StatusBadRequest)
		return
	}

	_, dlurl, err := client.FS.Download(item)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		log.Println("getHandelr Err:", err)
		return
	}

	http.Redirect(rw, req, dlurl.String(), http.StatusFound)
}

var client *goBoom.Client

func init() {
	client = goBoom.NewClient(nil)

	code, _, err := client.User.Login("email", "clearPassword")
	check(err)

	log.Println("Login Response: ", code)
}

func main() {
	var err error

	http.HandleFunc("/", listHandler)
	http.HandleFunc("/get", getHandler)

	err = http.ListenAndServe(":3002", nil)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
