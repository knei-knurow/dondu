package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	ApiUrl = "https://pracownia.knei.pl/akcje.php"
)

var (
	socketNumber int
	socketState  int
	minutes      int
	changeState  int
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("dondu: ")

	flag.IntVar(&socketNumber, "socket", 0, "number of the socket to manipulate")
	flag.IntVar(&socketState, "state", 0, "state to set on socket ")
	flag.IntVar(&minutes, "minutes", 0, "time")

}

func main() {
	/*
		gniazdo_nr: 5
		gniazdo_stan: 0
		minuty:
		rozdzielnia_zmien_stan: 1
	*/

	data := url.Values{}
	data.Set("gniazdo_nr", strconv.Itoa(socketNumber))
	data.Set("gniazdo_stan", strconv.Itoa(socketState))
	data.Set("rozdzielnia_zmien_stan", strconv.Itoa(1))

	res, err := http.Post(ApiUrl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalf("failed to make HTTP POST request: %v\n", err)
	}

	log.Printf("res: %v\n", res)
}
