package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

const (
	BaseURL    = "https://pracownia.knei.pl"
	LoginURL   = BaseURL + "/" + "?m=wyloguj" // login subpage address
	ActionsURL = BaseURL + "/" + "akcje.php"  // actions subpage address
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("dondu: ")

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("failed to load .env file")
	}

	http.DefaultClient.Jar, err = cookiejar.New(nil)
	if err != nil {
		log.Fatalln("failed to create cookie jar")
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
			Name:    "minutes",
			Aliases: []string{"m"},
			Value:   30,
			Usage:   "time after which to disable the socket",
		},
	},
	Action: func(c *cli.Context) error {
		err := login()
		if err != nil {
			return fmt.Errorf("login: %v", err)
		}

		socket := c.Int("socket")
		minutes := c.Int("minutes")

		err = update(socket, true, minutes)
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
		err := login()
		if err != nil {
			return fmt.Errorf("login: %v", err)
		}

		socket := c.Int("socket")

		err = update(socket, false, 0)
		if err != nil {
			return fmt.Errorf("disable socket %d: %v", socket, err)
		}

		log.Printf("disabled socket %d\n", socket)
		return nil
	},
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

// login performs user authentication (sets the PHPSESSID cookie).
func login() error {
	username := os.Getenv("DONDU_USERNAME")
	password := os.Getenv("DONDU_PASSWORD")

	data := url.Values{}
	data.Set("login", username)
	data.Set("haslo", password)

	res, err := http.Post(LoginURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to make HTTP POST login request: %v", err)
	}

	log.Println("auth status:", res.Status)

	return nil
}

// update modifies state of single socket on the switchboard.
func update(socket int, enabled bool, minutes int) error {
	state := 0
	if enabled {
		state = 1
	}

	data := url.Values{}
	data.Set("gniazdo_nr", strconv.Itoa(socket))
	data.Set("gniazdo_stan", strconv.Itoa(state))
	data.Set("minuty", strconv.Itoa(minutes))
	data.Set("rozdzielnia_zmien_stan", strconv.Itoa(1))

	_, err := http.Post(ActionsURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to make HTTP POST request: %v", err)
	}

	return nil
}
