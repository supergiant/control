package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/supergiant/supergiant/pkg/model"
)

type Client struct {
	BaseURL   string
	AuthType  string // token, session
	AuthToken string

	httpClient *http.Client

	Sessions         *Sessions
	Users            *Users
	CloudAccounts    *CloudAccounts
	Kubes            *Kubes
	Apps             *Apps
	Components       *Components
	Releases         *Releases
	Instances        *Instances
	Volumes          *Volumes
	PrivateImageKeys *PrivateImageKeys
	Entrypoints      *Entrypoints
	Nodes            *Nodes
}

func New(url string, authType string, authToken string, certFile string) *Client {
	client := &Client{
		BaseURL:   url,
		AuthType:  authType,
		AuthToken: authToken,
	}

	transport := new(http.Transport)

	if certFile != "" {
		pem, err := ioutil.ReadFile(certFile)
		if err != nil {
			panic(err)
		}
		roots := x509.NewCertPool()
		if ok := roots.AppendCertsFromPEM(pem); !ok {
			panic("failed to parse root certificate")
		}
		transport.TLSClientConfig = &tls.Config{
			RootCAs: roots,
		}
	}

	client.httpClient = &http.Client{
		Transport: transport,
	}

	client.Sessions = &Sessions{Collection{client, "sessions"}}
	client.Users = &Users{Collection{client, "users"}}
	client.CloudAccounts = &CloudAccounts{Collection{client, "cloud_accounts"}}
	client.Kubes = &Kubes{Collection{client, "kubes"}}
	client.Apps = &Apps{Collection{client, "apps"}}
	client.Components = &Components{Collection{client, "components"}}
	client.Releases = &Releases{Collection{client, "releases"}}
	client.Instances = &Instances{Collection{client, "instances"}}
	client.Volumes = &Volumes{Collection{client, "volumes"}}
	client.PrivateImageKeys = &PrivateImageKeys{Collection{client, "private_image_keys"}}
	client.Entrypoints = &Entrypoints{Collection{client, "entrypoints"}}
	client.Nodes = &Nodes{Collection{client, "nodes"}}

	return client
}

func (c *Client) request(method string, path string, in interface{}, out interface{}, queryValues map[string][]string) error {
	body := new(bytes.Buffer)
	if in != nil {
		if err := json.NewEncoder(body).Encode(in); err != nil {
			return err
		}
	}

	requestURL, err := url.Parse(c.BaseURL + "/" + path)
	if err != nil {
		return err
	}

	q := requestURL.Query()
	for key, values := range queryValues {
		for _, value := range values {
			q.Add(key, value)
		}
	}
	requestURL.RawQuery = q.Encode()

	req, err := http.NewRequest(method, requestURL.String(), body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf(`SGAPI %s="%s"`, c.AuthType, c.AuthToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.Status[:2] != "20" {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		errModel := new(model.Error)
		if err := json.Unmarshal(body, errModel); err != nil {
			return err
		}
		return errModel
	}

	if out != nil {
		defer resp.Body.Close()
		if err = json.NewDecoder(resp.Body).Decode(out); err != nil {
			return err
		}
	}

	return nil
}
