package npm

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/razshare/go-implicits/files"
	messages2 "github.com/razshare/go-implicits/tui/messages"
)

func Install(bun string, app string, pkgs ...string) error {
	if len(pkgs) == 0 {
		return nil
	}

	if !files.IsDirectory(app) {
		return fmt.Errorf("directory %s not found", app)
	}

	ok := 0
	for _, pkg := range pkgs {
		cmd := exec.Command(bun, "add", pkg)
		cmd.Dir = app
		cmd.Env = os.Environ()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err != nil {
			messages2.Error(fmt.Sprintf("failed to install %s: %v", pkg, err))
			continue
		}

		messages2.Success(fmt.Sprintf("installed %s to %s/node_modules", pkg, app))
		ok++
	}

	if ok > 0 {
		messages2.Success(fmt.Sprintf("successfully installed %d package(s) to %s/node_modules", ok, app))
	}

	return nil
}
