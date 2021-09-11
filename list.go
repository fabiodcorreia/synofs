package synofs

import (
	"errors"
	"fmt"
	"net/url"
)

const (
	listEndpoint    = `webapi/entry.cgi`
	listAPI         = `SYNO.FileStation.List`
	listVersion     = `2`
	listShareMethod = `list_share`
	listMethod      = `list`
)

type listService struct {
	client *Client
}

func newListService(client *Client) *listService {
	return &listService{
		client: client,
	}
}

func (s *listService) Shares() ([]ListShare, error) {
	if !s.client.authenticated {
		return nil, fmt.Errorf("client shares error: not authenticated")
	}

	params := url.Values{}
	params.Add(paramAPI, listAPI)
	params.Add(paramMethod, listShareMethod)
	params.Add(paramVersion, listVersion)

	resp, err := s.client.httpClient.Get(genEndpoint(s.client.address, listEndpoint, params))
	if err != nil {
		return nil, fmt.Errorf("client shares error: %q", errors.Unwrap(err))
	}

	var listShareResp listSharesResponse
	if err := handleResponse(resp, &listShareResp); err != nil {
		return nil, fmt.Errorf("client shares error: %q", err)
	}

	if !listShareResp.Success {
		listShareResp.Error.API = listAPI
		return nil, fmt.Errorf("client shares error: %w", listShareResp.Error)
	}

	return listShareResp.Data.Shares, nil
}

func (s *listService) Files(path string) ([]ListFile, error) {
	if !s.client.authenticated {
		return nil, fmt.Errorf("client list error: not authenticated")
	}

	params := url.Values{}
	params.Add(paramAPI, listAPI)
	params.Add(paramMethod, listMethod)
	params.Add(paramVersion, listVersion)
	params.Add("folder_path", path)
	params.Add("additional", `["size","time","type"]`)

	resp, err := s.client.httpClient.Get(genEndpoint(s.client.address, listEndpoint, params))
	if err != nil {
		return nil, fmt.Errorf("client list error: %q", errors.Unwrap(err))
	}

	var listResp listResponse
	if err := handleResponse(resp, &listResp); err != nil {
		return nil, fmt.Errorf("client list error: %q", err)
	}

	if !listResp.Success {
		listResp.Error.API = listAPI
		return nil, fmt.Errorf("client list error: %w", listResp.Error)
	}

	return listResp.Data.Files, nil
}

type listResponse struct {
	Data listData `json:"data"`
	APIResponse
}

type listData struct {
	Offset int        `json:"offset"`
	Total  int        `json:"total"`
	Files  []ListFile `json:"files"`
}

type ListFile struct {
	IsDir      bool       `json:"isdir"`
	Name       string     `json:"name"`
	Path       string     `json:"path"`
	Additional Additional `json:"additional"`
}

type Additional struct {
	Type string  `json:"type"`
	Size int64   `json:"size"`
	Time APITime `json:"time"`
}

type listSharesResponse struct {
	Data listSharesData `json:"data"`
	APIResponse
}

type listSharesData struct {
	Offset int         `json:"offset"`
	Total  int         `json:"total"`
	Shares []ListShare `json:"shares"`
}

type ListShare struct {
	IsDir bool   `json:"isdir"`
	Name  string `json:"name"`
	Path  string `json:"path"`
}
