package main

import (
	"bufio"
	"fmt"
	"github.com/Braden-Griebel/cloxgo/vm"
	"io"
	"os"
)

func main() {
	machine := vm.InitVM()

	if len(os.Args) == 1 {
		repl(&machine)
	} else if len(os.Args) == 2 {
		runFile(&machine, os.Args[1])
	} else {
		_, err := os.Stderr.WriteString("Usage: cloxgo [path]\n")
		if err != nil {
			machine.FreeVM()
			panic(err)
		}
	}
	machine.FreeVM()
}

func repl(machine *vm.VM) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			machine.FreeVM()
			panic(err)
		}
		_ = machine.Interpret(line)
	}
}

func runFile(machine *vm.VM, filename string) {
	program, err := os.ReadFile(filename)
	if err != nil {
		panic("Couldn't read file: " + filename)
	}
	programStr := string(program)
	result := machine.Interpret(programStr)
	if result == vm.INTERPRET_COMPILE_ERROR {
		machine.FreeVM()
		os.Exit(65)
	}
	if result == vm.INTERPRET_RUNTIME_ERROR {
		machine.FreeVM()
		os.Exit(70)
	}

}
