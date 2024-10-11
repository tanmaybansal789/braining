package Parser

import (
	"braining/AST"
	"braining/Lexer"
	"braining/Logging"
	"fmt"
)

type Parser struct {
	lexer      *Lexer.Lexer
	blockStack []*AST.BlockNode
	macros     map[string]*AST.MacroNode
	Ast        AST.Ast
	logger     *Logging.Logger
	pos        int
}

func NewParser(source string, logger *Logging.Logger) *Parser {
	if logger == nil {
		logger = Logging.NewLoggerWithDefaultColors("braining_parser", Logging.ERROR)
	}
	p := &Parser{
		lexer:      Lexer.NewLexer(source, logger),
		blockStack: []*AST.BlockNode{},
		macros:     make(map[string]*AST.MacroNode),
		logger:     logger,
	}
	p.blockStack = append(p.blockStack, &p.Ast.Root)
	return p
}

func (p *Parser) copyToken(t AST.Token, tbl map[AST.IdentToken]AST.Token) AST.Token {
	switch t.Type() {
	case AST.T_IDENT:
		id := t.(*AST.IdentToken)
		if val, ok := tbl[*id]; ok {
			return val
		}
		return id
	case AST.T_LIT:
		return t
	default:
		return nil
	}
}

func parseLit(lit string) AST.Token {
	if len(lit) == 3 && lit[0] == '\'' && lit[2] == '\'' {
		return &AST.LitToken{Value: fmt.Sprintf("%d", lit[1])}
	}
	return &AST.LitToken{Value: lit}
}

func (p *Parser) appendNode(n AST.Node) {
	p.blockStack[len(p.blockStack)-1].Nodes = append(p.blockStack[len(p.blockStack)-1].Nodes, n)
	switch n.Type() {
	case AST.N_BLOCK:
		p.blockStack = append(p.blockStack, n.(*AST.BlockNode))
	case AST.N_IF:
		p.blockStack = append(p.blockStack, &n.(*AST.IfNode).Block)
	case AST.N_IFNOT:
		p.blockStack = append(p.blockStack, &n.(*AST.IfNotNode).Block)
	case AST.N_WHILE:
		p.blockStack = append(p.blockStack, &n.(*AST.WhileNode).Block)
	case AST.N_WHILENOT:
		p.blockStack = append(p.blockStack, &n.(*AST.WhileNotNode).Block)
	case AST.N_MACRO:
		p.blockStack = append(p.blockStack, &n.(*AST.MacroNode).Block)
	}
}

func (p *Parser) popBlock() {
	if len(p.blockStack) < 2 {
		err := Logging.InvalidEndParserError{Line: p.lexer.Line()}
		p.logger.Error(err.Error())
	}
	p.blockStack = p.blockStack[:len(p.blockStack)-1]
}

func (p *Parser) copyNode(n AST.Node, tbl map[AST.IdentToken]AST.Token) AST.Node {
	switch n.Type() {
	case AST.N_BLOCK:
		b := n.(*AST.BlockNode)
		res := &AST.BlockNode{Nodes: make([]AST.Node, len(b.Nodes))}
		for i, node := range b.Nodes {
			res.Nodes[i] = p.copyNode(node, tbl)
		}
		return res
	case AST.N_ASSIGN:
		a := n.(*AST.AssignNode)
		return &AST.AssignNode{
			Left:  *p.copyToken(&a.Left, tbl).(*AST.IdentToken),
			Right: p.copyToken(a.Right, tbl),
		}
	case AST.N_ADD:
		a := n.(*AST.AddNode)
		return &AST.AddNode{
			Left:  *p.copyToken(&a.Left, tbl).(*AST.IdentToken),
			Right: p.copyToken(a.Right, tbl),
		}
	case AST.N_SUB:
		s := n.(*AST.SubNode)
		return &AST.SubNode{
			Left:  *p.copyToken(&s.Left, tbl).(*AST.IdentToken),
			Right: p.copyToken(s.Right, tbl),
		}
	case AST.N_IF:
		i := n.(*AST.IfNode)
		b := p.copyNode(&i.Block, tbl).(*AST.BlockNode)
		return &AST.IfNode{
			Id:    *p.copyToken(&i.Id, tbl).(*AST.IdentToken),
			Block: *b,
		}
	case AST.N_IFNOT:
		i := n.(*AST.IfNotNode)
		b := p.copyNode(&i.Block, tbl).(*AST.BlockNode)
		return &AST.IfNotNode{
			Id:    *p.copyToken(&i.Id, tbl).(*AST.IdentToken),
			Block: *b,
		}
	case AST.N_WHILE:
		w := n.(*AST.WhileNode)
		b := p.copyNode(&w.Block, tbl).(*AST.BlockNode)
		return &AST.WhileNode{
			Id:    *p.copyToken(&w.Id, tbl).(*AST.IdentToken),
			Block: *b,
		}
	case AST.N_WHILENOT:
		w := n.(*AST.WhileNotNode)
		b := p.copyNode(&w.Block, tbl).(*AST.BlockNode)
		return &AST.WhileNotNode{
			Id:    *p.copyToken(&w.Id, tbl).(*AST.IdentToken),
			Block: *b,
		}
	case AST.N_WRITE:
		w := n.(*AST.WriteNode)
		return &AST.WriteNode{Value: p.copyToken(w.Value, tbl)}
	case AST.N_READ:
		r := n.(*AST.ReadNode)
		return &AST.ReadNode{Value: *p.copyToken(&r.Value, tbl).(*AST.IdentToken)}
	case AST.N_FREE:
		f := n.(*AST.FreeNode)
		return &AST.FreeNode{Value: *p.copyToken(&f.Value, tbl).(*AST.IdentToken)}
	case AST.N_MACRO:
		m := n.(*AST.MacroNode)
		b := p.copyNode(&m.Block, tbl).(*AST.BlockNode)
		return &AST.MacroNode{
			Name:   *p.copyToken(&m.Name, tbl).(*AST.IdentToken),
			Params: m.Params,
			Block:  *b,
		}
	case AST.N_MACRO_CALL:
		mc := n.(*AST.MacroCallNode)

		macro := p.macros[mc.Name.Name]
		newTbl := make(map[AST.IdentToken]AST.Token)
		for _, param := range macro.Params {
			newTbl[param] = p.copyToken(mc.Args[param], tbl)
		}

		macroBlock := p.copyNode(&macro.Block, newTbl).(*AST.BlockNode)
		return macroBlock
	case AST.N_BREAKPOINT:
		return &AST.BreakpointNode{}
	}
	return nil
}

func (p *Parser) parseNext() bool {
	t := p.lexer.Advance()
	switch t.Type {
	case Lexer.T_IF:
		id := p.lexer.Advance()
		if id.Type == Lexer.T_NOT {
			id = p.lexer.Advance()
			if id.Type != Lexer.T_IDENT {
				err := Logging.InvalidIdentifierParserError{Name: id.Value, Line: p.lexer.Line()}
				p.logger.Error(err.Error())
			}
			i := AST.IfNotNode{Id: AST.IdentToken{Name: id.Value}}
			p.appendNode(&i)
		} else if id.Type != Lexer.T_IDENT {
			err := Logging.InvalidIdentifierParserError{Name: id.Value, Line: p.lexer.Line()}
			p.logger.Error(err.Error())
		} else {
			i := AST.IfNode{Id: AST.IdentToken{Name: id.Value}}
			p.appendNode(&i)
		}

	case Lexer.T_WHILE:
		id := p.lexer.Advance()
		if id.Type == Lexer.T_NOT {
			id = p.lexer.Advance()
			if id.Type != Lexer.T_IDENT {
				err := Logging.InvalidIdentifierParserError{Name: id.Value, Line: p.lexer.Line()}
				p.logger.Error(err.Error())
			}
			w := AST.WhileNotNode{Id: AST.IdentToken{Name: id.Value}}
			p.appendNode(&w)
		} else if id.Type != Lexer.T_IDENT {
			err := Logging.InvalidIdentifierParserError{Name: id.Value, Line: p.lexer.Line()}
			p.logger.Error(err.Error())
		} else {
			w := AST.WhileNode{Id: AST.IdentToken{Name: id.Value}}
			p.appendNode(&w)
		}

	case Lexer.T_END:
		p.popBlock()

	case Lexer.T_DONE:
		if len(p.blockStack) != 1 {
			err := Logging.InvalidEndParserError{Line: p.lexer.Line()}
			p.logger.Error(err.Error())
		}
		return false

	case Lexer.T_IDENT:
		op := p.lexer.Advance()
		r := p.lexer.Advance()

		var rt AST.Token
		if r.Type == Lexer.T_IDENT {
			rt = &AST.IdentToken{Name: r.Value}
		} else if r.Type == Lexer.T_LIT {
			rt = parseLit(r.Value)
		} else {
			err := Logging.InvalidRightParserError{Line: p.lexer.Line()}
			p.logger.Error(err.Error())
		}

		switch op.Type {
		case Lexer.T_ASSIGN:
			p.appendNode(&AST.AssignNode{
				Left:  AST.IdentToken{Name: t.Value},
				Right: rt,
			})
		case Lexer.T_ADD:
			p.appendNode(&AST.AddNode{
				Left:  AST.IdentToken{Name: t.Value},
				Right: rt,
			})
		case Lexer.T_SUB:
			p.appendNode(&AST.SubNode{
				Left:  AST.IdentToken{Name: t.Value},
				Right: rt,
			})
		default:
			err := Logging.InvalidOperatorParserError{Line: p.lexer.Line()}
			p.logger.Error(err.Error())
		}

	case Lexer.T_WRITE:
		r := p.lexer.Advance()
		var rt AST.Token
		if r.Type == Lexer.T_IDENT {
			rt = &AST.IdentToken{Name: r.Value}
		} else if r.Type == Lexer.T_LIT {
			rt = parseLit(r.Value)
		} else {
			err := Logging.InvalidRightParserError{Line: p.lexer.Line()}
			p.logger.Error(err.Error())
		}
		p.appendNode(&AST.WriteNode{Value: rt})

	case Lexer.T_READ:
		id := p.lexer.Advance()
		if id.Type != Lexer.T_IDENT {
			err := Logging.InvalidIdentifierParserError{Name: id.Value, Line: p.lexer.Line()}
			p.logger.Error(err.Error())
		}
		p.appendNode(&AST.ReadNode{
			Value: AST.IdentToken{Name: id.Value},
		})

	case Lexer.T_FREE:
		id := p.lexer.Advance()
		if id.Type != Lexer.T_IDENT {
			err := Logging.InvalidIdentifierParserError{Name: id.Value, Line: p.lexer.Line()}
			p.logger.Error(err.Error())
		}
		p.appendNode(&AST.FreeNode{
			Value: AST.IdentToken{Name: id.Value},
		})

	case Lexer.T_BREAKPOINT:
		p.appendNode(&AST.BreakpointNode{})

	case Lexer.T_MACRO_BEGIN:
		id := p.lexer.Advance()
		if id.Type != Lexer.T_IDENT {
			err := Logging.InvalidIdentifierParserError{Name: id.Value, Line: p.lexer.Line()}
			p.logger.Error(err.Error())
		}
		prms := p.lexer.Advance()
		if prms.Type != Lexer.T_MACRO_SET_PARAMS {
			err := Logging.InvalidIdentifierParserError{Line: p.lexer.Line()}
			p.logger.Error(err.Error())
		}
		m := AST.MacroNode{Name: AST.IdentToken{Name: id.Value}}
		for p.lexer.Peek().Type != Lexer.T_MACRO_DEFINE {
			paramId := p.lexer.Advance()
			if paramId.Type != Lexer.T_IDENT {
				err := Logging.InvalidIdentifierParserError{Name: paramId.Value, Line: p.lexer.Line()}
				p.logger.Error(err.Error())
			}
			m.Params = append(m.Params, AST.IdentToken{Name: paramId.Value})
		}
		p.lexer.Advance()
		p.appendNode(&m)
		p.macros[m.Name.Name] = &m

	case Lexer.T_MACRO_END:
		if len(p.blockStack) < 2 {
			err := Logging.InvalidEndParserError{Line: p.lexer.Line()}
			p.logger.Error(err.Error())
		}
		p.blockStack = p.blockStack[:len(p.blockStack)-1]

	case Lexer.T_MACRO_CALL:
		id := p.lexer.Advance()
		if id.Type != Lexer.T_IDENT {
			err := Logging.InvalidIdentifierParserError{Name: id.Value, Line: p.lexer.Line()}
			p.logger.Error(err.Error())
		}

		mc := &AST.MacroCallNode{
			Name: AST.IdentToken{Name: id.Value},
			Args: make(map[AST.IdentToken]AST.Token),
		}
		macro, found := p.macros[mc.Name.Name]
		if !found {
			err := Logging.InvalidIdentifierParserError{Name: mc.Name.Name, Line: p.lexer.Line()}
			p.logger.Error(err.Error())
		}

		for i := range macro.Params {
			arg := p.lexer.Advance()
			if arg.Type == Lexer.T_IDENT {
				mc.Args[macro.Params[i]] = &AST.IdentToken{Name: arg.Value}
			} else if arg.Type == Lexer.T_LIT {
				mc.Args[macro.Params[i]] = parseLit(arg.Value)
			} else {
				err := Logging.InvalidLiteralParserError{Line: p.lexer.Line()}
				p.logger.Error(err.Error())
			}
		}

		expandedBlock := p.copyNode(mc, nil).(*AST.BlockNode)
		p.blockStack[len(p.blockStack)-1].Nodes = append(p.blockStack[len(p.blockStack)-1].Nodes, expandedBlock)

	default:
		err := Logging.InvalidIdentifierParserError{Name: t.Value, Line: p.lexer.Line()}
		p.logger.Error(err.Error())
	}
	return true
}

func (p *Parser) Parse() AST.Ast {
	for p.parseNext() {
	}
	return p.Ast
}
