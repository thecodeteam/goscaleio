package goscaleio

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/thecodeteam/goscaleio/api"
	types "github.com/thecodeteam/goscaleio/types/v1"
)

var (
	accHeader string
	conHeader string

	errNilReponse = errors.New("nil response from API")
	errBodyRead   = errors.New("error reading body")
	errNoLink     = errors.New("Error: problem finding link")

	debug, _    = strconv.ParseBool(os.Getenv("GOSCALEIO_DEBUG"))
	showHTTP, _ = strconv.ParseBool(os.Getenv("GOSCALEIO_SHOWHTTP"))
)

type Client struct {
	configConnect *ConfigConnect
	api           api.Client
}

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

func (c *Client) getVersion() (string, error) {

	resp, err := c.api.DoAndGetResponseBody(
		context.Background(), http.MethodGet, "/api/version", nil, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// parse the response
	switch {
	case resp == nil:
		return "", errNilReponse
	case !(resp.StatusCode >= 200 && resp.StatusCode <= 299):
		return "", c.api.ParseJSONError(resp)
	}

	version, err := extractString(resp)
	if err != nil {
		return "", err
	}

	versionRX := regexp.MustCompile(`^(\d+?\.\d+?).*$`)
	if m := versionRX.FindStringSubmatch(version); len(m) > 0 {
		return m[1], nil
	}
	return version, nil
}

func (c *Client) updateVersion() error {

	version, err := c.getVersion()
	if err != nil {
		return err
	}
	c.configConnect.Version = version

	accHeader = api.HeaderValContentTypeJSON
	if version != "" {
		accHeader = accHeader + ";version=" + version
	}
	conHeader = accHeader

	return nil
}

func (c *Client) Authenticate(configConnect *ConfigConnect) (Cluster, error) {

	configConnect.Version = c.configConnect.Version
	c.configConnect = configConnect

	c.api.SetToken("")

	headers := make(map[string]string, 1)
	headers["Authorization"] = "Basic " + basicAuth(
		configConnect.Username, configConnect.Password)

	resp, err := c.api.DoAndGetResponseBody(
		context.Background(), http.MethodGet, "api/login", headers, nil)
	if err != nil {
		doLog(log.WithError(err).Error, "")
		return Cluster{}, err
	}
	defer resp.Body.Close()

	// parse the response
	switch {
	case resp == nil:
		return Cluster{}, errNilReponse
	case !(resp.StatusCode >= 200 && resp.StatusCode <= 299):
		return Cluster{}, c.api.ParseJSONError(resp)
	}

	token, err := extractString(resp)
	if err != nil {
		return Cluster{}, nil
	}

	c.api.SetToken(token)

	if c.configConnect.Version == "" {
		err = c.updateVersion()
		if err != nil {
			return Cluster{}, errors.New("error getting version of ScaleIO")
		}
	}

	return Cluster{}, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (c *Client) getJSONWithRetry(
	method, uri string,
	body, resp interface{}) error {

	headers := make(map[string]string, 2)
	headers[api.HeaderKeyAccept] = accHeader
	headers[api.HeaderKeyContentType] = conHeader

	err := c.api.DoWithHeaders(
		context.Background(), method, uri, headers, body, resp)
	if err == nil {
		return nil
	}

	// check if we need to authenticate
	if e, ok := err.(*types.Error); ok {
		doLog(log.WithError(err).Debug, fmt.Sprintf("Got JSON error: %+v", e))
		if e.HTTPStatusCode == 401 {
			doLog(log.Info, "Need to re-auth")
			// Authenticate then try again
			if _, err := c.Authenticate(c.configConnect); err != nil {
				return fmt.Errorf("Error Authenticating: %s", err)
			}
			return c.api.Do(
				context.Background(),
				method, uri, nil, resp)
		}
	}
	doLog(log.WithError(err).Error, "returning error")

	return err
}

func extractString(resp *http.Response) (string, error) {
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errBodyRead
	}

	s := string(bs)

	s = strings.TrimRight(s, `"`)
	s = strings.TrimLeft(s, `"`)

	return s, nil
}

func (c *Client) getStringWithRetry(
	method, uri string,
	body interface{}) (string, error) {

	headers := make(map[string]string, 2)
	headers[api.HeaderKeyAccept] = accHeader
	headers[api.HeaderKeyContentType] = conHeader

	checkResponse := func(resp *http.Response) (string, bool, error) {
		defer resp.Body.Close()

		// parse the response
		switch {
		case resp == nil:
			return "", false, errNilReponse
		case resp.StatusCode == 401:
			return "", true, c.api.ParseJSONError(resp)
		case !(resp.StatusCode >= 200 && resp.StatusCode <= 299):
			return "", false, c.api.ParseJSONError(resp)
		}

		s, err := extractString(resp)
		if err != nil {
			return "", false, err
		}

		return s, false, nil
	}

	resp, err := c.api.DoAndGetResponseBody(
		context.Background(), method, uri, headers, body)
	if err != nil {
		return "", err
	}
	s, retry, httpErr := checkResponse(resp)
	if httpErr != nil {
		if retry {
			doLog(log.Info, "need to re-auth")
			// Authenticate then try again
			if _, err = c.Authenticate(c.configConnect); err != nil {
				return "", fmt.Errorf("Error Authenticating: %s", err)
			}
			resp, err = c.api.DoAndGetResponseBody(
				context.Background(), method, uri, headers, body)
			if err != nil {
				return "", err
			}
			s, _, err = checkResponse(resp)
		} else {
			return "", httpErr
		}
	}

	return s, nil
}

func (c *Client) SetToken(token string) {
	c.api.SetToken(token)
}

func (c *Client) GetToken() string {
	return c.api.GetToken()
}

func NewClient() (client *Client, err error) {
	return NewClientWithArgs(
		os.Getenv("GOSCALEIO_ENDPOINT"),
		os.Getenv("GOSCALEIO_VERSION"),
		os.Getenv("GOSCALEIO_INSECURE") == "true",
		os.Getenv("GOSCALEIO_USECERTS") == "true")
}

func NewClientWithArgs(
	endpoint string,
	version string,
	insecure,
	useCerts bool) (client *Client, err error) {

	if showHTTP {
		debug = true
	}

	fields := map[string]interface{}{
		"endpoint": endpoint,
		"insecure": insecure,
		"useCerts": useCerts,
		"version":  version,
		"debug":    debug,
		"showHTTP": showHTTP,
	}

	doLog(log.WithFields(fields).Debug, "goscaleio client init")

	if endpoint == "" {
		doLog(log.WithFields(fields).Error, "endpoint is required")
		return nil,
			withFields(fields, "endpoint is required")
	}

	opts := api.ClientOptions{
		Insecure: insecure,
		UseCerts: useCerts,
		ShowHTTP: showHTTP,
	}

	ac, err := api.New(context.Background(), endpoint, opts, debug)
	if err != nil {
		doLog(log.WithError(err).Error, "Unable to create HTTP client")
		return nil, err
	}

	client = &Client{
		api: ac,
		configConnect: &ConfigConnect{
			Version: version,
		},
	}

	accHeader = api.HeaderValContentTypeJSON
	if version != "" {
		accHeader = accHeader + ";version=" + version
	}
	conHeader = accHeader

	return client, nil
}

func GetLink(links []*types.Link, rel string) (*types.Link, error) {
	for _, link := range links {
		if link.Rel == rel {
			return link, nil
		}
	}

	return nil, errNoLink
}

func withFields(fields map[string]interface{}, message string) error {
	return withFieldsE(fields, message, nil)
}

func withFieldsE(
	fields map[string]interface{}, message string, inner error) error {

	if fields == nil {
		fields = make(map[string]interface{})
	}

	if inner != nil {
		fields["inner"] = inner
	}

	x := 0
	l := len(fields)

	var b bytes.Buffer
	for k, v := range fields {
		if x < l-1 {
			b.WriteString(fmt.Sprintf("%s=%v,", k, v))
		} else {
			b.WriteString(fmt.Sprintf("%s=%v", k, v))
		}
		x = x + 1
	}

	return fmt.Errorf("%s %s", message, b.String())
}

func doLog(
	l func(args ...interface{}),
	msg string) {

	if debug {
		l(msg)
	}
}
