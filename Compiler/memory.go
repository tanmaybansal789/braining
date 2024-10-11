package Compiler

import (
	"strings"
)

const (
	BF_INC        = "+"
	BF_DEC        = "-"
	BF_PTR_L      = "<"
	BF_PTR_R      = ">"
	BF_WRITE      = "."
	BF_READ       = ","
	BF_OPEN       = "["
	BF_CLOSE      = "]"
	BF_CLEAR      = "[-]"
	BF_BREAKPOINT = "!"
)

type MemoryManager struct {
	UsedMemory  map[int]bool
	Variables   map[string]int
	FreedMemory []int
	NextLoc     int

	pointer int
}

func NewMemoryManager() *MemoryManager {
	return &MemoryManager{
		UsedMemory:  make(map[int]bool),
		Variables:   make(map[string]int),
		FreedMemory: make([]int, 0),
		NextLoc:     0,
		pointer:     0,
	}
}

func (m *MemoryManager) GetMemoryLoc(name string) (int, string) {
	if loc, ok := m.Variables[name]; ok {
		return loc, m.MovePointer(loc)
	}
	loc := m.getNextLoc()
	m.Variables[name] = loc
	m.UsedMemory[loc] = true
	return loc, m.MovePointer(loc)
}

func (m *MemoryManager) IdentifierExists(name string) bool {
	_, ok := m.Variables[name]
	return ok
}

func (m *MemoryManager) GetClearLoc(name string) (int, string) {
	loc, code := m.GetMemoryLoc(name)
	return loc, code + BF_CLEAR
}

func (m *MemoryManager) GetTempLoc() (int, string) {
	loc := m.getNextLoc()
	m.UsedMemory[loc] = true
	return loc, m.MovePointer(loc)
}

func (m *MemoryManager) FreeMemoryLoc(name string) string {
	loc := m.Variables[name]
	output := m.FreeTempLoc(loc)
	delete(m.Variables, name)
	return output
}

func (m *MemoryManager) FreeTempLoc(loc int) string {
	m.FreedMemory = append(m.FreedMemory, loc)
	m.UsedMemory[loc] = false
	out := m.MovePointer(loc)
	out += BF_CLEAR
	return out
}

func (m *MemoryManager) MovePointer(pos int) string {
	var out string
	if m.pointer > pos {
		out = strings.Repeat(BF_PTR_L, m.pointer-pos)
	} else {
		out = strings.Repeat(BF_PTR_R, pos-m.pointer)
	}
	m.pointer = pos
	return out
}

func (m *MemoryManager) getNextLoc() int {
	if len(m.FreedMemory) > 0 {
		loc := m.FreedMemory[0]
		m.FreedMemory = m.FreedMemory[1:]
		return loc
	}
	loc := m.NextLoc
	m.NextLoc++
	return loc
}
