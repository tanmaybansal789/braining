package Logging

import "fmt"

type ErrorType int

const (
	E_PARSER = iota
	E_COMPILER
)

type Error interface {
	Error() string
	Type() ErrorType
}

// Errors for invalid literals

type InvalidLiteralParserError struct {
	Value string
	Line  int
}

func (e *InvalidLiteralParserError) Error() string {
	return "(PARSER) Invalid Literal: " + e.Value + " at line " + string(e.Line)
}

func (e *InvalidLiteralParserError) Type() ErrorType {
	return E_PARSER
}

type InvalidLiteralCompilerError struct {
	Value string
}

func (e *InvalidLiteralCompilerError) Error() string {
	return "(COMPILER) Invalid Literal: " + e.Value
}

func (e *InvalidLiteralCompilerError) Type() ErrorType {
	return E_COMPILER
}

// Errors for invalid identifiers

type InvalidIdentifierParserError struct {
	Name string
	Line int
}

func (e *InvalidIdentifierParserError) Error() string {
	return fmt.Sprintf("(PARSER) Invalid Identifier: %s at line %d", e.Name, e.Line)
}

func (e *InvalidIdentifierParserError) Type() ErrorType {
	return E_PARSER
}

type InvalidIdentifierCompilerError struct {
	Name string
}

func (e *InvalidIdentifierCompilerError) Error() string {
	return fmt.Sprintf("(COMPILER) Invalid identifier: %s", e.Name)
}

func (e *InvalidIdentifierCompilerError) Type() ErrorType {
	return E_COMPILER
}

// Errors for invalid end statements (only applicable to the parser)

type InvalidEndParserError struct {
	Line int
}

func (e *InvalidEndParserError) Error() string {
	return fmt.Sprintf("(PARSER) Mismatched end/done statement at line %d", e.Line)
}

func (e *InvalidEndParserError) Type() ErrorType {
	return E_PARSER
}

// Errors for invalid operators

type InvalidOperatorParserError struct {
	Line int
}

func (e *InvalidOperatorParserError) Error() string {
	return fmt.Sprintf("(PARSER) Invalid operator at line %d", e.Line)
}

func (e *InvalidOperatorParserError) Type() ErrorType {
	return E_PARSER
}

// Errors for invalid right-hand side of an assignment

type InvalidRightParserError struct {
	Line int
}

func (e *InvalidRightParserError) Error() string {
	return fmt.Sprintf("(PARSER) Invalid right-hand side of an assignment at line %d", e.Line)
}
