package goscaleio

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	types "github.com/emccode/goscaleio/types/v1"
)

type SdcMappedVolume struct {
	MdmID     string
	VolumeID  string
	SdcDevice string
}

type Volume struct {
	Volume *types.Volume
	client *Client
}

func NewVolume(client *Client) *Volume {
	return &Volume{
		Volume: new(types.Volume),
		client: client,
	}
}

func (storagePool *StoragePool) GetVolume(storagepoolhref string) (volumes []*types.Volume, err error) {

	endpoint := storagePool.client.SIOEndpoint

	if storagepoolhref == "" {
		link, err := GetLink(storagePool.StoragePool.Links, "/api/StoragePool/relationship/Volume")
		if err != nil {
			return []*types.Volume{}, errors.New("Error: problem finding link")
		}
		endpoint.Path = link.HREF
	} else {
		endpoint.Path = storagepoolhref
	}

	req := storagePool.client.NewRequest(map[string]string{}, "GET", endpoint, nil)
	req.SetBasicAuth("", storagePool.client.Token)
	req.Header.Add("Accept", "application/json;version=1.0")

	resp, err := checkResp(storagePool.client.Http.Do(req))
	if err != nil {
		return []*types.Volume{}, fmt.Errorf("problem getting response: %v", err)
	}
	defer resp.Body.Close()

	if storagepoolhref == "" {
		if err = decodeBody(resp, &volumes); err != nil {
			return []*types.Volume{}, fmt.Errorf("error decoding storage pool response: %s", err)
		}
	} else {
		storagePool := &types.Volume{}
		if err = decodeBody(resp, &storagePool); err != nil {
			return []*types.Volume{}, fmt.Errorf("error decoding instances response: %s", err)
		}
		volumes = append(volumes, storagePool)
	}
	return volumes, nil
}

func (storagePool *StoragePool) FindVolume(id, name, href string) (volume *types.Volume, err error) {
	// volumes, err := storagePool.GetVolume(href)
	// if err != nil {
	// 	return &types.Volume{}, errors.New("Error getting volumes")
	// }
	//
	// for _, volume = range volumes {
	// 	if volume.ID == id || volume.Name == name || href != "" {
	// 		return volume, nil
	// 	}
	// }
	//
	// return &types.Volume{}, errors.New("Couldn't find volumes")
	return &types.Volume{}, nil
}

func GetLocalVolumeMap() (mappedVolumes []*SdcMappedVolume, err error) {

	// get sdc kernel guid
	// /bin/emc/scaleio/drv_cfg --query_guid
	// sdcKernelGuid := "271bad82-08ee-44f2-a2b1-7e2787c27be1"

	mappedVolumesMap := make(map[string]*SdcMappedVolume)

	out, err := exec.Command("/bin/emc/scaleio/drv_cfg", "--query_vols").Output()
	if err != nil {
		return []*SdcMappedVolume{}, fmt.Errorf("Error querying volumes: ", err)
	}

	result := string(out)
	lines := strings.Split(result, "\n")

	for _, line := range lines {
		split := strings.Split(line, " ")
		if split[0] == "VOL-ID" {
			mappedVolume := &SdcMappedVolume{MdmID: split[3], VolumeID: split[1]}
			mdmVolumeID := fmt.Sprintf("%s-%s", mappedVolume.MdmID, mappedVolume.VolumeID)
			mappedVolumesMap[mdmVolumeID] = mappedVolume
		}
	}

	diskIDPath := "/dev/disk/by-id"
	files, _ := ioutil.ReadDir(diskIDPath)
	r, _ := regexp.Compile(`^emc-vol-\w*-\w*$`)
	for _, f := range files {
		matched := r.MatchString(f.Name())
		if matched {
			mdmVolumeID := strings.Replace(f.Name(), "emc-vol-", "", 1)
			devPath, _ := filepath.EvalSymlinks(fmt.Sprintf("%s/%s", diskIDPath, f.Name()))
			mappedVolumesMap[mdmVolumeID].SdcDevice = devPath
		}
	}

	keys := make([]string, 0, len(mappedVolumesMap))
	for key := range mappedVolumesMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		mappedVolumes = append(mappedVolumes, mappedVolumesMap[key])
	}

	return mappedVolumes, nil
}
