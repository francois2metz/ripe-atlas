// common.go
//
// This file implements the configuration part for when you need the API
// key to modify things in the Atlas configuration and manage measurements.

package atlas

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	"github.com/pkg/errors"
)

const (
	apiEndpoint = "https://atlas.ripe.net/api/v2"
)

// getPageNum returns the value of the page= parameter
func getPageNum(url string) (page string) {
	re := regexp.MustCompile(`page=(\d+)`)
	if m := re.FindStringSubmatch(url); len(m) >= 1 {
		return m[1]
	}
	return ""
}

// AddQueryParameters adds query parameters to the URL.
func AddQueryParameters(baseURL string, queryParams map[string]string) string {
	if len(queryParams) == 0 {
		return baseURL
	}
	baseURL += "?"
	params := url.Values{}
	for key, value := range queryParams {
		params.Add(key, value)
	}
	return baseURL + params.Encode()
}

// addAPIKey insert the key into options if needed
func (c *Client) addAPIKey(opts map[string]string) map[string]string {
	key, ok := c.HasAPIKey()
	// Insert key
	if ok {
		opts["key"] = key
	}
	return opts
}

// prepareRequest insert all pre-defined stuff
func (c *Client) prepareRequest(method, what string, opts map[string]string) (req *http.Request) {
	var endPoint string

	// This is a hack to fetch direct urls for results
	if method == "FETCH" {
		endPoint = what
		method = "GET"
	} else {
		if c.config.endpoint != "" {
			endPoint = fmt.Sprintf("%s/%s", c.config.endpoint, what)
		}
	}

	c.mergeGlobalOptions(opts)
	c.debug("Options:\n%v", opts)
	baseURL := AddQueryParameters(endPoint, opts)

	req, err := http.NewRequest(method, baseURL, nil)
	if err != nil {
		c.log.Printf("error parsing %s: %v", baseURL, err)
		return &http.Request{}
	}

	c.debug("req.url=%s", baseURL)
	// We need these when we POST
	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
	}

	return
}

// decodeAPIError does the deserialisation
func decodeAPIError(body []byte) (*APIError, error) {
	var e APIError

	err := json.Unmarshal(body, &e)
	return &e, err
}

// client.handleAPIResponse check status code & return undecoded APIError
func (c *Client) handleAPIResponse(r *http.Response) ([]byte, error) {
	var (
		err  error
		body []byte
	)

	if r == nil {
		return []byte{}, fmt.Errorf("handleAPIResponse/r is nil")
	}

	// If there is a body, get it
	if r.Body != nil {
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return []byte{}, errors.Wrap(err, "read body")
		}
		defer r.Body.Close()
	}

	// Everything is fine
	if r.StatusCode == http.StatusOK || r.StatusCode == 0 {
		c.debug("OK")
		return body, nil
	}

	// Everything is fine too (200-2xx)
	if r.StatusCode >= http.StatusOK && r.StatusCode < http.StatusMultipleChoices {
		c.debug("OKbis")
		return body, nil
	}

	// Check this condition (3xx are handled directly)
	if r.StatusCode >= http.StatusMultipleChoices && r.StatusCode < http.StatusBadRequest {
		c.debug("mostlyOK")
		return body, nil
	}

	c.debug("err=%s", string(body))

	apie, err := decodeAPIError(body)
	if err != nil {
		return body, errors.Wrap(err, "decodeAPIError")
	}

	c.debug("the end")
	return body, fmt.Errorf(apie.Error())
}

func (c *Client) mergeGlobalOptions(opts map[string]string) {
	opts = mergeOptions(opts, c.opts)
}

func mergeOptions(from, with map[string]string) map[string]string {
	m := from
	for i, opt := range with {
		// "" means delete
		if opt != "" {
			m[i] = opt
		} else {
			delete(m, i)
		}
	}
	return m
}
