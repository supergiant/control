// +build integration

package pki

import (
	"net/http"
	"testing"

	"bytes"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/stretchr/testify/require"
	"github.com/supergiant/supergiant/pkg/pki"
	"github.com/supergiant/supergiant/pkg/storage"
	"github.com/supergiant/supergiant/pkg/testutils/assert"
	"net/http/httptest"
)

const (
	defaultETCDHost = "http://127.0.0.1:2379"
)

var defaultConfig clientv3.Config

const validCert = `-----BEGIN CERTIFICATE-----
MIIDWjCCAkICCQDL2Tw+Yt3GVTANBgkqhkiG9w0BAQUFADBvMRAwDgYDVQQDDAdx
Ym94LmlvMQ0wCwYDVQQKDARhc2RmMRAwDgYDVQQHDAdhc2Rmc2RmMRAwDgYDVQQI
DAdBbGFiYW1hMQswCQYDVQQGEwJVUzEbMBkGCSqGSIb3DQEJARYMYXNkZkBhc2Qu
Y29tMB4XDTE4MDgwNjE0MTc1N1oXDTE5MDgwNjE0MTc1N1owbzEQMA4GA1UEAwwH
cWJveC5pbzENMAsGA1UECgwEYXNkZjEQMA4GA1UEBwwHYXNkZnNkZjEQMA4GA1UE
CAwHQWxhYmFtYTELMAkGA1UEBhMCVVMxGzAZBgkqhkiG9w0BCQEWDGFzZGZAYXNk
LmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALl1lRL97v6c+v3u
BTdX/hpHjuXHVgM9IIVpekPF3hq1tI/tKBL/TShn0zKrIvSF65w1VEfezwPQuIad
XMMQQgM4gaGkEwQ/lDkZGvZ9z61S3ydxcOpKBs1iDr9gKYcln+Hgnpxnhw9mUh/7
/SCayggRDfzNa6KlSDYs5J1+hbdsPlcVcQdIqdb8kbPyp65skb2a71iKY5ZhZ5kH
WvJp9YtZpm+QMCP2Jpfh0lmN62XpxkMZvwtBncPJueeoFNL0u3iSd+x5AFautC8t
EiRPyA241z/xkZdU/hqV55knk/h0tNJmGcPKLvZ+RsptRRhTPjKY+1/tICjCCj1O
MTIgwh0CAwEAATANBgkqhkiG9w0BAQUFAAOCAQEAaLcKaxtBJKMqXAkEv891ESka
2hpWdW+I7b75O25VMJRy74fvaoNFric532zj3hTqaRWG17TeoAR8XE4G08UrcNIk
Z+4jfJBFxYgF3MJkx9M/zadfVIFxllM4BPw0yEOUwITVT6IlzgnT9z56OtTUhaRV
rBgYCvnkY0C+pklGmd/Lvdr603umavbyyXXAJ4n5z2lBo3f4dElDPMr0hOF/KcwM
AFrC2HIoROyJc07utw6pgNpZQGbx3B5sjX/bm8BS3NXMR1DDX2AAluKMbToSTioO
hc4whUkBq2OT/W37+SK/nOKFUn6LDG8yFRxeJ5uScpfJgUAqxi0/k0fnXodd9A==
-----END CERTIFICATE-----`

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAuXWVEv3u/pz6/e4FN1f+GkeO5cdWAz0ghWl6Q8XeGrW0j+0o
Ev9NKGfTMqsi9IXrnDVUR97PA9C4hp1cwxBCAziBoaQTBD+UORka9n3PrVLfJ3Fw
6koGzWIOv2AphyWf4eCenGeHD2ZSH/v9IJrKCBEN/M1roqVINizknX6Ft2w+VxVx
B0ip1vyRs/KnrmyRvZrvWIpjlmFnmQda8mn1i1mmb5AwI/Yml+HSWY3rZenGQxm/
C0Gdw8m556gU0vS7eJJ37HkAVq60Ly0SJE/IDbjXP/GRl1T+GpXnmSeT+HS00mYZ
w8ou9n5Gym1FGFM+Mpj7X+0gKMIKPU4xMiDCHQIDAQABAoIBAGWt1KSL+ms379gm
lk+Ie7U2xF6wUjUGX30lnjXoFuR3+N3r/UulE01y1vTxpQGBJvMGvgWFX+RMm86a
GhCMKlUPturDRPXQUdiYLhM0WRdC1zwN0wVwvpf+Ce3csAf7ldPGTc+cZw0HYUFN
67Ljip6vkwamLTwH+DZTmKfMhU1RKVI+85wm3yE7OuPT7C6iFTAOzVoZCso7Zt9m
NW3iXcSrPn6JBpezoX/0ZrkloHy5L68cZiPwNoQg4EkMTBOZ6nIdUxNxN3RKfE5Q
Tm+LVWJJJF9+ZE+MQBj++cUtH8hdZwIjDfrswTmozf/DXnBG8/+WjTnVMLvjt+h9
iZ0RgSECgYEA209zJBnWBen79li5gTrJaG9LmEzFfuUkBZndmmAGclcFsVXfLKo1
bHOojLdruRUWcBdgH8YTBub5dcSh6mUbBFOmEf3AjleHMX+eNDVqTykWi7KNMsW3
YdqekAZ4rT+kdrBNG7L0m2AZ33mnd7xMYTUF2fXH3tAB+5VaVm8pUpkCgYEA2Hxf
EKpfMA4Srni/5p3xPF7q9qCRKxcYFXpoH6iG3jrs51RqziJoopNikcdoaeo3iyYA
2L5LG0RMSLYh0rz6aR7lw38hW7TxdS+OqX73yl/o+amvqpOu5caUvKpLodYRARHy
AUAsUcmMPZa+DItgDmWdIL3+cyJMriKm8SLboiUCgYBTFjWcHsGr+erAePrGz/vQ
OiIcsDE+kxdjm9iODQVEOl3owozLwix9SxA3R6JjO28FxoVfZE5/FfC6wmVJhUaI
DBzlwgo6o0SP5zaLtxTwqrNk959w9eE1DHt4O0tq76qiYMbF0LXFS9JhjRh6T3ds
eIcf/XLcolet9faEupagOQKBgHY3lII5RzmqtbDo54I8BZv+CTkcfamWNuSjr3B6
SwvYCb5ZbumaCKGe8ljBF9eeuy4VVqkFYWZGaZHbQ6Uc5XG6GaYkKkc2DBT+H12X
pCCzNzn+25q+guefBWHxbNO3XhnDfvAH5yvSb+7B/o5DHfU+sAtNNUISHOWKrrdH
XcCxAoGBAITb5Y49LlCesvEwoHQfKiLs2eh29tjCGBA6y62hO8EzIReibfMPiSxn
uZHHNeMDu9+DXXLGkf6K4AxCm8bW7iO9iovJbdGUiqEo2wHxmcKoSAdAIwvUlf+0
5WHLGFYvTkR+AV9AhtA7pR6db2zYRj6yqhT2iGvd6rOBZPRoJ8eF
-----END RSA PRIVATE KEY-----`

func init() {
	assert.MustRunETCD(defaultETCDHost)
	defaultConfig = clientv3.Config{
		Endpoints: []string{defaultETCDHost},
	}
}

func TestGenerateCertificates(t *testing.T) {
	kv := storage.NewETCDRepository(defaultConfig)

	h := pki.NewHandler(pki.NewService("/test/cert", kv))

	handler := http.HandlerFunc(h.Generate)

	model := pki.CARequest{
		DNSDomain: "test",
		IPs:       []string{"127.0.0.1"},
	}
	caReq, _ := json.Marshal(model)

	req, _ := http.NewRequest(http.MethodPost, "", bytes.NewReader(caReq))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, rr.Body.String())
	fmt.Println(rr.Body.String())
}
