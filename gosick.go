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
