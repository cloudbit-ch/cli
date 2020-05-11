package output

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
)

const (
	AnsiReset     = "\033[0m"
	AnsiBold      = "\u001B[1m"
	AnsiItalic    = "\u001B[3m"
	AnsiUnderline = "\u001B[4m"
)

type AnsiColor int

const (
	AnsiBlack AnsiColor = 30
	AnsiRed
	AnsiGreen
	AnsiYellow
	AnsiBlue
	AnsiMagenta
	AnsiCyan
	AnsiWhite

	AnsiBackground AnsiColor = 10
	AnsiBright     AnsiColor = 60
)

type Output struct {
	Writer       io.Writer
	EnableColors bool
}

func NewConsoleOutput(writer *os.File) *Output {
	return &Output{
		Writer:       writer,
		EnableColors: terminal.IsTerminal(int(writer.Fd())),
	}
}

func (o *Output) AnsiSequence(sequence string) *Output {
	if o.EnableColors {
		_, _ = fmt.Fprintf(o.Writer, sequence)
	}
	return o
}

func (o *Output) Color(color AnsiColor) *Output {
	return o.AnsiSequence(fmt.Sprintf("\033[%dm", color))
}

func (o *Output) Printf(format string, a ...interface{}) *Output {
	_, _ = fmt.Fprintf(o.Writer, format, a...)
	return o
}

func (o *Output) Print(a ...interface{}) *Output {
	_, _ = fmt.Fprint(o.Writer, a...)
	return o
}

func (o *Output) Println(a ...interface{}) *Output {
	_, _ = fmt.Fprintln(o.Writer, a...)
	return o
}

func (o *Output) Reset() *Output {
	return o.AnsiSequence(AnsiReset)
}

func (o *Output) Bold(format string, a ...interface{}) *Output {
	return o.AnsiSequence(AnsiBold).
		Printf(format, a...).
		Reset()
}

func (o *Output) Errorf(format string, a ...interface{}) *Output {
	return o.Color(AnsiRed).
		Printf(format, a...).
		Reset()
}
