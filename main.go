package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/kyawphyothu/momo/program"
)

func main() {

	m := program.InitialModel()
	p := tea.NewProgram(m)
	_, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
