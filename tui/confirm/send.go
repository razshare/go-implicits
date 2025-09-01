package confirm

import (
	"fmt"

	"github.com/razshare/go-implicits/tui/program"
)

func Send(defaultValue bool, message string) (yes bool, err error) {
	var model *Model
	if model, err = program.Run(&Model{Prompt: message, DefaultValue: defaultValue}); err != nil {
		return
	}

	yes = model.Confirmed

	return
}

func Sendf(defaultValue bool, format string, vars ...any) (yes bool, err error) {
	return Send(defaultValue, fmt.Sprintf(format, vars...))
}
