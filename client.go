package goscaleio

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	types "github.com/emccode/goscaleio/types/v1"
)

type Client struct {
	Token         string
	SIOEndpoint   url.URL
	Http          http.Client
	Insecure      string
	showHttp      bool
	configConnect *ConfigConnect
}

//Cluster deprecate
type Cluster struct {
}

type ConfigConnect struct {
	Endpoint string
	Version  string
	Username string
	Password string
}

type ClientPersistent struct {
	configConnect *ConfigConnect
	client        *Client
}

// getVersion returns the API version supported by the ScaleIO API gateway.
func (c *Client) getVersion() (string, error) {
	// ver is populated when new Client created.
	if c.configConnect.Version != "" {
		return c.configConnect.Version, nil
	}

	endpoint := c.SIOEndpoint
	endpoint.Path = "/api/version"

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request for", endpoint)
		return "", err
	}
	req.Header.Add("User-Agent", "go-scaleio")
	log.WithField("url", req.URL).Debug("sending request")
	resp, err := c.Http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("error reading body")
	}

	version := string(bs)
	version = strings.TrimRight(version, `"`)
	version = strings.TrimLeft(version, `"`)
	c.configConnect.Version = version

	log.WithField("version", version).Debug("scaleio api version")

	return version, nil
}

// Authenticate connects and login to the Scaleio server.
func (c *Client) Authenticate(configConnect *ConfigConnect) error {

	ver := c.configConnect.Version
	c.configConnect = configConnect
	c.configConnect.Version = ver

	endpoint := c.SIOEndpoint
	endpoint.Path += "/login"

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request for", endpoint)
		return err
	}
	req.SetBasicAuth(configConnect.Username, configConnect.Password)
	resp, err := c.send(req)
	if err != nil {
		return err
	}

	if err := c.validateResponse(resp); err != nil {
		return err
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	token := string(bs)
	token = strings.TrimRight(token, `"`)
	token = strings.TrimLeft(token, `"`)
	c.Token = token
	log.WithField("token", token).Debug("received api token")

	return nil
}

// send sends request to remote server and returns an http respose or error.
func (c *Client) send(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", "go-scaleio")
	req.Header.Add("Accept", "application/json;version="+c.configConnect.Version)
	if req.ContentLength > 0 {
		req.Header.Add("Content-Type", "application/json;version="+c.configConnect.Version)
	}

	log.WithField("url", req.URL).Debug("sent request")
	resp, err := c.Http.Do(req)
	if log.GetLevel() == log.DebugLevel && c.showHttp {
		log.WithField("REQ", req).Debug("http request sent")
		log.WithField("RESP", resp).Debug("http response received")
	}
	// TODO(vladimirvivien) insert timeout/failure/backoff logic (eventually)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// validateResponse ensures the response is not error, if so, parse and return it.
func (c *Client) validateResponse(resp *http.Response) error {
	log.WithField("status", resp.StatusCode).Debug("status code rcvd")
	switch resp.StatusCode {
	case 200, 201, 202, 204:
		return nil
	default:
		errBody, _ := c.parseErr(resp)
		return errBody
	}
}

// decodeBody extracts the JSON body from a response.
func (c *Client) decodeBody(resp *http.Response, out interface{}) error {

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if log.GetLevel() == log.DebugLevel && c.showHttp {
		var prettyJSON bytes.Buffer
		_ = json.Indent(&prettyJSON, body, "", "  ")
		log.WithField("body", prettyJSON.String()).Debug(
			"decoded response body")
	}

	if err = json.Unmarshal(body, &out); err != nil {
		return err
	}

	return nil
}

func (c *Client) parseErr(resp *http.Response) (*types.Error, error) {
	//TODO(vladimirvivien) update to only return only error
	errBody := new(types.Error)

	// if there was an error decoding the body, just return that
	if err := c.decodeBody(resp, errBody); err != nil {
		return &types.Error{}, fmt.Errorf("error parsing error body for non-200 request: %s", err)
	}

	return errBody, nil //fmt.Errorf("API (%d) Error: %d: %s", resp.StatusCode, errBody.MajorErrorCode, errBody.Message)
}

// NewClient returns a pointer to a Client or error.
func NewClient() (c *Client, err error) {
	return NewClientWithArgs(
		os.Getenv("GOSCALEIO_ENDPOINT"),
		os.Getenv("GOSCALEIO_VERSION"),
		os.Getenv("GOSCALEIO_INSECURE") == "true",
		os.Getenv("GOSCALEIO_USECERTS") == "true")
}

// NewClientWithArgs uses the passed params to return a new Client
func NewClientWithArgs(
	endpoint string,
	version string,
	insecure,
	useCerts bool) (c *Client, err error) {

	debugEnabled := os.Getenv("GOSCALEIO_DEBUG") == "true"
	if debugEnabled {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	fields := map[string]interface{}{
		"endpoint": endpoint,
		"insecure": insecure,
		"useCerts": useCerts,
		"version":  version,
	}

	var uri *url.URL
	uri, err = url.ParseRequestURI(endpoint)
	if err != nil {
		log.WithFields(fields).Errorf("error parsing endpoint")
		return &Client{}, errors.New("error parsing endpoint")
	}

	c = &Client{
		SIOEndpoint: *uri,
		Http: http.Client{
			Transport: &http.Transport{
				TLSHandshakeTimeout: 120 * time.Second,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: insecure,
				},
			},
		},
	}

	if useCerts {
		log.Debug("Setting up secure client")

		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(pemCerts)

		c.Http.Transport = &http.Transport{
			TLSHandshakeTimeout: 120 * time.Second,
			TLSClientConfig: &tls.Config{
				RootCAs:            pool,
				InsecureSkipVerify: insecure,
			},
		}
	}

	c.showHttp = os.Getenv("GOSCALEIO_SHOW_HTTP") == "true"
	c.configConnect = &ConfigConnect{}

	ver, err := c.getVersion()
	if err != nil {
		log.WithError(err).Error("failed to get API version")
		return nil, err
	}

	if version != "" && version != ver {
		log.WithField("version", version).WithField("api", ver).
			Errorf("expecting api version %s, got %s", version, ver)
		return nil, fmt.Errorf("api version mismatched")
	}

	c.configConnect.Version = ver

	return c, nil
}
