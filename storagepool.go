package goscaleio

import (
	"fmt"

	types "github.com/emccode/goscaleio/types/v1"
)

func (protectionDomain *ProtectionDomain) GetStoragePool() (storagePools []*types.StoragePool, err error) {
	endpoint := protectionDomain.client.SIOEndpoint
	endpoint.Path = fmt.Sprintf("/api/instances/ProtectionDomain::%v/relationships/StoragePool", protectionDomain.ProtectionDomain.ID)

	req := protectionDomain.client.NewRequest(map[string]string{}, "GET", endpoint, nil)
	req.SetBasicAuth("", protectionDomain.client.Token)
	req.Header.Add("Accept", "application/json;version=1.0")

	resp, err := checkResp(protectionDomain.client.Http.Do(req))
	if err != nil {
		return []*types.StoragePool{}, fmt.Errorf("problem getting response: %v", err)
	}
	defer resp.Body.Close()

	if err = decodeBody(resp, &storagePools); err != nil {
		return []*types.StoragePool{}, fmt.Errorf("error decoding storage pool response: %s", err)
	}

	return storagePools, nil
}

// func (system *System) FindProtectionDomain(id, name string) (protectionDomain *types.ProtectionDomain, err error) {
// 	protectionDomains, err := system.GetProtectionDomain()
// 	if err != nil {
// 		return &types.ProtectionDomain{}, errors.New("Error getting protection domains")
// 	}
//
// 	for _, protectionDomain = range protectionDomains {
// 		if protectionDomain.ID == id || protectionDomain.Name == name {
// 			return protectionDomain, nil
// 		}
// 	}
//
// 	return &types.ProtectionDomain{}, errors.New("Couldn't find protection domain")
// }
//
