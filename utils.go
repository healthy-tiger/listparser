package listparser

type stack struct {
	entries []*ListElement
	len     int
}

func newStack() *stack {
	return &stack{make([]*ListElement, 0), 0}
}

func (s *stack) push(v *ListElement) {
	if len(s.entries) >= s.len+1 {
		s.entries[s.len] = v
	} else {
		s.entries = append(s.entries, v)
	}
	s.len = s.len + 1
}

func (s *stack) pop() *ListElement {
	if s.len <= 0 {
		return nil
	}
	v := s.entries[s.len-1]
	s.len = s.len - 1
	return v
}

func (s *stack) peek() *ListElement {
	if s.len <= 0 {
		return nil
	}
	return s.entries[s.len-1]
}
