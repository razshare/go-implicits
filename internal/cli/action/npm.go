package action

import (
	"os/exec"
	"path/filepath"

	"github.com/razshare/go-implicits/files"
	"github.com/razshare/go-implicits/internal/cli/npm"
	"github.com/razshare/go-implicits/tui/npmselect"
)

func Npm(options NpmOptions) (err error) {
	var pkgs []string
	if pkgs, err = npmselect.Send(); err != nil {
		return
	}

	if len(pkgs) == 0 {
		return nil
	}

	var bun string
	if files.IsFile(options.Bun) {
		if bun, err = filepath.Rel(options.App, options.Bun); err != nil {
			return
		}
	} else if bun, err = exec.LookPath(options.Bun); err != nil {
		bun = options.Bun
	}

	return npm.Install(bun, options.App, pkgs...)
}
