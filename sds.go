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

type Sds struct {
	Sds    *types.Sds
	client *Client
}

func NewSds(client *Client) *Sds {
	return &Sds{
		Sds:    new(types.Sds),
		client: client,
	}
}

func NewSdsEx(client *Client, sds *types.Sds) *Sds {
	return &Sds{
		Sds:    sds,
		client: client,
	}
}

func (protectionDomain *ProtectionDomain) CreateSds(name string, ipList []string, ipMode []string, faultSetId string) (string, error) {
	if len(ipList) == 0 {
		return "", fmt.Errorf("Must provide at least 1 SDS IP")
	}
	if len(ipList) != len(ipMode) {
		return "", fmt.Errorf("The length of the ipList must match the length of the ipMode.")
	}

	endpoint := protectionDomain.client.SIOEndpoint

	sdsParam := &types.SdsParam{}
	sdsParam.Name = name
	sdsParam.ProtectionDomainID = protectionDomain.ProtectionDomain.ID
	if len(faultSetId) > 0 {
		sdsParam.FaultSetID = faultSetId
	}

	for i := 0; i < len(ipList); i++ {
		sdsIP := types.SdsIp{IP: ipList[0], Role: ipMode[0]}
		sdsIPList := &types.SdsIpList{sdsIP}
		sdsParam.IPList = append(sdsParam.IPList, sdsIPList)
	}

	jsonOutput, err := json.Marshal(&sdsParam)
	if err != nil {
		return "", fmt.Errorf("error marshaling: %s", err)
	}
	endpoint.Path = fmt.Sprintf("/api/types/Sds/instances")

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

	var sds types.SdsResp
	err = json.Unmarshal(bs, &sds)
	if err != nil {
		return "", err
	}

	return sds.ID, nil
}

func (protectionDomain *ProtectionDomain) GetSds() (sdss []types.Sds, err error) {
	endpoint := protectionDomain.client.SIOEndpoint
	endpoint.Path = fmt.Sprintf("/api/instances/ProtectionDomain::%v/relationships/Sds", protectionDomain.ProtectionDomain.ID)

	req := protectionDomain.client.NewRequest(map[string]string{}, "GET", endpoint, nil)
	req.SetBasicAuth("", protectionDomain.client.Token)
	req.Header.Add("Accept", "application/json;version="+protectionDomain.client.configConnect.Version)

	resp, err := protectionDomain.client.retryCheckResp(&protectionDomain.client.Http, req)
	if err != nil {
		return []types.Sds{}, fmt.Errorf("problem getting response: %v", err)
	}
	defer resp.Body.Close()

	if err = protectionDomain.client.decodeBody(resp, &sdss); err != nil {
		return []types.Sds{}, fmt.Errorf("error decoding instances response: %s", err)
	}

	return sdss, nil
}

func (protectionDomain *ProtectionDomain) FindSds(field, value string) (sds *types.Sds, err error) {
	sdss, err := protectionDomain.GetSds()
	if err != nil {
		return &types.Sds{}, nil
	}

	for _, sds := range sdss {
		valueOf := reflect.ValueOf(sds)
		switch {
		case reflect.Indirect(valueOf).FieldByName(field).String() == value:
			return &sds, nil
		}
	}

	return &types.Sds{}, errors.New("Couldn't find SDS")
}
