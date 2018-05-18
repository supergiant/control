package core

import (
	"strings"
	"testing"
)

func TestReleaseConfigAsFlagValue(t *testing.T) {
	testCases := []struct {
		config   map[string]interface{}
		expected string
	}{
		{
			config: map[string]interface{}{
				"nested": map[string]interface{}{
					"key": "value",
				},
				"not-nested": 6,
			},
			expected: "nested.key='value',not-nested=6",
		},
		{
			config: map[string]interface{}{
				"nested": map[string]interface{}{
					"key": "value",
				},
				"akey": 6,
			},
			expected: "akey=6,nested.key='value'",
		},
	}

	for _, testCase := range testCases {
		actual := releaseConfigAsFlagValue(testCase.config, "")

		if !strings.EqualFold(actual, testCase.expected) {
			t.Errorf("Wrong keys string expected %s actual %s", testCase.expected, actual)
		}
	}
}
