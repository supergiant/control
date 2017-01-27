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
	Version string

	BaseURL   string
	AuthType  string // token, session
	AuthToken string

	httpClient *http.Client

	Sessions      SessionsInterface
	Users         UsersInterface
	CloudAccounts CloudAccountsInterface
	Kubes         KubesInterface
	KubeResources KubeResourcesInterface
	Nodes         NodesInterface
	LoadBalancers LoadBalancersInterface
	HelmRepos     HelmReposInterface
	HelmCharts    HelmChartsInterface
	HelmReleases  HelmReleasesInterface
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
	client.KubeResources = &KubeResources{Collection{client, "kube_resources"}}
	client.Nodes = &Nodes{Collection{client, "nodes"}}
	client.LoadBalancers = &LoadBalancers{Collection{client, "load_balancers"}}
	client.HelmRepos = &HelmRepos{Collection{client, "helm_repos"}}
	client.HelmCharts = &HelmCharts{Collection{client, "helm_charts"}}
	client.HelmReleases = &HelmReleases{Collection{client, "helm_releases"}}

	return client
}

func (c *Client) request(method string, path string, in interface{}, out interface{}, queryValues map[string][]string) error {
	body := new(bytes.Buffer)
	if in != nil {
		if err := json.NewEncoder(body).Encode(in); err != nil {
			return err
		}
	}

	requestURL, err := url.Parse(c.BaseURL + "/api/v0/" + path)
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

	req.Close = true

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
			// If unmarshalling failed, we have to fallback to capturing the full text
			errModel.Message = string(body)
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
