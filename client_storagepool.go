package goscaleio

import (
	"errors"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	types "github.com/emccode/goscaleio/types/v1"
)

// GetStoragePoolByID returns the a storage pool  based on ID
// API endpoint - /api/instances/StoragePool::<id>
func (c *Client) GetStoragePoolByID(id string) (*types.StoragePool, error) {
	endpoint := c.SIOEndpoint
	endpoint.Path += fmt.Sprintf("/instances/StoragePool::%s", id)

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request for", endpoint)
		return nil, err
	}
	req.SetBasicAuth("", c.Token)
	resp, err := c.send(req)
	if err = c.validateResponse(resp); err != nil {
		log.WithError(err).Error("failed to get storage pool")
		return nil, err
	}
	defer resp.Body.Close()

	var pool = new(types.StoragePool)
	if err = c.decodeBody(resp, pool); err != nil {
		log.WithError(err).Error("json decoding failed")
		return nil, err
	}
	return pool, nil
}

//GetStoragePools return a list of storage poools for entire cluster
//API endpoint - /api/types/StoragePool/instances
func (c *Client) GetStoragePools() ([]*types.StoragePool, error) {
	endpoint := c.SIOEndpoint
	endpoint.Path = "/api/types/StoragePool/instances"

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request for", endpoint)
		return nil, err
	}
	req.SetBasicAuth("", c.Token)

	resp, err := c.send(req)
	if err = c.validateResponse(resp); err != nil {
		log.WithError(err).Error("failed to get storage pools")
		return nil, err
	}
	defer resp.Body.Close()

	if err != nil {
		log.WithError(err).Error("error reading body")
		return nil, err
	}

	var pools []*types.StoragePool
	if err = c.decodeBody(resp, &pools); err != nil {
		log.WithError(err).Error("failed to decode storage pools")
		return nil, err
	}

	log.WithField("count", len(pools)).Debug("received storage pools")
	return pools, nil
}

//GetStoragePoolsByProtectionDomainID returns storage poools for a protection domain
//API endpoint - /api/instances/ProtectionDomain::{id}/relationships/StoragePool
func (c *Client) GetStoragePoolsByProtectionDomainID(id string) ([]*types.StoragePool, error) {
	endpoint := c.SIOEndpoint
	endpoint.Path = fmt.Sprintf(
		"/api/instances/ProtectionDomain::%s/relationships/StoragePool",
		id,
	)

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request for", endpoint)
		return nil, err
	}
	req.SetBasicAuth("", c.Token)

	resp, err := c.send(req)
	if err = c.validateResponse(resp); err != nil {
		log.WithError(err).Error("failed to get storage pools")
		return nil, err
	}
	defer resp.Body.Close()

	if err != nil {
		log.WithError(err).Error("error reading body")
		return nil, err
	}

	var pools []*types.StoragePool
	if err = c.decodeBody(resp, &pools); err != nil {
		log.WithError(err).Error("failed to decode storage pools")
		return nil, err
	}

	log.WithField("count", len(pools)).Debug("received storage pools")
	return pools, nil
}

//GetStoragePoolByName returns a storage pool matching specified name
func (c *Client) GetStoragePoolByName(name string) (*types.StoragePool, error) {
	pools, err := c.GetStoragePools()
	if err != nil {
		return nil, err
	}
	for _, pool := range pools {
		if pool.Name == name {
			return pool, nil
		}
	}
	log.WithField("name", name).Error("storage pool not found")
	return nil, errors.New("not found")
}
