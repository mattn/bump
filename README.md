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

### package.json

```json
{
  "name": "my-app",
  "version": "1.2.3"
}
```

```bash
bump patch -w -f package.json -p '"version":\s*"(\d+\.\d+\.\d+)"'
```

### pyproject.toml

```toml
[project]
name = "my-app"
version = "1.2.3"
```

```bash
bump minor -w -f pyproject.toml -p 'version\s*=\s*"(\d+\.\d+\.\d+)"'
```

### Cargo.toml

```toml
[package]
name = "my-app"
version = "1.2.3"
```

```bash
bump major -w -f Cargo.toml -p 'version\s*=\s*"(\d+\.\d+\.\d+)"'
```

### Go source

```go
const version = "1.2.3"
```

```bash
bump patch -w -f version.go -p 'version\s*=\s*"(\d+\.\d+\.\d+)"'
```

### Interactive prompt

```bash
bump up -w -f package.json -p '"version":\s*"(\d+\.\d+\.\d+)"'
```

### Set exact version

```bash
bump set 2.0.0 -w -f pyproject.toml -p 'version\s*=\s*"(\d+\.\d+\.\d+)"'
```

### Show current version

```bash
bump show -f Cargo.toml -p 'version\s*=\s*"(\d+\.\d+\.\d+)"'
```

## License

MIT

## Author

mattn
