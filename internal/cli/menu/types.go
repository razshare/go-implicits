package menu

import (
	"github.com/razshare/go-implicits/tui/search"
)

type Menu struct {
	Items []Item
}

type Item struct {
	Choice  search.Choice
	Hidden  bool
	Active  func() bool // Active returns true if the user has typed the choice directly in the terminal.
	Handler func() error
}
