package goscaleio

import (
	"errors"
	"fmt"
	"net/http"

	types "github.com/thecodeteam/goscaleio/types/v1"
)

type ProtectionDomain struct {
	ProtectionDomain *types.ProtectionDomain
	client           *Client
}

func NewProtectionDomain(client *Client) *ProtectionDomain {
	return &ProtectionDomain{
		ProtectionDomain: &types.ProtectionDomain{},
		client:           client,
	}
}

func NewProtectionDomainEx(client *Client, pd *types.ProtectionDomain) *ProtectionDomain {
	return &ProtectionDomain{
		ProtectionDomain: pd,
		client:           client,
	}
}

func (s *System) CreateProtectionDomain(name string) (string, error) {

	protectionDomainParam := &types.ProtectionDomainParam{
		Name: name,
	}

	path := fmt.Sprintf("/api/types/ProtectionDomain/instances")

	pd := types.ProtectionDomainResp{}
	err := s.client.getJSONWithRetry(
		http.MethodPost, path, protectionDomainParam, &pd)
	if err != nil {
		return "", err
	}

	return pd.ID, nil
}

func (s *System) GetProtectionDomain(
	pdhref string) ([]*types.ProtectionDomain, error) {

	var (
		err error
		pd  = &types.ProtectionDomain{}
		pds []*types.ProtectionDomain
	)

	if pdhref == "" {
		var link *types.Link
		link, err = GetLink(
			s.System.Links,
			"/api/System/relationship/ProtectionDomain")
		if err != nil {
			return nil, err
		}

		err = s.client.getJSONWithRetry(
			http.MethodGet, link.HREF, nil, &pds)
	} else {
		err = s.client.getJSONWithRetry(
			http.MethodGet, pdhref, nil, pd)
	}
	if err != nil {
		return nil, err
	}

	if pdhref != "" {
		pds = append(pds, pd)
	}
	return pds, nil
}

func (s *System) FindProtectionDomain(
	id, name, href string) (*types.ProtectionDomain, error) {

	pds, err := s.GetProtectionDomain(href)
	if err != nil {
		return nil, fmt.Errorf("Error getting protection domains %s", err)
	}

	for _, pd := range pds {
		if pd.ID == id || pd.Name == name || href != "" {
			return pd, nil
		}
	}

	return nil, errors.New("Couldn't find protection domain")
}
