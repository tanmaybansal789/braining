package Lexer

import (
	"braining/Logging"
	"regexp"
	"unicode"
)

type Lexer struct {
	source           []rune
	pos              int
	line             int
	compiledPatterns []*regexp.Regexp
	logger           *Logging.Logger
}

const COMMENT_SYMBOL = '|'

func NewLexer(source string, logger *Logging.Logger) *Lexer {
	compiledPatterns := make([]*regexp.Regexp, len(TokenPatternJmp))

	if logger == nil {
		logger = Logging.NewLoggerWithDefaultColors("braining_lexer", Logging.ERROR)
	}

	for i, pattern := range TokenPatternJmp {
		compiledPatterns[i] = regexp.MustCompile(string(pattern))
	}

	return &Lexer{source: []rune(source), pos: 0, compiledPatterns: compiledPatterns, logger: logger}
}

func (l *Lexer) matchToken(pattern int) (bool, string) {
	compiled := l.compiledPatterns[pattern]
	match := compiled.FindString(string(l.source[l.pos:]))
	return match != "", match
}

func (l *Lexer) passWhitespace() {
	for l.pos < len(l.source) && unicode.IsSpace(l.source[l.pos]) {
		if l.source[l.pos] == '\n' {
			l.line++
		}
		l.pos++
	}
}

func (l *Lexer) passCommentsAndWhitespace() {
	for l.pos < len(l.source) {
		if unicode.IsSpace(l.source[l.pos]) {
			if l.source[l.pos] == '\n' {
				l.line++
			}
			l.pos++
			continue
		}

		if l.source[l.pos] == COMMENT_SYMBOL {
			l.pos++
			for l.pos < len(l.source) && l.source[l.pos] != COMMENT_SYMBOL {
				if l.source[l.pos] == '\n' {
					l.line++
				}
				l.pos++
			}
			if l.pos < len(l.source) && l.source[l.pos] == COMMENT_SYMBOL {
				l.pos++
			}
			continue
		}

		break
	}
}

func (l *Lexer) Advance() Token {
	l.passCommentsAndWhitespace()

	if l.pos >= len(l.source) {
		return Token{T_DONE, string(P_DONE)}
	}

	for i := range len(TokenPatternJmp) {
		if ok, match := l.matchToken(i); ok {
			if TokenType(i) == T_DONE {
				l.pos = len(l.source)
				return Token{T_DONE, match}
			}
			l.pos += len(match)
			return Token{TokenType(i), match}
		}
	}

	var unrec string
	for l.pos < len(l.source) && !unicode.IsSpace(l.source[l.pos]) {
		unrec += string(l.source[l.pos])
		l.pos++
	}
	err := Logging.InvalidIdentifierParserError{Name: unrec, Line: l.line}
	l.logger.Error(err.Error())
	return Token{T_DONE, string(P_DONE)}
}

func (l *Lexer) Peek() Token {
	pos := l.pos
	line := l.line
	token := l.Advance()
	l.pos = pos
	l.line = line
	return token
}

func (l *Lexer) Line() int {
	return l.line
}
