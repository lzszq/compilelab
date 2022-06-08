package utils

import "fmt"

type Type int32

const (
	VARIABLE Type = 0
	CONSTANT Type = 1
	OTHERS   Type = 2
)

type Node struct {
	Data     Element
	Value    float64
	Variable string
	Function string
	Parent   *Node
	NodeType Type
	Children []*Node
}

func InitNode(data Element) *Node {
	return &Node{Data: data, Value: 0, Variable: "", Parent: nil, NodeType: OTHERS, Children: make([]*Node, 0)}
}

func InitNodeChar(data Element) *Node {
	return &Node{Data: data, Value: 0, Variable: "", Parent: nil, NodeType: OTHERS, Children: make([]*Node, 0)}
}

func InitNodeValue(data Element, value float64) *Node {
	return &Node{Data: data, Value: value, Variable: "", Parent: nil, NodeType: CONSTANT, Children: make([]*Node, 0)}
}
func InitNodeVariable(data Element, variable string) *Node {
	return &Node{Data: data, Value: 0, Variable: variable, Parent: nil, NodeType: VARIABLE, Children: make([]*Node, 0)}
}

func (node *Node) Add(n *Node) {
	node.Children = append(node.Children, n)
}

func (node *Node) DeleteChildren() {
	node.Children = nil
}

func DFS(root *Node, ret *[]string) {
	if root == nil {
		return
	} else if len((*root).Children) == 0 {
		if (*root).NodeType == VARIABLE {
			*ret = append(*ret, (*root).Variable)
		} else if (*root).NodeType == CONSTANT {
			*ret = append(*ret, fmt.Sprintf("%f", (*root).Value))
		} else if (*root).NodeType == OTHERS {
			*ret = append(*ret, (*root).Function)
		}
	}
	for _, i := range (*root).Children {
		DFS(i, ret)
	}

}

func PrintExpression(root *Node) {
	var ret []string
	DFS(root, &ret)
	fmt.Printf("value: \"")
	for i := 1; i <= len(ret); i++ {
		fmt.Printf("%s", ret[len(ret)-i])
	}
	fmt.Println("\"")
}
