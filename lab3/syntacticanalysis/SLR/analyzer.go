package SLR

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
	GramTable, SLRTable, GoToTable := GetGramTableAndSLRTableAndGoToTable()
	LAStream := convertLAStream(keys, tmps)
	LAStream = append(LAStream, item{Data: "#", ID: "#"})

	var symbolStack utils.Stack
	symbolStack.Init()
	symbolStack.Push("#")

	var stateStack utils.Stack
	stateStack.Init()
	stateStack.Push(0)

	root := utils.InitNode("S")
	var treeStack utils.Stack
	treeStack.Init()

	curIndex := 0
	parent := utils.InitNode("")
	top := utils.InitNode("")

	for curIndex < len(LAStream) {
		row := stateStack.Top().(int)
		col := Term[LAStream[curIndex].ID]
		action := SLRTable[row][col]
		if action.Right == 44 {
			root = treeStack.Top().(*utils.Node)
			treeStack.Pop()
			utils.DrawAST(root)
			return true
		} else if action.Left == "s" {
			symbolStack.Push(LAStream[curIndex].ID)
			stateStack.Push(action.Right)
			curIndex += 1
		} else if action.Left == "r" {
			g := GramTable[action.Right]
			parent = utils.InitNode(g.Left)
			for g.Right != 0 {
				reduce := symbolStack.Top().(string)
				symbolStack.Pop()
				stateStack.Pop()
				if !treeStack.IsEmpty() && ((treeStack.Top().(*utils.Node)).Data).(string) == reduce {
					top = treeStack.Top().(*utils.Node)
					treeStack.Pop()
					parent.Add(top)
					(*top).Parent = parent
				} else {
					top = utils.InitNode(reduce)
					parent.Add(top)
					top.Parent = parent
				}
				g.Right -= 1
			}
			top = parent
			treeStack.Push(top)
			row = stateStack.Top().(int)
			col = NonTerm[g.Left]
			symbolStack.Push(g.Left)
			if GoToTable[row][col].Right > 0 {
				stateStack.Push(GoToTable[row][col].Right)
			} else {
				fmt.Printf("%d and %d query GoToTable fail\n", row, col)
				return false
			}
		} else {
			fmt.Printf("%d and %d query SLRTable fail\n", row, col)
			return false
		}
	}
	return false
}
