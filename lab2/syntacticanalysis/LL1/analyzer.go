package LL1

import (
	"fmt"
	"syntacticanalysis/utils"
)

type item struct {
	Data string
	ID   string
}

func convertLAStream(keys, tmps []string) []item {
	var LAStream []item
	for i := range keys {
		if keys[i] == "Variable" || keys[i] == "Constant" {
			LAStream = append(LAStream, item{Data: tmps[i], ID: keys[i]})
		} else if tmps[i] != "\t" && tmps[i] != "\n" && tmps[i] != " " {
			LAStream = append(LAStream, item{Data: tmps[i], ID: tmps[i]})
		}
	}
	return LAStream
}

func Analyzer(keys, tmps []string) bool {
	Term := GetIDTerm()
	NonTerm := GetIDNonTerm()
	PredictTable := GetPredictTable()
	LAStream := convertLAStream(keys, tmps)
	LAStream = append(LAStream, item{Data: "#", ID: "#"})

	var stack utils.Stack
	stack.Init()
	stack.Push("#")
	stack.Push("S")
	var stackTree utils.Stack
	stackTree.Init()
	root := utils.InitNode("S")
	stackTree.Push(utils.InitNode("#"))
	stackTree.Push(root)
	curIndex := 0

loop:
	for !stack.IsEmpty() && curIndex < len(LAStream) {
		str := stack.Top().(string)
		stack.Pop()
		node := stackTree.Top().(*utils.Node)
		stackTree.Pop()
		if _, ok := NonTerm[str]; ok {
			next := PredictTable[NonTerm[str]][Term[LAStream[curIndex].ID]]
			if len(next) == 0 {
				for i := 0; i <= curIndex; i++ {
					fmt.Printf("%s", LAStream[i].Data)
				}
				fmt.Printf("\nsyntax error 1: %s and %s parse fail\n", str, LAStream[curIndex].Data)
				return false
			} else {
				if next[0] == "ε" {
					continue loop
				}
				for _, i := range next {
					t := utils.InitNode(i)
					node.Add(t)
					stack.Push(i)
					stackTree.Push(t)
				}
			}
		} else if _, ok := Term[str]; ok {
			if str == LAStream[curIndex].ID {
				fmt.Printf("Match: %s\n", LAStream[curIndex].Data)
				curIndex += 1
			} else {
				for i := 0; i <= curIndex; i++ {
					fmt.Printf("%s", LAStream[i].Data)
				}
				fmt.Printf("\nsyntax error 2: Term %s  parse fail\n", str)
				return false
			}
		} else {
			for i := 0; i <= curIndex; i++ {
				fmt.Printf("%s", LAStream[i].Data)
			}
			fmt.Printf("\nsyntax error 3: %s is not Term or NonTerm\n", str)
			return false
		}
	}
	utils.DrawAST(root)
	return true
}

