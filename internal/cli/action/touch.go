package action

import (
	"os"
	"path/filepath"

	"github.com/razshare/go-implicits/files"
)

func Touch(options TouchOptions) (err error) {
	touch := func(name string) (err error) {
		dir := filepath.Dir(name)

		if !files.IsDirectory(dir) {
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {
				return
			}
		}

		var file *os.File
		if file, err = os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0666); err != nil {
			return
		}

		if err = file.Close(); err != nil {
			return
		}

		return
	}

	if err = os.MkdirAll(filepath.Join(options.App, "dist"), os.ModePerm); err != nil {
		return
	}

	if err = touch(filepath.Join(options.App, "dist", "app.server.js")); err != nil {
		return
	}

	if err = touch(filepath.Join(options.App, "dist", "client", "index.html")); err != nil {
		return
	}

	if files.IsDirectory(filepath.Join("lib", "core", "svelte", "ssr")) && !files.IsDirectory(filepath.Join("lib", "core", "svelte", "ssr", "app")) {
		if err = files.CopyDirectory(filepath.Join(options.App, "dist"), filepath.Join("lib", "core", "svelte", "ssr", "app", "dist")); err != nil {

		}
	}

	return
}
