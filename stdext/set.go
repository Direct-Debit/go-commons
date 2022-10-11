package stdext

type StrSet map[string]struct{}

func (ss StrSet) Add(v string) {
	ss[v] = struct{}{}
}

func (ss StrSet) Has(v string) bool {
	_, ok := ss[v]
	return ok
}

func (ss StrSet) Remove(v string) {
	if _, ok := ss[v]; ok {
		delete(ss, v)
	}
}

func (ss StrSet) ToSlice() []string {
	s := make([]string, len(ss))
	idx := 0
	for v := range ss {
		s[idx] = v
		idx++
	}
	return s
}

type Set[T comparable] struct {
	values map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		values: make(map[T]struct{}),
	}
}

func SetFromSlice[T comparable](s []T) *Set[T] {
	set := &Set[T]{
		values: make(map[T]struct{}),
	}
	for _, v := range s {
		set.values[v] = struct{}{}
	}
	return set
}

func (s *Set[T]) Add(v T) {
	s.values[v] = struct{}{}
}

func (s *Set[T]) Has(v T) bool {
	_, ok := s.values[v]
	return ok
}

func (s *Set[T]) Remove(v T) {
	if _, ok := s.values[v]; ok {
		delete(s.values, v)
	}
}

func (s *Set[T]) ToSlice() []T {
	slice := make([]T, len(s.values))
	idx := 0
	for v := range s.values {
		slice[idx] = v
		idx++
	}
	return slice
}
