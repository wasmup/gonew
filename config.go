package main

import (
	"flag"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	OutputDir  string
	Start      uint64
	Prefix     string
	Compiler   string
	Editor     string
	Git        bool
	Version    bool
	Log        slog.LevelVar
	JSONLogger bool
	TimeFormat
}

var Version = "dev"

func parseFlags(dir string) *Config {
	outDir := os.Getenv(dir)
	if outDir == "" {
		outDir = "."
	}

	cfg := &Config{}
	cfg.Log.Set(slog.LevelError)
	flag.Uint64Var(&cfg.Start, "n", 101, "starting file number")
	flag.BoolVar(&cfg.Version, "version", false, "print version and exit")
	flag.StringVar(&cfg.Prefix, "prefix", "go", "prefix")
	flag.StringVar(&cfg.Compiler, "c", "go", "compiler path")
	flag.StringVar(&cfg.Editor, "e", "code -n", "editor path")
	flag.StringVar(&cfg.OutputDir, "o", outDir, "output dir")
	flag.BoolVar(&cfg.Git, "git", false, "init git too")
	flag.BoolVar(&cfg.JSONLogger, "json", false, "use JSON logger")
	flag.TextVar(&cfg.Log, "log", &cfg.Log, "log level: debug|info|warn|error")

	flag.StringVar((*string)(&cfg.TimeFormat), "log-time", "unix",
		"log time format: unix|rfc3339")

	flag.Parse()

	return cfg
}

type TimeFormat string

const (
	TimeUnix   TimeFormat = "unix"
	TimeRFC333 TimeFormat = "rfc3339"
)

func newLogger(cfg *Config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: &cfg.Log,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key != slog.TimeKey {
				return a
			}

			t := a.Value.Time()

			switch cfg.TimeFormat {
			case TimeUnix:
				return slog.Int64(slog.TimeKey, t.Unix())
			case TimeRFC333:
				return slog.String(slog.TimeKey, t.Format(time.RFC3339))
			default:
				return a
			}
		},
	}

	if cfg.JSONLogger {
		return slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}

	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}
