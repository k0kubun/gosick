package main

import (
	"fmt"
	"github.com/GeertJohan/go.linenoise"
	"github.com/jessevdk/go-flags"
	"github.com/k0kubun/gosick/scheme"
	"io/ioutil"
	"log"
	"strings"
)

type Options struct {
	FileName string `short:"f" long:"file" description:"interpret selected scheme source file"`
	DumpAST  bool   `short:"a" long:"ast" default:"false" description:"whether leaf nodes are plotted"`
}

func main() {
	options := new(Options)
	if _, err := flags.Parse(options); err != nil {
		return
	}

	if len(options.FileName) > 0 {
		executeSourceCode(options)
	} else {
		invokeInteractiveShell(options)
	}
}

func executeSourceCode(options *Options) {
	buffer, err := ioutil.ReadFile(options.FileName)
	if err != nil {
		log.Fatal(err)
	}

	interpreter := scheme.NewInterpreter(string(buffer))
	interpreter.Eval(options.DumpAST)
}

func invokeInteractiveShell(options *Options) {
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
				interpreter.Eval(options.DumpAST)
				break
			} else if indentLevel < 0 {
				log.Println("Error: extra close parentheses")
				expression = ""
				indentLevel = 0
			}
		}
	}
}

func shellPrompt(indentLevel int) string {
	if indentLevel == 0 {
		return "gosick> "
	} else if indentLevel > 0 {
		return fmt.Sprintf("gosick* %s", strings.Repeat("  ", indentLevel))
	} else {
		panic("Negative indent level")
	}
}
