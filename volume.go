package goscaleio

import (
	"errors"
	"fmt"

	types "github.com/emccode/goscaleio/types/v1"
)

type Volume struct {
	Volume *types.Volume
	client *Client
}

func NewVolume(client *Client) *Volume {
	return &Volume{
		Volume: new(types.Volume),
		client: client,
	}
}

func (storagePool *StoragePool) GetVolume(storagepoolhref string) (volumes []*types.Volume, err error) {

	endpoint := storagePool.client.SIOEndpoint

	if storagepoolhref == "" {
		link, err := GetLink(storagePool.StoragePool.Links, "/api/StoragePool/relationship/Volume")
		if err != nil {
			return []*types.Volume{}, errors.New("Error: problem finding link")
		}
		endpoint.Path = link.HREF
	} else {
		endpoint.Path = storagepoolhref
	}

	req := storagePool.client.NewRequest(map[string]string{}, "GET", endpoint, nil)
	req.SetBasicAuth("", storagePool.client.Token)
	req.Header.Add("Accept", "application/json;version=1.0")

	resp, err := checkResp(storagePool.client.Http.Do(req))
	if err != nil {
		return []*types.Volume{}, fmt.Errorf("problem getting response: %v", err)
	}
	defer resp.Body.Close()

	if storagepoolhref == "" {
		if err = decodeBody(resp, &volumes); err != nil {
			return []*types.Volume{}, fmt.Errorf("error decoding storage pool response: %s", err)
		}
	} else {
		storagePool := &types.Volume{}
		if err = decodeBody(resp, &storagePool); err != nil {
			return []*types.Volume{}, fmt.Errorf("error decoding instances response: %s", err)
		}
		volumes = append(volumes, storagePool)
	}
	return volumes, nil
}

func (storagePool *StoragePool) FindVolume(id, name, href string) (volume *types.Volume, err error) {
	// volumes, err := storagePool.GetVolume(href)
	// if err != nil {
	// 	return &types.Volume{}, errors.New("Error getting volumes")
	// }
	//
	// for _, volume = range volumes {
	// 	if volume.ID == id || volume.Name == name || href != "" {
	// 		return volume, nil
	// 	}
	// }
	//
	// return &types.Volume{}, errors.New("Couldn't find volumes")
	return &types.Volume{}, nil
}
