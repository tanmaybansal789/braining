package AST

import (
	"fmt"
	"strings"
)

type Ast struct {
	Root BlockNode
}

func (a *Ast) Display() {
	fmt.Println("Program")
	displayNode(&a.Root, 0)
}

func getTokenString(token Token) string {
	switch token.Type() {
	case T_IDENT:
		return token.(*IdentToken).Name
	case T_LIT:
		return token.(*LitToken).Value
	default:
		return "Unknown Token"
	}
}

func displayNode(node Node, indent int) {
	indentation := strings.Repeat("  ", indent)
	switch node.Type() {
	case N_BLOCK:
		fmt.Println(indentation + "{")
		n := node.(*BlockNode)
		for _, child := range n.Nodes {
			displayNode(child, indent+1)
		}
		fmt.Println(indentation + "}")
	case N_ASSIGN:
		n := node.(*AssignNode)
		fmt.Println(indentation + "AssignNode")
		fmt.Println(indentation + "  Left: " + n.Left.Name)
		fmt.Println(indentation + "  Right: " + getTokenString(n.Right))
	case N_ADD:
		n := node.(*AddNode)
		fmt.Println(indentation + "AddNode")
		fmt.Println(indentation + "  Left: " + n.Left.Name)
		fmt.Println(indentation + "  Right: " + getTokenString(n.Right))

	case N_SUB:
		n := node.(*SubNode)
		fmt.Println(indentation + "SubNode")
		fmt.Println(indentation + "  Left: " + n.Left.Name)
		fmt.Println(indentation + "  Right: " + getTokenString(n.Right))
	case N_IF:
		n := node.(*IfNode)
		fmt.Println(indentation + "IfNode")
		fmt.Println(indentation + "  Id: " + n.Id.Name)
		displayNode(&n.Block, indent+1)
	case N_IFNOT:
		n := node.(*IfNotNode)
		fmt.Println(indentation + "IfNotNode")
		fmt.Println(indentation + "  Id: " + n.Id.Name)
		displayNode(&n.Block, indent+1)
	case N_WHILE:
		n := node.(*WhileNode)
		fmt.Println(indentation + "WhileNode")
		fmt.Println(indentation + "  Id: " + n.Id.Name)
		displayNode(&n.Block, indent+1)
	case N_WHILENOT:
		n := node.(*WhileNotNode)
		fmt.Println(indentation + "WhileNotNode")
		fmt.Println(indentation + "  Id: " + n.Id.Name)
		displayNode(&n.Block, indent+1)
	case N_WRITE:
		n := node.(*WriteNode)
		fmt.Println(indentation + "WriteNode")
		fmt.Println(indentation + "  Value/Id: " + getTokenString(n.Value))
	case N_READ:
		n := node.(*ReadNode)
		fmt.Println(indentation + "ReadNode")
		fmt.Println(indentation + "  Id: " + n.Value.Name)
	case N_FREE:
		n := node.(*FreeNode)
		fmt.Println(indentation + "FreeNode")
		fmt.Println(indentation + "  Id: " + n.Value.Name)
	case N_MACRO:
		n := node.(*MacroNode)
		fmt.Println(indentation + "MacroNode")
		fmt.Println(indentation + "  Name: " + n.Name.Name)
		for _, param := range n.Params {
			fmt.Println(indentation + "  Param: " + param.Name)
		}
		displayNode(&n.Block, indent+1)
	case N_MACRO_CALL:
		n := node.(*MacroCallNode)
		fmt.Println(indentation + "MacroCallNode")
		fmt.Println(indentation + "  Name: " + n.Name.Name)
		for _, arg := range n.Args {
			fmt.Println(indentation + "  Arg: " + getTokenString(arg))
		}
	case N_BREAKPOINT:
		fmt.Println(indentation + "BreakpointNode")
	default:
		fmt.Println(indentation + "Unknown Node")
	}
}
