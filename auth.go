package synofs

import (
	"errors"
	"fmt"
	"net/url"
)

const (
	authEndpoint     = `webapi/auth.cgi`
	authAPI          = `SYNO.API.Auth`
	authVersion      = `3`
	authLoginMethod  = `login`
	authLogoutMethod = `logout`
	authSession      = `FileStation`
)

type authService struct {
	client *Client
}

func newAuthService(client *Client) *authService {
	return &authService{
		client: client,
	}
}

func (s *authService) Login(username, password string) (string, error) {
	params := url.Values{}
	params.Add(paramAPI, authAPI)
	params.Add(paramMethod, authLoginMethod)
	params.Add(paramVersion, authVersion)
	params.Add("account", username)
	params.Add("passwd", password)
	params.Add("session", authSession)

	resp, err := s.client.httpClient.Get(genEndpoint(s.client.address, authEndpoint, params))
	if err != nil {
		return "", fmt.Errorf("client login error: %q", errors.Unwrap(err))
	}

	var authRes authResponse
	if err := handleResponse(resp, &authRes); err != nil {
		return "", fmt.Errorf("client login error: %q", err)
	}

	if !authRes.Success {
		authRes.Error.API = authAPI
		return "", fmt.Errorf("client login error:  %w", authRes.Error)
	}

	s.client.authenticated = true

	return authRes.Data.SID, nil
}

func (s *authService) Logout() error {
	params := url.Values{}
	params.Add(paramAPI, authAPI)
	params.Add(paramMethod, authLogoutMethod)
	params.Add(paramVersion, authVersion)
	params.Add("session", authSession)

	resp, err := s.client.httpClient.Get(genEndpoint(s.client.address, authEndpoint, params))
	if err != nil {
		return fmt.Errorf("client logout error: %q", errors.Unwrap(err))
	}

	var authRes authResponse
	if err := handleResponse(resp, &authRes); err != nil {
		return fmt.Errorf("client logout error: %q", err)
	}

	if !authRes.Success {
		authRes.Error.API = authAPI
		return fmt.Errorf("client logout error: %w", authRes.Error)
	}

	s.client.authenticated = false

	return nil
}

type authResponse struct {
	APIResponse
	Data authData `json:"data"`
}

type authData struct {
	SID string `json:"sid"`
}
