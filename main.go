package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	_ "embed"
)

func main() {
	cfg := parseFlags(`GO_NEW_DIR`)
	if cfg.Version {
		fmt.Println(`gonew`, Version)
		return
	}
	logger := newLogger(cfg)
	slog.SetDefault(logger)

	dir := cfg.OutputDir
	next, err := findNextIndex(cfg.Start, cfg.Prefix, dir)
	if err != nil {
		slog.Error(`find next number failed`, `err`, err)
		os.Exit(1)
	}

	err = os.Chdir(dir)
	if err != nil {
		slog.Error(`change dir failed`, `err`, err)
		os.Exit(1)
	}
	slog.Info("change directory", "path", dir)

	name, err := makeDir(cfg.Prefix, next)
	if err != nil {
		slog.Error(`make dir failed`, `err`, err)
		os.Exit(1)
	}
	slog.Info("directory created", "name", name)

	err = os.Chdir(name)
	if err != nil {
		slog.Error(`change dir failed`, `err`, err)
		os.Exit(1)
	}
	slog.Info("change directory", "path", name)

	err = run(cfg.Compiler, "mod", "init", name)
	if err != nil {
		slog.Error(`go mod init failed`, `name`, name, `err`, err)
		os.Exit(1)
	}
	slog.Info("module initialized", "module", name)

	if cfg.Git {
		err = run("git", "init")
		if err != nil {
			slog.Error(`git init failed`, `err`, err)
			os.Exit(1)
		}
		slog.Info(`git init`)
	}

	// dir = filepath.Join(dir, name)
	// filename := filepath.Join(dir, "main.go")
	filename := "main.go"
	err = os.WriteFile(filename, []byte(mainTemplate), 0o644)
	if err != nil {
		slog.Error(`write file failed`, `filename`, filename, `err`, err)
		os.Exit(1)
	}

	// err = exec.Command(cfg.Editor, `-n`,`.`, filename).Run()
	args := strings.Fields(cfg.Editor)
	if len(args) == 0 {
		slog.Error("invalid editor command")
		os.Exit(1)
	}
	cmd := exec.Command(args[0], append(args[1:], ".", filename)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // the editor becomes its own process group, so signals to gonew won’t kill it.
	// err = cmd.Run()
	if err := cmd.Start(); err != nil {
		slog.Error("start editor failed", "err", err)
		os.Exit(1)
	}
	if err := cmd.Wait(); err != nil {
		slog.Error("editor exited with error", "err", err)
		os.Exit(1)
	}
	slog.Info("project created", "name", name, "path", dir+"/"+name)
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func makeDir(prefix string, next uint64) (name string, err error) {
	for {
		name = fmt.Sprintf("%s%d", prefix, next)
		err = os.Mkdir(name, 0o755)
		if err == nil {
			return
		}
		if os.IsExist(err) {
			next++
			continue
		}
		return
	}
}

// findNextIndex finds max next number
// gaps are never reused.
func findNextIndex(start uint64, prefix, dir string) (next uint64, err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		name := e.Name()
		if !strings.HasPrefix(name, prefix) {
			continue
		}

		n, err := strconv.ParseUint(name[len(prefix):], 10, 64)
		if err != nil {
			continue
		}

		if n > next {
			next = n
		}
	}
	return max(start, next+1), nil
}

//go:embed main.tpl
var mainTemplate string
