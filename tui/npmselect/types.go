package npmselect

import (
	"time"

	"github.com/razshare/go-implicits/internal/cli/npm"
	"github.com/razshare/go-implicits/tui/search"
	"github.com/razshare/go-implicits/tui/viewport"
)

type Model struct {
	Selected  []string
	Prompt    string
	LastQuery string
	Search    *search.Search
	Viewport  *viewport.Viewport
	Packages  []npm.PackageInfo
	Debounce  time.Duration
	Debouncer *time.Timer
	Error     error
	Loading   bool
	Quitting  bool
	Confirmed bool
}

type SearchResultMsg struct {
	Packages []npm.PackageInfo
	Error    error
}

type DebouncedSearchMsg struct {
	Query string
}
