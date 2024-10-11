package AST

type NodeType int

const (
	N_BLOCK NodeType = iota
	N_ASSIGN
	N_ADD
	N_SUB
	N_IF
	N_IFNOT
	N_WHILE
	N_WHILENOT
	N_WRITE
	N_READ
	N_FREE
	N_MACRO
	N_MACRO_CALL
	N_BREAKPOINT
)

type Node interface {
	Type() NodeType
}

type BlockNode struct {
	Nodes []Node
}

type AssignNode struct {
	Left  IdentToken
	Right Token
}

type AddNode struct {
	Left  IdentToken
	Right Token
}

type SubNode struct {
	Left  IdentToken
	Right Token
}

type IfNode struct {
	Id    IdentToken
	Block BlockNode
}

type IfNotNode struct {
	Id    IdentToken
	Block BlockNode
}

type WhileNode struct {
	Id    IdentToken
	Block BlockNode
}

type WhileNotNode struct {
	Id    IdentToken
	Block BlockNode
}

type WriteNode struct {
	Value Token
}

type ReadNode struct {
	Value IdentToken
}

type FreeNode struct {
	Value IdentToken
}

type MacroNode struct {
	Name   IdentToken
	Params []IdentToken
	Block  BlockNode
}

type MacroCallNode struct {
	Name IdentToken
	Args map[IdentToken]Token
}

type BreakpointNode struct{}

func (n *BlockNode) Type() NodeType {
	return N_BLOCK
}

func (n *AssignNode) Type() NodeType {
	return N_ASSIGN
}

func (n *AddNode) Type() NodeType {
	return N_ADD
}

func (n *SubNode) Type() NodeType {
	return N_SUB
}

func (n *IfNode) Type() NodeType {
	return N_IF
}

func (n *IfNotNode) Type() NodeType {
	return N_IFNOT
}

func (n *WhileNode) Type() NodeType {
	return N_WHILE
}

func (n *WhileNotNode) Type() NodeType {
	return N_WHILENOT
}

func (n *WriteNode) Type() NodeType {
	return N_WRITE
}

func (n *ReadNode) Type() NodeType {
	return N_READ
}

func (n *FreeNode) Type() NodeType {
	return N_FREE
}

func (n *MacroNode) Type() NodeType {
	return N_MACRO
}

func (n *MacroCallNode) Type() NodeType {
	return N_MACRO_CALL
}

func (n *BreakpointNode) Type() NodeType {
	return N_BREAKPOINT
}
