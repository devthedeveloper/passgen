# passgen

> Cryptographically secure CLI password generator for macOS, Linux & Windows.

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey?style=flat)](https://github.com/devthedeveloper/passgen/releases)
[![Release](https://img.shields.io/github/v/release/devthedeveloper/passgen?style=flat&color=orange)](https://github.com/devthedeveloper/passgen/releases)

```
nJ0aK-6V96F-zMtQs-o8kuh-sZty9   ← passgen -
X7&kP2!qL9mR@wZ#                 ← passgen -length 16
mrwlN_XfZNo_awT1y_wDrPk_JY6X9   ← passgen _
```

---

## Features

- **`crypto/rand`** — genuinely unpredictable, not `math/rand`
- **Quick mode** — `passgen -` or `passgen _` for an instant segmented password
- **Interactive mode** — guided step-by-step prompts when run with no flags
- **Segmented passwords** — configurable segments, length & separator (`-` / `_`)
- **Auto clipboard** — every generated password is copied instantly
- **Zero dependencies** — pure Go stdlib, single static binary

---

## Install

### Option 1 — Download binary (recommended)

Go to **[Releases](https://github.com/devthedeveloper/passgen/releases)** and download the binary for your platform.

**macOS / Linux:**
```sh
# Apple Silicon
tar xzf passgen-1.0.0-macos-arm64.tar.gz
sudo mv passgen /usr/local/bin/

# Intel Mac
tar xzf passgen-1.0.0-macos-intel.tar.gz
sudo mv passgen /usr/local/bin/

# Linux x64
tar xzf passgen-1.0.0-linux-amd64.tar.gz
sudo mv passgen /usr/local/bin/
```

**Windows:**
Extract `passgen-1.0.0-windows-amd64.zip` and move `passgen.exe` to a folder in your `PATH`.

---

### Option 2 — Build from source

Requires [Go 1.21+](https://go.dev/dl/).

```sh
git clone https://github.com/devthedeveloper/passgen.git
cd passgen
go install .
```

Or install directly without cloning:
```sh
go install github.com/devthedeveloper/passgen@latest
```

> Make sure `$GOPATH/bin` (usually `~/go/bin`) is in your `PATH`.
> Add this to your `~/.zshrc` or `~/.bashrc`:
> ```sh
> export PATH="$PATH:$HOME/go/bin"
> ```

---

## Usage

### Quick mode
```sh
passgen -     # instant segmented: nJ0aK-6V96F-zMtQs-o8kuh-sZty9
passgen _     # instant segmented: mrwlN_XfZNo_awT1y_wDrPk_JY6X9
```
Generates 5 segments × 5 chars, uppercase + lowercase + digits, copies to clipboard.

---

### Interactive mode
```sh
passgen
```
Walks you through all options step by step — type, length, character sets, count.

---

### Random passwords
```sh
passgen -length 32
passgen -count 5
passgen -no-symbols
passgen -no-upper -no-symbols          # lowercase + digits only
passgen -exclude "0OIl1"              # strip ambiguous characters
passgen -length 24 -count 3 -no-copy  # no clipboard copy
```

---

### Segmented passwords
```sh
passgen -type segment                              # ab3k-n9xQ-r7mW  (default)
passgen -type segment -segments 4 -seg-length 6   # ab3k2f-n9xQt1-r7mWp0-02yLs8
passgen -type segment -separator _                # ab3k_n9xQ_r7mW
passgen -type segment -no-upper                   # lowercase + digits only
```

---

## All flags

| Flag | Default | Description |
|---|---|---|
| `-type` | `random` | Password type: `random` or `segment` |
| `-length` | `16` | Password length (random mode) |
| `-count` | `1` | Number of passwords to generate |
| `-no-upper` | `false` | Exclude uppercase A–Z |
| `-no-lower` | `false` | Exclude lowercase a–z |
| `-no-digits` | `false` | Exclude digits 0–9 |
| `-no-symbols` | `false` | Exclude symbols `!@#$...` (random only) |
| `-exclude` | `""` | Specific characters to exclude |
| `-segments` | `3` | Number of segments (segment mode) |
| `-seg-length` | `4` | Characters per segment (segment mode) |
| `-separator` | `-` | Segment separator: `-` or `_` |
| `-no-copy` | `false` | Skip copying to clipboard |

---

## Building releases

```sh
./build.sh 1.1.0     # cross-compiles for all platforms into dist/
```

Produces binaries for:
- `macos-arm64` (Apple Silicon)
- `macos-intel` (x86_64)
- `linux-amd64`
- `linux-arm64`
- `windows-amd64`

---

## License

MIT © [devthedeveloper](https://github.com/devthedeveloper)
