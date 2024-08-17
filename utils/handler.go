// The fastest unofficial Schwab TraderAPI wrapper
// Copyright (C) 2024 Samuel Troyer <samjtro.com>
// See the GNU General Public License for more details
package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Handler is the general purpose request function for the td-ameritrade api, all functions will be routed through this handler function, which does all of the API calling work
// It performs a GET request after adding the apikey found in the config.env file in the same directory as the program calling the function,
// then returns the body of the GET request's return.
// It takes one parameter:
// req = a request of type *http.Request
func (agent *Agent) Handler(req *http.Request) (*http.Response, error) {
	if (&Agent{}) == agent {
		log.Fatal("[ERR] empty agent - call 'Agent.Initiate' before making any API function calls.")
	}
	if !time.Now().Before(agent.tokens.BearerExpiration) {
		agent.refresh()
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", agent.tokens.Bearer))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode == 401 {
		err := os.Remove(fmt.Sprintf("%s/.trade", homeDir()))
		check(err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		check(err)
		log.Fatalf("[ERR] %d, %s", resp.StatusCode, body)
	}
	return resp, nil
}
