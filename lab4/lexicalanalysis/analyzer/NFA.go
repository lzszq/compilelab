package analyzer

import (
	"lexicalanalysis/utils"
)

type NFA struct {
	Name             string
	Regex            string
	RegexWithContact string
	SuffixRegex      string
	StateNum         int
	StartAndEnd      State
	Graph            [][]Node
	NFA              []Part
}

type Node struct {
	Number int
	Char   rune
}

type State struct {
	StartNum int
	EndNum   int
}

func (nfa *NFA) init(name string, regex string) {
	nfa.Name = name
	nfa.Regex = regex
	nfa.SuffixRegex = ""
	nfa.SuffixRegex = ""
	nfa.StartAndEnd = State{}
	nfa.StateNum = 0
	nfa.Graph = append(nfa.Graph, make([]Node, 0))
}

func (nfa *NFA) priority(char rune) int {
	flag := 0
	switch char {
	case '~':
		flag = 3
	case '`':
		flag = 2
	case '|':
		flag = 1
	}
	return flag
}

func (nfa *NFA) insertContact() {
	for i := 0; i < len(nfa.Regex)-1; i++ {
		nfa.RegexWithContact += string(nfa.Regex[i])
		if nfa.Regex[i] != '[' && nfa.Regex[i+1] != ']' && nfa.Regex[i] != '|' && nfa.Regex[i+1] != '|' && nfa.Regex[i+1] != '~' {
			nfa.RegexWithContact += "`"
		}
	}
	nfa.RegexWithContact += string(nfa.Regex[len(nfa.Regex)-1])
}

func (nfa *NFA) regex2SuffixRegex() {
	var stack utils.Stack
	stack.Init()
	for _, c := range nfa.RegexWithContact {
		if c == '[' {
			stack.Push(c)
		} else if c == ']' {
			for stack.Top().(rune) != '[' {
				nfa.SuffixRegex += string(stack.Top().(rune))
				stack.Pop()
			}
			stack.Pop()
		} else if c == '~' || c == '`' || c == '|' {
		loop:
			for stack.Len() != 0 {
				if nfa.priority(c) <= nfa.priority(stack.Top().(rune)) {
					nfa.SuffixRegex += string(stack.Top().(rune))
					stack.Pop()
				} else {
					break loop
				}
			}
			stack.Push(c)
		} else {
			nfa.SuffixRegex += string(c)
		}
	}
	for stack.Len() != 0 {
		nfa.SuffixRegex += string(stack.Top().(rune))
		stack.Pop()
	}
}

func (nfa *NFA) newState() int {
	nfa.Graph = append(nfa.Graph, make([]Node, 0))
	nfa.StateNum += 1
	return nfa.StateNum
}

func (nfa *NFA) suffixRegex2NFA() {
	var states utils.Stack
	states.Init()
	var startNum, endNum int
loop:
	for _, c := range nfa.SuffixRegex {
		if c != '~' && c != '`' && c != '|' {
			startNum = nfa.newState()
			endNum = nfa.newState()
			states.Push(State{StartNum: startNum, EndNum: endNum})
			nfa.Graph[startNum] = append(nfa.Graph[startNum], Node{Number: endNum, Char: c})
			continue loop
		}
		switch c {
		case '~':
			origin := states.Top().(State)
			states.Pop()
			startNum = nfa.newState()
			endNum = nfa.newState()
			states.Push(State{StartNum: startNum, EndNum: endNum})
			nfa.Graph[startNum] = append(nfa.Graph[startNum], Node{Number: origin.StartNum, Char: utils.EPSILON})
			nfa.Graph[startNum] = append(nfa.Graph[startNum], Node{Number: endNum, Char: utils.EPSILON})
			nfa.Graph[origin.EndNum] = append(nfa.Graph[origin.EndNum], Node{Number: endNum, Char: utils.EPSILON})
			nfa.Graph[origin.EndNum] = append(nfa.Graph[origin.EndNum], Node{Number: origin.StartNum, Char: utils.EPSILON})
		case '`':
			right := states.Top().(State)
			states.Pop()
			left := states.Top().(State)
			states.Pop()
			states.Push(State{StartNum: left.StartNum, EndNum: right.EndNum})
			nfa.Graph[left.EndNum] = append(nfa.Graph[left.EndNum], Node{Number: right.StartNum, Char: utils.EPSILON})
		case '|':
			down := states.Top().(State)
			states.Pop()
			up := states.Top().(State)
			states.Pop()
			startNum = nfa.newState()
			endNum = nfa.newState()
			states.Push(State{StartNum: startNum, EndNum: endNum})
			nfa.Graph[startNum] = append(nfa.Graph[startNum], Node{Number: up.StartNum, Char: utils.EPSILON})
			nfa.Graph[startNum] = append(nfa.Graph[startNum], Node{Number: down.StartNum, Char: utils.EPSILON})
			nfa.Graph[up.EndNum] = append(nfa.Graph[up.EndNum], Node{Number: endNum, Char: utils.EPSILON})
			nfa.Graph[down.EndNum] = append(nfa.Graph[down.EndNum], Node{Number: endNum, Char: utils.EPSILON})
		}
	}
	tmp := states.Top().(State)
	nfa.StartAndEnd = State{StartNum: tmp.StartNum, EndNum: tmp.EndNum}
}

func regex2NFA(name string, regex string) NFA {
	var nfa NFA
	nfa.init(name, regex)
	nfa.insertContact()
	nfa.regex2SuffixRegex()
	nfa.suffixRegex2NFA()
	for i := 1; i <= nfa.StateNum; i++ {
	loop:
		for _, item := range nfa.Graph[i] {
			if item.Char == '@' {
				nfa.NFA = append(nfa.NFA, Part{Src: i, Edge: 'ε', Dst: item.Number})
				continue loop
			}
			nfa.NFA = append(nfa.NFA, Part{Src: i, Edge: item.Char, Dst: item.Number})
		}
	}

	return nfa
}

func GetNFA() []NFA {
	var list []NFA
	list = append(list, regex2NFA("Keyword", "sin|cos|tg|ctg|log|lg|ln"))
	list = append(list, regex2NFA("Operator", "+|-|*|/|^|="))
	list = append(list, regex2NFA("Constant", "[[0|"+utils.NONZERODIGIT+"]["+utils.DIGIT+"]~][.["+utils.DIGIT+"]~["+utils.NONZERODIGIT+"]|@]|PI|E"))
	list = append(list, regex2NFA("Variable", "[_|"+utils.LETTER+"][_|"+utils.LETTER+"|"+utils.DIGIT+"]~"))
	list = append(list, regex2NFA("Delimiter", ";|(|)|\n|\t|,| |?"))
	return list
}
