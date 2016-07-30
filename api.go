package goscaleio

//https://github.com/chrislusf/teeproxy/blob/master/teeproxy.go
import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	types "github.com/emccode/goscaleio/types/v1"
)

func (c *Client) updateVersion() error {

	version, err := c.getVersion()
	if err != nil {
		return err
	}
	c.configConnect.Version = version

	return nil
}

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

func (c *Client) retryCheckResp(httpClient *http.Client, req *http.Request) (*http.Response, error) {

	req1, req2 := DuplicateRequest(req)
	resp, errBody, err := c.checkResp(httpClient.Do(req1))
	if errBody == nil && err != nil {
		return &http.Response{}, err
	} else if errBody != nil && err != nil {
		if resp == nil {
			return nil, errors.New("Problem getting response from endpoint")
		}

		if resp.StatusCode == 401 {
			err := c.Authenticate(c.configConnect)
			if err != nil {
				return nil, fmt.Errorf("Error re-authenticating: %s", err)
			}

			ioutil.ReadAll(resp.Body)
			resp.Body.Close()

			req2.SetBasicAuth("", c.Token)
			resp, errBody, err = c.checkResp(httpClient.Do(req2))
			if err != nil {
				return &http.Response{}, errors.New(errBody.Message)
			}
		} else {
			return &http.Response{}, errors.New(errBody.Message)
		}
	}

	return resp, nil
}

func (c *Client) checkResp(resp *http.Response, err error) (*http.Response, *types.Error, error) {
	if err != nil {
		return resp, &types.Error{}, err
	}

	switch i := resp.StatusCode; {
	// Valid request, return the response.
	case i == 200 || i == 201 || i == 202 || i == 204:
		return resp, &types.Error{}, nil
	// Invalid request, parse the XML error returned and return it.
	case i == 400 || i == 401 || i == 403 || i == 404 || i == 405 || i == 406 || i == 409 || i == 415 || i == 500 || i == 503 || i == 504:
		errBody, err := c.parseErr(resp)
		return resp, errBody, err
	// Unhandled response.
	default:
		return nil, &types.Error{}, fmt.Errorf("unhandled API response, please report this issue, status code: %s", resp.Status)
	}
}

// NewRequest returns a point to a http.Request value constructed.
func (c *Client) NewRequest(params map[string]string, method string, u url.URL, body io.Reader) *http.Request {

	if log.GetLevel() == log.DebugLevel && c.showHttp && body != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(body)
		log.WithField("body", buf.String()).Debug("print new request body")
	}

	p := url.Values{}

	for k, v := range params {
		p.Add(k, v)
	}

	u.RawQuery = p.Encode()

	req, _ := http.NewRequest(method, u.String(), body)

	return req

}

func GetLink(links []*types.Link, rel string) (*types.Link, error) {
	for _, link := range links {
		if link.Rel == rel {
			return link, nil
		}
	}

	return &types.Link{}, errors.New("Couldn't find link")
}
