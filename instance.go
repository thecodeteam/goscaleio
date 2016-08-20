package goscaleio

import (
	"errors"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	types "github.com/emccode/goscaleio/types/v1"
)

func (c *Client) GetInstance(systemhref string) ([]*types.System, error) {
	endpoint := c.SIOEndpoint
	if systemhref == "" {
		endpoint.Path += "/types/System/instances"
	} else {
		endpoint.Path = systemhref
	}

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request for", endpoint)
		return nil, err
	}
	req.SetBasicAuth("", c.Token)
	resp, err := c.send(req)
	if err = c.validateResponse(resp); err != nil {
		log.WithError(err).Error("failed to get instances")
		return nil, err
	}
	defer resp.Body.Close()

	var systems []*types.System
	if systemhref == "" {
		if err = c.decodeBody(resp, &systems); err != nil {
			return nil, fmt.Errorf("error decoding instances response: %s", err)
		}
	} else {
		system := &types.System{} // get info agbout specified system
		if err = c.decodeBody(resp, &system); err != nil {
			return nil, fmt.Errorf("error decoding instances response: %s", err)
		}
		systems = append(systems, system)
	}
	if log.GetLevel() == log.DebugLevel {
		log.WithField("instances", fmt.Sprintf("%#v", systems)).Debug("retrieved system instances")
	}
	return systems, nil
}

func (c *Client) FindVolumeID(volumename string) (volumeID string, err error) {
	return c.GetVolumeID(volumename)
}

func (c *Client) GetVolume(volumehref, volumeid, ancestorvolumeid, volumename string, getSnapshots bool) (volumes []*types.Volume, err error) {

	endpoint := c.SIOEndpoint

	if volumename != "" {
		volumeid, err = c.FindVolumeID(volumename)
		if err != nil && err.Error() == "Not found" {
			return nil, err
		}
		if err != nil {
			return nil, fmt.Errorf("Error: problem finding volume: %s", err)
		}
	}

	if volumeid != "" {
		endpoint.Path = fmt.Sprintf("/api/instances/Volume::%s", volumeid)
	} else if volumehref == "" {
		endpoint.Path = "/api/types/Volume/instances"
	} else {
		endpoint.Path = volumehref
	}

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request for", endpoint)
		return nil, err
	}
	req.SetBasicAuth("", c.Token)

	resp, err := c.send(req)
	if err = c.validateResponse(resp); err != nil {
		log.WithError(err).Error("failed to get volumes")
		return nil, err
	}
	defer resp.Body.Close()

	if volumehref == "" && volumeid == "" {
		if err = c.decodeBody(resp, &volumes); err != nil {
			return nil, fmt.Errorf("error decoding storage pool response: %s", err)
		}
		var volumesNew []*types.Volume
		for _, volume := range volumes {
			if (!getSnapshots && volume.AncestorVolumeID == ancestorvolumeid) || (getSnapshots && volume.AncestorVolumeID != "") {
				volumesNew = append(volumesNew, volume)
			}
		}
		volumes = volumesNew
	} else {
		volume := &types.Volume{}
		if err = c.decodeBody(resp, &volume); err != nil {
			return nil, fmt.Errorf("error decoding instances response: %s", err)
		}
		volumes = append(volumes, volume)
	}
	return volumes, nil
}

func (c *Client) GetStoragePool(storagepoolhref string) (storagePools []*types.StoragePool, err error) {

	endpoint := c.SIOEndpoint

	if storagepoolhref == "" {
		endpoint.Path = "/api/types/StoragePool/instances"
	} else {
		endpoint.Path = storagepoolhref
	}

	req := c.NewRequest(map[string]string{}, "GET", endpoint, nil)
	req.SetBasicAuth("", c.Token)
	req.Header.Add("Accept", "application/json;version="+c.configConnect.Version)

	resp, err := c.retryCheckResp(&c.Http, req)
	if err != nil {
		return []*types.StoragePool{}, fmt.Errorf("problem getting response: %v", err)
	}
	defer resp.Body.Close()

	if storagepoolhref == "" {
		if err = c.decodeBody(resp, &storagePools); err != nil {
			return []*types.StoragePool{}, fmt.Errorf("error decoding storage pool response: %s", err)
		}
	} else {
		storagePool := &types.StoragePool{}
		if err = c.decodeBody(resp, &storagePool); err != nil {
			return []*types.StoragePool{}, fmt.Errorf("error decoding instances response: %s", err)
		}
		storagePools = append(storagePools, storagePool)
	}
	return storagePools, nil
}

func (c *Client) FindStoragePool(id, name, href string) (storagePool *types.StoragePool, err error) {
	storagePools, err := c.GetStoragePool(href)
	if err != nil {
		return &types.StoragePool{}, fmt.Errorf("Error getting storage pool %s", err)
	}

	for _, storagePool = range storagePools {
		if storagePool.ID == id || storagePool.Name == name || href != "" {
			return storagePool, nil
		}
	}

	return &types.StoragePool{}, errors.New("Couldn't find storage pool")

}
