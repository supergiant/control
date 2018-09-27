package awssdk

import "testing"

func TestNew(t *testing.T) {
	testCases := []struct{
		region string
		key string
		token string
		secret string
		expectedError error
	}{
		{
			region: "",
			expectedError:ErrEmptyRegion,
		},
		{
			region: "us-west-1",
			key: "key",
			expectedError: ErrInvalidCreds,
		},
		{
			region: "us-west-1",
			secret: "key",
			expectedError: ErrInvalidCreds,
		},
		{
			region: "us-west-1",
			key: "key",
			secret: "key",
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		sdk, err := New(testCase.region, testCase.key, testCase.secret, testCase.token)

		if err != testCase.expectedError {
			t.Errorf("expected error %v actual %v", testCase.expectedError, err)
		}

		if testCase.expectedError == nil && sdk == nil {
			t.Errorf("sdk should not be nil")
		}
	}
}
