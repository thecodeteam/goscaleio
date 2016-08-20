package goscaleio

import (
	"errors"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	types "github.com/emccode/goscaleio/types/v1"
)

// GetSystemByID returns a system instance based on provided id.
// API endpoint - /instances/System::{id}
func (c *Client) GetSystemByID(id string) (*types.System, error) {
	endpoint := c.SIOEndpoint
	endpoint.Path += fmt.Sprintf("/instances/System::%s", id)

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request for", endpoint)
		return nil, err
	}
	req.SetBasicAuth("", c.Token)
	resp, err := c.send(req)
	if err = c.validateResponse(resp); err != nil {
		log.WithError(err).Error("failed to get instance")
		return nil, err
	}
	defer resp.Body.Close()

	var system = new(types.System)
	if err = c.decodeBody(resp, system); err != nil {
		log.WithError(err).Error("json decoding failed")
		return nil, err
	}
	return system, nil
}

//GetSystems retrieved a list of available systems filtered by nameFilter
// API endpoint - /types/System/instances
func (c *Client) GetSystems() ([]*types.System, error) {
	endpoint := c.SIOEndpoint
	endpoint.Path += "/types/System/instances"

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request for", endpoint)
		return nil, err
	}
	req.SetBasicAuth("", c.Token)
	resp, err := c.send(req)
	if err = c.validateResponse(resp); err != nil {
		log.WithError(err).Error("failed to get instance")
		return nil, err
	}
	defer resp.Body.Close()

	var systems []*types.System
	if err = c.decodeBody(resp, &systems); err != nil {
		log.WithError(err).Error("json decoding failed")
		return nil, err
	}

	log.WithField("count", len(systems)).Debug("retrieved system instances")

	return systems, nil
}

//GetSystemByName returns list of filtered systems by name
//See GetSystems
func (c *Client) GetSystemByName(name string) (*types.System, error) {
	systems, err := c.GetSystems()
	if err != nil {
		return nil, err
	}

	for _, sys := range systems {
		if sys.Name == name {
			return sys, nil
		}
	}
	log.WithField("name", name).Error("system not found")
	return nil, errors.New("system not found")
}

// GetProtectionDomainByID returns the a protection domain based on ID
// API endpoint - /api/instances/ProtectionDomain::<id>
func (c *Client) GetProtectionDomainByID(id string) (*types.ProtectionDomain, error) {
	endpoint := c.SIOEndpoint
	endpoint.Path += fmt.Sprintf("/instances/ProtectionDomain::%s", id)

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request for", endpoint)
		return nil, err
	}
	req.SetBasicAuth("", c.Token)
	resp, err := c.send(req)
	if err = c.validateResponse(resp); err != nil {
		log.WithError(err).Error("failed to get protection domain")
		return nil, err
	}
	defer resp.Body.Close()

	var pd = new(types.ProtectionDomain)
	if err = c.decodeBody(resp, pd); err != nil {
		log.WithError(err).Error("json decoding failed")
		return nil, err
	}
	return pd, nil
}

//GetProtectionDomains retrieved a list of protection domains
// API endpoint - /api/types/ProtectionDomain/instances
func (c *Client) GetProtectionDomains() ([]*types.ProtectionDomain, error) {
	endpoint := c.SIOEndpoint
	endpoint.Path += "/types/ProtectionDomain/instances"

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request for", endpoint)
		return nil, err
	}
	req.SetBasicAuth("", c.Token)
	resp, err := c.send(req)
	if err = c.validateResponse(resp); err != nil {
		log.WithError(err).Error("failed to get protection domains")
		return nil, err
	}
	defer resp.Body.Close()

	var domains []*types.ProtectionDomain
	if err = c.decodeBody(resp, &domains); err != nil {
		log.WithError(err).Error("json decoding failed")
		return nil, err
	}

	log.WithField("count", len(domains)).Debug("retrieved protection domains")

	return domains, nil
}

//GetProtectionDomainByName returns a protection with name
func (c *Client) GetProtectionDomainByName(name string) (*types.ProtectionDomain, error) {
	domains, err := c.GetProtectionDomains()
	if err != nil {
		return nil, err
	}
	for _, domain := range domains {
		if domain.Name == name {
			return domain, nil
		}
	}
	log.WithField("name", name).Error("protection domain not found")
	return nil, errors.New("protection domain not found")
}

// GetSdcByID returns the a protection domain based on ID
// API endpoint - /api/instances/Sdc::<id>
func (c *Client) GetSdcByID(id string) (*types.Sdc, error) {
	endpoint := c.SIOEndpoint
	endpoint.Path += fmt.Sprintf("/instances/Sdc::%s", id)

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request for", endpoint)
		return nil, err
	}
	req.SetBasicAuth("", c.Token)
	resp, err := c.send(req)
	if err = c.validateResponse(resp); err != nil {
		log.WithError(err).Error("failed to get Sdc")
		return nil, err
	}
	defer resp.Body.Close()

	var sdc = new(types.Sdc)
	if err = c.decodeBody(resp, sdc); err != nil {
		log.WithError(err).Error("json decoding failed")
		return nil, err
	}
	return sdc, nil
}

// GetSdcs returns the SDCs for storage system
// API endpoint - /api/types/Sdc/instances
func (c *Client) GetSdcs() ([]*types.Sdc, error) {
	endpoint := c.SIOEndpoint
	endpoint.Path += "/types/Sdc/instances"

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request for", endpoint)
		return nil, err
	}
	req.SetBasicAuth("", c.Token)
	resp, err := c.send(req)
	if err = c.validateResponse(resp); err != nil {
		log.WithError(err).Error("failed to get SDCs")
		return nil, err
	}
	defer resp.Body.Close()

	var sdcs []*types.Sdc
	if err = c.decodeBody(resp, &sdcs); err != nil {
		log.WithError(err).Error("json decoding failed")
		return nil, err
	}

	log.WithField("count", len(sdcs)).Debug("retrieved SDCs")

	return sdcs, nil
}

//GetSdcByGUID returns a protection with name
func (c *Client) GetSdcByGUID(guid string) (*types.Sdc, error) {
	sdcs, err := c.GetSdcs()
	if err != nil {
		return nil, err
	}
	for _, sdc := range sdcs {
		if sdc.SdcGuid == guid {
			return sdc, nil
		}
	}
	log.WithField("guid", guid).Error("sdc not found")
	return nil, errors.New("sdc not found")
}
