package xray

import (
	"fmt"
)

type HashedSet struct {
	items    map[string]struct{}
	hashHigh uint32
	hashLow  uint32
}

func NewHashedSet() *HashedSet {
	return &HashedSet{
		items:    make(map[string]struct{}),
		hashHigh: 0,
		hashLow:  0,
	}
}

func Djb2Dual(str string) (high, low uint32) {
	var h int32 = 5381
	var l int32 = 5387

	for i := 0; i < len(str); i++ {
		c := int32(str[i])
		h = ((h << 5) + h + c)    // h = h * 33 + char
		l = ((l << 6) + l + c*37) // l = l * 65 + char * 37
	}

	return uint32(h), uint32(l)
}

func (s *HashedSet) Add(str string) {
	if _, exists := s.items[str]; !exists {
		s.items[str] = struct{}{}
		high, low := Djb2Dual(str)
		s.hashHigh ^= high
		s.hashLow ^= low
	}
}

func (s *HashedSet) Delete(str string) {
	if _, exists := s.items[str]; exists {
		delete(s.items, str)
		high, low := Djb2Dual(str)
		s.hashHigh ^= high
		s.hashLow ^= low
	}
}

func (s *HashedSet) Has(str string) bool {
	_, exists := s.items[str]
	return exists
}

func (s *HashedSet) Size() int {
	return len(s.items)
}

func (s *HashedSet) Clear() {
	s.items = make(map[string]struct{})
	s.hashHigh = 0
	s.hashLow = 0
}

func (s *HashedSet) Hash64String() string {
	return fmt.Sprintf("%08x%08x", s.hashHigh, s.hashLow)
}

func (s *HashedSet) Items() []string {
	result := make([]string, 0, len(s.items))
	for item := range s.items {
		result = append(result, item)
	}
	return result
}
