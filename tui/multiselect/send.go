package multiselect

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/razshare/go-implicits/tui/program"
	"github.com/razshare/go-implicits/tui/search"
	"github.com/razshare/go-implicits/tui/viewport"
)

func Send(choices []search.Choice, message string) (selected []string, err error) {
	// Initialize the search input
	input := textinput.New()
	input.Width = 80
	var model *Model
	model, err = program.Run(&Model{
		Prompt:   message,
		Viewport: &viewport.Viewport{Visible: 6},
		Selected: make([]string, 0),
		Search: &search.Search{
			Choices:  choices,
			Filtered: choices,
			Input:    input,
		},
	})

	if err != nil {
		return
	}

	selected = model.Selected

	return
}

func Sendf(choices []search.Choice, format string, vars ...any) (selected []string, err error) {
	return Send(choices, fmt.Sprintf(format, vars...))
}
