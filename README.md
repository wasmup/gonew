# gonew

gonew creates numbered Go project directories, runs go mod init, optionally initialises git, writes a starter main.go, and opens your editor.


It is designed for developers who frequently create small, sequential Go projects (e.g., exercises, experiments, daily coding tasks).

---

## Why This Exists

This tool is especially useful for:

- coding exercises
- daily Go practice
- kata directories
- quick prototypes
- sequential experiment folders

## ✨ Features

- ✅ Automatically finds the next available numbered directory
- ✅ Directory numbers are monotonically increasing (Existing gaps are never reused).
- ✅ Runs `go mod init`
- ✅ Optional `git init`
- ✅ Creates a starter `main.go`
- ✅ Opens your editor automatically
- ✅ Structured logging with `slog`
- ✅ Supports JSON or text logs
- ✅ Configurable timestamp format
- ✅ zero dependencies (stdlib only)

---

## Requirements

- Go 1.21+

## 📦 Installation

```bash
go install github.com/wasmup/gonew@latest

# go install -trimpath -ldflags "-s -w -X main.Version=$(git describe --tags --always)" github.com/wasmup/gonew@latest

# locally:
go install -trimpath -ldflags "-s -X main.Version=$(git describe --tags --always)"
```

Or build locally:

```bash
git clone https://github.com/wasmup/gonew.git
cd gonew
go build -trimpath -ldflags "-X main.Version=$(git describe --tags --always)"
```

---

## 🚀 Simple Usage

```bash
gonew
```

## Help

```bash
gonew --help
```

### Example

```bash
gonew -prefix day -n 101 -git
```

---

## ⚙️ Flags

| Flag | Default | Description |
|------|---------|------------|
| `-o` | `$GO_NEW_DIR` or `.` | Output directory |
| `-prefix` | `go` | Directory name prefix |
| `-n` | `101` | Starting number |
| `-c` | `go` | Compiler/Go binary |
| `-e` | `code` | Editor command |
| `-git` | `false` | Initialize git repository |
| `-json` | `false` | Use JSON logger |
| `-log` | `error` | Log level: `debug`, `info`, `warn`, `error` |
| `-log-time` | `unix` | Time format: `unix`, `rfc3339` |
| `-version` | `false` | Print version and exit |

---

## 🔢 Numbering Behavior

`gonew` scans existing directories matching:

```
<prefix><number>
```

Example:

```
go1
go2
go5
```

If you run:

```bash
gonew
```

It will create:

```
go6
```

> Gaps are never reused.

The next number is calculated as:
```
next = max(start, existing+1)
```

---

## 📝 Generated `main.go`

Each new project contains:

```go
package main

import "fmt"

func main() {
	fmt.Println("Hi")
}
```

---

## 🧾 Logging

Uses Go's structured logger (`log/slog`).

### Text logger (default)

```bash
gonew -log debug
```

### JSON logger

```bash
gonew -json -log info
```

### Time formats

- `unix` (default)
- `rfc3339`

```bash
gonew -log-time rfc3339
```

---

## 🌍 Environment Variables

### `GO_NEW_DIR`

If set, it becomes the default for `-o`.

```bash
export GO_NEW_DIR=~/projects
gonew
```

---

## 🛠 Typical Workflow

For daily coding practice:

```bash
gonew -prefix day -git
```

Produces:

```
day1/
day2/
day3/
...
```

Each ready with:

- `go.mod`
- `main.go`
- Git repository (optional)
- Editor opened

---

## 📌 Version

Print version:

```bash
gonew -version
```

Version is injected at build time.

---

## 🧠 Design Goals

- Minimal dependencies (stdlib only)
- Deterministic directory numbering
- Fast execution
- Developer workflow automation
- Clean structured logging

---

## 📄 License

MIT
