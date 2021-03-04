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
	}
	return s
}
