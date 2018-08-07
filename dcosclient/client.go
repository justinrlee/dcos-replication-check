package dcosclient

import (
	"bytes"
	// "crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	// "os"
	// "time"

	"github.com/Sirupsen/logrus"
)

type Client struct {
	Host       string
	UserAgent  string
	httpClient *http.Client
	Token      string
}

func getClient() {

}

func (c *Client) newRequest(method string, port int, path string, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	// Generate URL from parameters.  
	// TODO adapt to try both http and https
	url := "https://" + c.Host + ":" + string(port) + "/" + path

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "token="+c.Token)
	} else {
		req.Header.Del("Authorization")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	return req, err
}

func (c *Client) do(req *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if 200 != resp.StatusCode {
		if 401 == resp.StatusCode {
			logrus.Infoln("Authentication expired. Re-authorizing account")
			c.auth()
		} else {
			return nil, fmt.Errorf("%s", body)
		}
	}
	return body, err

}