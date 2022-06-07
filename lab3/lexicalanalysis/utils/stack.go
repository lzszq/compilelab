package utils

type Stack struct {
	List []Element
}

func (s *Stack) Init() {
	s.List = make([]Element, 0)
}

func (s *Stack) Len() int {
	return len(s.List)
}

func (s *Stack) IsEmpty() bool {
	return len(s.List) == 0
}

func (s *Stack) Push(x interface{}) {
	s.List = append(s.List, x)
}

func (s *Stack) Pop() Element {
	if len(s.List) <= 0 {
		return nil
	} else {
		ret := s.List[len(s.List)-1]
		s.List = s.List[:len(s.List)-1]
		return ret
	}
}

func (s *Stack) Top() Element {
	if s.IsEmpty() {
		return nil
	} else {
		return s.List[len(s.List)-1]
	}
}

func (s *Stack) Clear() {
	if len(s.List) == 0 {
		return
	}
	for i := 0; i < s.Len(); i++ {
		s.List[i] = nil
	}
	s.List = make([]Element, 0)
}
