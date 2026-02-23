package main

import (
	"bufio"
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	charUppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charLowercase = "abcdefghijklmnopqrstuvwxyz"
	charDigits    = "0123456789"
	charSymbols   = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

// ── Crypto helpers ────────────────────────────────────────────────────────────

func randInt(max int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

func filterChars(s, exclude string) string {
	if exclude == "" {
		return s
	}
	var sb strings.Builder
	for _, ch := range s {
		if !strings.ContainsRune(exclude, ch) {
			sb.WriteRune(ch)
		}
	}
	return sb.String()
}

// ── Configs ───────────────────────────────────────────────────────────────────

type RandomConfig struct {
	Length    int
	NoUpper   bool
	NoLower   bool
	NoDigits  bool
	NoSymbols bool
	Exclude   string
}

type SegmentConfig struct {
	Segments  int
	SegLength int
	Separator string // "-" or "_"
	NoUpper   bool
	NoLower   bool
	NoDigits  bool
	Exclude   string
}

// ── Generators ────────────────────────────────────────────────────────────────

func buildSets(noUpper, noLower, noDigits bool, noSymbols *bool, exclude string) (sets []string, fullCharset string) {
	var cs strings.Builder
	addSet := func(chars string) {
		filtered := filterChars(chars, exclude)
		if filtered != "" {
			sets = append(sets, filtered)
			cs.WriteString(filtered)
		}
	}
	if !noUpper { addSet(charUppercase) }
	if !noLower { addSet(charLowercase) }
	if !noDigits { addSet(charDigits) }
	if noSymbols != nil && !*noSymbols { addSet(charSymbols) }
	return sets, cs.String()
}

func shuffleBytes(b []byte) error {
	for i := len(b) - 1; i > 0; i-- {
		j, err := randInt(i + 1)
		if err != nil {
			return err
		}
		b[i], b[j] = b[j], b[i]
	}
	return nil
}

func generateRandom(cfg RandomConfig) (string, error) {
	noSym := cfg.NoSymbols
	sets, charset := buildSets(cfg.NoUpper, cfg.NoLower, cfg.NoDigits, &noSym, cfg.Exclude)
	if charset == "" {
		return "", fmt.Errorf("no characters available — all sets excluded")
	}

	password := make([]byte, cfg.Length)

	// Guarantee at least one char from each active set
	pos := 0
	for _, set := range sets {
		if pos >= cfg.Length {
			break
		}
		idx, err := randInt(len(set))
		if err != nil {
			return "", err
		}
		password[pos] = set[idx]
		pos++
	}

	// Fill remainder from full charset
	for i := pos; i < cfg.Length; i++ {
		idx, err := randInt(len(charset))
		if err != nil {
			return "", err
		}
		password[i] = charset[idx]
	}

	if err := shuffleBytes(password); err != nil {
		return "", err
	}
	return string(password), nil
}

func generateSegmented(cfg SegmentConfig) (string, error) {
	_, charset := buildSets(cfg.NoUpper, cfg.NoLower, cfg.NoDigits, nil, cfg.Exclude)
	if charset == "" {
		return "", fmt.Errorf("no characters available — all sets excluded")
	}

	parts := make([]string, cfg.Segments)
	for i := range parts {
		seg := make([]byte, cfg.SegLength)
		for j := range seg {
			idx, err := randInt(len(charset))
			if err != nil {
				return "", err
			}
			seg[j] = charset[idx]
		}
		parts[i] = string(seg)
	}
	return strings.Join(parts, cfg.Separator), nil
}

// ── Clipboard ─────────────────────────────────────────────────────────────────

func copyToClipboard(text string) error {
	var cmd *exec.Cmd
	switch {
	case commandExists("pbcopy"):                          // macOS
		cmd = exec.Command("pbcopy")
	case commandExists("xclip"):                           // Linux (X11)
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case commandExists("xsel"):                            // Linux alt
		cmd = exec.Command("xsel", "--clipboard", "--input")
	case commandExists("wl-copy"):                         // Wayland
		cmd = exec.Command("wl-copy")
	case commandExists("clip"):                            // Windows
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("no clipboard utility found (pbcopy / xclip / xsel / wl-copy / clip)")
	}
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// ── Interactive mode ──────────────────────────────────────────────────────────

var reader = bufio.NewReader(os.Stdin)

func ask(prompt string) string {
	fmt.Print(prompt)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func askDefault(prompt, def string) string {
	fmt.Printf("%s [%s]: ", prompt, def)
	line, _ := reader.ReadString('\n')
	val := strings.TrimSpace(line)
	if val == "" {
		return def
	}
	return val
}

func askInt(prompt string, def int) int {
	for {
		raw := askDefault(prompt, strconv.Itoa(def))
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 {
			fmt.Println("  ✗  Please enter a number >= 1.")
			continue
		}
		return n
	}
}

func askYesNo(prompt string, def bool) bool {
	defStr := "y"
	if !def {
		defStr = "n"
	}
	raw := strings.ToLower(askDefault(prompt+" (y/n)", defStr))
	return raw == "y" || raw == "yes"
}

func askChoice(prompt string, choices []string, def string) string {
	for {
		raw := askDefault(prompt, def)
		for _, c := range choices {
			if strings.EqualFold(raw, c) {
				return c
			}
		}
		fmt.Printf("  ✗  Please choose one of: %s\n", strings.Join(choices, " / "))
	}
}

func printDivider() { fmt.Println("  " + strings.Repeat("─", 44)) }

func runInteractive() {
	fmt.Println()
	fmt.Println("  ┌──────────────────────────────────────────┐")
	fmt.Println("  │         passgen · Password Generator       │")
	fmt.Println("  └──────────────────────────────────────────┘")
	fmt.Println()

	// ── Type ──
	fmt.Println("  Password type:")
	fmt.Println("    1  Random    e.g. X7&kP2!qL9mR@wZ")
	fmt.Println("    2  Segmented e.g. ab12-cd34-ef56")
	fmt.Println()
	typeChoice := askChoice("  Choose type", []string{"1", "2"}, "1")
	fmt.Println()

	count := askInt("  How many passwords to generate", 1)
	fmt.Println()

	var passwords []string

	switch typeChoice {
	case "1":
		// ── Random ──
		printDivider()
		fmt.Println("  Random password options")
		printDivider()
		length    := askInt("  Length", 16)
		noUpper   := !askYesNo("  Include uppercase  (A-Z)", true)
		noLower   := !askYesNo("  Include lowercase  (a-z)", true)
		noDigits  := !askYesNo("  Include digits     (0-9)", true)
		noSymbols := !askYesNo("  Include symbols    (!@#$...)", true)
		excludeRaw := askDefault("  Exclude characters (leave blank to skip)", "")
		fmt.Println()

		cfg := RandomConfig{
			Length:    length,
			NoUpper:   noUpper,
			NoLower:   noLower,
			NoDigits:  noDigits,
			NoSymbols: noSymbols,
			Exclude:   excludeRaw,
		}
		for i := 0; i < count; i++ {
			p, err := generateRandom(cfg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  error: %v\n", err)
				os.Exit(1)
			}
			passwords = append(passwords, p)
		}

	case "2":
		// ── Segmented ──
		printDivider()
		fmt.Println("  Segmented password options")
		printDivider()
		segments  := askInt("  Number of segments", 3)
		segLength := askInt("  Characters per segment", 4)
		separator := askChoice("  Separator  (- or _)", []string{"-", "_"}, "-")
		fmt.Println()
		noUpper  := !askYesNo("  Include uppercase  (A-Z)", false)
		noLower  := !askYesNo("  Include lowercase  (a-z)", true)
		noDigits := !askYesNo("  Include digits     (0-9)", true)
		excludeRaw := askDefault("  Exclude characters (leave blank to skip)", "")
		fmt.Println()

		cfg := SegmentConfig{
			Segments:  segments,
			SegLength: segLength,
			Separator: separator,
			NoUpper:   noUpper,
			NoLower:   noLower,
			NoDigits:  noDigits,
			Exclude:   excludeRaw,
		}
		for i := 0; i < count; i++ {
			p, err := generateSegmented(cfg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  error: %v\n", err)
				os.Exit(1)
			}
			passwords = append(passwords, p)
		}
	}

	// ── Output ──
	printDivider()
	if len(passwords) == 1 {
		fmt.Printf("  Password: %s\n", passwords[0])
	} else {
		fmt.Println("  Generated passwords:")
		for i, p := range passwords {
			fmt.Printf("  %2d. %s\n", i+1, p)
		}
	}
	printDivider()

	// Copy to clipboard (last password if multiple)
	toCopy := passwords[len(passwords)-1]
	if err := copyToClipboard(toCopy); err != nil {
		fmt.Fprintf(os.Stderr, "  (clipboard unavailable: %v)\n", err)
	} else {
		if len(passwords) == 1 {
			fmt.Println("  Copied to clipboard!")
		} else {
			fmt.Printf("  Password #%d copied to clipboard!\n", len(passwords))
		}
	}
	fmt.Println()
}

// ── Flag mode ─────────────────────────────────────────────────────────────────

func runQuickSegment(separator string) {
	cfg := SegmentConfig{
		Segments:  5,
		SegLength: 5,
		Separator: separator,
		NoUpper:   false,
		NoLower:   false,
		NoDigits:  false,
	}
	p, err := generateSegmented(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(p)
	if err := copyToClipboard(p); err != nil {
		fmt.Fprintf(os.Stderr, "(clipboard unavailable: %v)\n", err)
	} else {
		fmt.Fprintln(os.Stderr, "Copied to clipboard.")
	}
}

func main() {
	// No args → interactive
	if len(os.Args) == 1 {
		runInteractive()
		return
	}

	// Quick segmented mode: passgen - or passgen _
	if len(os.Args) == 2 && (os.Args[1] == "-" || os.Args[1] == "_") {
		runQuickSegment(os.Args[1])
		return
	}

	fs := flag.NewFlagSet("passgen", flag.ExitOnError)

	mode      := fs.String("type",      "random", "Password type: random or segment")
	length    := fs.Int("length",       16,       "Password length (random mode)")
	count     := fs.Int("count",        1,        "Number of passwords to generate")
	noUpper   := fs.Bool("no-upper",    false,    "Exclude uppercase letters (A-Z)")
	noLower   := fs.Bool("no-lower",    false,    "Exclude lowercase letters (a-z)")
	noDigits  := fs.Bool("no-digits",   false,    "Exclude digits (0-9)")
	noSymbols := fs.Bool("no-symbols",  false,    "Exclude symbols (random mode only)")
	exclude   := fs.String("exclude",   "",       "Characters to exclude")
	segments  := fs.Int("segments",     3,        "Number of segments (segment mode)")
	segLen    := fs.Int("seg-length",   4,        "Characters per segment (segment mode)")
	separator := fs.String("separator", "-",      "Segment separator: - or _ (segment mode)")
	noCopy    := fs.Bool("no-copy",     false,    "Skip copying to clipboard")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "passgen — Cryptographically secure password generator")
		fmt.Fprintln(os.Stderr, "\nUsage:")
		fmt.Fprintln(os.Stderr, "  passgen                                   Interactive mode")
		fmt.Fprintln(os.Stderr, "  passgen [options]                         Flag mode")
		fmt.Fprintln(os.Stderr, "\nOptions:")
		fs.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, `  passgen -length 32`)
		fmt.Fprintln(os.Stderr, `  passgen -count 5 -no-symbols`)
		fmt.Fprintln(os.Stderr, `  passgen -exclude "0OIl1"`)
		fmt.Fprintln(os.Stderr, `  passgen -type segment -segments 4 -seg-length 5`)
		fmt.Fprintln(os.Stderr, `  passgen -type segment -separator _`)
		fmt.Fprintln(os.Stderr, `  passgen -type segment -segments 3 -seg-length 6 -no-copy`)
	}

	fs.Parse(os.Args[1:])

	var passwords []string

	switch strings.ToLower(*mode) {
	case "random":
		if *length < 1 {
			fmt.Fprintln(os.Stderr, "error: -length must be >= 1")
			os.Exit(1)
		}
		cfg := RandomConfig{
			Length:    *length,
			NoUpper:   *noUpper,
			NoLower:   *noLower,
			NoDigits:  *noDigits,
			NoSymbols: *noSymbols,
			Exclude:   *exclude,
		}
		for i := 0; i < *count; i++ {
			p, err := generateRandom(cfg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			passwords = append(passwords, p)
		}

	case "segment":
		if *separator != "-" && *separator != "_" {
			fmt.Fprintln(os.Stderr, "error: -separator must be - or _")
			os.Exit(1)
		}
		if *segments < 1 {
			fmt.Fprintln(os.Stderr, "error: -segments must be >= 1")
			os.Exit(1)
		}
		if *segLen < 1 {
			fmt.Fprintln(os.Stderr, "error: -seg-length must be >= 1")
			os.Exit(1)
		}
		cfg := SegmentConfig{
			Segments:  *segments,
			SegLength: *segLen,
			Separator: *separator,
			NoUpper:   *noUpper,
			NoLower:   *noLower,
			NoDigits:  *noDigits,
			Exclude:   *exclude,
		}
		for i := 0; i < *count; i++ {
			p, err := generateSegmented(cfg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			passwords = append(passwords, p)
		}

	default:
		fmt.Fprintf(os.Stderr, "error: unknown type %q — use random or segment\n", *mode)
		os.Exit(1)
	}

	for _, p := range passwords {
		fmt.Println(p)
	}

	if !*noCopy && len(passwords) > 0 {
		toCopy := passwords[len(passwords)-1]
		if err := copyToClipboard(toCopy); err != nil {
			fmt.Fprintf(os.Stderr, "(clipboard unavailable: %v)\n", err)
		} else {
			if len(passwords) == 1 {
				fmt.Fprintln(os.Stderr, "Copied to clipboard.")
			} else {
				fmt.Fprintf(os.Stderr, "Password #%d copied to clipboard.\n", len(passwords))
			}
		}
	}
}
