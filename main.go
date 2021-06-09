package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	BaseURL    = "https://pracownia.knei.pl"
	LoginURL   = BaseURL + "/" + "?m=wyloguj" // login subpage address
	ActionsURL = BaseURL + "/" + "akcje.php"  // actions subpage address
)

var (
	username string
	password string
)

var (
	socketNumber int
	socketState  int
	minutes      int
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("dondu: ")

	flag.IntVar(&socketNumber, "socket", 0, "number of the socket to manipulate")
	flag.IntVar(&socketState, "state", 0, "state to set on socket ")
	flag.IntVar(&minutes, "minutes", 30, "time")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file")
	}

	username = os.Getenv("DONDU_USERNAME")
	password = os.Getenv("DONDU_PASSWORD")
}

func main() {
	var err error
	http.DefaultClient.Jar, err = cookiejar.New(nil)
	if err != nil {
		log.Fatalln("failed to create cookie jar")
	}

	err = auth()
	if err != nil {
		log.Fatalln("failed to authencticate:", err)
	}

	err = update()
	if err != nil {
		log.Fatalln("failed to update the switchboard:", err)
	}
}

// auth performs user authentication (sets the PHPSESSID cookie).
func auth() error {
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
func update() error {
	data := url.Values{}
	data.Set("gniazdo_nr", strconv.Itoa(socketNumber))
	data.Set("gniazdo_stan", strconv.Itoa(socketState))
	data.Set("minuty", strconv.Itoa(socketState))
	data.Set("rozdzielnia_zmien_stan", strconv.Itoa(1))

	_, err := http.Post(ActionsURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to make HTTP POST request: %v", err)
	}

	return nil
}
