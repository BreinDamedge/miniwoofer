// main file for μWoofer experiment
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"miniwoofer/tui"
)

func main() {
	p := tea.NewProgram(tui.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Woof! Something broke: %v", err)
		os.Exit(1)
	}
}
