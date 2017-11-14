package goscaleio

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	types "github.com/thecodeteam/goscaleio/types/v1"
)

type Sds struct {
	Sds    *types.Sds
	client *Client
}

func NewSds(client *Client) *Sds {
	return &Sds{
		Sds:    &types.Sds{},
		client: client,
	}
}

func NewSdsEx(client *Client, sds *types.Sds) *Sds {
	return &Sds{
		Sds:    sds,
		client: client,
	}
}

func (pd *ProtectionDomain) CreateSds(
	name string, ipList []string) (string, error) {

	sdsParam := &types.SdsParam{
		Name:               name,
		ProtectionDomainID: pd.ProtectionDomain.ID,
	}

	if len(ipList) == 0 {
		return "", fmt.Errorf("Must provide at least 1 SDS IP")
	} else if len(ipList) == 1 {
		sdsIP := types.SdsIp{IP: ipList[0], Role: "all"}
		sdsIPList := &types.SdsIpList{SdsIP: sdsIP}
		sdsParam.IPList = append(sdsParam.IPList, sdsIPList)
	} else if len(ipList) >= 2 {
		sdsIP1 := types.SdsIp{IP: ipList[0], Role: "sdcOnly"}
		sdsIP2 := types.SdsIp{IP: ipList[1], Role: "sdsOnly"}
		sdsIPList1 := &types.SdsIpList{SdsIP: sdsIP1}
		sdsIPList2 := &types.SdsIpList{SdsIP: sdsIP2}
		sdsParam.IPList = append(sdsParam.IPList, sdsIPList1)
		sdsParam.IPList = append(sdsParam.IPList, sdsIPList2)
	}

	path := fmt.Sprintf("/api/types/Sds/instances")

	sds := types.SdsResp{}
	err := pd.client.getJSONWithRetry(
		http.MethodPost, path, sdsParam, &sds)
	if err != nil {
		return "", err
	}

	return sds.ID, nil
}

func (pd *ProtectionDomain) GetSds() ([]types.Sds, error) {

	path := fmt.Sprintf("/api/instances/ProtectionDomain::%v/relationships/Sds",
		pd.ProtectionDomain.ID)

	var sdss []types.Sds
	err := pd.client.getJSONWithRetry(
		http.MethodGet, path, nil, &sdss)
	if err != nil {
		return nil, err
	}

	return sdss, nil
}

func (pd *ProtectionDomain) FindSds(
	field, value string) (*types.Sds, error) {

	sdss, err := pd.GetSds()
	if err != nil {
		return nil, err
	}

	for _, sds := range sdss {
		valueOf := reflect.ValueOf(sds)
		switch {
		case reflect.Indirect(valueOf).FieldByName(field).String() == value:
			return &sds, nil
		}
	}

	return nil, errors.New("Couldn't find SDS")
}
