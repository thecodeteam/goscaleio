package goscaleio

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	types "github.com/thecodeteam/goscaleio/types/v1"
)

type Device struct {
	Device *types.Device
	client *Client
}

func NewDevice(client *Client) *Device {
	return &Device{
		Device: &types.Device{},
		client: client,
	}
}

func NewDeviceEx(client *Client, device *types.Device) *Device {
	return &Device{
		Device: device,
		client: client,
	}
}

func (sp *StoragePool) AttachDevice(
	path string,
	sdsID string) (string, error) {

	deviceParam := &types.DeviceParam{
		Name: path,
		DeviceCurrentPathname: path,
		StoragePoolID:         sp.StoragePool.ID,
		SdsID:                 sdsID,
		TestMode:              "testAndActivate"}

	dev := types.DeviceResp{}
	err := sp.client.getJSONWithRetry(
		http.MethodPost, "/api/types/Device/instances",
		deviceParam, &dev)
	if err != nil {
		return "", err
	}

	return dev.ID, nil
}

func (sp *StoragePool) GetDevice() ([]types.Device, error) {

	path := fmt.Sprintf(
		"/api/instances/StoragePool::%v/relationships/Device",
		sp.StoragePool.ID)

	var devices []types.Device
	err := sp.client.getJSONWithRetry(
		http.MethodGet, path, nil, &devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func (sp *StoragePool) FindDevice(
	field, value string) (*types.Device, error) {

	devices, err := sp.GetDevice()
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		valueOf := reflect.ValueOf(device)
		switch {
		case reflect.Indirect(valueOf).FieldByName(field).String() == value:
			return &device, nil
		}
	}

	return nil, errors.New("Couldn't find DEV")
}
