package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/Masterminds/semver/v3"
	"github.com/manifoldco/promptui"
	"github.com/mattn/go-tty/v2"
)

const version = "0.0.5"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: bump (major|minor|patch|set <version>|show) -f <file> -p <pattern>

Commands:
  major             bump major version up
  minor             bump minor version up
  patch             bump patch version up
  up                bump up with prompt
  set <version>     set exact version (no increments)
  show              only show the current version

Flags:
`)
	flag.PrintDefaults()
}

func run(argv []string) error {
	if len(argv) < 1 {
		usage()
		return errors.New("please specify subcommand")
	}

	var (
		majorDelta uint64
		minorDelta uint64
		patchDelta uint64
		exact      string
		showOnly   bool
		prompt     bool
	)

	parseOffset := 1
	switch argv[0] {
	case "major":
		majorDelta = 1
	case "minor":
		minorDelta = 1
	case "patch":
		patchDelta = 1
	case "up":
		prompt = true
	case "set":
		if len(argv) < 2 {
			return errors.New("please specify a version to set")
		}
		exact = argv[1]
		parseOffset = 2
	case "show":
		showOnly = true
	case "-v", "-version", "--version":
		fmt.Println(version)
		return nil
	case "-h", "-help", "--help":
		usage()
		return nil
	default:
		return fmt.Errorf("unknown subcommand %q", argv[0])
	}

	fs := flag.NewFlagSet("bump", flag.ContinueOnError)
	var (
		file    string
		pattern string
		write   bool
	)
	var yes bool
	fs.StringVar(&file, "f", "", "target file")
	fs.StringVar(&pattern, "p", "", "regexp pattern with a capture group for the version")
	fs.BoolVar(&write, "w", false, "write result to file instead of stdout")
	fs.BoolVar(&yes, "y", false, "skip prompt and use patch (for non-interactive environments)")
	if err := fs.Parse(argv[parseOffset:]); err != nil {
		return err
	}

	if file == "" {
		return errors.New("-f flag is required")
	}
	if pattern == "" {
		return errors.New("-p flag is required")
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern: %w", err)
	}
	if re.NumSubexp() < 1 {
		return errors.New("pattern must contain at least one capture group for the version")
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	loc := re.FindSubmatchIndex(content)
	if loc == nil {
		return errors.New("pattern did not match")
	}

	// loc[2], loc[3] are the start and end of the first capture group
	currentVersion := string(content[loc[2]:loc[3]])

	if showOnly {
		fmt.Println(currentVersion)
		return nil
	}

	if prompt {
		result, err := promptTarget(currentVersion, file)
		if err != nil {
			if !yes {
				return err
			}
			patchDelta = 1
		} else {
			switch result {
			case promptResultPatch:
				patchDelta = 1
			case promptResultMinor:
				minorDelta = 1
			case promptResultMajor:
				majorDelta = 1
			}
		}
	}

	newVersion, err := bumpVersion(currentVersion, majorDelta, minorDelta, patchDelta, exact)
	if err != nil {
		return fmt.Errorf("version bump failed: %w", err)
	}

	result := make([]byte, 0, len(content)+len(newVersion)-len(currentVersion))
	result = append(result, content[:loc[2]]...)
	result = append(result, []byte(newVersion)...)
	result = append(result, content[loc[3]:]...)

	if write {
		if err := os.WriteFile(file, result, 0644); err != nil {
			return err
		}
		fmt.Println(newVersion)
	} else {
		os.Stdout.Write(result)
	}

	return nil
}

func bumpVersion(version string, majorDelta, minorDelta, patchDelta uint64, exact string) (string, error) {
	if exact != "" {
		ev, err := semver.StrictNewVersion(exact)
		if err != nil {
			return "", fmt.Errorf("invalid version %q: %w", exact, err)
		}
		if v, err := semver.StrictNewVersion(version); err == nil {
			if !ev.GreaterThan(v) {
				return "", fmt.Errorf("version %s is not greater than the current version %s", ev, v)
			}
		}
		return ev.String(), nil
	}

	v, err := semver.StrictNewVersion(version)
	if err != nil {
		return "", fmt.Errorf("invalid current version %q: %w", version, err)
	}

	if majorDelta > 0 {
		for i := uint64(0); i < majorDelta; i++ {
			*v = v.IncMajor()
		}
	} else if minorDelta > 0 {
		for i := uint64(0); i < minorDelta; i++ {
			*v = v.IncMinor()
		}
	} else if patchDelta > 0 {
		for i := uint64(0); i < patchDelta; i++ {
			*v = v.IncPatch()
		}
	}

	return v.String(), nil
}

type escInterceptReader struct {
	r   io.ReadCloser
	buf []byte
}

func newEscInterceptReader(r io.ReadCloser) *escInterceptReader {
	return &escInterceptReader{r: r}
}

func (e *escInterceptReader) Read(p []byte) (int, error) {
	if len(e.buf) > 0 {
		n := copy(p, e.buf)
		e.buf = e.buf[n:]
		return n, nil
	}

	n, err := e.r.Read(p)
	if n > 0 && p[0] == 0x1b { // ESC
		if n == 1 {
			// ESC alone → Ctrl+C (interrupt)
			p[0] = 0x03
		}
		// n > 1: part of escape sequence (e.g. arrows), pass through
	}
	return n, err
}

func (e *escInterceptReader) Close() error {
	return e.r.Close()
}

type promptResult int

const (
	promptResultNone promptResult = iota
	promptResultPatch
	promptResultMinor
	promptResultMajor
)

func promptTarget(currentVersion, target string) (promptResult, error) {
	t, err := tty.Open()
	if err != nil {
		return promptResultNone, err
	}
	defer t.Close()

	candidates := []struct {
		name   string
		delta  [3]uint64 // major, minor, patch
		result promptResult
	}{
		{"patch", [3]uint64{0, 0, 1}, promptResultPatch},
		{"minor", [3]uint64{0, 1, 0}, promptResultMinor},
		{"major", [3]uint64{1, 0, 0}, promptResultMajor},
	}

	items := make([]string, len(candidates))
	for i, c := range candidates {
		newVersion, err := bumpVersion(currentVersion, c.delta[0], c.delta[1], c.delta[2], "")
		if err != nil {
			return promptResultNone, err
		}
		items[i] = fmt.Sprintf("%s (%s -> %s)", c.name, currentVersion, newVersion)
	}

	stdin := newEscInterceptReader(t.Input())
	p := promptui.Select{
		Label:    "Bump up " + target,
		HideHelp: true,
		Items:    items,
		Stdin:    stdin,
		Stdout:   t.Output(),
	}

	index, _, err := p.Run()
	if err != nil {
		return promptResultNone, err
	}

	return candidates[index].result, nil
}
