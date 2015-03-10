package goscaleio

import (
	"fmt"

	types "github.com/emccode/goscaleio/types/v1"
)

func (client *Client) GetInstance() (systems []types.System, err error) {

	endpoint := client.SIOEndpoint
	endpoint.Path += "/types/System/instances"

	req := client.NewRequest(map[string]string{}, "GET", endpoint, nil)
	req.SetBasicAuth("", client.Token)
	req.Header.Add("Accept", "application/json;version=1.0")

	resp, err := checkResp(client.Http.Do(req))
	if err != nil {
		return []types.System{}, fmt.Errorf("problem getting response: %v", err)
	}
	defer resp.Body.Close()

	if err = decodeBody(resp, &systems); err != nil {
		return []types.System{}, fmt.Errorf("error decoding instances response: %s", err)
	}

	// bs, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return types.Systems{}, errors.New("error reading body")
	// }

	return systems, nil
}
