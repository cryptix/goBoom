package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/cryptix/goBoom"
)

var client *goBoom.Client

func init() {
	client = goBoom.NewClient(nil)

	code, _, err := client.User.Login("email", "clearPassword")
	check(err)

	log.Println("Login Response: ", code)
}

func main() {
	app := cli.NewApp()
	app.Name = "boomTool"
	app.Commands = []cli.Command{
		{
			Name:  "ls",
			Usage: "list...",
			Action: func(c *cli.Context) {
				_, ls, err := client.Info.Ls(c.Args().First())
				check(err)
				for _, item := range ls.Items {
					log.Printf("%8s - %s\n", item.ID, item.Name)

				}
			},
		},
		{
			Name:      "put",
			ShortName: "p",
			Usage:     "put a file",
			Action: func(c *cli.Context) {
				println("putting:", c.Args().First())
			},
		},
		{
			Name:      "get",
			ShortName: "g",
			Usage:     "get a file",
			Action: func(c *cli.Context) {

				item := c.Args().First()
				if item == "" {
					println("no item id")
					os.Exit(1)
				}
				_, url, err := client.FS.Download(item)
				check(err)
				println(url.String())
			},
		},
	}

	app.Run(os.Args)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
