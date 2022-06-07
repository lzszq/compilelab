package utils

import "fmt"

type Node struct {
	Data     Element
	Parent   *Node
	Children []*Node
}

func InitNode(data Element) *Node {
	return &Node{Data: data, Parent: nil, Children: make([]*Node, 0)}
}

func (node *Node) Add(n *Node) {
	node.Children = append(node.Children, n)
}

func DFS(root *Node, level int, ret *[]string) {
	if len((*root).Children) == 0 {
		return
	} else {
		for _, i := range (*root).Children {
			if len((*i).Children) != 0 {
				DFS(i, level+1, ret)
			}
			*ret = append(*ret, " ")
			for i := 0; i < level-1; i++ {
				(*ret)[len(*ret)-1] += " "
			}
			(*ret)[len(*ret)-1] += (*i).Data.(string)
		}
	}
}

func DrawAST(root *Node) {
	var ret []string
	DFS(root, 0, &ret)
	max := 0
	for index := range ret {
		if t := len([]rune(ret[index])); max < t {
			max = t
		}
	}
	for i := range ret {
		for j := 0; j < max; j++ {
			ret[i] += " "
		}
	}
	length := len(ret)
	for j := 0; j < max; j++ {
		for i := 0; i < length; i++ {
			fmt.Printf("%c  ", ret[length-i-1][j])
		}
		fmt.Println()
	}
}
