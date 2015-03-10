package goscaleio

import (
	"fmt"

	types "github.com/emccode/goscaleio/types/v1"
)

func (system *System) GetProtectionDomain() (protectionDomains []types.ProtectionDomain, err error) {
	endpoint := system.client.SIOEndpoint
	endpoint.Path = fmt.Sprintf("/api/instances/System::%v/relationships/ProtectionDomain", system.System.ID)

	req := system.client.NewRequest(map[string]string{}, "GET", endpoint, nil)
	req.SetBasicAuth("", system.client.Token)
	req.Header.Add("Accept", "application/json;version=1.0")

	resp, err := checkResp(system.client.Http.Do(req))
	if err != nil {
		return []types.ProtectionDomain{}, fmt.Errorf("problem getting response: %v", err)
	}
	defer resp.Body.Close()

	if err = decodeBody(resp, &protectionDomains); err != nil {
		return []types.ProtectionDomain{}, fmt.Errorf("error decoding instances response: %s", err)
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
