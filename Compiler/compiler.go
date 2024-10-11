package Compiler

import (
	"braining/AST"
	"braining/Logging"
	"os"
	"strconv"
	"strings"
)

// Compiler struct that holds the memory manager, logger, AST, and generated code
type Compiler struct {
	memoryManager *MemoryManager
	logger        *Logging.Logger
	Ast           AST.Ast
	Code          string
}

func NewCompiler(ast AST.Ast, logger *Logging.Logger) *Compiler {
	if logger == nil {
		logger = Logging.NewLoggerWithDefaultColors("braining_compiler", Logging.ERROR) // Crash on error by default
	}

	return &Compiler{
		memoryManager: NewMemoryManager(),
		logger:        logger,
		Ast:           ast,
		Code:          "",
	}
}

// ----------------------------------------------------
// Entry Point
// ----------------------------------------------------

// Compile starts compiling from the root of the AST
func (c *Compiler) Compile() {
	c.compileNode(&c.Ast.Root)
}

// ----------------------------------------------------
// Operations
// ----------------------------------------------------

// Increment the value at the given memory location by a specified amount
func (c *Compiler) inc(loc, val int) {
	c.inject(c.memoryManager.MovePointer(loc))
	c.inject(strings.Repeat(BF_INC, val))
}

// Decrement the value at the given memory location by a specified amount
func (c *Compiler) dec(loc, val int) {
	c.inject(c.memoryManager.MovePointer(loc))
	c.inject(strings.Repeat(BF_DEC, val))
}

// Clear the value at the given memory location
func (c *Compiler) clear(loc int) {
	c.inject(c.memoryManager.MovePointer(loc))
	c.inject(BF_CLEAR)
}

// Open and close loops in Brainf***
func (c *Compiler) open() {
	c.inject(BF_OPEN)
}

func (c *Compiler) openAt(loc int) {
	c.inject(c.memoryManager.MovePointer(loc))
	c.open()
}

func (c *Compiler) close() {
	c.inject(BF_CLOSE)
}

func (c *Compiler) closeAt(loc int) {
	c.inject(c.memoryManager.MovePointer(loc))
	c.close()
}

func (c *Compiler) write(loc int) {
	c.inject(c.memoryManager.MovePointer(loc))
	c.inject(BF_WRITE)
}

func (c *Compiler) read(loc int) {
	c.inject(c.memoryManager.MovePointer(loc))
	c.inject(BF_READ)
}

func (c *Compiler) inject(code string) {
	c.Code += code
}

// ----------------------------------------------------
// Memory Management Helper Functions
// ----------------------------------------------------

func (c *Compiler) getLoc(name string) int {
	loc, code := c.memoryManager.GetMemoryLoc(name)
	c.inject(code)
	return loc
}

func (c *Compiler) getClearLoc(name string) int {
	loc, code := c.memoryManager.GetClearLoc(name)
	c.inject(code)
	return loc
}

func (c *Compiler) getTemp() int {
	loc, code := c.memoryManager.GetTempLoc()
	c.inject(code)
	return loc
}

func (c *Compiler) freeTemp(loc int) {
	c.inject(c.memoryManager.FreeTempLoc(loc))
}

func (c *Compiler) free(name string) {
	c.inject(c.memoryManager.FreeMemoryLoc(name))
}

// ----------------------------------------------------
// High Level Operations
// ----------------------------------------------------

func (c *Compiler) copy(from, to int) {
	tmp := c.getTemp()
	c.clear(to)

	c.openAt(from)
	c.inc(tmp, 1)
	c.inc(to, 1)
	c.dec(from, 1)
	c.closeAt(from)

	c.openAt(tmp)
	c.inc(from, 1)
	c.dec(tmp, 1)
	c.closeAt(tmp)

	c.freeTemp(tmp)
}

func (c *Compiler) add(left, right int) {
	tmp := c.getTemp()

	c.copy(right, tmp)

	c.openAt(tmp)
	c.inc(left, 1)
	c.dec(tmp, 1)
	c.closeAt(tmp)

	c.freeTemp(tmp)
}

func (c *Compiler) sub(left, right int) {
	tmp := c.getTemp()

	c.copy(right, tmp)

	c.openAt(tmp)
	c.dec(left, 1)
	c.dec(tmp, 1)
	c.closeAt(tmp)

	c.freeTemp(tmp)
}

// ----------------------------------------------------
// AST Node Compilation
// ----------------------------------------------------

func (c *Compiler) compileNode(node AST.Node) {
	switch node.Type() {
	case AST.N_BLOCK:
		n := node.(*AST.BlockNode)
		for _, child := range n.Nodes {
			c.compileNode(child)
		}

	case AST.N_ASSIGN:
		n := node.(*AST.AssignNode)
		left := c.getClearLoc(n.Left.Name)
		if n.Right.Type() == AST.T_LIT {
			val, err := strconv.Atoi(n.Right.(*AST.LitToken).Value)
			if err != nil {
				err := Logging.InvalidLiteralCompilerError{Value: n.Right.(*AST.LitToken).Value}
				c.logger.Error(err.Error())
			}
			c.inc(left, val)
		} else {
			if !c.memoryManager.IdentifierExists(n.Right.(*AST.IdentToken).Name) {
				err := Logging.InvalidIdentifierCompilerError{Name: n.Right.(*AST.IdentToken).Name}
				c.logger.Error(err.Error())
			}
			c.copy(c.getLoc(n.Right.(*AST.IdentToken).Name), left)
		}

	case AST.N_ADD:
		n := node.(*AST.AddNode)
		left := c.getLoc(n.Left.Name)
		if n.Right.Type() == AST.T_LIT {
			val, err := strconv.Atoi(n.Right.(*AST.LitToken).Value)
			if err != nil {
				err := Logging.InvalidLiteralCompilerError{Value: n.Right.(*AST.LitToken).Value}
				c.logger.Error(err.Error())
			}
			c.inc(left, val)
		} else {
			if !c.memoryManager.IdentifierExists(n.Right.(*AST.IdentToken).Name) {
				err := Logging.InvalidIdentifierCompilerError{Name: n.Right.(*AST.IdentToken).Name}
				c.logger.Error(err.Error())
			}
			c.add(left, c.getLoc(n.Right.(*AST.IdentToken).Name))
		}

	case AST.N_SUB:
		n := node.(*AST.SubNode)
		left := c.getLoc(n.Left.Name)
		if n.Right.Type() == AST.T_LIT {
			val, err := strconv.Atoi(n.Right.(*AST.LitToken).Value)
			if err != nil {
				err := Logging.InvalidLiteralCompilerError{Value: n.Right.(*AST.LitToken).Value}
				c.logger.Error(err.Error())
			}
			c.dec(left, val)
		} else {
			if !c.memoryManager.IdentifierExists(n.Right.(*AST.IdentToken).Name) {
				err := Logging.InvalidIdentifierCompilerError{Name: n.Right.(*AST.IdentToken).Name}
				c.logger.Error(err.Error())
			}
			c.sub(left, c.getLoc(n.Right.(*AST.IdentToken).Name))
		}

	case AST.N_WRITE:
		n := node.(*AST.WriteNode)
		if n.Value.Type() == AST.T_LIT {
			val, err := strconv.Atoi(n.Value.(*AST.LitToken).Value)
			if err != nil {
				err := Logging.InvalidLiteralCompilerError{Value: n.Value.(*AST.LitToken).Value}
				c.logger.Error(err.Error())
			}
			tmp := c.getTemp()
			c.inc(tmp, val)
			c.write(tmp)
			c.freeTemp(tmp)
		} else {
			if !c.memoryManager.IdentifierExists(n.Value.(*AST.IdentToken).Name) {
				err := Logging.InvalidIdentifierCompilerError{Name: n.Value.(*AST.IdentToken).Name}
				c.logger.Error(err.Error())
			}
			c.write(c.getLoc(n.Value.(*AST.IdentToken).Name))
		}

	case AST.N_READ:
		n := node.(*AST.ReadNode)
		if !c.memoryManager.IdentifierExists(n.Value.Name) {
			err := Logging.InvalidIdentifierCompilerError{Name: n.Value.Name}
			c.logger.Error(err.Error())
		}
		c.read(c.getLoc(n.Value.Name))

	case AST.N_FREE:
		n := node.(*AST.FreeNode)
		if !c.memoryManager.IdentifierExists(n.Value.Name) {
			err := Logging.InvalidIdentifierCompilerError{Name: n.Value.Name}
			c.logger.Error(err.Error())
		}
		c.free(n.Value.Name)

	case AST.N_IF:
		n := node.(*AST.IfNode)
		id := c.getLoc(n.Id.Name)
		tmp := c.getTemp()
		c.copy(id, tmp)

		c.openAt(tmp)
		c.compileNode(&n.Block)
		c.clear(tmp)
		c.closeAt(tmp)

		c.freeTemp(tmp)

	case AST.N_IFNOT:
		n := node.(*AST.IfNotNode)
		id := c.getLoc(n.Id.Name)

		tmp := c.getTemp()
		tmp2 := c.getTemp()

		c.copy(id, tmp)
		c.inc(tmp2, 1)

		c.openAt(tmp)
		c.clear(tmp)
		c.dec(tmp2, 1)
		c.closeAt(tmp)

		c.openAt(tmp2)
		c.dec(tmp2, 1)
		c.compileNode(&n.Block)
		c.closeAt(tmp2)

		c.freeTemp(tmp)
		c.freeTemp(tmp2)

	case AST.N_WHILE:
		n := node.(*AST.WhileNode)
		id := c.getLoc(n.Id.Name)

		c.openAt(id)
		c.compileNode(&n.Block)
		c.closeAt(id)

	case AST.N_WHILENOT:
		n := node.(*AST.WhileNotNode)
		id := c.getLoc(n.Id.Name)

		tmp := c.getTemp()
		tmp2 := c.getTemp()

		c.copy(id, tmp)
		c.inc(tmp2, 1)

		c.openAt(tmp)
		c.clear(tmp)
		c.dec(tmp2, 1)
		c.closeAt(tmp)

		c.openAt(tmp2)
		c.dec(tmp2, 1)

		tmp3 := c.getTemp()
		c.inc(tmp3, 1)
		c.openAt(tmp3)

		c.compileNode(&n.Block)

		tmp4 := c.getTemp()
		c.copy(id, tmp4)

		c.openAt(tmp4)
		c.clear(tmp4)
		c.dec(tmp3, 1)

		c.closeAt(tmp4)

		c.closeAt(tmp3)

		c.closeAt(tmp2)

		c.freeTemp(tmp)
		c.freeTemp(tmp2)
		c.freeTemp(tmp3)
		c.freeTemp(tmp4)

	case AST.N_BREAKPOINT:
		c.inject(BF_BREAKPOINT)
	}
}

func (c *Compiler) WriteToFile(path string) {
	f, err := os.Create(path)
	if err != nil {
		c.logger.Error(err.Error())
	}
	defer f.Close()

	n, err := f.WriteString(c.Code)
	if err != nil {
		c.logger.Error(err.Error())
	}

	c.logger.Info("Wrote " + strconv.Itoa(n) + " bytes to " + path)

	err = f.Sync()
	if err != nil {
		c.logger.Error(err.Error())
	}
}
