package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	Number   = 0
	Operator = 1
)

type Node struct {
	Type  int
	Value string
	Left  *Node
	Right *Node
}

// input: 1 + 4 - 2
// result:
//     		-
//        /   \
//    	 +     2
//      / \
//     1   4
func getAst(expr string) *Node {

	operator := make(map[string]int)
	operator["+"] = Operator
	operator["-"] = Operator

	nodeList := make([]Node, 0)
	var root *Node

	expr = strings.Trim(expr, " ")
	words := strings.Split(expr, " ")
	for _, word := range words {

		var node Node

		if _, ok := operator[word]; ok {
			node.Type = Operator
		} else {
			node.Type = Number
		}
		node.Value = word
		nodeList = append(nodeList, node)
	}

	for i := 0; i < len(nodeList); i++ {
		if root == nil {
			root = &nodeList[i]
			continue
		}

		switch nodeList[i].Type {
		case Operator:
			nodeList[i].Left = root
			root = &nodeList[i]
		case Number:
			root.Right = &nodeList[i]
		}
	}

	return root
}

func getResult(node *Node) string {
	switch node.Type {
	case Number:
		return node.Value
	case Operator:
		return calc(getResult(node.Left), getResult(node.Right), node.Value)
	}
	return ""
}

func calc(left, right string, operator string) string {
	leftVal, _ := TransToInt(left)
	rightVal, _ := TransToInt(right)
	val := 0
	switch operator {
	case "+":
		val = leftVal + rightVal
	case "-":
		val = leftVal - rightVal
	}
	return TransToString(val)
}

func main() {
	expr := `1 + 4 - 2 + 100 - 20 + 12 `
	//expr := ` 1 + 4 `

	ast := getAst(expr)
	result := getResult(ast)

	fmt.Println(result)
}

func TransToString(data interface{}) (res string) {
	val := reflect.ValueOf(data)
	return strconv.FormatInt(val.Int(), 10)
}

func TransToInt(data interface{}) (res int, err error) {
	return strconv.Atoi(strings.TrimSpace(data.(string)))
}
