package goscaleio

import (
	"errors"
	"fmt"

	types "github.com/emccode/goscaleio/types/v1"
)

type ProtectionDomain struct {
	ProtectionDomain *types.ProtectionDomain
	client           *Client
}

func NewProtectionDomain(client *Client) *ProtectionDomain {
	return &ProtectionDomain{
		ProtectionDomain: new(types.ProtectionDomain),
		client:           client,
	}
}

func (system *System) GetProtectionDomain(protectiondomainhref string) (protectionDomains []*types.ProtectionDomain, err error) {

	// if systemhref == "" {
	// 	if err = decodeBody(resp, &systems); err != nil {
	// 		return []*types.System{}, fmt.Errorf("error decoding instances response: %s", err)
	// 	}
	// } else {
	// 	system := &types.System{}
	// 	if err = decodeBody(resp, &system); err != nil {
	// 		return []*types.System{}, fmt.Errorf("error decoding instances response: %s", err)
	// 	}
	// 	systems = append(systems, system)
	// }

	endpoint := system.client.SIOEndpoint

	if protectiondomainhref == "" {
		link, err := GetLink(system.System.Links, "/api/System/relationship/ProtectionDomain")
		if err != nil {
			return []*types.ProtectionDomain{}, errors.New("Error: problem finding link")
		}

		endpoint.Path = link.HREF
	} else {
		endpoint.Path = protectiondomainhref
	}

	req := system.client.NewRequest(map[string]string{}, "GET", endpoint, nil)
	req.SetBasicAuth("", system.client.Token)
	req.Header.Add("Accept", "application/json;version=1.0")

	resp, err := checkResp(system.client.Http.Do(req))
	if err != nil {
		return []*types.ProtectionDomain{}, fmt.Errorf("problem getting response: %v", err)
	}
	defer resp.Body.Close()

	if protectiondomainhref == "" {
		if err = decodeBody(resp, &protectionDomains); err != nil {
			return []*types.ProtectionDomain{}, fmt.Errorf("error decoding instances response: %s", err)
		}
	} else {
		protectionDomain := &types.ProtectionDomain{}
		if err = decodeBody(resp, &protectionDomain); err != nil {
			return []*types.ProtectionDomain{}, fmt.Errorf("error decoding instances response: %s", err)
		}
		protectionDomains = append(protectionDomains, protectionDomain)

	}
	//
	// bs, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return []types.ProtectionDomain{}, errors.New("error reading body")
	// }
	//
	// fmt.Println(string(bs))
	// log.Fatalf("here")
	// return []types.ProtectionDomain{}, nil
	return protectionDomains, nil
}

func (system *System) FindProtectionDomain(id, name, href string) (protectionDomain *types.ProtectionDomain, err error) {
	protectionDomains, err := system.GetProtectionDomain(href)
	if err != nil {
		return &types.ProtectionDomain{}, errors.New("Error getting protection domains")
	}

	for _, protectionDomain = range protectionDomains {
		if protectionDomain.ID == id || protectionDomain.Name == name || href != "" {
			return protectionDomain, nil
		}
	}

	return &types.ProtectionDomain{}, errors.New("Couldn't find protection domain")

}
