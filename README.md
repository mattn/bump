# bump

Bump semantic version in any file using a regular expression pattern.

Inspired by [gobump](https://github.com/x-motemen/gobump), but works with any file format.

## Installation

```
go install github.com/mattn/bump@latest
```

## Usage

```
bump (major|minor|patch|up|set <version>|show) -f <file> -p <pattern> [-w]
```

### Commands

| Command | Description |
|---|---|
| `major` | Bump major version up (e.g. 1.2.3 → 2.0.0) |
| `minor` | Bump minor version up (e.g. 1.2.3 → 1.3.0) |
| `patch` | Bump patch version up (e.g. 1.2.3 → 1.2.4) |
| `up` | Bump up with interactive prompt |
| `set <version>` | Set exact version (no increments) |
| `show` | Only show the current version |

### Flags

| Flag | Description |
|---|---|
| `-f <file>` | Target file (required) |
| `-p <pattern>` | Regexp pattern with a capture group for the version (required) |
| `-w` | Write result to file instead of stdout |

## Examples

Suppose you have a source file `version.go`:

```go
package main

const version = "1.2.3"
```

```bash
# Show current version
bump show -f version.go -p 'version\s*=\s*"(\d+\.\d+\.\d+)"'
# Output: 1.2.3

# Bump patch version (prints modified file to stdout)
bump patch -f version.go -p 'version\s*=\s*"(\d+\.\d+\.\d+)"'

# Bump minor version and write to file
bump minor -w -f version.go -p 'version\s*=\s*"(\d+\.\d+\.\d+)"'

# Interactive prompt
bump up -w -f version.go -p 'version\s*=\s*"(\d+\.\d+\.\d+)"'

# Set exact version
bump set 2.0.0 -w -f version.go -p 'version\s*=\s*"(\d+\.\d+\.\d+)"'
```

Works with any file format — JSON, YAML, TOML, Markdown, etc.:

```bash
# package.json
bump patch -w -f package.json -p '"version":\s*"(\d+\.\d+\.\d+)"'

# Cargo.toml
bump patch -w -f Cargo.toml -p 'version\s*=\s*"(\d+\.\d+\.\d+)"'
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
