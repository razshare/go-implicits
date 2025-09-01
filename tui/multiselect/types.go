package multiselect

import (
	"github.com/razshare/go-implicits/tui/search"
	"github.com/razshare/go-implicits/tui/viewport"
)

type Model struct {
	Selected []string
	Prompt   string
	Search   *search.Search
	Viewport *viewport.Viewport
}
