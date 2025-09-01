package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/razshare/go-implicits/files"
	"github.com/razshare/go-implicits/tui/confirm"
	"github.com/razshare/go-implicits/tui/messages"
	"github.com/razshare/go-implicits/tui/spinner"
)

//go:embed version.mirror
var version string

func init() {
	version = strings.TrimSpace(version)
}

func fetch() {
	spin := spinner.New(fmt.Sprintf("fetching frizzante %s", version))
	go spinner.Start(spin)
	defer spinner.Stop(spin)

	if files.IsDirectory(fmt.Sprintf("frizzante-%s", version)) {
		var overwrite bool
		var err error
		if overwrite, err = confirm.Sendf(true, "frizzante-%s already exists. Overwrite?", version); err != nil {
			messages.Fatal(err)
		}
		if !overwrite {
			messages.Info("skipping download")
		} else if err = os.Remove(fmt.Sprintf("frizzante-%s", version)); err != nil {
			messages.Fatal(err)
		}
	}
	if !files.IsDirectory(fmt.Sprintf("frizzante-%s", version)) {
		if files.IsFile(fmt.Sprintf("frizzante-%s.zip", version)) {
			var overwrite bool
			var err error
			if overwrite, err = confirm.Sendf(true, "frizzante-%s.zip already exists. Overwrite?", version); err != nil {
				messages.Fatal(err)
			}
			if !overwrite {
				messages.Info("skipping download")
			} else if err = os.Remove(fmt.Sprintf("frizzante-%s.zip", version)); err != nil {
				messages.Fatal(err)
			}
		}
	}
	url := fmt.Sprintf("https://github.com/razshare/frizzante/archive/refs/tags/%s.zip", version)
	if err := files.DownloadFile(url, fmt.Sprintf("frizzante-%s.zip", version)); err != nil {
		messages.Fatal(err)
	}
	if err := files.UnzipFile(fmt.Sprintf("frizzante-%s.zip", version), "."); err != nil {
		messages.Fatal(err)
	}
	if err := os.Rename(
		fmt.Sprintf("frizzante-%s", strings.Replace(version, "v", "", 1)),
		fmt.Sprintf("frizzante-%s", version),
	); err != nil {
		return
	}
	if err := os.RemoveAll(fmt.Sprintf("frizzante-%s.zip", version)); err != nil {
		messages.Fatal(err)
	}
}

func fixmod() {
	spin := spinner.New("fixing go.mod")
	go spinner.Start(spin)
	defer spinner.Stop(spin)

	name := filepath.Join(fmt.Sprintf("frizzante-%s", version), "internal", "project", "go.mod")
	var data []byte
	var err error
	if data, err = os.ReadFile(name); err != nil {
		messages.Fatal(err)
	}
	var builder strings.Builder
	lines := strings.Split(string(data), "\n")
	count := len(lines)
	for i := 0; i < count; i++ {
		line := lines[i]
		if strings.Contains(line, "=>") {
			// ignore replacements
			continue
		}

		if strings.Contains(line, "github.com/razshare/go-implicits v") {
			builder.WriteString("\tgithub.com/razshare/go-implicits " + version + "\n")
			continue
		}

		builder.WriteString(line + "\n")
	}
	if err = os.WriteFile(name, []byte(builder.String()), os.ModePerm); err != nil {
		messages.Fatal(err)
	}
}

func zip() {
	spin := spinner.New("archiving project")
	go spinner.Start(spin)
	defer spinner.Stop(spin)

	if err := files.ZipDirectory(
		filepath.Join(fmt.Sprintf("frizzante-%s", version), "internal", "project"),
		filepath.Join(fmt.Sprintf("frizzante-%s", version), "internal", "project.zip"),
	); err != nil {
		messages.Fatal(err)
	}
}

func update() {
	spin := spinner.New("updating")
	go spinner.Start(spin)
	defer spinner.Stop(spin)

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = fmt.Sprintf("frizzante-%s", version)
	cmd.Env = append(os.Environ())
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		messages.Fatal(err)
	}
}

func install() {
	spin := spinner.New(fmt.Sprintf("installing frizzante %s", version))
	go spinner.Start(spin)
	defer spinner.Stop(spin)

	cmd := exec.Command("go", "install", ".")
	cmd.Dir = fmt.Sprintf("frizzante-%s", version)
	cmd.Env = append(os.Environ())
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		messages.Fatal(err)
	}
}

func clean() {
	if err := os.RemoveAll(fmt.Sprintf("frizzante-%s", version)); err != nil {
		messages.Fatal(err)
	}
}

func main() {
	fetch()
	fixmod()
	zip()
	update()
	install()
	clean()
	messages.Successf("frizzante %s installed", version)
}
