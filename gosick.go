package main

import (
	"bufio"
	"fmt"
	"github.com/k0kubun/gosick/scheme"
	"log"
	"os"
)

func printShellPrompt(indentLevel int) {
	if indentLevel == 0 {
		fmt.Print("gosick> ")
	} else if indentLevel > 0 {
		fmt.Print("gosick* ")

		for length := indentLevel; length > 0; length-- {
			fmt.Print("  ")
		}
	} else {
		panic("Negative indent level")
	}
}

func invokeInteractiveShell() {
	commandLine := bufio.NewReader(os.Stdin)

	for {
		indentLevel := 0
		expression := ""

		for {
			printShellPrompt(indentLevel)

			currentLine, err := commandLine.ReadString('\n')
			if err != nil {
				log.Fatal(err)
				return
			}
			expression += currentLine

			interpreter := scheme.NewInterpreter(expression)
			if indentLevel = interpreter.IndentLevel(); indentLevel == 0 {
				interpreter.Eval()
				break
			}
		}
	}
}

func main() {
	invokeInteractiveShell()
}
