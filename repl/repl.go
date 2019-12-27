package repl

import (
	"dvc/compiler"
	"dvc/evaluator"
	"dvc/lexer"
	"dvc/object"
	"dvc/parser"
	"dvc/vm"
	"fmt"
	"github.com/chzyer/readline"
	"io"
	"os"
	"strings"
)

func Start(engine *string) {
	var constants []object.Object
	globals := make([]object.Object, vm.GlobalsSize)

	symbolTable := compiler.NewSymbolTable()
	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	readLine, err := readline.NewEx(&readline.Config{
		Prompt:      ">> ",
		HistoryFile: "/tmp/readline.tmp",
		InterruptPrompt: "^C",
		HistorySearchFold: true,
	})

	if err != nil {
		panic(err)
	}

	defer func() { _ = readLine.Close() }()

	env := object.NewEnvironment()

	for {
		line, err := readLine.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)

		if strings.ToLower(line) == "bye" || strings.ToLower(line) == "exit" {
			fmt.Println("bye")
			os.Exit(0)
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(p.Errors())
			continue
		}

		if *engine == "vm" {

			comp := compiler.NewWithState(symbolTable, constants)
			err = comp.Compile(program)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				continue
			}

			code := comp.Bytecode()
			constants = code.Constants

			machine := vm.NewWithGlobalsStore(code, globals)
			err = machine.Run()
			if err != nil {
				fmt.Printf("Bytecode execution failed: %s\n", err)
				continue
			}

			lastPopped := machine.LastPoppedStackElem()
			fmt.Println(lastPopped.Inspect())
		} else {
			result := evaluator.Eval(program, env)
			if result != nil {
				o := result.(object.Object)
				fmt.Println(o.Inspect())
			}
		}
	}
}

func printParserErrors(errors []string) {
	for _, msg := range errors {
		fmt.Println(msg)
	}
}
