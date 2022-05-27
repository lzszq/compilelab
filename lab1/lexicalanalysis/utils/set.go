package utils

type Set struct {
	Data map[Element]bool
}

func (s *Set) Init() {
	s.Data = make(map[Element]bool)
}

func (s *Set) GetData() map[Element]bool {
	return s.Data
}

func (s *Set) Add(x interface{}) {
	s.Data[x] = true
}

func (s *Set) AddSet(x Set) {
	for k := range x.Data {
		s.Data[k] = true
	}
}

func (s *Set) IsIn(x Set) bool {
	flag := true
	var tmp Set
	tmp.Init()
	tmp.AddSet(x)
	lastLen := tmp.Len()
	for k := range s.GetData() {
		tmp.Data[k] = true
	}
	currentLen := tmp.Len()
	if lastLen != currentLen {
		flag = false
	}
	return flag
}

func (s *Set) IsSame(x Set) bool {
	flag := true
	for k := range s.Data {
		if !x.Data[k] {
			flag = false
		}
	}
	return flag
}

func (s *Set) Delete(x interface{}) {
	delete(s.Data, x)
}

func (s *Set) Clear() {
	for k := range s.Data {
		delete(s.Data, k)
	}
}

func (s *Set) Len() int {
	return len(s.Data)
}
