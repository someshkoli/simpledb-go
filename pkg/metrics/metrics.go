package metrics

import "fmt"

type Statistics interface {
	Get() map[string]int
	Increment(name string, n int)
	Update(name string, value int)
}

type InternalStatistics struct {
	name  string
	store map[string]int
}

func NewInternalStatistics(name string) *InternalStatistics {
	return &InternalStatistics{
		name:  name,
		store: map[string]int{},
	}
}

func (s *InternalStatistics) Get() map[string]int {
	return s.store
}

func (s *InternalStatistics) buildKey(name string) string {
	return fmt.Sprintf("%s_%s", s.name, name)
}

func (s *InternalStatistics) Increment(name string, n int) {
	key := s.buildKey(name)
	v, ok := s.store[key]
	if !ok {
		s.store[key] = n
	} else {
		s.store[key] = v + n
	}
}

func (s *InternalStatistics) Update(name string, n int) {
	key := s.buildKey(name)
	s.store[key] = n
}
