package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"
)

func runStatus(opt Opt, command string, args []string) {
	_, t := targeting(opt.Config)
	info, err := t.GetInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	fmt.Printf("@G{Config}\n")
	fmt.Printf("  %-10s @Y{%s}\n", "", opt.Config)
	fmt.Printf("\n")
	fmt.Printf("@G{Director}\n")
	fmt.Printf("  %-10s @Y{%s}\n", "Name", info.Name)
	fmt.Printf("  %-10s @Y{%s}\n", "URL", t.URL)
	fmt.Printf("  %-10s @Y{%s}\n", "Version", info.Version)
	fmt.Printf("  %-10s @Y{%s}\n", "User", t.Username)
	fmt.Printf("  %-10s @Y{%s}\n", "UUID", info.UUID)
	fmt.Printf("  %-10s @Y{%s}\n", "CPI", info.CPI)
}
