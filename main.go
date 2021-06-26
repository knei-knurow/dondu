package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"

	"github.com/knei-knurow/dondu/api"
	"github.com/urfave/cli/v2"
)

var (
	username string
	password string
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("dondu: ")

	username = os.Getenv("DONDU_USERNAME")
	password = os.Getenv("DONDU_PASSWORD")

	var err error
	http.DefaultClient.Jar, err = cookiejar.New(nil)
	if err != nil {
		log.Fatalln("failed to create cookie jar")
	}
}

func main() {
	var err error
	app := &cli.App{
		Name:  "dondu",
		Usage: "Easily control our switchboard from the command line.",
		Action: func(c *cli.Context) error {
			log.Println("no commands passed")
			return nil
		},
		Commands: []*cli.Command{
			&enableCommand,
			&disableCommand,
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

var enableCommand = cli.Command{
	Name:  "enable",
	Usage: "enable socket",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "socket",
			Aliases: []string{"s"},
			Value:   -1,
			Usage:   "number of the socket to enable",
		},
		&cli.IntFlag{
			Name:    "time",
			Aliases: []string{"t"},
			Value:   30,
			Usage:   "time after which to disable the socket (minutes)",
		},
	},
	Action: func(c *cli.Context) error {
		err := api.Login(username, password)
		if err != nil {
			return fmt.Errorf("login: %v", err)
		}

		socket := c.Int("socket")
		minutes := c.Int("time")

		err = api.Update(socket, true, minutes)
		if err != nil {
			return fmt.Errorf("enable socket %d for %d minutes: %v", socket, minutes, err)
		}

		log.Printf("enabled socket %d for %d minutes\n", socket, minutes)
		return nil
	},
}

var disableCommand = cli.Command{
	Name:  "disable",
	Usage: "disable socket",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "socket",
			Aliases: []string{"s"},
			Value:   -1,
			Usage:   "number of the socket to disable",
		},
	},
	Action: func(c *cli.Context) error {
		err := api.Login(username, password)
		if err != nil {
			return fmt.Errorf("login: %v", err)
		}

		socket := c.Int("socket")

		err = api.Update(socket, false, 0)
		if err != nil {
			return fmt.Errorf("disable socket %d: %v", socket, err)
		}

		log.Printf("disabled socket %d\n", socket)
		return nil
	},
}
