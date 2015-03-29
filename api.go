package goscaleio

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	types "github.com/emccode/goscaleio/types/v1"
)

type Client struct {
	Token       string
	SIOEndpoint url.URL
	Http        http.Client
}

type Cluster struct {
}

type ConfigConnect struct {
	Endpoint string
	Username string
	Password string
}

func (client *Client) Authenticate(configConnect *ConfigConnect) (Cluster, error) {

	endpoint := client.SIOEndpoint
	endpoint.Path += "/login"

	req := client.NewRequest(map[string]string{}, "GET", endpoint, nil)
	req.SetBasicAuth(configConnect.Username, configConnect.Password)
	req.Header.Add("Accept", "application/json;version=1.0")

	resp, err := checkResp(client.Http.Do(req))
	if err != nil {
		return Cluster{}, fmt.Errorf("problem getting response: %v", err)
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Cluster{}, errors.New("error reading body")
	}

	token := string(bs)

	if os.Getenv("SCALEIO_SHOW_BODY") == "true" {
		fmt.Printf("%+v\n", token)
	}

	token = strings.TrimRight(token, `"`)
	token = strings.TrimLeft(token, `"`)
	client.Token = token

	return Cluster{}, nil
}

func checkResp(resp *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return resp, err
	}

	switch i := resp.StatusCode; {
	// Valid request, return the response.
	case i == 200 || i == 201 || i == 202 || i == 204:
		return resp, nil
	// Invalid request, parse the XML error returned and return it.
	case i == 400 || i == 401 || i == 403 || i == 404 || i == 405 || i == 406 || i == 409 || i == 415 || i == 500 || i == 503 || i == 504:
		return nil, parseErr(resp)
	// Unhandled response.
	default:
		return nil, fmt.Errorf("unhandled API response, please report this issue, status code: %s", resp.Status)
	}
}

func decodeBody(resp *http.Response, out interface{}) error {

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if os.Getenv("SCALEIO_SHOW_BODY") == "true" {
		var prettyJSON bytes.Buffer
		_ = json.Indent(&prettyJSON, body, "", "    ")
		fmt.Printf("%+v\n", prettyJSON.String())
	}

	if err = json.Unmarshal(body, &out); err != nil {
		return err
	}

	return nil
}

func parseErr(resp *http.Response) error {

	errBody := new(types.Error)

	// if there was an error decoding the body, just return that
	if err := decodeBody(resp, errBody); err != nil {
		return fmt.Errorf("error parsing error body for non-200 request: %s", err)
	}

	return fmt.Errorf("API Error: %d: %s", errBody.MajorErrorCode, errBody.Message)
}

func (c *Client) NewRequest(params map[string]string, method string, u url.URL, body io.Reader) *http.Request {

	debug := os.Getenv("SCALEIO_SHOW_BODY")
	if debug == "true" && body != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(body)
		fmt.Printf("\n\nXML DEBUG: \n%s\n\n", buf.String())
	}

	p := url.Values{}

	// Build up our request parameters
	for k, v := range params {
		p.Add(k, v)
	}

	// Add the params to our URL
	u.RawQuery = p.Encode()

	// Build the request, no point in checking for errors here as we're just
	// passing a string version of an url.URL struct and http.NewRequest returns
	// error only if can't process an url.ParseRequestURI().
	req, _ := http.NewRequest(method, u.String(), body)

	// if c.VCDAuthHeader != "" && c.VCDToken != "" {
	// 	// Add the authorization header
	// 	req.Header.Add(c.VCDAuthHeader, c.VCDToken)
	// 	// Add the Accept header for VCD
	// 	req.Header.Add("Accept", "application/*+xml;version=5.6")
	// }

	return req

}

func NewClient() (client *Client, err error) {

	var uri *url.URL

	if os.Getenv("GOSCALEIO_ENDPOINT") != "" {
		uri, err = url.ParseRequestURI(os.Getenv("GOSCALEIO_ENDPOINT"))
		if err != nil {
			return &Client{}, fmt.Errorf("cannot parse endpoint coming from VCLOUDAIR_ENDPOINT")
		}
	} else {
		return &Client{}, errors.New("missing GOSCALEIO_ENDPOINT")

	}

	var insecureSkipVerify bool
	if os.Getenv("GOSCALEIO_INSECURE") == "true" {
		insecureSkipVerify = true
	} else {
		insecureSkipVerify = false
	}

	client = &Client{
		SIOEndpoint: *uri,
		Http: http.Client{
			Transport: &http.Transport{
				TLSHandshakeTimeout: 120 * time.Second,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: insecureSkipVerify,
				},
			},
		},
	}

	if os.Getenv("GOSCALEIO_USECERTS") == "true" {
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(pemCerts)

		client.Http.Transport = &http.Transport{
			TLSHandshakeTimeout: 120 * time.Second,
			TLSClientConfig: &tls.Config{
				RootCAs:            pool,
				InsecureSkipVerify: insecureSkipVerify,
			},
		}
	}

	return client, nil
}
