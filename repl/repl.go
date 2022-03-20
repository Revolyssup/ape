package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/Revolyssup/ape/compiler"
	"github.com/Revolyssup/ape/lexer"
	"github.com/Revolyssup/ape/parser"
	"github.com/Revolyssup/ape/vm"
)

func PrintParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
func CloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal. Ape says bye!")
		os.Exit(0)
	}()
}
func StartRepl(in io.Reader, out io.Writer) {
	buf := bufio.NewScanner(in)
	CloseHandler()
	for {
		fmt.Printf("\n[APE]>>")
		scanned := buf.Scan()
		if !scanned {
			return
		}

		input := buf.Text()

		l := lexer.New(input)

		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			PrintParserErrors(out, p.Errors())
			continue
		}

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			continue
		}
		machine := vm.New(comp.ByteCode())
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
			continue
		}
		stackTop := machine.StackTop()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")
	}
}
