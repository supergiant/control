package amazon

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/sgerrors"
)

func TestFindOutboundIPCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := FindOutboundIP(ctx, func() (string, error) {
		return "", sgerrors.ErrNotFound
	})

	if err != context.Canceled {
		t.Errorf("Expected error %v actual %v", context.Canceled, err)
	}
}

func TestFindOutboundIPDeadline(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Nanosecond*1)
	time.Sleep(time.Nanosecond * 2)

	_, err := FindOutboundIP(ctx, func() (string, error) {
		return "", sgerrors.ErrNotFound
	})

	if err != context.DeadlineExceeded {
		t.Errorf("Expected error %v actual %v", context.DeadlineExceeded, err)
	}
}

func TestFindOutboundIPSuccess(t *testing.T) {
	expectedAddr := "10.20.30.40"

	addr, _ := FindOutboundIP(context.Background(), func() (string, error) {
		return expectedAddr, nil
	})

	if addr != expectedAddr {
		t.Errorf("Expected addr %v actual %v", expectedAddr, addr)
	}
}

func TestFindOutboundIP(t *testing.T) {
	testCases := []struct {
		description string

		awsUrl  string
		awsResp []byte
		awsCode int

		awsReadErr error

		myExtUrl  string
		myExtResp []byte
		myExtCode int

		expectedIP string
		hasErr     bool
	}{
		{
			description: "aws success",
			awsUrl:      serviceURLs[0],
			awsCode:     http.StatusOK,
			awsResp:     []byte(`10.20.30.40`),
			awsReadErr:  io.EOF,
			expectedIP:  "10.20.30.40",
		},
		{
			description: "aws failure",
			awsUrl:      serviceURLs[0],
			awsCode:     http.StatusGatewayTimeout,

			myExtUrl:   serviceURLs[1],
			myExtCode:  http.StatusOK,
			myExtResp:  []byte(`11.20.30.40`),
			expectedIP: "11.20.30.40",
		},
		{
			description: "aws bad body",
			awsUrl:      serviceURLs[0],
			awsCode:     http.StatusGatewayTimeout,
			awsReadErr:  errors.New("error"),

			myExtUrl:   serviceURLs[1],
			myExtCode:  http.StatusOK,
			myExtResp:  []byte(`11.20.30.40`),
			expectedIP: "11.20.30.40",
		},
		{
			description: "failure",
			awsUrl:      serviceURLs[0],
			awsCode:     http.StatusGatewayTimeout,

			myExtUrl:  serviceURLs[1],
			myExtCode: http.StatusGatewayTimeout,
			hasErr:    true,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.description)
		httpmock.Activate()
		httpmock.RegisterResponder(http.MethodGet, testCase.awsUrl,
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewBytesResponse(testCase.awsCode,
					testCase.awsResp)

				if testCase.awsCode != http.StatusOK {
					return resp, errors.New("error")
				}

				return resp, nil
			})

		httpmock.RegisterResponder(http.MethodGet, testCase.myExtUrl,
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewBytesResponse(testCase.myExtCode,
					testCase.myExtResp)

				if testCase.myExtCode != http.StatusOK {
					return resp, errors.New("error")
				}

				return resp, nil
			})

		ip, err := findOutBoundIP()
		httpmock.DeactivateAndReset()

		if testCase.hasErr && err == nil {
			t.Errorf("Error must not be nil")
		}

		if testCase.expectedIP != ip {
			t.Errorf("Wrong ip expected %s actual %s",
				testCase.expectedIP, ip)
		}
	}
}
