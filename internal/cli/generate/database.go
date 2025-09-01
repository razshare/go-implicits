package generate

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/razshare/go-implicits/files"
	"github.com/razshare/go-implicits/tui/confirm"
	messages2 "github.com/razshare/go-implicits/tui/messages"
	"github.com/razshare/go-implicits/tui/search"
	"github.com/razshare/go-implicits/tui/singleselect"
	spinner2 "github.com/razshare/go-implicits/tui/spinner"
)

func Database(options DatabaseOptions) (err error) {
	var choice string
	if choice, err = singleselect.Send([]search.Choice{{Id: "sqlite"}}, "what type of database would you like to setup?"); err != nil {
		return
	}

	choice = strings.ToLower(choice)

	lib := filepath.Join("lib", "database", choice)

	if files.IsDirectory(lib) {
		if !options.Auto {
			var overwrite bool
			if overwrite, err = confirm.Sendf(true, "%s already exists. Overwrite?", lib); err != nil {
				return
			}

			if !overwrite {
				messages2.Infof("skipping %s", lib)
				return nil
			}
		}

		if err = os.RemoveAll(lib); err != nil {
			return
		}
	}

	if err = Copy(CopyOptions{
		From: "lib/database/" + choice,
		To:   lib,
		Auto: options.Auto,
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	if choice == "sqlite" {
		spin := spinner2.New("adding github.com/mattn/go-sqlite3")

		go spinner2.Start(spin)
		install := exec.Command(options.Go, "get", "github.com/mattn/go-sqlite3")
		install.Env = append(os.Environ())
		install.Stderr = os.Stderr
		install.Stdout = os.Stdout
		install.Stdin = os.Stdin
		err = install.Run()
		spinner2.Stop(spin)

		if err != nil {
			return
		}

		spin = spinner2.New("updating go dependencies")

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

		messages2.Success("sqlite database is ready")

		if strings.Contains(strings.ToLower(options.Generate), "queries") {
			var queries bool
			if queries, err = confirm.Send(true, "would you like to also generate your queries?"); err != nil {
				return
			}

			yaml := "lib/core/database/sqlite/sqlc.yaml"

			//if choice == "sqlite" {
			//	yaml = "lib/core/database/sqlite/sqlc.yaml"
			//}

			if queries {
				if err = Queries(QueriesOptions{
					Auto:     options.Auto,
					Sqlc:     options.Sqlc,
					Platform: options.Platform,
					SqlcYaml: yaml,
				}); err != nil {
					return
				}
			}
		}
	}

	return
}
