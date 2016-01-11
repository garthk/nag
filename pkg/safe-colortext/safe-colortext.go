package sc

import (
	"github.com/daviddengcn/go-colortext"
	"github.com/mattn/go-isatty"
	"os"
)

func IsStdoutColorSafe() bool {
	term := os.Getenv("TERM")
	return term != "dumb" && isatty.IsTerminal(os.Stdout.Fd())
}

func Foreground(cl ct.Color, bright bool) {
	if IsStdoutColorSafe() {
		ct.Foreground(cl, bright)
	}
}

func ResetColor() {
	if IsStdoutColorSafe() {
		ct.ResetColor()
	}
}
