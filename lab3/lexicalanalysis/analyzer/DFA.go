package analyzer

import (
	"lexicalanalysis/utils"
)

func initFlags(nfa NFA) []rune {
	var flags utils.Set
	flags.Init()
	for i := range nfa.Graph {
		for _, j := range nfa.Graph[i] {
			flags.Add(j.Char)
		}
	}
	flags.Delete(utils.EPSILON)
	flags_new := make([]rune, 0)
	for k := range flags.Data {
		flags_new = append(flags_new, k.(rune))
	}
	return flags_new
}

type StateTables struct {
	Data []utils.Set
}

func (stateTables *StateTables) add(content utils.Set) {
	stateTables.Data = append(stateTables.Data, utils.Set{})
	t := len(stateTables.Data) - 1
	stateTables.Data[t].Init()
	stateTables.Data[t].AddSet(content)
}

func (stateTables *StateTables) isExist(q utils.Set) bool {
	flag := false
	for _, i := range stateTables.Data {
		if i.IsSame(q) {
			flag = true
		}
	}
	return flag
}

func generateStateTable(nfa NFA) (StateTables, []rune) {
	flags := initFlags(nfa)
	cnt := 0

	var closure Closure
	closure.initClosure(nfa)
	closure.calculateClosure(nfa.StartAndEnd.StartNum)
	q := closure.Closure

	var stateTables StateTables
	stateTables.add(q)
	cnt += 1

	var workList []utils.Set
	workList = append(workList, q)
	for len(workList) != 0 {
		q := workList[0]
		workList = workList[1:]
		for _, flag := range flags {
			closure.initClosure(nfa)
			for _, i := range closure.move(q, flag) {
				closure.calculateClosure(i)
			}
			t := closure.Closure
			if !stateTables.isExist(t) && t.Len() != 0 {
				workList = append(workList, t)
				stateTables.add(t)
				cnt += 1
			}
		}
	}
	return stateTables, flags
}

type Closure struct {
	IsArrive utils.Set
	Content  NFA
	Closure  utils.Set
}

func (closure *Closure) calculateClosure(StartNum int) {
	closure.IsArrive.Add(StartNum)
	closure.Closure.Add(StartNum)
	for _, item := range closure.Content.Graph[StartNum] {
		if item.Char == utils.EPSILON && !closure.IsArrive.Data[item.Number] {
			closure.calculateClosure(item.Number)
		}
	}
}

func (closure *Closure) move(q utils.Set, flag rune) []int {
	var list []int
	for k := range q.Data {
		for _, i := range closure.Content.Graph[k.(int)] {
			if i.Char == flag {
				list = append(list, i.Number)
			}
		}
	}
	return list
}

func (closure *Closure) initClosure(nfa NFA) {
	closure.Content = nfa
	closure.Closure.Init()
	closure.IsArrive.Init()
}

type DFA struct {
	Name        string
	StateTables StateTables
	Flags       []rune
	StateNum    int
	NonEndNums  []int
	EndNums     []int
	Graph       [][]Node
	DFA         []Part
}

type Part struct {
	Src  int
	Edge rune
	Dst  int
}

func (dfa *DFA) init(nfa NFA) {
	dfa.Name = nfa.Name
	stateTables, flags := generateStateTable(nfa)
	dfa.StateTables = stateTables
	dfa.Flags = flags
	dfa.StateNum = len(stateTables.Data)
	dfa.nfa2dfa(nfa)
	dfa.generateEndNumsAndNonEndNums(nfa)
	dfa.minimizeDFA()
}

func (dfa *DFA) nfa2dfa(nfa NFA) {
	var closure Closure
	dfa.Graph = make([][]Node, dfa.StateNum)
	for index, item := range dfa.StateTables.Data {
		for _, flag := range dfa.Flags {
			closure.initClosure(nfa)
			for _, i := range closure.move(item, flag) {
				closure.calculateClosure(i)
			}
			t := closure.Closure
			for j, tt := range dfa.StateTables.Data {
				if tt.IsSame(t) {
					dfa.Graph[index] = append(dfa.Graph[index], Node{Number: j, Char: flag})
				}
			}
		}
	}
	for index, i := range dfa.Graph {
		for _, item := range i {
			dfa.DFA = append(dfa.DFA, Part{Src: index, Edge: item.Char, Dst: item.Number})
		}
	}
}

func (dfa *DFA) generateEndNumsAndNonEndNums(nfa NFA) {
	endNum := nfa.StartAndEnd.EndNum
	for index, i := range dfa.StateTables.Data {
		if i.Data[endNum] {
			dfa.EndNums = append(dfa.EndNums, index)
		} else {
			dfa.NonEndNums = append(dfa.NonEndNums, index)
		}
	}
}

func move(DFA []Part, src int, edge rune) int {
	for _, i := range DFA {
		if i.Src == src && i.Edge == edge {
			return i.Dst
		}
	}
	return -1
}

func getSourceSet(targetSet utils.Set, char rune, DFA []Part, allNums utils.Set) utils.Set {
	var sourceSet utils.Set
	sourceSet.Init()
	for k := range allNums.Data {
		if targetSet.Data[move(DFA, k.(int), char)] {
			sourceSet.Add(k.(int))
		}
	}
	return sourceSet
}

func (dfa *DFA) minimizeDFA() {
	var P, W []utils.Set
	var NonEndNums, EndNums, allNums utils.Set
	NonEndNums.Init()
	EndNums.Init()
	allNums.Init()
	for _, i := range dfa.NonEndNums {
		NonEndNums.Add(i)
	}
	for _, i := range dfa.EndNums {
		EndNums.Add(i)
	}

	allNums.AddSet(NonEndNums)
	allNums.AddSet(EndNums)

	P = append(P, EndNums, NonEndNums)
	W = append(W, EndNums, NonEndNums)

	for len(W) != 0 {
		A := W[0]
		W = W[1:]

		for _, char := range dfa.Flags {
			X := getSourceSet(A, char, dfa.DFA, allNums)
			var P_ []utils.Set

			for _, Y := range P {
				var S, S1 utils.Set
				S.Init()
				S1.Init()
				S.AddSet(X)
				S.AddSet(Y)
				S1.AddSet(Y)
				for k := range X.Data {
					S1.Delete(k)
				}

				if S.Len() == 0 && S1.Len() == 0 {
					P_ = append(P_, S)
					P_ = append(P_, S1)

					flag := true
					for index, i := range W {
						if Y.IsSame(i) && flag {
							W = append(W[:index], W[index+1:]...)
							W = append(W, S)
							W = append(W, S1)
							flag = false
						}
					}
					if flag {
						if S.Len() <= S1.Len() {
							W = append(W, S)
						} else {
							W = append(W, S1)
						}
					}
				} else {
					P_ = append(P_, Y)
				}
			}
			P = nil
			P = append(P, P_...)
		}
	}
}

type EndTable struct {
	Name    string
	EndNums []int
}

func mergeDFA(dfas []DFA) ([]Part, []EndTable) {
	var newDFA []Part
	var endTables []EndTable
	cnt := 0
	for _, dfa := range dfas {
		for _, part := range dfa.DFA {
			if part.Src == 0 {
				newDFA = append(newDFA, Part{Src: part.Src, Edge: part.Edge, Dst: part.Dst + cnt})
			} else {
				newDFA = append(newDFA, Part{Src: part.Src + cnt, Edge: part.Edge, Dst: part.Dst + cnt})
			}
		}
		endTables = append(endTables, EndTable{Name: dfa.Name, EndNums: dfa.EndNums})
		for i := range endTables[len(endTables)-1].EndNums {
			endTables[len(endTables)-1].EndNums[i] += cnt
		}
		cnt += dfa.StateNum
	}
	return newDFA, endTables
}

func GetDFA(nfas []NFA) ([]DFA, []Part, []EndTable) {
	var dfas []DFA
	for i := range nfas {
		var dfa DFA
		dfa.init(nfas[i])
		dfas = append(dfas, dfa)
	}
	newDFA, endTables := mergeDFA(dfas)
	return dfas, newDFA, endTables
}
