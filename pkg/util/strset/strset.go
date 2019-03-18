package strset

import "sort"

// Set holds a unique collection of strings.
type Set struct {
	data map[string]struct{}
}

// New returns a ne and initialized Set.
func New() *Set {
	return &Set{
		data: make(map[string]struct{}),
	}
}

// Add adds a string to the set.
func (s *Set) Add(items ...string) {
	if s.data == nil {
		s.data = make(map[string]struct{})
	}
	for _, i := range items {
		s.data[i] = struct{}{}
	}
}

// Contains checks if an item exits in the set.
func (s *Set) Contains(item string) bool {
	if s.data == nil {
		return false
	}
	_, ok := s.data[item]
	return ok
}

// ToSlice is used to represent the set as a slice of strings.
func (s *Set) ToSlice() []string {
	if s == nil {
		return nil
	}
	out := make([]string, len(s.data))
	i := 0
	for k := range s.data {
		out[i] = k
		i++
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
