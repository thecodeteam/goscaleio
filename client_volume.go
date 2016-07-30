package goscaleio

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	types "github.com/emccode/goscaleio/types/v1"
)

// RemoveMode exposes Volume removal modes
var RemoveMode = struct {
	OnlyMe             string
	IncludeDescendants string
	DescendantsOnly    string
	WholeTree          string
}{
	OnlyMe:             "ONLY_ME",
	IncludeDescendants: "INCLUDE_DESCENDANTS",
	DescendantsOnly:    "DESCENDANTS_ONLY",
	WholeTree:          "WHOLE_VTREE",
}

//GetVolumeID returns the ID for volume with name volumename.
//API endpoint: /api/types/Volume/instances/action/queryIdByKey.
func (c *Client) GetVolumeID(volumename string) (volumeID string, err error) {
	endpoint := c.SIOEndpoint
	endpoint.Path = "/api/types/Volume/instances/action/queryIdByKey"

	volumeQeryIdByKeyParam := &types.VolumeQeryIdByKeyParam{Name: volumename}
	jsonOutput, err := json.Marshal(&volumeQeryIdByKeyParam)
	if err != nil {
		log.WithError(err).Error("failed to marshal request parameter")
		return "", err
	}

	req, err := http.NewRequest("POST", endpoint.String(), bytes.NewReader(jsonOutput))
	if err != nil {
		log.WithError(err).Error("Failed to create request for", endpoint)
		return "", err
	}
	req.SetBasicAuth("", c.Token)
	resp, err := c.send(req)
	if err = c.validateResponse(resp); err != nil {
		log.WithError(err).Error("failed to get volume id")
		return "", err
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("error reading body")
		return "", err
	}

	volumeID = string(bs)
	volumeID = strings.TrimRight(volumeID, `"`)
	volumeID = strings.TrimLeft(volumeID, `"`)
	log.WithField("volumeID", volumeID).Debug("received volumeID")

	return volumeID, nil
}

//GetVolumeByID returns a volume instance retrieved by id.
//API endpoint: /api/instances/Volume::{is}
func (c *Client) GetVolumeByID(id string) (*types.Volume, error) {
	endpoint := c.SIOEndpoint
	endpoint.Path = fmt.Sprintf("/api/instances/Volume::%s", id)

	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.WithError(err).Error("Failed to create request for", endpoint)
		return nil, err
	}
	req.SetBasicAuth("", c.Token)

	resp, err := c.send(req)
	if err = c.validateResponse(resp); err != nil {
		log.WithError(err).Error("failed to get volume with id", id)
		return nil, err
	}
	defer resp.Body.Close()

	if err != nil {
		log.WithError(err).Error("error reading body")
		return nil, err
	}

	var volume *types.Volume
	if err = c.decodeBody(resp, &volume); err != nil {
		log.WithError(err).Error("failed to decode body")
		return nil, err
	}
	log.WithField("volueID", volume.ID).WithField("name", volume.Name).Debug("received volume")
	return volume, nil
}

// GetVolumes returns all volume instances for the scaleio cluster.
// API endpoint - /api/types/Volume/instances or /api/instances/Volume::id
func (c *Client) GetVolumes() ([]*types.Volume, error) {
	endpoint := c.SIOEndpoint
	endpoint.Path = "/api/types/Volume/instances"

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

	if err != nil {
		log.WithError(err).Error("error reading body")
		return nil, err
	}

	var volumes []*types.Volume
	if err = c.decodeBody(resp, &volumes); err != nil {
		log.WithError(err).Error("failed to decode body")
		return nil, err
	}

	log.WithField("count", len(volumes)).Debug("received volumes")
	return volumes, nil
}

// GetVolumesByName returns a list of Volumes filtered by name
// See GetVolumes
func (c *Client) GetVolumesByName(name string) (*types.Volume, error) {
	volumes, err := c.GetVolumes()
	if err != nil {
		return nil, err
	}

	for _, vol := range volumes {
		if vol.Name == name {
			return vol, nil
		}
	}
	log.WithField("name", name).Error("volume not found")
	return nil, errors.New("volume not found")
}

// GetVolumesByStoragePoolID returns all volumes for specified pool storage id.
// API endpoint - /api/instances/StoragePool::id/relationships/Volume
func (c *Client) GetVolumesByStoragePoolID(id string) ([]*types.Volume, error) {
	endpoint := c.SIOEndpoint
	endpoint.Path = fmt.Sprintf("/api/instances/StoragePool::%s/relationships/Volume", id)

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

	if err != nil {
		log.WithError(err).Error("error reading body")
		return nil, err
	}

	var volumes []*types.Volume
	if err = c.decodeBody(resp, &volumes); err != nil {
		log.WithError(err).Error("failed to decode body")
		return nil, err
	}

	log.WithField("count", len(volumes)).Debug("received storage pool volumes")
	return volumes, nil
}

//CreateVolume creates a new volume and returns the id for the created volume
//API endpoint - /api/types/Volume/instances
func (c *Client) CreateVolume(volume *types.VolumeParam) (string, error) {
	endpoint := c.SIOEndpoint
	endpoint.Path = "/api/types/Volume/instances"

	jsonOutput, err := json.Marshal(&volume)
	if err != nil {
		log.WithError(err).Error("failed to marshal VolumeParam")
		return "", err
	}

	req, err := http.NewRequest("POST", endpoint.String(), bytes.NewReader(jsonOutput))
	if err != nil {
		log.WithError(err).Error("Failed to create request for", endpoint)
		return "", err
	}
	req.SetBasicAuth("", c.Token)
	resp, err := c.send(req)
	if err = c.validateResponse(resp); err != nil {
		log.WithError(err).Error("failed to create volume")
		return "", err
	}
	defer resp.Body.Close()

	var volumeResp *types.VolumeResp
	if err = c.decodeBody(resp, &volumeResp); err != nil {
		return "", err
	}

	return volumeResp.ID, nil
}

// RemoveVolume removes the volume with specified id using mode.
// API endpoint - /api/instances/Volume::<id>/action/removeVolume
func (c *Client) RemoveVolume(id, mode string) error {
	endpoint := c.SIOEndpoint
	endpoint.Path = fmt.Sprintf("/api/instances/Volume::%s/action/removeVolume", id)

	param := &types.RemoveVolumeParam{
		RemoveMode: mode,
	}

	jsonOutput, err := json.Marshal(param)
	if err != nil {
		log.WithError(err).Error("failed to marshal RemoveVolumeParam")
		return err
	}

	req, err := http.NewRequest("POST", endpoint.String(), bytes.NewReader(jsonOutput))
	if err != nil {
		log.WithError(err).Error("Failed to create request for", endpoint)
		return err
	}
	req.SetBasicAuth("", c.Token)
	resp, err := c.send(req)
	if err = c.validateResponse(resp); err != nil {
		log.WithError(err).Error("failed to remove volue")
		return err
	}
	return nil
}

//AddMappedSdc maps an SDC to specified volume
//API endpoint - /api/instances/Volume::<id>/action/addMappedSdc
func (c *Client) AddMappedSdc(volumeID string, param *types.MapVolumeSdcParam) error {
	endpoint := c.SIOEndpoint
	endpoint.Path = fmt.Sprintf("/api/instances/Volume::%s/action/addMappedSdc", volumeID)

	jsonOutput, err := json.Marshal(param)
	if err != nil {
		log.WithError(err).Error("failed to marshal MapVolumeSdcParam")
		return err
	}

	req, err := http.NewRequest("POST", endpoint.String(), bytes.NewReader(jsonOutput))
	if err != nil {
		log.WithError(err).Error("Failed to create request for", endpoint)
		return err
	}
	req.SetBasicAuth("", c.Token)
	resp, err := c.send(req)
	if err = c.validateResponse(resp); err != nil {
		log.WithError(err).Error("failed to add mapped sdc")
		return err
	}
	return nil
}

//RemoveMappedSdc removes the sdc mapping for specified volume
//API endpoint - /api/instances/Volume::<id>/action/removeMappedSdc
func (c *Client) RemoveMappedSdc(volumeID string, param *types.UnmapVolumeSdcParam) error {
	endpoint := c.SIOEndpoint
	endpoint.Path = fmt.Sprintf("/api/instances/Volume::%s/action/removeMappedSdc", volumeID)

	jsonOutput, err := json.Marshal(param)
	if err != nil {
		log.WithError(err).Error("failed to marshal UnmapVolumeSdcParam")
		return err
	}

	req, err := http.NewRequest("POST", endpoint.String(), bytes.NewReader(jsonOutput))
	if err != nil {
		log.WithError(err).Error("Failed to create request for", endpoint)
		return err
	}
	req.SetBasicAuth("", c.Token)
	resp, err := c.send(req)
	if err = c.validateResponse(resp); err != nil {
		log.WithError(err).Error("failed to add mapped sdc")
		return err
	}
	return nil
}
