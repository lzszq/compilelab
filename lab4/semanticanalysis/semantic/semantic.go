package semantic

import (
	"fmt"
	"math"
	"semanticanalysis/utils"
	"strconv"
)

type Type int32

const (
	ASSIGN    Type = 0
	CALCULATE Type = 1
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

type Node utils.Node

type SemanticAnalyzer struct {
	Root      *utils.Node
	SAType    Type
	Symbol    map[string]float64
	Term      map[string]int
	NonTerm   map[string]int
	GramTable []Pair
	SLRTable  [][]Pair
	GoToTable [][]Pair
}

func Init() SemanticAnalyzer {
	var SA SemanticAnalyzer
	SA.Root = utils.InitNode("S")
	SA.Symbol = make(map[string]float64)
	SA.Symbol["PI"] = math.Pi
	SA.Symbol["E"] = math.E
	SA.Term = GetIDTerm()
	SA.NonTerm = GetIDNonTerm()
	SA.GramTable, SA.SLRTable, SA.GoToTable = GetGramTableAndSLRTableAndGoToTable()
	return SA
}

func Analyzer(SA SemanticAnalyzer, keys, tmps []string) (SemanticAnalyzer, bool) {
	LAStream := convertLAStream(keys, tmps)
	LAStream = append(LAStream, item{Data: "#", ID: "#"})

	var symbolStack utils.Stack
	symbolStack.Init()
	symbolStack.Push("#")

	var stateStack utils.Stack
	stateStack.Init()
	stateStack.Push(0)

	var treeStack utils.Stack
	treeStack.Init()

	curIndex := 0
	parent := utils.InitNode("")
	top := utils.InitNode("")

	var reduceVal string
	var valueStack utils.Stack
	valueStack.Init()
	valueStack.Push("#")
	var reduce string

	for curIndex < len(LAStream) {
		row := stateStack.Top().(int)
		col := SA.Term[LAStream[curIndex].ID]
		action := SA.SLRTable[row][col]
		if action.Right == 44 {
			if SA.SAType == ASSIGN {
				SA.Root = treeStack.Top().(*utils.Node)
			}
			treeStack.Pop()

			return SA, true
		} else if action.Left == "s" {
			symbolStack.Push(LAStream[curIndex].ID)
			stateStack.Push(action.Right)
			valueStack.Push(LAStream[curIndex].Data)
			curIndex += 1
		} else if action.Left == "r" {
			g := SA.GramTable[action.Right]
			parent = utils.InitNode(g.Left)
			t := g.Right
			for t != 0 {
				reduce = symbolStack.Top().(string)
				symbolStack.Pop()
				stateStack.Pop()
				reduceVal = valueStack.Top().(string)
				valueStack.Pop()
				if !treeStack.IsEmpty() && ((treeStack.Top().(*utils.Node)).Data).(string) == reduce {
					top = treeStack.Top().(*utils.Node)
					treeStack.Pop()
					parent.Add(top)
					(*top).Parent = parent
				} else {
					if reduce == "Variable" {
						if value, ok := SA.Symbol[reduceVal]; ok {
							top = utils.InitNodeValue(reduce, value)
						} else {
							top = utils.InitNodeVariable(reduce, reduceVal)
						}
					} else if reduce == "Constant" {
						if reduceVal == "PI" || reduceVal == "E" {
							top = utils.InitNodeValue(reduce, SA.Symbol[reduceVal])
						} else {
							f, _ := strconv.ParseFloat(reduceVal, 64)
							top = utils.InitNodeValue(reduce, f)
						}
					} else {
						top = utils.InitNodeChar(reduce)
						(*top).Function = reduce
					}
					(*parent).Add(top)
					(*top).Parent = parent
				}
				t -= 1
			}
			SA = calculate(SA, parent, action.Right)
			top = parent
			treeStack.Push(top)
			row = stateStack.Top().(int)
			col = SA.NonTerm[g.Left]
			symbolStack.Push(g.Left)
			valueStack.Push(LAStream[curIndex].Data)
			if SA.GoToTable[row][col].Right > 0 {
				stateStack.Push(SA.GoToTable[row][col].Right)
			} else {
				fmt.Printf("%d and %d query GoToTable fail\n", row, col)
				return SA, false
			}
		} else {
			fmt.Printf("%d and %d query SLRTable fail\n", row, col)
			return SA, false
		}
	}
	return SA, false
}

func PrintExpression(SA SemanticAnalyzer) {
	utils.PrintExpression(SA.Root)
}

func calculate(SA SemanticAnalyzer, root *utils.Node, reduceId int) SemanticAnalyzer {
	switch reduceId {
	case 2:
		SA.SAType = ASSIGN
		t3 := *((*root).Children[3])
		t1 := *((*root).Children[1])

		if t1.NodeType == utils.CONSTANT {
			SA.Symbol[t3.Variable] = t1.Value
		}
	case 3:
		SA.Root = (*root).Children[1]
		SA.SAType = CALCULATE
	case 4:
		t0 := *((*root).Children[0])
		t2 := *((*root).Children[2])

		if t0.NodeType == utils.CONSTANT && t2.NodeType == utils.CONSTANT {
			(*root).Value = t0.Value + t2.Value
			(*root).NodeType = utils.CONSTANT
			(*root).DeleteChildren()
		}
	case 5:
		t0 := *((*root).Children[0])
		t2 := *((*root).Children[2])
		if t0.NodeType == utils.CONSTANT && t2.NodeType == utils.CONSTANT {
			(*root).Value = t0.Value - t2.Value
			(*root).NodeType = utils.CONSTANT
			(*root).DeleteChildren()
		}
	case 6, 9, 14:
		t0 := *((*root).Children[0])
		if t0.NodeType == utils.CONSTANT {
			(*root).Value = t0.Value
			(*root).NodeType = utils.CONSTANT
			(*root).DeleteChildren()
		}
	case 7:
		t0 := *((*root).Children[0])
		t2 := *((*root).Children[2])
		if t0.NodeType == utils.CONSTANT && t2.NodeType == utils.CONSTANT {
			(*root).Value = t0.Value * t2.Value
			(*root).NodeType = utils.CONSTANT
			(*root).DeleteChildren()
		}
	case 8:
		t0 := *((*root).Children[0])
		t2 := *((*root).Children[2])
		if t0.NodeType == utils.CONSTANT && t2.NodeType == utils.CONSTANT {
			(*root).Value = t2.Value / t0.Value
			(*root).NodeType = utils.CONSTANT
			(*root).DeleteChildren()
		}
	case 10:
		t2 := *((*root).Children[2])
		t4 := *((*root).Children[4])
		if t2.NodeType == utils.CONSTANT && t4.NodeType == utils.CONSTANT && t2.Value > 0 && t4.Value > 0 {
			(*root).Value = math.Log(t4.Value) / math.Log(t2.Value)
			(*root).NodeType = utils.CONSTANT
			(*root).DeleteChildren()
		}
	case 11:
		t1 := *((*root).Children[1])
		if t1.NodeType == utils.CONSTANT && t1.Value > 0 {
			(*root).Value = math.Log2(t1.Value)
			(*root).NodeType = utils.CONSTANT
			(*root).DeleteChildren()
		}
	case 12:
		t0 := *((*root).Children[0])
		t2 := *((*root).Children[2])
		if t0.NodeType == utils.CONSTANT && t2.NodeType == utils.CONSTANT {
			(*root).Value = math.Pow(t0.Value, t2.Value)
			(*root).NodeType = utils.CONSTANT
			(*root).DeleteChildren()
		}
	case 13:
		t3 := *((*root).Children[3])
		t1 := *((*root).Children[1])
		if t1.NodeType == utils.CONSTANT {
			if t3.Function == "sin" {
				(*root).Value = math.Sin(t1.Value)
			} else if t3.Function == "cos" {
				(*root).Value = math.Cos(t1.Value)
			} else if t3.Function == "tg" && math.Cos(t1.Value) != 0 {
				(*root).Value = math.Tan(t1.Value)
			} else if t3.Function == "ctg" && math.Sin(t1.Value) != 0 {
				(*root).Value = 1 / math.Tan(t1.Value)
			} else if t3.Function == "lg" && t1.Value > 0 {
				(*root).Value = math.Log10(t1.Value)
			} else if t3.Function == "ln" && t1.Value > 0 {
				(*root).Value = math.Log(t1.Value)
			}
			(*root).NodeType = utils.CONSTANT
			(*root).DeleteChildren()
		}
	case 15, 17:
		t0 := *((*root).Children[0])
		if t0.NodeType == utils.CONSTANT {
			(*root).Value = t0.Value
			(*root).NodeType = utils.CONSTANT
			(*root).DeleteChildren()
		}
	case 16, 18:
		t0 := *((*root).Children[0])
		if t0.NodeType == utils.CONSTANT {
			(*root).Value = -t0.Value
			(*root).NodeType = utils.CONSTANT
			(*root).DeleteChildren()
		}
	case 19, 20:
		t0 := *((*root).Children[0])
		if t0.NodeType == utils.CONSTANT {
			(*root).Value = t0.Value
			(*root).NodeType = utils.CONSTANT
			(*root).DeleteChildren()
		}
	case 21:
		t1 := *((*root).Children[1])
		if t1.NodeType == utils.CONSTANT {
			(*root).Value = t1.Value
			(*root).NodeType = utils.CONSTANT
			(*root).DeleteChildren()
		}
	}
	return SA
}
