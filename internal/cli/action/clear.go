package action

import (
	"github.com/razshare/go-implicits/tui/text"
)

func Clear(_ ClearOptions) error {
	text.Clrscr()
	return nil
}
