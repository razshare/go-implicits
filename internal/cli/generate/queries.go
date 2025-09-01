package generate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/go-implicits/files"
	messages2 "github.com/razshare/go-implicits/tui/messages"
	"github.com/razshare/go-implicits/tui/search"
	"github.com/razshare/go-implicits/tui/singleselect"
	spinner2 "github.com/razshare/go-implicits/tui/spinner"
)

func Queries(options QueriesOptions) (err error) {
	if options.SqlcYaml == "" {
		choices := make([]search.Choice, 0)

		if files.IsFile(filepath.Join("lib", "database", "sqlite", "sqlc.yaml")) {
			choices = append(choices, search.Choice{Id: "lib/database/sqlite/sqlc.yaml", Description: "lib/database/sqlite/sqlc.yaml"})
		}

		choices = append(choices, search.Choice{Id: "other", Description: "other"})
		options.SqlcYaml, err = singleselect.Sendf(choices, "where is your sqlc.yaml file located?")
	}

	lib := filepath.Dir(options.SqlcYaml)

	if _, err = exec.LookPath(options.Sqlc); err != nil && !files.IsFile(options.Sqlc) {
		if err = Sqlc(SqlcOptions{Sqlc: options.Sqlc, Platform: options.Platform, Auto: options.Auto}); err != nil {
			return
		}
	}

	yaml := filepath.Join(lib, "sqlc.yaml")

	if !files.IsFile(yaml) {
		return fmt.Errorf("%s not found", yaml)
	}

	spin := spinner2.New("generating queries")

	var sqlc string
	if files.IsFile(options.Sqlc) {
		if sqlc, err = filepath.Rel(lib, options.Sqlc); err != nil {
			return err
		}
	} else if sqlc, err = exec.LookPath(options.Sqlc); err != nil {
		sqlc = options.Sqlc
	}

	go spinner2.Start(spin)
	generate := exec.Command(sqlc, "generate")
	generate.Dir = lib
	generate.Env = append(os.Environ())
	generate.Stderr = os.Stderr
	generate.Stdout = os.Stdout
	generate.Stdin = os.Stdin
	err = generate.Run()
	spinner2.Stop(spin)

	if err != nil {
		return
	}

	messages2.Success(
		"queries generated at database.Queries.*\n",
		lib+"/queries.go",
	)
	messages2.Tip(
		"## usage example\n",
		"func(c *client.Client){\n",
		"    u, _ := database.Queries.FindUsers(c.Request.Context())\n",
		"    send.Json(c, u)\n",
		"}",
	)

	return
}
