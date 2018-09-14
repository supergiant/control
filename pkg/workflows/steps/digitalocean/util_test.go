package digitalocean

import (
	"context"
	"strings"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
)

type mockKeyService struct {
	key *godo.Key
	err error
}

func (m *mockKeyService) Create(context.Context, *godo.KeyCreateRequest) (*godo.Key, *godo.Response, error) {
	return m.key, nil, m.err
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

func TestCreateKey(t *testing.T) {
	testCases := []struct {
		key *godo.Key
		err error
	}{
		{
			key: &godo.Key{
				1,
				"name",
				"fingerprint",
				"",
			},
			err: nil,
		},
		{
			&godo.Key{},
			errors.New("error create key"),
		},
	}

	for _, testCase := range testCases {
		keyService := &mockKeyService{
			key: testCase.key,
			err: testCase.err,
		}
		_, err := createKey(context.Background(), keyService, testCase.key.Name, testCase.key.PublicKey)

		if err != testCase.err {
			t.Errorf("Unexpected err value expected %v actual %v", testCase.err, err)
		}
	}
}
