package action

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/go-implicits/files"
	"github.com/razshare/go-implicits/tui/messages"
	spinner2 "github.com/razshare/go-implicits/tui/spinner"
)

func Update(options UpdateOptions) (err error) {
	if err = Touch(TouchOptions{App: options.App}); err != nil {
		return
	}

	spin := spinner2.New("updating go dependencies")

	go spinner2.Start(spin)
	get := exec.Command(options.Go, "get", "-u", "./...")
	get.Env = append(os.Environ())
	get.Stderr = os.Stderr
	get.Stdout = os.Stdout
	get.Stdin = os.Stdin
	err = get.Run()
	spinner2.Stop(spin)

	if err != nil {
		return
	}

	var bun string
	if files.IsFile(options.Bun) {
		if bun, err = filepath.Rel(options.App, options.Bun); err != nil {
			return
		}
	} else if bun, err = exec.LookPath(options.Bun); err != nil {
		bun = options.Bun
	}

	pretty := exec.Command(bun, "update")
	pretty.Dir = options.App
	pretty.Env = append(os.Environ())
	pretty.Stderr = os.Stderr
	pretty.Stdout = os.Stdout
	pretty.Stdin = os.Stdin
	if err = pretty.Run(); err != nil {
		return
	}

	messages.Success("project dependencies updated")

	return
}
