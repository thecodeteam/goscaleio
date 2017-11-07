package goscaleio

import (
	"errors"
	"fmt"
	"net/http"

	types "github.com/thecodeteam/goscaleio/types/v1"
)

type StoragePool struct {
	StoragePool *types.StoragePool
	client      *Client
}

func NewStoragePool(client *Client) *StoragePool {
	return &StoragePool{
		StoragePool: &types.StoragePool{},
		client:      client,
	}
}

func NewStoragePoolEx(client *Client, pool *types.StoragePool) *StoragePool {
	return &StoragePool{
		StoragePool: pool,
		client:      client,
	}
}

func (pd *ProtectionDomain) CreateStoragePool(name string) (string, error) {

	storagePoolParam := &types.StoragePoolParam{
		Name:               name,
		ProtectionDomainID: pd.ProtectionDomain.ID,
	}

	path := fmt.Sprintf("/api/types/StoragePool/instances")

	sp := types.StoragePoolResp{}
	err := pd.client.getJSONWithRetry(
		http.MethodPost, path, storagePoolParam, &sp)
	if err != nil {
		return "", err
	}

	return sp.ID, nil
}

func (pd *ProtectionDomain) GetStoragePool(
	storagepoolhref string) ([]*types.StoragePool, error) {

	var (
		err error
		sp  = &types.StoragePool{}
		sps []*types.StoragePool
	)

	if storagepoolhref == "" {
		var link *types.Link
		link, err := GetLink(pd.ProtectionDomain.Links,
			"/api/ProtectionDomain/relationship/StoragePool")
		if err != nil {
			return nil, err
		}
		err = pd.client.getJSONWithRetry(
			http.MethodGet, link.HREF, nil, &sps)
	} else {
		err = pd.client.getJSONWithRetry(
			http.MethodGet, storagepoolhref, nil, sp)
	}
	if err != nil {
		return nil, err
	}

	if storagepoolhref != "" {
		sps = append(sps, sp)
	}
	return sps, nil
}

func (pd *ProtectionDomain) FindStoragePool(
	id, name, href string) (*types.StoragePool, error) {

	sps, err := pd.GetStoragePool(href)
	if err != nil {
		return nil, fmt.Errorf("Error getting protection domains %s", err)
	}

	for _, sp := range sps {
		if sp.ID == id || sp.Name == name || href != "" {
			return sp, nil
		}
	}

	return nil, errors.New("Couldn't find storage pool")

}

func (sp *StoragePool) GetStatistics() (*types.Statistics, error) {

	link, err := GetLink(sp.StoragePool.Links,
		"/api/StoragePool/relationship/Statistics")
	if err != nil {
		return nil, err
	}

	stats := types.Statistics{}
	err = sp.client.getJSONWithRetry(
		http.MethodGet, link.HREF, nil, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
