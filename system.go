package goscaleio

import (
	"errors"
	"fmt"

	types "github.com/emccode/goscaleio/types/v1"
)

type System struct {
	System *types.System
	client *Client
}

func NewSystem(client *Client) *System {
	return &System{
		System: new(types.System),
		client: client,
	}
}

func (client *Client) FindSystem(instanceID, href string) (*System, error) {
	systems, err := client.GetInstance(href)
	if err != nil {
		return &System{}, fmt.Errorf("err: problem getting instances: %s", err)
	}

	for _, system := range systems {
		if system.ID == instanceID || href != "" {
			outSystem := NewSystem(client)
			outSystem.System = system
			return outSystem, nil
		}
	}
	return &System{}, fmt.Errorf("err: systemid not found")
}

func (system *System) GetStatistics() (statistics *types.Statistics, err error) {
	endpoint := system.client.SIOEndpoint
	// endpoint.Path = fmt.Sprintf("/api/instances/System::%v/relationships/Statistics", system.System.ID)

	link, err := GetLink(system.System.Links, "/api/System/relationship/Statistics")
	if err != nil {
		return &types.Statistics{}, errors.New("Error: problem finding link")
	}

	endpoint.Path = link.HREF

	req := system.client.NewRequest(map[string]string{}, "GET", endpoint, nil)
	req.SetBasicAuth("", system.client.Token)
	req.Header.Add("Accept", "application/json;version=1.0")

	resp, err := checkResp(system.client.Http.Do(req))
	if err != nil {
		return &types.Statistics{}, fmt.Errorf("problem getting response: %v", err)
	}
	defer resp.Body.Close()

	if err = decodeBody(resp, &statistics); err != nil {
		return &types.Statistics{}, fmt.Errorf("error decoding instances response: %s", err)
	}

	// bs, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return errors.New("error reading body")
	// }
	//
	// fmt.Println(string(bs))
	return statistics, nil
}
