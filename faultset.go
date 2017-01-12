package goscaleio

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"

	types "github.com/codedellemc/goscaleio/types/v1"
)

type FaultSet struct {
	FaultSet *types.FaultSet
	client   *Client
}

func NewFaultSet(client *Client) *FaultSet {
	return &FaultSet{
		FaultSet: new(types.FaultSet),
		client:   client,
	}
}

func NewFaultSetEx(client *Client, fs *types.FaultSet) *FaultSet {
	return &FaultSet{
		FaultSet: fs,
		client:   client,
	}
}

func (protectionDomain *ProtectionDomain) CreateFaultSet(name string) (string, error) {
	endpoint := protectionDomain.client.SIOEndpoint

	faultSetParam := &types.FaultSetParam{}
	faultSetParam.Name = name
	faultSetParam.ProtectionDomainID = protectionDomain.ProtectionDomain.ID

	jsonOutput, err := json.Marshal(&faultSetParam)
	if err != nil {
		return "", fmt.Errorf("error marshaling: %s", err)
	}
	endpoint.Path = fmt.Sprintf("/api/types/FaultSet/instances")

	req := protectionDomain.client.NewRequest(map[string]string{}, "POST", endpoint, bytes.NewBufferString(string(jsonOutput)))
	req.SetBasicAuth("", protectionDomain.client.Token)
	req.Header.Add("Accept", "application/json;version="+protectionDomain.client.configConnect.Version)
	req.Header.Add("Content-Type", "application/json;version="+protectionDomain.client.configConnect.Version)

	resp, err := protectionDomain.client.retryCheckResp(&protectionDomain.client.Http, req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("error reading body")
	}

	var sds types.FaultSetResp
	err = json.Unmarshal(bs, &sds)
	if err != nil {
		return "", err
	}

	return sds.ID, nil
}

func (protectionDomain *ProtectionDomain) GetFaultSets() (faultSets []types.FaultSet, err error) {
	endpoint := protectionDomain.client.SIOEndpoint
	endpoint.Path = fmt.Sprintf("/api/instances/ProtectionDomain::%v/relationships/FaultSet", protectionDomain.ProtectionDomain.ID)

	req := protectionDomain.client.NewRequest(map[string]string{}, "GET", endpoint, nil)
	req.SetBasicAuth("", protectionDomain.client.Token)
	req.Header.Add("Accept", "application/json;version="+protectionDomain.client.configConnect.Version)

	resp, err := protectionDomain.client.retryCheckResp(&protectionDomain.client.Http, req)
	if err != nil {
		return []types.FaultSet{}, fmt.Errorf("problem getting response: %v", err)
	}
	defer resp.Body.Close()

	if err = protectionDomain.client.decodeBody(resp, &faultSets); err != nil {
		return []types.FaultSet{}, fmt.Errorf("error decoding instances response: %s", err)
	}

	return faultSets, nil
}

func (protectionDomain *ProtectionDomain) FindFaultSet(field, value string) (faultSet *types.FaultSet, err error) {
	faultSets, err := protectionDomain.GetFaultSets()
	if err != nil {
		return &types.FaultSet{}, nil
	}

	for _, faultSet := range faultSets {
		valueOf := reflect.ValueOf(faultSet)
		switch {
		case reflect.Indirect(valueOf).FieldByName(field).String() == value:
			return &faultSet, nil
		}
	}

	return &types.FaultSet{}, errors.New("Couldn't find FaultSets")
}
