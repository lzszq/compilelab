package analyzer

import (
	"fmt"
)

func GetType(endTables []EndTable, t int) string {
	name := ""
	for _, i := range endTables {
		for _, j := range i.EndNums {
			if t == j {
				name = i.Name
			}
		}
	}
	return name
}

type GobData struct {
	Nfas      []NFA
	Dfas      []DFA
	NewDFA    []Part
	EndTables []EndTable
}

func GetNFADFA() GobData {
	nfas := GetNFA()
	dfas, newDFA, endTables := GetDFA(nfas)
	data := GobData{
		Nfas:      nfas,
		Dfas:      dfas,
		NewDFA:    newDFA,
		EndTables: endTables,
	}
	return data
}

func checkKey(key string, tmp string) bool {
	flag := false
	if key == "Variable" && len(tmp) > 32 {
		fmt.Printf("parse fail, Variable is too long\n\n")
		flag = false
	} else if key == "Constant" && tmp[0] == '0' && len(tmp) > 1 {
		if tmp[1] >= '0' && tmp[1] != '.' {
			fmt.Printf("parse fail, Constant exist 0 before positive number\n\n")
			flag = false
		} else if tmp[1] == '.' {
			// fmt.Printf("{%q, %q}\n", key, tmp)
			flag = true
		}
	} else if key != "" {
		// fmt.Printf("{%q, %q}\n", key, tmp)
		flag = true
	} else {
		fmt.Printf("parse fail\n\n")
		flag = false
	}
	return flag
}

func Match(test string, data GobData) ([]string, []string) {
	// fmt.Println(data.EndTables)
	fmt.Printf("\ntest string: %q\n", test)
	fmt.Printf("lexical analysis: ")
	keys := make([]string, 0)
	tmps := make([]string, 0)
	tmp := ""
	t := 0
	key := ""
	for _, c := range test {
		flag := true
		for _, item := range data.NewDFA {
			if item.Src == t && item.Edge == c && flag {
				t = item.Dst
				// fmt.Printf("%d -> %q -> %d\n", item.Src, item.Edge, item.Dst)
				flag = false
				tmp += string(c)
			}
		}
		if flag {
			key = GetType(data.EndTables, t)
			if !checkKey(key, tmp) {
				return nil, nil
			}
			keys = append(keys, key)
			tmps = append(tmps, tmp)
			tmp = ""
			t = 0
			for _, item := range data.NewDFA {
				if item.Src == t && item.Edge == c && flag {
					t = item.Dst
					// fmt.Printf("%d -> %q -> %d\n", item.Src, item.Edge, item.Dst)
					flag = false
					tmp += string(c)
				}
			}
		}
	}
	key = GetType(data.EndTables, t)
	if !checkKey(key, tmp) {
		return nil, nil
	}
	keys = append(keys, key)
	tmps = append(tmps, tmp)
	fmt.Printf("parse success\n\n")
	return keys, tmps
}
