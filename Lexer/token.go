package Lexer

type TokenType int

const (
	T_IF TokenType = iota
	T_NOT
	T_WHILE
	T_END
	T_DONE
	T_WRITE
	T_READ
	T_FREE
	T_MACRO_BEGIN
	T_MACRO_SET_PARAMS
	T_MACRO_DEFINE
	T_MACRO_END
	T_MACRO_CALL
	T_BREAKPOINT

	T_IDENT
	T_LIT

	T_ASSIGN
	T_ADD
	T_SUB
)

type TokenPattern string

const (
	// Keywords
	P_IF               TokenPattern = `^if`
	P_NOT              TokenPattern = `^not`
	P_WHILE            TokenPattern = `^while`
	P_END              TokenPattern = `^end`
	P_DONE             TokenPattern = `^done`
	P_WRITE            TokenPattern = `^write`
	P_READ             TokenPattern = `^read`
	P_FREE             TokenPattern = `^free`
	P_MACRO_BEGIN      TokenPattern = `^macro`
	P_MACRO_SET_PARAMS TokenPattern = `^takes`
	P_MACRO_DEFINE     TokenPattern = `^define`
	P_MACRO_END        TokenPattern = `^emcro`
	P_MACRO_CALL       TokenPattern = `^call`
	P_BREAKPOINT       TokenPattern = `^breakpoint`

	// Identifiers and literals
	P_IDENT TokenPattern = `^[a-zA-Z_][a-zA-Z0-9_]*`
	P_LIT   TokenPattern = `^\d+|^'.'`

	// Operators
	P_ASSIGN TokenPattern = `^=`
	P_ADD    TokenPattern = `^\+=`
	P_SUB    TokenPattern = `^-=`
)

var TokenPatternJmp = [...]TokenPattern{P_IF, P_NOT, P_WHILE, P_END, P_DONE, P_WRITE, P_READ, P_FREE, P_MACRO_BEGIN, P_MACRO_SET_PARAMS, P_MACRO_DEFINE, P_MACRO_END, P_MACRO_CALL, P_BREAKPOINT, P_IDENT, P_LIT, P_ASSIGN, P_ADD, P_SUB}

type Token struct {
	Type  TokenType
	Value string
}
