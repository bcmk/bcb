package lib

import (
	"fmt"
	"net/http"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

var offlineDiv = cascadia.MustCompile("div.status-off")
var bannedDiv = cascadia.MustCompile("div.banned")

// CheckModelStripchat checks Stripchat model status
func CheckModelStripchat(client *Client, modelID string, userAgent string, dbg bool) StatusKind {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://stripchat.com/%s", modelID), nil)
	CheckErr(err)
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}
	resp, err := client.Client.Do(req)
	if err != nil {
		Lerr("[%v] cannot send a query, %v", client.Addr, err)
		return StatusUnknown
	}
	defer func() {
		CheckErr(resp.Body.Close())
	}()
	if resp.StatusCode == 404 {
		return StatusNotFound
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		Lerr("[%v] cannot parse body for model %s, %v", client.Addr, modelID, err)
		return StatusUnknown
	}

	if offlineDiv.MatchFirst(doc) != nil {
		return StatusOffline
	}

	if bannedDiv.MatchFirst(doc) != nil {
		return StatusDenied
	}

	return StatusOnline
}
