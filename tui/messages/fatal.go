package messages

import (
	"fmt"
	"os"
	"strings"

	"github.com/razshare/go-implicits/tui/config"
)

func Fatal(args ...any) {
	length := len(args)
	entries := make([]string, length)
	for i := 0; i < length; i++ {
		entries[i] = fmt.Sprintf("%s", args[i])
	}

	Status("ERROR", strings.Join(entries, ""), config.Colors.Error, "233", config.Colors.Error)
	os.Exit(1)
}

func Fatalf(format string, vars ...any) {
	Fatal(fmt.Sprintf(format, vars...))
}
