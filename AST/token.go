package AST

type TokenType int

const (
	T_IDENT TokenType = iota
	T_LIT
)

type Token interface {
	Type() TokenType
}

type IdentToken struct {
	Name string
}

type LitToken struct {
	Value string
}

func (t *IdentToken) Type() TokenType {
	return T_IDENT
}

func (t *LitToken) Type() TokenType {
	return T_LIT
}
