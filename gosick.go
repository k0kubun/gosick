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
	FileName   string   `short:"f" long:"file" description:"interpret selected scheme source file"`
	Expression []string `short:"e" long:"expression" description:"excecute given expression"`
	DumpAST    bool     `short:"a" long:"ast" default:"false" description:"whether leaf nodes are plotted"`
}

func main() {
	options := new(Options)
	if _, err := flags.Parse(options); err != nil {
		return
	}

	if len(options.FileName) > 0 {
		executeSourceCode(options)
	} else if len(options.Expression) > 0 {
		executeExpression(strings.Join(options.Expression, " "), options.DumpAST)
	} else {
		invokeInteractiveShell(options)
	}
}

func executeSourceCode(options *Options) {
	buffer, err := ioutil.ReadFile(options.FileName)
	if err != nil {
		log.Fatal(err)
	}

	executeExpression(string(buffer), options.DumpAST)
}

func executeExpression(expression string, dumpAST bool) {
	interpreter := scheme.NewInterpreter(expression)
	interpreter.PrintResult(dumpAST)
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
			linenoise.AddHistory(currentLine)
			expression += " "
			expression += currentLine

			interpreter := scheme.NewInterpreter(expression)
			indentLevel = interpreter.IndentLevel()
			if indentLevel == 0 {
				executeExpression(expression, options.DumpAST)
				break
			} else if indentLevel < 0 {
				fmt.Println("*** ERROR: extra close parentheses")
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
