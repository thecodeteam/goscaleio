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
	Insecure    string
}

type Cluster struct {
}

type ConfigConnect struct {
	Endpoint string
	Username string
	Password string
}

type ClientPersistent struct {
	configConnect *ConfigConnect
	client        *Client
}

var clientPersistentGlobal ClientPersistent

func (client *Client) Authenticate(configConnect *ConfigConnect) (Cluster, error) {
	clientPersistentGlobal.configConnect = configConnect
	clientPersistentGlobal.client = client

	endpoint := client.SIOEndpoint
	endpoint.Path += "/login"

	req := client.NewRequest(map[string]string{}, "GET", endpoint, nil)
	req.SetBasicAuth(configConnect.Username, configConnect.Password)
	req.Header.Add("Accept", "application/json;version=1.0")

	resp, err := retryCheckResp(&client.Http, req)
	if err != nil {
		return Cluster{}, fmt.Errorf("problem getting response: %v", err)
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Cluster{}, errors.New("error reading body")
	}

	token := string(bs)

	if os.Getenv("GOSCALEIO_SHOW_BODY") == "true" {
		fmt.Printf("%+v\n", token)
	}

	token = strings.TrimRight(token, `"`)
	token = strings.TrimLeft(token, `"`)
	client.Token = token

	return Cluster{}, nil
}

//https://github.com/chrislusf/teeproxy/blob/master/teeproxy.go
type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func DuplicateRequest(request *http.Request) (request1 *http.Request, request2 *http.Request) {
	request1 = &http.Request{
		Method:        request.Method,
		URL:           request.URL,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        request.Header,
		Host:          request.Host,
		ContentLength: request.ContentLength,
	}
	request2 = &http.Request{
		Method:        request.Method,
		URL:           request.URL,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        request.Header,
		Host:          request.Host,
		ContentLength: request.ContentLength,
	}

	if request.Body != nil {
		b1 := new(bytes.Buffer)
		b2 := new(bytes.Buffer)
		w := io.MultiWriter(b1, b2)
		io.Copy(w, request.Body)
		request1.Body = nopCloser{b1}
		request2.Body = nopCloser{b2}

		defer request.Body.Close()
	}

	return
}

func retryCheckResp(httpClient *http.Client, req *http.Request) (*http.Response, error) {

	req1, req2 := DuplicateRequest(req)
	resp, errBody, err := checkResp(httpClient.Do(req1))
	if errBody == nil && err != nil {
		return &http.Response{}, err
	} else if errBody != nil && err != nil {
		if resp != nil && resp.StatusCode == 401 && errBody.MajorErrorCode == 0 {
			_, err := clientPersistentGlobal.client.Authenticate(clientPersistentGlobal.configConnect)
			if err != nil {
				return nil, fmt.Errorf("Error re-authenticating: %s", err)
			}

			ioutil.ReadAll(resp.Body)
			resp.Body.Close()

			req2.SetBasicAuth("", clientPersistentGlobal.client.Token)
			resp, errBody, err = checkResp(httpClient.Do(req2))
			if err != nil {
				return &http.Response{}, errors.New(errBody.Message)
			}
		} else {
			return &http.Response{}, errors.New(errBody.Message)
		}
	}

	return resp, nil
}

func checkResp(resp *http.Response, err error) (*http.Response, *types.Error, error) {
	if err != nil {
		return resp, &types.Error{}, err
	}

	switch i := resp.StatusCode; {
	// Valid request, return the response.
	case i == 200 || i == 201 || i == 202 || i == 204:
		return resp, &types.Error{}, nil
	// Invalid request, parse the XML error returned and return it.
	case i == 400 || i == 401 || i == 403 || i == 404 || i == 405 || i == 406 || i == 409 || i == 415 || i == 500 || i == 503 || i == 504:
		errBody, err := parseErr(resp)
		return resp, errBody, err
	// Unhandled response.
	default:
		return nil, &types.Error{}, fmt.Errorf("unhandled API response, please report this issue, status code: %s", resp.Status)
	}
}

func decodeBody(resp *http.Response, out interface{}) error {

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if os.Getenv("GOSCALEIO_SHOW_BODY") == "true" {
		var prettyJSON bytes.Buffer
		_ = json.Indent(&prettyJSON, body, "", "    ")
		fmt.Printf("%+v\n", prettyJSON.String())
	}

	if err = json.Unmarshal(body, &out); err != nil {
		return err
	}

	return nil
}

func parseErr(resp *http.Response) (*types.Error, error) {

	errBody := new(types.Error)

	// if there was an error decoding the body, just return that
	if err := decodeBody(resp, errBody); err != nil {
		return &types.Error{}, fmt.Errorf("error parsing error body for non-200 request: %s", err)
	}

	return errBody, fmt.Errorf("API (%d) Error: %d: %s", resp.StatusCode, errBody.MajorErrorCode, errBody.Message)
}

func (c *Client) NewRequest(params map[string]string, method string, u url.URL, body io.Reader) *http.Request {

	debug := os.Getenv("SCALEIO_SHOW_BODY")
	if debug == "true" && body != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(body)
		fmt.Printf("\n\nDEBUG: \n%s\n\n", buf.String())
	}

	p := url.Values{}

	for k, v := range params {
		p.Add(k, v)
	}

	u.RawQuery = p.Encode()

	req, _ := http.NewRequest(method, u.String(), body)

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

func GetLink(links []*types.Link, rel string) (*types.Link, error) {
	for _, link := range links {
		if link.Rel == rel {
			return link, nil
		}
	}

	return &types.Link{}, errors.New("Couldn't find link")
}
