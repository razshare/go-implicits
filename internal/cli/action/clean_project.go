package action

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/go-implicits/tui/messages"
)

func CleanProject(options CleanProjectOptions) (err error) {
	clean := exec.Command(options.Go, "clean")
	clean.Env = append(os.Environ())
	clean.Stderr = os.Stderr
	clean.Stdout = os.Stdout
	clean.Stdin = os.Stdin
	if err = clean.Run(); err != nil {
		return
	}

	if err = os.RemoveAll(".gen"); err != nil {
		return
	}

	if err = os.RemoveAll(filepath.Join(options.App, "dist")); err != nil {
		return
	}

	if err = os.RemoveAll(filepath.Join(options.App, "node_modules")); err != nil {
		return
	}

	if err = os.RemoveAll(filepath.Join(".vite")); err != nil {
		return
	}

	if err = os.RemoveAll(filepath.Join("lib", "core", "svelte", "ssr", "app")); err != nil {
		return
	}

	if err = Touch(TouchOptions{App: options.App}); err != nil {
		return
	}

	messages.Success("project cleaned")

	return
}
