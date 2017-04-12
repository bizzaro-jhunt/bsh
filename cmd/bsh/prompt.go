package main

import (
	"bufio"
	"os"
	"strings"

	fmt "github.com/jhunt/go-ansi"

	"github.com/mattn/go-isatty"
	"golang.org/x/crypto/ssh/terminal"
)

func prompt(label string, hide bool) string {
	if isatty.IsTerminal(os.Stdin.Fd()) {
		fmt.Fprintf(os.Stderr, label)
		if hide {
			b, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
			fmt.Fprintf(os.Stderr, "\n")
			return string(b)
		}
	}

	s, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.TrimSuffix(s, "\n")
}

func confirm(label string) bool {
	s := prompt(label, false)
	switch strings.ToLower(s) {
	case "yes", "y", "yup", "sure", "yeah", "yea":
		return true
	default:
		return false
	}
}
