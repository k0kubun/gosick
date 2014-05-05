package main

import (
	"fmt"
	"github.com/GeertJohan/go.linenoise"
	"github.com/k0kubun/gosick/scheme"
	"log"
	"strings"
)

func shellPrompt(indentLevel int) string {
	if indentLevel == 0 {
		return "gosick> "
	} else if indentLevel > 0 {
		return fmt.Sprintf("gosick* %s", strings.Repeat("  ", indentLevel))
	} else {
		panic("Negative indent level")
	}
}

func invokeInteractiveShell() {
	for {
		indentLevel := 0
		expression := ""

		for {
			currentLine, err := linenoise.Line(shellPrompt(indentLevel))
			if err != nil {
				log.Fatal(err)
				return
			}
			if currentLine == "exit" {
				return
			}
			expression += " "
			expression += currentLine

			interpreter := scheme.NewInterpreter(expression)
			indentLevel = interpreter.IndentLevel()
			if indentLevel == 0 {
				// Because the IndentLevel() method changes its reading position,
				// recreate interpreter to initialize the position.
				interpreter = scheme.NewInterpreter(expression)
				interpreter.Eval()
				break
			} else if indentLevel < 0 {
				log.Println("Error: extra close parentheses")
				expression = ""
				indentLevel = 0
			}
		}
	}
}

func main() {
	invokeInteractiveShell()
}
