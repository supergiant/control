package digitalocean

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"github.com/digitalocean/godo"
	"golang.org/x/crypto/ssh"
	"strings"
	"testing"
)

var (
	privateKeyBytes = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAwdnJH7+a1ckkqu6fc4hLfsGlNvGhRqPVTRq/YKLbQu3XsNjf
oIC4+w81FABCNvNPajwcTRTfQv6wtmj/m60vI3LyJ23DMThyv0M8czkg40EDn6ub
UOOdR7IZAqkRbc+XCEYK5/8i2j2fuHa/oLKbILvKr0xmSG/Mwso1T2bH4YO9GfEJ
ftZJ/tSk5Yzpb0Oo5cWb2V96OvJDG/MEb5BXPKbOMU8ngnx2gqs+jh1p1fQwrb0G
JjJgA83sySjr8AHkelH9QNS8HvgfHlL8bX2nmFGeaVrRxRXzspJndN15NKUWAMLU
gvWVno3ulqKMkzH+ehK6SFCuKoF5icxtqBx4wwIDAQABAoIBAAwbSulVsRja+BRI
1OKFR5nCBEx7KMRdpQusuPkTEriKXCcqVEUU5PihCYKXRYtjBLmwyV+zBwKLH4Q0
6InTdhczrZXyz/b5/IifbV4Q2lH3FH/bWtbhcEgzAkbdQj5mcZtNrI6yq32PzbLa
j7s8jF2t/MmX7udlPBeKQ2wTEjauRKy5yAS3nshHtsZ1yTa5wFdtrdv3Ulav8w+O
zcoQFWP5CM63MKi4vQ8wN/Nl+wYFlFztVJMbqfUhcWHkkQwyZRha3PfjyrbShwDS
rq7v9UoJ8Zuun0FYQkkwavWBvrG+euhFrDJSV/77m2rkvWJyymfzlNDo4rNESjqs
U+RSppECgYEA8daWHxlou0F2iSEAU3m2femj3a+MrDSgYN5Gdg4MH12+X5WC6hvs
JAvqEi7NkBqozei+FhaKuB4mkONBBwjx7mIwT/J/QFy9Ejr4Lg3BJeyachrB5WxG
6PcZn8IzPXXpyxBzZUMeo17Gif5/dQwTI8CHyFz7jGT0BKnJfpPOf0kCgYEAzTPQ
tJNb83HYAbziKzlacGZFUDrjHDTV/wr1NnMpj551jJVrFTBvz99wWv4uvMXTtiyH
QZftY+wyLbfgFGPAQnp0x7c8/YC656bEehYEa4ecrPy14seLk/4KAI6IMiUbspe1
iuUdSlsS3FaUSE1BM35QCpgaTPs3rNsZGjH626sCgYEA2BLtYG34aE6uFQl6XBsE
VW26LmkaHAaNQN94PyR/6kp8vLQ+GuPF0dMfWQ2eNuHK7ubDZ8LOQIEX3h5dzGZO
mrn6BoRY8+2oNLCha6x4ZWUH/Wkw0sYyeRXGPDpsQ76lm/xfzhrxNfCJHWRZBwA5
3Zi4+OkzC5Zre/sjf8eaGZkCgYBgqo0h091YNIQWZX2B8TW6h2MVpXgBfJ5m1Cmp
6dxlTLeBb44PYE779P0/0EgCI4tVYWqiKsjo7obA5MMJt+gFKRzETHzNywvBPt2F
ycNxSGQ1VaL1Xx1QrTbXBk4AmVyP6EncUYxXz8l1xM97s/EIKfPY2chiBWI36srL
fUn4mwKBgQCpELvH91vYlzECe+uM5SPHEKWXZf27+I+ABOttVlxMYomjSfWnuKRp
aJTJFasfzhumpKl9+T3HZ3WNx2YzzggkDs3l/CuBeJcuHmWBZ/99VSAcQ4v7bJez
ftubvUfbIYGZDXxwqpBifsVm22U3lTggt1MSzSFzUaS6RkvOw+wW2g==
-----END RSA PRIVATE KEY-----`
	expectedPublicKey = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDB2ckfv5rVySSq7p9ziEt+waU28aFGo9VNGr9gottC7dew2N+ggLj7DzUUAEI2809qPBxNFN9C/rC2aP+brS8jcvInbcMxOHK/QzxzOSDjQQOfq5tQ451HshkCqRFtz5cIRgrn/yLaPZ+4dr+gspsgu8qvTGZIb8zCyjVPZsfhg70Z8Ql+1kn+1KTljOlvQ6jlxZvZX3o68kMb8wRvkFc8ps4xTyeCfHaCqz6OHWnV9DCtvQYmMmADzezJKOvwAeR6Uf1A1Lwe+B8eUvxtfaeYUZ5pWtHFFfOykmd03Xk0pRYAwtSC9ZWeje6WooyTMf56ErpIUK4qgXmJzG2oHHjD`
)

func TestFingerPrint(t *testing.T) {
	expected := "ed:79:fd:40:e6:a1:05:64:ce:84:40:94:72:eb:9c:ee"
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCpJTpadNp+c8MMQ/cyiWMjio5WxsFklDxG4RCuP6tgUEWqnANelNxT/lkIO5hUCfCS8a4wGPiOWIJpYMmmQRz7lysqm3hGGLVSv1H8m9XY//t/Xd+On7M/FZtr1AB/WV/11YBU8jW0TWk/pgPHjUUYnbbPAK1iilQS1ULx/Wen6EmjzVqD8XDLl82/cQgfT6UF1ZVQd+7qPmdeK4her+Otg/rTwIqjQI7DObhThpn7ZHehclTULw0jtAGw7/3Bek/DAuKSG3yQ+hMg+0xqO1t6zo12kYlRwpGTiCW2zLAuVw7PW7nz3SGvOTAjXAzKYcVdCn9rSs6UqufP4FV2BlbW3ZoQJY2KoEuDFbgmyhP8Z/+A6EXVkQBY/jHHsJGWIZS1QGpSAbYEGubb/lKryw0k1nr4X+bmFeymuOSWdipYOv/b4nXUrI+qIAZIza7heSM5BuRqkvVO/SSqyNbrypWHmL8x+EVb0WiSLQqFh/VZKiW0cgZ2gWL+qYyHuKlTPXCa+vO3SpPVFyIKV6WlblrSeCpwC6dj94RSkQejOojXvUJ1eT504dU8zyDYgE5nAgxeJecnM5+5Kowb/Zi5ByIjAmRE8e7ST4C9g73sue3t5foJ6IItJtlVgIoP5W3GLbRJ8p8T5SQY7fIVR6BiUmWU9BR2XdWVi2sH/x1IW9meoQ=="

	fg, err := fingerprint(publicKey)

	if err != nil {
		t.Error(err)
	}

	if !strings.EqualFold(fg, expected) {
		t.Errorf("Wrong fingerprint expected %s actual %s", expected, fg)
	}
}

func TestGeneratePublicKey(t *testing.T) {
	pk, _ := ssh.ParseRawPrivateKey([]byte(privateKeyBytes))

	privateKey, _ := pk.(*rsa.PrivateKey)
	publicKey, err := generatePublicKey(&privateKey.PublicKey)

	if err != nil {
		t.Errorf("generate public key %v", err)
	}

	if !strings.EqualFold(string(publicKey[:len(publicKey)-1]), expectedPublicKey) {
		t.Errorf("Wrong public key expected \n %s actual \n %s",
			expectedPublicKey, publicKey)
	}
}

func TestGenerateKeyPair(t *testing.T) {
	privateKey, _, err := generateKeyPair(keySize)

	if err != nil {
		t.Errorf("generate key pair %v", err)
	}

	privateKeyRSA, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))

	if err != nil {
		t.Errorf("parse private key %v", err)
	}

	privateKeyRSA.Validate()
}

func TestGetPublicIpAddr(t *testing.T) {
	testCases := []struct {
		networks []godo.NetworkV4
		expected string
	}{
		{
			networks: []godo.NetworkV4{
				{
					IPAddress: "10.0.0.2",
				},
				{
					IPAddress: "103.208.177.11",
					Type:      "public",
				},
			},
			expected: "103.208.177.11",
		},
		{
			networks: []godo.NetworkV4{
				{
					IPAddress: "10.0.0.2",
				},
				{
					IPAddress: "172.16.0.5",
				},
			},
			expected: "",
		},
	}

	for _, testCase := range testCases {
		ipAddr := getPublicIpPort(testCase.networks)

		if !strings.EqualFold(ipAddr, testCase.expected) {
			t.Errorf("Wrong public ip address expected %s actual %s", testCase.expected, ipAddr)
		}
	}
}

func TestGetPrivateIpAddr(t *testing.T) {
	testCases := []struct {
		networks []godo.NetworkV4
		expected string
	}{
		{
			networks: []godo.NetworkV4{
				{
					IPAddress: "10.0.0.2",
					Type:      "private",
				},
				{
					IPAddress: "103.208.177.11",
					Type:      "public",
				},
			},
			expected: "10.0.0.2",
		},
		{
			networks: []godo.NetworkV4{},
			expected: "",
		},
	}

	for _, testCase := range testCases {
		ipAddr := getPrivateIpPort(testCase.networks)

		if !strings.EqualFold(ipAddr, testCase.expected) {
			t.Errorf("Wrong private ip address expected %s actual %s", testCase.expected, ipAddr)
		}
	}
}
