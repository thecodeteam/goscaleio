package goscaleio

import (
	"errors"
	"fmt"
	"net/http"

	types "github.com/thecodeteam/goscaleio/types/v1"
)

func (c *Client) GetInstance(systemhref string) ([]*types.System, error) {

	var (
		err     error
		system  = &types.System{}
		systems []*types.System
	)

	if systemhref == "" {
		err = c.getJSONWithRetry(
			http.MethodGet, "api/types/System/instances", nil, &systems)
	} else {
		err = c.getJSONWithRetry(
			http.MethodGet, systemhref, nil, system)
	}
	if err != nil {
		return nil, err
	}

	if systemhref != "" {
		systems = append(systems, system)
	}

	return systems, nil
}

func (c *Client) GetVolume(
	volumehref, volumeid, ancestorvolumeid, volumename string,
	getSnapshots bool) ([]*types.Volume, error) {

	var (
		err     error
		path    string
		volume  = &types.Volume{}
		volumes []*types.Volume
	)

	if volumename != "" {
		volumeid, err = c.FindVolumeID(volumename)
		if err != nil && err.Error() == "Not found" {
			return nil, nil
		}
		if err != nil {
			return nil, fmt.Errorf("Error: problem finding volume: %s", err)
		}
	}

	if volumeid != "" {
		path = fmt.Sprintf("/api/instances/Volume::%s", volumeid)
	} else if volumehref == "" {
		path = "/api/types/Volume/instances"
	} else {
		path = volumehref
	}

	if volumehref == "" && volumeid == "" {
		err = c.getJSONWithRetry(
			http.MethodGet, path, nil, &volumes)
	} else {
		err = c.getJSONWithRetry(
			http.MethodGet, path, nil, volume)

	}
	if err != nil {
		return nil, err
	}

	if volumehref == "" && volumeid == "" {
		var volumesNew []*types.Volume
		for _, volume := range volumes {
			if (!getSnapshots && volume.AncestorVolumeID == ancestorvolumeid) || (getSnapshots && volume.AncestorVolumeID != "") {
				volumesNew = append(volumesNew, volume)
			}
		}
		volumes = volumesNew
	} else {
		volumes = append(volumes, volume)
	}
	return volumes, nil
}

func (c *Client) FindVolumeID(volumename string) (string, error) {

	volumeQeryIdByKeyParam := &types.VolumeQeryIdByKeyParam{
		Name: volumename,
	}

	path := fmt.Sprintf("/api/types/Volume/instances/action/queryIdByKey")

	volumeID, err := c.getStringWithRetry(http.MethodPost, path,
		volumeQeryIdByKeyParam)
	if err != nil {
		return "", err
	}

	return volumeID, nil
}

func (c *Client) CreateVolume(
	volume *types.VolumeParam,
	storagePoolName string) (*types.VolumeResp, error) {

	path := "/api/types/Volume/instances"

	storagePool, err := c.FindStoragePool("", storagePoolName, "")
	if err != nil {
		return nil, err
	}

	volume.StoragePoolID = storagePool.ID
	volume.ProtectionDomainID = storagePool.ProtectionDomainID

	vol := &types.VolumeResp{}
	err = c.getJSONWithRetry(
		http.MethodPost, path, volume, vol)
	if err != nil {
		return nil, err
	}

	return vol, nil
}

func (c *Client) GetStoragePool(
	storagepoolhref string) ([]*types.StoragePool, error) {

	var (
		err          error
		storagePool  = &types.StoragePool{}
		storagePools []*types.StoragePool
	)

	if storagepoolhref == "" {
		err = c.getJSONWithRetry(
			http.MethodGet, "/api/types/StoragePool/instances",
			nil, &storagePools)
	} else {
		err = c.getJSONWithRetry(
			http.MethodGet, storagepoolhref, nil, storagePool)
	}
	if err != nil {
		return nil, err
	}

	if storagepoolhref != "" {
		storagePools = append(storagePools, storagePool)
	}
	return storagePools, nil
}

func (c *Client) FindStoragePool(
	id, name, href string) (*types.StoragePool, error) {

	storagePools, err := c.GetStoragePool(href)
	if err != nil {
		return nil, fmt.Errorf("Error getting storage pool %s", err)
	}

	for _, storagePool := range storagePools {
		if storagePool.ID == id || storagePool.Name == name || href != "" {
			return storagePool, nil
		}
	}

	return nil, errors.New("Couldn't find storage pool")
}
