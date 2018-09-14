package digitalocean

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"strings"
	"context"
	"encoding/base64"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
)

// Returns private ip
func getPrivateIpPort(networks []godo.NetworkV4) string {
	for _, network := range networks {
		if network.Type == "private" {
			return network.IPAddress
		}
	}

	return ""
}

// Returns public ip
func getPublicIpPort(networks []godo.NetworkV4) string {
	for _, network := range networks {
		if network.Type == "public" {
			return network.IPAddress
		}
	}

	return ""
}

func fingerprint(key string) (string, error) {
	parts := strings.Fields(string(key))

	k, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", errors.Wrap(err, "fingerprint decode string")
	}

	fp := md5.Sum(k)
	buffer := &bytes.Buffer{}

	for i, b := range fp {
		fmt.Fprintf(buffer, "%02x", b)

		if i < len(fp)-1 {
			fmt.Fprint(buffer, ":")
		}
	}

	return buffer.String(), nil
}

func createKey(ctx context.Context, keyService KeyService, name, publicKey string) (*godo.Key, error) {
	req := &godo.KeyCreateRequest{
		Name:      name,
		PublicKey: publicKey,
	}

	key, _, err := keyService.Create(ctx, req)

	if err != nil {
		return nil, err
	}

	return key, err
}
