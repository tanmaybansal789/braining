package main

import (
	"braining/Compiler"
	"braining/Parser"
	"os"
)

func main() {
	f, err := os.ReadFile("test2.br")
	if err != nil {
		panic(err)
	}

	src := string(f)

	p := Parser.NewParser(src, nil)
	a := p.Parse()

	a.Display()

	c := Compiler.NewCompiler(a, nil)
	c.Compile()

	c.WriteToFile("test.b")
}
