package synofs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/fabiodcorreia/httpcg"
)

const (
	paramAPI     = `api`
	paramVersion = `version`
	paramMethod  = `method`
)

type Client struct {
	httpClient    *http.Client
	address       string
	authenticated bool
	Auth          *authService
	List          *listService
	Upload        *uploadService
}

func New(address string) *Client {
	//TODO handle this error
	httpClient, _ := httpcg.NewBuilder().
		WithCookies().
		ResponseHeaderTimeout(0 * time.Second).
		Build()
	return NewWithHTTP(address, httpClient)
}

func NewWithHTTP(address string, httpClient *http.Client) *Client {
	c := &Client{
		address:       address,
		authenticated: false,
		httpClient:    httpClient,
	}

	c.Auth = newAuthService(c)
	c.List = newListService(c)
	c.Upload = newUploadService(c)
	return c
}

type APIResponse struct {
	Error   ClientError `json:"error"`
	Success bool        `json:"success"`
}

type APITime struct {
	AccessTime   int `json:"atime"`
	CreationTime int `json:"crtime"`
	ChangeTime   int `json:"ctime"`
	ModifiedTime int `json:"mtime"`
}

func genEndpoint(address, endpoint string, params url.Values) string {
	target := fmt.Sprintf("%s/%s?%s", address, endpoint, params.Encode())
	fmt.Println(target)
	return target
}

func handleResponse(resp *http.Response, v interface{}) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code %v is not %v", resp.StatusCode, http.StatusOK)
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(v)
}
