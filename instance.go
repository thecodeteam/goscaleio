package goscaleio

import (
	"fmt"

	types "github.com/emccode/goscaleio/types/v1"
)

func (client *Client) GetInstance(systemhref string) (systems []*types.System, err error) {

	endpoint := client.SIOEndpoint
	if systemhref == "" {
		endpoint.Path += "/types/System/instances"
	} else {
		endpoint.Path = systemhref
	}

	req := client.NewRequest(map[string]string{}, "GET", endpoint, nil)
	req.SetBasicAuth("", client.Token)
	req.Header.Add("Accept", "application/json;version=1.0")

	resp, err := client.retryCheckResp(&client.Http, req)
	if err != nil {
		return []*types.System{}, fmt.Errorf("problem getting response: %v", err)
	}
	defer resp.Body.Close()

	if systemhref == "" {
		if err = client.decodeBody(resp, &systems); err != nil {
			return []*types.System{}, fmt.Errorf("error decoding instances response: %s", err)
		}
	} else {
		system := &types.System{}
		if err = client.decodeBody(resp, &system); err != nil {
			return []*types.System{}, fmt.Errorf("error decoding instances response: %s", err)
		}
		systems = append(systems, system)
	}

	// bs, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return types.Systems{}, errors.New("error reading body")
	// }

	return systems, nil
}
