package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var initialMsg string
	if len(os.Args) > 1 {
		initialMsg = os.Args[1]
	}
	var p *tea.Program
	if strings.ToLower(initialMsg) == "configure" {
		p = tea.NewProgram(initialConfigureModel())
	} else {
		p = tea.NewProgram(initialChatModel(initialMsg))
	}
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
