package pki

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testCAKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEA0UfRTFxIjp4+uNWvwe15U4PsQJX4ZAaN3RBvZYOQM23Kw8E1
aWimpmfON6V/hlpTidrhVHdNuD5UotuyYAcQUfSSKG/E4RWnunDMKkNq9wOz1lx0
lTuSqLAY4W589qQGDOg6PX/28eAglfcaNrwC+ANiu8cPZBA4CFn67cmDkXp7VhQT
ERmV+DXH0AM7gWpyUe8cEypxGV/ZJwckh1lKDZdiPwUMCa9DkPsWRHoM4MD/PlmG
ucjX/GDUH+Qdfl7Z726FJ4TdhBuHLYgMBcGeNmu/KWk7HVfDMrB/y/uY876e998S
YVA1/V8JcRaPA8x+orHtiwRIHj5rNNBANEWTJQIDAQABAoIBAE112Ij76zsmZS7C
oOAVhn+b55jsKVjUeDOnfqPDM19lttQHsj5GptAWweQk1HOlASCYLCi4U8LrByaU
TIxwcOD0thhTbjqlakR+tYK7G188NpcT92648wqOy1a9L3GWukqStePHdl6GR2la
YZB6vFqR3jyEbDTsL+EfdNoIaTMxyA8qXJz4jau42IBls/BRz4bBnpdjhE3S0IxJ
OL8L4ghLpF0jlbcjossGO5A0K5j2ky173U9CDMBW6vcilI1+tONRwzHDGoOIS3Gi
D4oIlAPBO6Jyrlb4jwRoyub3llzgzyhSBqzQnxWM9aL/Kl6DhP+BakE78VHj8v1w
DPtbvC0CgYEA67rAJlTm5zy9PBNYfeVKdQYRoHSa9/feu1KL8TD/KDK90P6KH1cd
Gixx79xOd44pFQ2J51u9juYjU0mjuGt+1U5x/nkOO18jKd+3qm/BsAMY+vmMgsNc
9IAPkBiLWRkbTGZfcCqi2f/SczOGngDMwb2cT1+BkRBHJmBdznf+VFMCgYEA40bU
sTWCrtslwrgNh6YXGtnlC9MNnRrNMaNd8gn9YBaAqAzSEVJWyZyAYmD0UN7GmknP
65LyHWwbt8HwyAYlBj5IIzqXeTM4OZHuuJxBvhfpsKfPy0yJEiAajaj2zhdfWJfW
4XlPURjkqmtmHQuBoyzjHE+aHgHY8xEOZtmvC6cCgYBkQKg3pSQOc+aHBjM8V6ey
3UHh27WMf/5Z7GFX0l6x2eKgX6Cec44M85oBSNCWR/9w1LExk/KqM3YSld7rL8xh
K1uPviwvU+bAiES0V5MoKCkXk8oOUsfVtCDqR4X7/pF9jIxKR9e6nvIBzIgT6oMq
Yll36EZSS3n2+ETs6ltfwQKBgAkjUOvrFd0H7KW+lrSsheNLfX0TOEnnyPZE9kME
Cc7yOKwJD+0oXVrv0u2hrlEOE/giHZ0AJIHwVdD2mELClHyCxo28DlkOKSWPa4S6
q54EAh5bMOygoCY9ajPl5j51DB1YxYf9Q6YkFRWRCeMDEmxIIr2BqdWpB1sGhYi3
GeWjAoGAIISevjT0InrKbEz3sFUAvb9BT5VvNbzqhqjN3b8Edg43qiqnJHLdwdBs
k663tkLdMfnpEIry8wWf6LKzUUCQp11TYcMIX7EhsVQwCX07wgpANCf/eokQvNw7
P3xlifXHcdOhmON1wA2CEmm4fID24iSVzUm6tIeM7BqxeprJr2E=
-----END RSA PRIVATE KEY-----`)
	testCACert = []byte(`-----BEGIN CERTIFICATE-----
MIIC9zCCAd+gAwIBAgIJAOA/9gyZIgEHMA0GCSqGSIb3DQEBCwUAMBIxEDAOBgNV
BAMMB2t1YmUtY2EwHhcNMTgxMDAzMTM0MTM2WhcNNDYwMjE4MTM0MTM2WjASMRAw
DgYDVQQDDAdrdWJlLWNhMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA
0UfRTFxIjp4+uNWvwe15U4PsQJX4ZAaN3RBvZYOQM23Kw8E1aWimpmfON6V/hlpT
idrhVHdNuD5UotuyYAcQUfSSKG/E4RWnunDMKkNq9wOz1lx0lTuSqLAY4W589qQG
DOg6PX/28eAglfcaNrwC+ANiu8cPZBA4CFn67cmDkXp7VhQTERmV+DXH0AM7gWpy
Ue8cEypxGV/ZJwckh1lKDZdiPwUMCa9DkPsWRHoM4MD/PlmGucjX/GDUH+Qdfl7Z
726FJ4TdhBuHLYgMBcGeNmu/KWk7HVfDMrB/y/uY876e998SYVA1/V8JcRaPA8x+
orHtiwRIHj5rNNBANEWTJQIDAQABo1AwTjAdBgNVHQ4EFgQUB19yMVhU5YVKYFHL
I8ouU1/qHtYwHwYDVR0jBBgwFoAUB19yMVhU5YVKYFHLI8ouU1/qHtYwDAYDVR0T
BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAEgVhmDU8zJK1goec+xaP61gjzcs1
tD/q4pFbF6h4vL39zueRITIqxT8apnftZsfkN/9NydRy8iFdZfrzBUTLRLi2d77Z
o4vJRIMnT2JetHQrPm22FvLU81isDInG/+3ytdBjOeR58OFCk35YEnQvMGnZ/+81
qOurhJ9euuP2ABoxPP2MVKtb6cERNiMY5dGqUtnNb60JWbl/MtqQRxEKbPAcF8rC
Tkas51LnxYGfF5M6hWRKdD/0kbnUbYFYb7rfy7xhiJiQBbtz35W8O52YJ/KBS7CR
ivLJLf7buMIFGDl+Dxf8ZYFtGBbko2ScVZKVcCimhSIQat/wcXVpgoMyZQ==
-----END CERTIFICATE-----`)
	testNonCAKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAw0AM9XKpZyED3BRWHrw201CwAFHkcuFfBJmoGU+LL53P5UYN
Apbs2XApbjuw+o0Abi6SAIMCpw1uw+K8z+9rNyPLnDP9rO2+a0Cf2Q7wxmeyZJ9D
QckzRDnuT7TJ/44tvFMbIh1H5jlExCsjhmcfIZa1eRxp8fvGnrsmYAernr32cjzl
1y8Ho+q6tga61wfy3oD8sumKXv4bnV7me5LYY6sQDU4D45I+TYd+VkjmxP5y1hn1
BK/wFcKo8KtpH9luRU81EilsaGrtsqOKkLZGYlTHx3emmA2jun7cbGU9Le083PHJ
33OBniIX8HFaxyTj36DXBUTeanadjq3AMKfecwIDAQABAoIBAGpE9EirbdR5vbYN
Q4pa+qQtKH4kFGkKjULdtwZ/qsvx4vNxGyerqaH3UDV1O7BDClbt9f1dURZHU2A+
inHKZ9qNGwcbrRhwFdIeTGJBpX73dbsai+rEbajWtwSe68qyQeZcfUomEboWtXzn
1eATRHhtpLwUYP2aAdLnHc3qntg4rFsZAi8pJsi/7exkxeJqWT1XbEcb+Tsk5lcw
FfE/IPn2DcjjzoV+jYCAEj9/bzPRoa7JdnX/K3YEe8nE5NnfskZCo7c9SBw2wbcK
O6sPTdjkTjhQCI2muBB9/1mjIS1ABSZOWgQW1AnwE4sFNEdXcPcy5Fm9T7tW4aWW
8LozuuECgYEA8OtagOGFIkNhoCvGJ9Fh7Sx/Oe6b89rhGSp1UEolZ7igdk2Ykf+Y
RQlKsSLRYOs9Cmdz4+nyqXnrcl4BkewqUQL/IiwlwsMxjGJTlstWC9XEfVHCzO/6
LQbM4RjH0iOLIwQam+ZGA7Bb8ZyjEhdYtnZK95gVaWTgyvEyUIO7cFsCgYEAz3jd
gT/vWrwFhH9aS1hTIaxRaatA/0pwVknt5imyjVMo+ZxxxVzmRewxwhR3B6Dv6vZx
F6Ops3z3T0YjcQpXVvnZJxwI67vh7/LpmsSmtrVBbFDjGzZtiSqlDNOXBlw7/R9r
9//5sClew1XwbBujwUBBgglBtF/zzedbYcXWpckCgYEAlGg2qRvDOlcNpXAxsceO
rl5xxQsScIZNkYYRHDOAlUMrPZURPiaX8zcFFtce6bgfMvCFeEleHT4oZpw4FV7I
tnzFE5TkcfRx6kuLuGFrkQDO+G/MMxhFIUWGIcd1GCKjDB/0EEMqsA0MpmpaHcPZ
9xQpnBnIXtMwknNADk8HwO8CgYAvmy5IgCEuEsK5Wnefnk7FBUNRGei6K5yHUEN0
ctDzuMdIL2uzu9Ni7AWm4QdHCtjCc3YT1IwWEXC2EgQD5jmQTZhUbwxk+yGm63hK
+SC//+tZLV5PWjfcJ5rjzJF09ikVteYSa/whPfzumYOnatgyecoOSo13FCVfc9z2
HG1acQKBgHTe8K4OiR88mlKVtgBCifFJYtbrnCRuSeJY+uBekpQVJpMR6pZmU8oQ
qmKBA+ZCCcM/PstVlOJ/SMGpVMro5sN8aoREC32153Np4//b9zmFd4dVPA0sEoby
s90+eaCT+U0+QWJ8CcHjClt8pMVqiaM88G/gaftWDz7//CyaUd6B
-----END RSA PRIVATE KEY-----`)
	testNonCACert = []byte(`-----BEGIN CERTIFICATE-----
MIIDPzCCAiegAwIBAgIJAOxqgOgmj2IFMA0GCSqGSIb3DQEBCwUAMBIxEDAOBgNV
BAMMB2t1YmUtY2EwHhcNMTgxMDAzMTM0NjE3WhcNMTkxMDAzMTM0NjE3WjAZMRcw
FQYDVQQDDA5rdWJlLWFwaXNlcnZlcjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
AQoCggEBAMNADPVyqWchA9wUVh68NtNQsABR5HLhXwSZqBlPiy+dz+VGDQKW7Nlw
KW47sPqNAG4ukgCDAqcNbsPivM/vazcjy5wz/aztvmtAn9kO8MZnsmSfQ0HJM0Q5
7k+0yf+OLbxTGyIdR+Y5RMQrI4ZnHyGWtXkcafH7xp67JmAHq5699nI85dcvB6Pq
urYGutcH8t6A/LLpil7+G51e5nuS2GOrEA1OA+OSPk2HflZI5sT+ctYZ9QSv8BXC
qPCraR/ZbkVPNRIpbGhq7bKjipC2RmJUx8d3ppgNo7p+3GxlPS3tPNzxyd9zgZ4i
F/BxWsck49+g1wVE3mp2nY6twDCn3nMCAwEAAaOBkDCBjTAJBgNVHRMEAjAAMAsG
A1UdDwQEAwIF4DBzBgNVHREEbDBqggprdWJlcm5ldGVzghJrdWJlcm5ldGVzLmRl
ZmF1bHSCFmt1YmVybmV0ZXMuZGVmYXVsdC5zdmOCHmt1YmVybmV0ZXMuZGVmYXVs
dC5zdmMuY2x1c3RlcocECgMAAYcEFB4oMocEChQeKDANBgkqhkiG9w0BAQsFAAOC
AQEAbMRIi2/3ZcUEB44uPWZ6mZOd2nKZUEuyRZnQLzq8MxcOMGjPjs3I2ktubYTF
0zET7hs/4B+rmbOa1sxMEq3Qq433oXFk7q1JuEwfh1/XAI2ITJBo4yNFw6yTOQsT
nxvlN5sBgE+iH5qd9ooJOGDZZrH6h0K7hU3Bvs0P4kvEHtgwJT5p4D3AYrtW95k2
sx8bzPoIaaUIohJsVqeX9yt2Zrq44q0TOf6PdTVtvOdyEaaRSN8f/t6+W1KD8UOM
ZelwQz0negCFYA3fAhXBCy/UpNlhwEnVg0ghGo5bGI5+6PEE6Z6f3FOZJzcmwEHQ
xUIn/rpHJeyLQdx+1S5dVrxzkg==
-----END CERTIFICATE-----`)
)

func TestGenerateSelfSignedCAKey(t *testing.T) {
	pki, err := NewCAPair(nil)
	require.NoError(t, err)

	require.NotNil(t, pki.Cert)
	require.NotNil(t, pki.Key)
}

func TestNewPKI(t *testing.T) {
	testCases := []struct {
		description string
		expectedErr error
		CA          []byte
	}{
		{
			description: "success self signed",
			CA:          nil,
			expectedErr: nil,
		},
		{
			description: "success provided",
			expectedErr: nil,
			CA:          testCACert,
		},
		{
			description: "error provided",
			expectedErr: ErrInvalidCA,
			CA:          testNonCACert,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		p, err := NewCAPair(testCase.CA)

		if err != testCase.expectedErr {
			t.Errorf("Wrong error expected %v actual %v",
				testCase.expectedErr, err)
		}

		if err == nil && p == nil {
			t.Errorf("pki bundle must not be nil")
		}
	}
}
