package strset

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	s := New()
	require.NotNil(t, s.data)
}

func TestSet_Add(t *testing.T) {
	for _, tc := range []struct {
		name        string
		input       []string
		expectedRes map[string]struct{}
	}{
		{
			name:        "empty input",
			expectedRes: map[string]struct{}{},
		},
		{
			name:  "single input",
			input: []string{"a"},
			expectedRes: map[string]struct{}{
				"a": {},
			},
		},
		{
			name:  "multiple input",
			input: []string{"a", "b", "c"},
			expectedRes: map[string]struct{}{
				"a": {},
				"b": {},
				"c": {},
			},
		},
	} {
		s := New()
		s.Add(tc.input...)

		require.Equalf(t, tc.expectedRes, s.data, "TC: %s", tc.name)
	}
}

func TestSet_Contains(t *testing.T) {
	for _, tc := range []struct {
		name        string
		s           Set
		item        string
		expectedRes bool
	}{
		{
			name: "empty set, empty input",
			s: Set{
				data: map[string]struct{}{},
			},
			item:        "",
			expectedRes: false,
		},
		{
			name: "empty set",
			s: Set{
				data: map[string]struct{}{},
			},
			item:        "a",
			expectedRes: false,
		},
		{
			name: "not found",
			s: Set{
				data: map[string]struct{}{
					"b": {},
				},
			},
			item:        "a",
			expectedRes: false,
		},
		{
			name: "found",
			s: Set{
				data: map[string]struct{}{
					"a": {},
					"b": {},
				},
			},
			item:        "a",
			expectedRes: true,
		},
	} {
		res := tc.s.Contains(tc.item)

		require.Equalf(t, tc.expectedRes, res, "TC: %s", tc.name)
	}
}

func TestSet_ToSlice(t *testing.T) {
	for _, tc := range []struct {
		name        string
		s           *Set
		expectedRes []string
	}{
		{
			name: "nil set",
		},
		{
			name: "empty set",
			s: &Set{
				data: map[string]struct{}{},
			},
			expectedRes: []string{},
		},
		{
			name: "set",
			s: &Set{
				data: map[string]struct{}{
					"a": {},
					"b": {},
				},
			},
			expectedRes: []string{"a", "b"},
		},
	} {
		res := tc.s.ToSlice()

		require.Equalf(t, tc.expectedRes, res, "TC: %s", tc.name)
	}
}
