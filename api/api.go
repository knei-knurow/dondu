package api

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	baseURL    = "https://pracownia.knei.pl"
	loginURL   = baseURL + "/" + "?m=wyloguj" // login subpage address
	actionsURL = baseURL + "/" + "akcje.php"  // actions subpage address
)

// Login performs user authentication (sets the PHPSESSID cookie).
func Login(username string, password string) error {
	data := url.Values{}
	data.Set("login", username)
	data.Set("haslo", password)

	res, err := http.Post(loginURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to make HTTP POST login request: %v", err)
	}

	log.Println("auth status:", res.Status)

	return nil
}

// Update modifies state of single socket on the switchboard.
func Update(socket int, enabled bool, minutes int) error {
	state := 0
	if enabled {
		state = 1
	}

	data := url.Values{}
	data.Set("gniazdo_nr", strconv.Itoa(socket))
	data.Set("gniazdo_stan", strconv.Itoa(state))
	data.Set("minuty", strconv.Itoa(minutes))
	data.Set("rozdzielnia_zmien_stan", strconv.Itoa(1))

	_, err := http.Post(actionsURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to make HTTP POST request: %v", err)
	}

	return nil
}
