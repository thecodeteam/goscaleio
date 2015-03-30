package goscaleio

import (
	"errors"
	"fmt"
	"os/exec"
	"reflect"
	"strings"

	types "github.com/emccode/goscaleio/types/v1"
)

type Sdc struct {
	Sdc    *types.Sdc
	client *Client
}

func NewSdc(client *Client, sdc *types.Sdc) *Sdc {
	return &Sdc{
		Sdc:    sdc,
		client: client,
	}
}

func (system *System) GetSdc() (sdcs []types.Sdc, err error) {
	endpoint := system.client.SIOEndpoint
	endpoint.Path = fmt.Sprintf("/api/instances/System::%v/relationships/Sdc", system.System.ID)

	req := system.client.NewRequest(map[string]string{}, "GET", endpoint, nil)
	req.SetBasicAuth("", system.client.Token)
	req.Header.Add("Accept", "application/json;version=1.0")

	resp, err := checkResp(system.client.Http.Do(req))
	if err != nil {
		return []types.Sdc{}, fmt.Errorf("problem getting response: %v", err)
	}
	defer resp.Body.Close()

	if err = decodeBody(resp, &sdcs); err != nil {
		return []types.Sdc{}, fmt.Errorf("error decoding instances response: %s", err)
	}

	// bs, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return []types.Sdc{}, errors.New("error reading body")
	// }
	//
	// fmt.Println(string(bs))
	// log.Fatalf("here")
	// return []types.Sdc{}, nil
	return sdcs, nil
}

func (system *System) FindSdc(field, value string) (sdc *Sdc, err error) {
	sdcs, err := system.GetSdc()
	if err != nil {
		return &Sdc{}, nil
	}

	for _, sdc := range sdcs {
		valueOf := reflect.ValueOf(sdc)
		switch {
		case reflect.Indirect(valueOf).FieldByName(field).String() == value:
			return NewSdc(system.client, &sdc), nil
		}
	}

	return &Sdc{}, errors.New("Couldn't find SDC")
}

func (sdc *Sdc) GetStatistics() (statistics types.Statistics, err error) {
	endpoint := sdc.client.SIOEndpoint
	endpoint.Path = fmt.Sprintf("/api/instances/Sdc::%v/relationships/Statistics", sdc.Sdc.ID)

	req := sdc.client.NewRequest(map[string]string{}, "GET", endpoint, nil)
	req.SetBasicAuth("", sdc.client.Token)
	req.Header.Add("Accept", "application/json;version=1.0")

	resp, err := checkResp(sdc.client.Http.Do(req))
	if err != nil {
		return types.Statistics{}, fmt.Errorf("problem getting response: %v", err)
	}
	defer resp.Body.Close()

	if err = decodeBody(resp, &statistics); err != nil {
		return types.Statistics{}, fmt.Errorf("error decoding instances response: %s", err)
	}

	// bs, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return errors.New("error reading body")
	// }
	//
	// fmt.Println(string(bs))
	return statistics, nil
}

func GetSdcLocalGUID() (sdcGUID string, err error) {

	// get sdc kernel guid
	// /bin/emc/scaleio/drv_cfg --query_guid
	// sdcKernelGuid := "271bad82-08ee-44f2-a2b1-7e2787c27be1"

	out, err := exec.Command("/bin/emc/scaleio/drv_cfg", "--query_guid").Output()
	if err != nil {
		return "", fmt.Errorf("Error querying volumes: ", err)
	}

	sdcGUID = strings.Replace(string(out), "\n", "", -1)

	return sdcGUID, nil

}
