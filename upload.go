package synofs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
)

const (
	uploadEndpoint = `webapi/entry.cgi`
	uploadAPI      = `SYNO.FileStation.Upload`
	uploadVersion  = 2
	uploadMethod   = `upload`
)

type uploadService struct {
	client *Client
}

func newUploadService(client *Client) *uploadService {
	return &uploadService{
		client: client,
	}
}

func (s *uploadService) Send(path, filename string, parents, overwrite bool, r io.Reader) error {
	if !s.client.authenticated {
		return fmt.Errorf("client upload error: not authenticated")
	}

	target := fmt.Sprintf(
		"%s/%s?api=%s&version=%d&method=%s",
		s.client.address,
		uploadEndpoint,
		uploadAPI,
		uploadVersion,
		uploadMethod,
	)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	err := writer.WriteField("path", path)
	if err != nil {
		return fmt.Errorf("client upload error adding path: %q", err)
	}

	err = writer.WriteField("create_parents", strconv.FormatBool(parents))
	if err != nil {
		return fmt.Errorf("client upload error adding create parents: %q", err)
	}

	writer.WriteField("overwrite", strconv.FormatBool(overwrite))
	if err != nil {
		return fmt.Errorf("client upload error adding overwrite: %q", err)
	}

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return fmt.Errorf("client upload error adding file: %q", err)
	}

	_, err = io.Copy(part, r)
	if err != nil {
		return fmt.Errorf("client upload error adding file body: %q", err)
	}

	resp, err := s.client.httpClient.Post(target, writer.FormDataContentType(), body)
	if err != nil {
		return fmt.Errorf("client upload error: %q", errors.Unwrap(err))
	}

	body.Reset() //TODO maybe a pool of this instead of create a new one on each call

	defer resp.Body.Close()

	var uploadResp APIResponse

	err = json.NewDecoder(resp.Body).Decode(&uploadResp)

	if err != nil {
		return fmt.Errorf("client upload error: %q", err)
	}

	if !uploadResp.Success {
		uploadResp.Error.API = uploadAPI
		return fmt.Errorf("client upload error: %w", uploadResp.Error)
	}

	return nil
}
