package singleselect

import (
	"github.com/razshare/go-implicits/tui/search"
	"github.com/razshare/go-implicits/tui/viewport"
)

// Model defines single selection options
type Model struct {
	Prompt   string
	Selected string
	Search   *search.Search
	Viewport *viewport.Viewport
}
