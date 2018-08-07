package dcos-client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

type DcosClient struct {
	BaseURL    string
	UserAgent  string
	httpClient *http.Client
	Token      string
}