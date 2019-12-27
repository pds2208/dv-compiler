package main

import (
	"dvc/repl"
	"flag"
	"fmt"
)

var engine = flag.String("engine", "vm",
	"use the virtual machine (vm) or interpreper (eval) engine")

func main() {

	flag.Parse()

	fmt.Println()
	fmt.Print("Derived Variables Language REPL - ")

	if *engine == "vm" {
		fmt.Println("[Virtual Machine]")
	} else {
		fmt.Println("[Interpreter]")
	}

	fmt.Println()
	fmt.Println("Enter commands. Type exit or bye to quit")
	fmt.Println()
	repl.Start(engine)
}
