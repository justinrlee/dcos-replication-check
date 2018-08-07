package dcosclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	// "os"
	"time"

	"github.com/Sirupsen/logrus"
)

type Client struct {
	Host       string
	UserAgent  string
	httpClient *http.Client
	Token      string
}

func (c *Client) Setup() {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	
	c.httpClient = &http.Client{Transport: tr, Timeout: time.Second * 10}
	c.auth()
}

// For external use
func (c *Client) Request(method string, port int, path string) ([]byte, error) {
		
	req, err := c.newRequest(method, port, path, nil)

	if err != nil {
		fmt.Println("Error generating request")
		logrus.Infoln(err)
	}

	body, err := c.do(req)

	if err != nil {
		fmt.Println("Error running request")
		logrus.Infoln(err)
	}

	return body, err

}

// Internal use
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
	url := "https://" + c.Host + ":" + strconv.Itoa(port) + "/" + path

	logrus.Infoln(url)

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