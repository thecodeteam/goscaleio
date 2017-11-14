package goscaleio

import (
	"fmt"
	"net/http"

	types "github.com/thecodeteam/goscaleio/types/v1"
)

func (s *System) GetScsiInitiator() ([]types.ScsiInitiator, error) {

	path := fmt.Sprintf(
		"/api/instances/System::%v/relationships/ScsiInitiator",
		s.System.ID)

	var si []types.ScsiInitiator
	err := s.client.getJSONWithRetry(
		http.MethodGet, path, nil, &si)
	if err != nil {
		return nil, err
	}

	return si, nil
}
