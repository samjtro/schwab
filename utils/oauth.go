// The fastest unofficial Schwab TraderAPI wrapper
// Copyright (C) 2024 Samuel Troyer <samjtro.com>
// See the GNU General Public License for more details
package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var (
	ctx = context.Background()
)

func init() {
	err := godotenv.Load("config.env")
	check(err)
	conf = &oauth2.Config{
		ClientID:     os.Getenv("APPKEY"),
		ClientSecret: os.Getenv("SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL: fmt.Sprintf("https://api.schwabapi.com/v1/oauth/authorize?client_id=%s&redirect_uri=%s", os.Getenv("APPKEY"), os.Getenv("CBURL")),
		},
	}
	verifier = oauth2.GenerateVerifier()
}

// Initiate the Schwab oAuth process to retrieve bearer/refresh tokens
func Initiate() *Agent {
	agent := Agent{}
	if _, err := os.Stat(fmt.Sprintf("%s/.trade", homeDir())); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(fmt.Sprintf("%s/.trade", homeDir()), os.ModePerm)
		check(err)
		// oAuth Leg 1 - Authorization Code
		openBrowser(fmt.Sprintf("https://api.schwabapi.com/v1/oauth/authorize?client_id=%s&redirect_uri=%s", os.Getenv("APPKEY"), os.Getenv("CBURL")))
		fmt.Printf("Log into your Schwab brokerage account. Copy Error404 URL and paste it here: ")
		var urlInput string
		fmt.Scanln(&urlInput)
		authCodeEncoded := getStringInBetween(urlInput, "?code=", "&session=")
		authCode, err := url.QueryUnescape(authCodeEncoded)
		check(err)
		// oAuth Leg 2 - Refresh, Bearer Tokens
		authStringLegTwo := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", os.Getenv("APPKEY"), os.Getenv("SECRET")))))
		client := http.Client{}
		payload := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s", string(authCode), os.Getenv("CBURL"))
		req, err := http.NewRequest("POST", "https://api.schwabapi.com/v1/oauth/token", bytes.NewBuffer([]byte(payload)))
		check(err)
		req.Header = http.Header{
			"Authorization": {authStringLegTwo},
			"Content-Type":  {"application/x-www-form-urlencoded"},
		}
		res, err := client.Do(req)
		check(err)
		defer res.Body.Close()
		bodyBytes, err := io.ReadAll(res.Body)
		check(err)
		agent.tokens = parseAccessTokenResponse(string(bodyBytes))
		tokensJson, err := json.Marshal(agent.tokens)
		check(err)
		err = os.WriteFile(fmt.Sprintf("%s/.trade/bar.json", homeDir()), tokensJson, 0777)
		check(err)
	} else {
		agent.tokens = readDB()
		if agent.tokens.Bearer == "" {
			err := os.RemoveAll(fmt.Sprintf("%s/.trade", homeDir()))
			check(err)
			log.Fatalf("[err] please reinitiate, something went wrong\n")
		}
	}
	return &agent
}

// Use refresh token to generate a new bearer token for authentication
func (agent *Agent) refresh() {
	oldTokens := readDB()
	authStringRefresh := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", os.Getenv("APPKEY"), os.Getenv("SECRET")))))
	client := http.Client{}
	req, err := http.NewRequest("POST", "https://api.schwabapi.com/v1/oauth/token", bytes.NewBuffer([]byte(fmt.Sprintf("grant_type=refresh_token&refresh_token=%s", oldTokens.Refresh))))
	check(err)
	req.Header = http.Header{
		"Authorization": {authStringRefresh},
		"Content-Type":  {"application/x-www-form-urlencoded"},
	}
	res, err := client.Do(req)
	check(err)
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	check(err)
	agent.tokens = parseAccessTokenResponse(string(bodyBytes))
}
