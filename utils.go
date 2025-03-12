/*
Copyright (C) 2025 github.com/go-schwab

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, see
<https://www.gnu.org/licenses/>.
*/

package trader

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	o "github.com/go-schwab/utils/oauth"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type Agent struct {
	Client *o.AuthorizedClient
	Tokens Token
}

type Token struct {
	RefreshExpiration time.Time
	Refresh           string
	BearerExpiration  time.Time
	Bearer            string
}

var (
	APPKEY string
	SECRET string
	CBURL  string
	PATH   string
)

// load env variables, check if you've run the program before
func init() {
	err := godotenv.Load(findAllEnvFiles()...)
	isErrNil(err)
	APPKEY = os.Getenv("APPKEY")
	SECRET = os.Getenv("SECRET")
	CBURL = os.Getenv("CBURL")
	homedir, err := os.UserHomeDir()
	isErrNil(err)
	PATH = homedir + "/.config/go-schwab/.json"
	if _, err := os.Stat(homedir + "/.config/go-schwab"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(homedir+"/.config/go-schwab", 0750)
		isErrNil(err)
	}
}

// find all env files
func findAllEnvFiles() []string {
	var files []string
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		split := strings.Split(d.Name(), ".")
		if len(split) > 1 {
			if split[1] == "env" {
				files = append(files, d.Name())
			}
		}
		return err
	})
	isErrNil(err)
	return files
}

// trim one FIRST & one LAST character in the string
func trimOneFirstOneLast(s string) string {
	if len(s) < 1 {
		return ""
	}
	return s[1 : len(s)-1]
}

// parse access token response
func parseAccessTokenResponse(s string) Token {
	token := Token{
		RefreshExpiration: time.Now().Add(time.Hour * 168),
		BearerExpiration:  time.Now().Add(time.Minute * 30),
	}
	for _, x := range strings.Split(s, ",") {
		for i1, x1 := range strings.Split(x, ":") {
			if trimOneFirstOneLast(x1) == "refresh_token" {
				token.Refresh = trimOneFirstOneLast(strings.Split(x, ":")[i1+1])
			} else if trimOneFirstOneLast(x1) == "access_token" {
				token.Bearer = trimOneFirstOneLast(strings.Split(x, ":")[i1+1])
			}
		}
	}
	return token
}

// Credit: https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		log.Fatal("Unsupported platform.")
	}
	isErrNil(err)
}

// Execute a command @ stdin, receive stdout
func execCommand(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}

// Credit: https://go.dev/play/p/C2sZRYC15XN
func getStringInBetween(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return
	}
	return str[s : s+e]
}

// read in tokens from PATH - linux
func readLinuxDB() Token {
	var tokens Token
	body, err := os.ReadFile(PATH)
	isErrNil(err)
	err = sonic.Unmarshal(body, &tokens)
	isErrNil(err)
	return tokens
}

// read in tokens from PATH - mac & windows
func readDB() Agent {
	var tok *oauth2.Token
	body, err := os.ReadFile(PATH)
	isErrNil(err)
	err = sonic.Unmarshal(body, &tok)
	isErrNil(err)
	conf := &oauth2.Config{
		ClientID:     APPKEY, // Schwab App Key
		ClientSecret: SECRET, // Schwab App Secret
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.schwabapi.com/v1/oauth/authorize",
			TokenURL: "https://api.schwabapi.com/v1/oauth/token",
		},
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}
	sslcli := &http.Client{Transport: tr}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, sslcli)
	return Agent{
		Client: &o.AuthorizedClient{
			conf.Client(ctx, tok),
			tok,
		},
	}
}

// create Agent - linux
func initiateLinux() Agent {
	var agent Agent
	// oAuth Leg 1 - Authorization Code
	openBrowser(fmt.Sprintf("https://api.schwabapi.com/v1/oauth/authorize?client_id=%s&redirect_uri=%s", os.Getenv("APPKEY"), os.Getenv("CBURL")))
	fmt.Printf("Log into your Schwab brokerage account. Copy Error404 URL and paste it here: ")
	var urlInput string
	fmt.Scanln(&urlInput)
	authCodeEncoded := getStringInBetween(urlInput, "?code=", "&session=")
	authCode, err := url.QueryUnescape(authCodeEncoded)
	isErrNil(err)
	// oAuth Leg 2 - Refresh, Bearer Tokens
	authStringLegTwo := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", os.Getenv("APPKEY"), os.Getenv("SECRET")))))
	client := http.Client{}
	payload := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s", string(authCode), os.Getenv("CBURL"))
	req, err := http.NewRequest("POST", "https://api.schwabapi.com/v1/oauth/token", bytes.NewBuffer([]byte(payload)))
	isErrNil(err)
	req.Header = http.Header{
		"Authorization": {authStringLegTwo},
		"Content-Type":  {"application/x-www-form-urlencoded"},
	}
	res, err := client.Do(req)
	isErrNil(err)
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	isErrNil(err)
	agent.Tokens = parseAccessTokenResponse(string(bodyBytes))
	bytes, err := sonic.Marshal(agent.Tokens)
	isErrNil(err)
	err = os.WriteFile(PATH, bytes, 0750)
	isErrNil(err)
	return agent
}

func initiateMacWindows() Agent {
	var agent Agent
	// execCommand("openssl req -x509 -out localhost.crt -keyout localhost.key   -newkey rsa:2048 -nodes -sha256   -subj '/CN=localhost' -extensions EXT -config <(;printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS.1:localhost,IP:127.0.0.1\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")")
	agent = Agent{Client: o.Initiate(APPKEY, SECRET)}
	bytes, err := sonic.Marshal(agent.Client.Token)
	isErrNil(err)
	err = os.WriteFile(PATH, bytes, 0750)
	isErrNil(err)
	return agent
}

// create Agent - mac & windows
func Initiate() *Agent {
	var agent Agent
	if runtime.GOOS == "linux" {
		if _, err := os.Stat(PATH); errors.Is(err, os.ErrNotExist) {
			agent = initiateLinux()
		} else {
			agent.Tokens = readLinuxDB()
		}
	} else {
		if _, err := os.Stat(PATH); errors.Is(err, os.ErrNotExist) {
			agent = initiateMacWindows()
		} else {
			agent = readDB()
		}
	}
	return &agent
}

func Reinitiate() *Agent {
	var agent Agent
	if _, err := os.Stat(PATH); !errors.Is(err, os.ErrNotExist) {
		err := os.Remove(PATH)
		isErrNil(err)
	}
	if runtime.GOOS == "linux" {
		agent = initiateLinux()
	} else {
		agent = initiateMacWindows()
	}
	return &agent
}

// use refresh to generate a new bearer token for authentication
func (agent *Agent) Refresh() {
	oldTokens := readLinuxDB()
	authStringRefresh := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", os.Getenv("APPKEY"), os.Getenv("SECRET")))))
	client := http.Client{}
	req, err := http.NewRequest("POST", "https://api.schwabapi.com/v1/oauth/token", bytes.NewBuffer([]byte(fmt.Sprintf("grant_type=refresh_token&refresh_token=%s", oldTokens.Refresh))))
	isErrNil(err)
	req.Header = http.Header{
		"Authorization": {authStringRefresh},
		"Content-Type":  {"application/x-www-form-urlencoded"},
	}
	res, err := client.Do(req)
	isErrNil(err)
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	isErrNil(err)
	agent.Tokens = parseAccessTokenResponse(string(bodyBytes))
}

// Handler is the general purpose request function for the td-ameritrade api, all functions will be routed through this handler function, which does all of the API calling work
// It performs a GET request after adding the apikey found in the config.env file in the same directory as the program calling the function,
// then returns the body of the GET request's return.
// It takes one parameter:
// req = a request of type *http.Request
func (agent *Agent) Handler(req *http.Request) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)
	if runtime.GOOS == "linux" {
		if !time.Now().Before(agent.Tokens.BearerExpiration) {
			agent.Refresh()
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", agent.Tokens.Bearer))
		client := http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			agent = Reinitiate()
		}
	} else {
		resp, err = agent.Client.Do(req)
		if err != nil {
			agent = Reinitiate()
		}
	}
	switch true {
	case resp.StatusCode == 200:
		return resp, nil
	case resp.StatusCode == 201:
		return resp, nil
	case resp.StatusCode == 401:
		body, err := io.ReadAll(resp.Body)
		isErrNil(err)
		if strings.Contains(string(body), "\"status\": 500") {
			return nil, WrapTraderError(ErrUnexpectedServer, resp)
		}
		return nil, WrapTraderError(ErrNeedReAuthorization, resp)
	case resp.StatusCode == 403:
		return nil, WrapTraderError(ErrForbidden, resp)
	case resp.StatusCode == 404:
		return nil, WrapTraderError(ErrNotFound, resp)
	case resp.StatusCode == 500:
		return nil, WrapTraderError(ErrUnexpectedServer, resp)
	case resp.StatusCode == 503:
		return nil, WrapTraderError(ErrTemporaryServer, resp)
	case resp.StatusCode == 400:
		body, err := io.ReadAll(resp.Body)
		isErrNil(err)
		if strings.Contains(string(body), "\"status\": 500") {
			return nil, WrapTraderError(ErrUnexpectedServer, resp)
		}
		// if io.ReadAll() fails:
		//     return nil, WrapTraderError(err, StatusCode, "could not read response", nil)
		// if sonic.Unmarshall() fails
		//     return nil, WrapTraderError(err, StatusCode, "could not unmarshall", nil)
		// Note: The two above situations would wrap the errors generated by io or sonic

		// otherwise okay but the API was unhappy with our request:
		// At this point we could populate an ErrorMessage struct based on Schwab definition
		//   which contains Message string; Error []string
		return nil, WrapTraderError(ErrValidation, resp)
	default:
		return nil, errors.New("error not defined")
	}
}
