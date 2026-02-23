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

// ── Word list (EFF short list style — common, easy-to-remember words) ────────

var wordList = []string{
	"acid", "acme", "aged", "also", "arch", "area", "army", "away",
	"back", "bail", "bake", "ball", "band", "bank", "barn", "base",
	"bath", "bead", "beam", "bear", "beat", "been", "bell", "belt",
	"bend", "best", "bike", "bird", "bite", "blow", "blue", "blur",
	"boat", "body", "bold", "bolt", "bomb", "bond", "bone", "book",
	"boot", "born", "boss", "bowl", "bulk", "bump", "burn", "busy",
	"cafe", "cage", "cake", "calm", "came", "camp", "cape", "card",
	"care", "cart", "cash", "cast", "cave", "chat", "chef", "chin",
	"chip", "chop", "cite", "city", "clad", "clam", "clan", "clap",
	"clay", "clip", "club", "clue", "coal", "coat", "code", "coil",
	"coin", "cold", "colt", "comb", "come", "cook", "cool", "cope",
	"copy", "cord", "core", "corn", "cost", "cozy", "crew", "crop",
	"crow", "cube", "cult", "cure", "curl", "cute", "dare", "dark",
	"dart", "dash", "data", "dawn", "deal", "dean", "dear", "debt",
	"deck", "deed", "deem", "deep", "deer", "demo", "deny", "desk",
	"dial", "dice", "diet", "dime", "dine", "dirt", "disc", "dish",
	"dock", "does", "dome", "done", "doom", "door", "dose", "dove",
	"down", "drag", "draw", "drip", "drop", "drum", "dual", "duck",
	"duel", "duke", "dull", "dumb", "dump", "dune", "dusk", "dust",
	"duty", "each", "earl", "earn", "ease", "east", "easy", "echo",
	"edge", "edit", "else", "emit", "ends", "epic", "euro", "even",
	"ever", "evil", "exam", "exec", "exit", "expo", "face", "fact",
	"fade", "fail", "fair", "fake", "fall", "fame", "fang", "fare",
	"farm", "fast", "fate", "fawn", "fear", "feat", "feed", "feel",
	"fell", "felt", "file", "fill", "film", "find", "fine", "fire",
	"firm", "fish", "fist", "flag", "flat", "flaw", "fled", "flew",
	"flex", "flip", "flow", "foam", "foil", "fold", "folk", "fond",
	"font", "food", "fool", "foot", "ford", "fork", "form", "fort",
	"foul", "four", "fowl", "free", "frog", "from", "fuel", "full",
	"fund", "funk", "fury", "fuse", "gain", "gale", "game", "gang",
	"gape", "garb", "gate", "gave", "gaze", "gear", "gene", "gift",
	"gild", "glad", "glow", "glue", "goat", "goes", "gold", "golf",
	"gone", "good", "grab", "gray", "grew", "grid", "grim", "grin",
	"grip", "grow", "gust", "guts", "hack", "hail", "hair", "hale",
	"half", "hall", "halt", "hand", "hang", "hard", "hare", "harm",
	"harp", "hash", "hate", "haul", "have", "hawk", "haze", "head",
	"heal", "heap", "hear", "heat", "heed", "heel", "held", "helm",
	"help", "herb", "herd", "here", "hero", "hike", "hill", "hilt",
	"hint", "hire", "hold", "hole", "holy", "home", "hood", "hook",
	"hope", "horn", "host", "hour", "howl", "huge", "hull", "hung",
	"hunt", "hurt", "hush", "hymn", "icon", "idea", "inch", "info",
	"into", "iron", "isle", "item", "jack", "jade", "jail", "jamb",
	"jaws", "jazz", "jean", "jerk", "jest", "jets", "jobs", "join",
	"joke", "jolt", "jump", "june", "jury", "just", "keen", "keep",
	"kelp", "kept", "kick", "kids", "kill", "kind", "king", "kiss",
	"kite", "knee", "knew", "knit", "knob", "knot", "know", "lace",
	"lack", "laid", "lake", "lamb", "lamp", "land", "lane", "lard",
	"lark", "lash", "last", "late", "lawn", "lead", "leaf", "leak",
	"lean", "leap", "left", "lend", "lens", "lent", "less", "levy",
	"liar", "lick", "lied", "life", "lift", "like", "limb", "lime",
	"limp", "line", "link", "lint", "lion", "lips", "list", "live",
	"load", "loaf", "loan", "lock", "loft", "logo", "lone", "long",
	"look", "loop", "lord", "lore", "lose", "loss", "lost", "love",
	"luck", "lump", "lung", "lure", "lurk", "lush", "made", "maid",
	"mail", "main", "make", "male", "malt", "mane", "many", "maps",
	"mare", "mark", "mars", "mash", "mask", "mass", "mast", "mate",
	"maze", "meal", "mean", "meat", "meld", "melt", "memo", "mend",
	"menu", "mere", "mesa", "mesh", "mild", "mile", "milk", "mill",
	"mime", "mind", "mine", "mint", "miss", "mist", "moan", "moat",
	"mock", "mode", "mold", "mole", "monk", "mood", "moon", "more",
	"moss", "most", "moth", "move", "much", "mule", "murk", "muse",
	"musk", "must", "myth", "nail", "name", "navy", "near", "neat",
	"neck", "need", "nest", "news", "next", "nice", "nine", "node",
	"none", "noon", "norm", "nose", "note", "noun", "null", "numb",
	"oath", "obey", "odds", "omit", "once", "only", "onto", "opal",
	"open", "oral", "orca", "oven", "over", "owed", "owls", "owns",
	"pace", "pack", "page", "paid", "pail", "pain", "pair", "pale",
	"palm", "pane", "park", "part", "pass", "past", "path", "pave",
	"peak", "pear", "peel", "perk", "pest", "pick", "pier", "pike",
	"pile", "pine", "pink", "pipe", "plan", "play", "plea", "plow",
	"plug", "plum", "plus", "poke", "pole", "poll", "polo", "pond",
	"pool", "poor", "pope", "pork", "port", "pose", "post", "pour",
	"pray", "prey", "prop", "pull", "pulp", "pump", "punk", "pure",
	"push", "quit", "quiz", "race", "rack", "raft", "rage", "raid",
	"rail", "rain", "rake", "ramp", "rang", "rank", "rare", "rash",
	"rate", "rave", "rays", "read", "real", "reap", "rear", "reed",
	"reef", "rein", "rely", "rent", "rest", "rice", "rich", "ride",
	"rift", "rims", "ring", "riot", "rise", "risk", "road", "roam",
	"robe", "rock", "rode", "role", "roll", "roof", "room", "root",
	"rope", "rose", "ruin", "rule", "rush", "rust", "sack", "safe",
	"sage", "said", "sail", "sake", "sale", "salt", "same", "sand",
	"sane", "sang", "sank", "save", "seal", "seam", "seed", "seek",
	"seen", "self", "sell", "semi", "send", "sent", "shed", "shin",
	"ship", "shop", "shot", "show", "shut", "sick", "side", "sift",
	"sign", "silk", "sink", "site", "size", "skin", "skip", "slab",
	"slam", "slap", "sled", "slew", "slid", "slim", "slip", "slot",
	"slow", "slug", "snap", "snow", "soak", "soap", "soar", "sock",
	"soft", "soil", "sold", "sole", "some", "song", "soon", "sore",
	"sort", "soul", "sour", "span", "spar", "spec", "sped", "spin",
	"spit", "spot", "spur", "star", "stay", "stem", "step", "stew",
	"stop", "stow", "stub", "such", "suit", "sulk", "sung", "sunk",
	"sure", "surf", "swan", "swap", "swim", "tabs", "tack", "tact",
	"tail", "take", "tale", "talk", "tall", "tame", "tank", "tape",
	"task", "taxi", "team", "tear", "tell", "tend", "tent", "term",
	"test", "text", "than", "them", "then", "they", "thin", "this",
	"tick", "tide", "tidy", "tied", "tier", "tile", "till", "tilt",
	"time", "tiny", "tire", "toad", "toga", "toil", "told", "toll",
	"tomb", "tone", "took", "tool", "tops", "tore", "torn", "toss",
	"tour", "town", "trap", "tray", "tree", "trek", "trim", "trio",
	"trip", "trot", "true", "tube", "tuck", "tuft", "tuna", "tune",
	"turn", "twin", "type", "ugly", "undo", "unit", "unto", "upon",
	"urge", "used", "user", "vale", "vane", "vary", "vast", "veil",
	"vein", "vent", "verb", "vest", "veto", "vial", "vice", "view",
	"vine", "visa", "void", "volt", "vote", "wade", "wage", "wait",
	"wake", "walk", "wall", "wand", "want", "ward", "warm", "warn",
	"warp", "wary", "wash", "vast", "wave", "wavy", "waxy", "ways",
	"weak", "wear", "weed", "week", "weep", "weld", "well", "went",
	"were", "west", "what", "when", "whom", "wick", "wide", "wife",
	"wild", "will", "wilt", "wily", "wind", "wine", "wing", "wink",
	"wipe", "wire", "wise", "wish", "wisp", "with", "woke", "wolf",
	"wood", "wool", "word", "wore", "work", "worm", "worn", "wrap",
	"wren", "yank", "yard", "yarn", "year", "yell", "yoga", "yoke",
	"your", "zeal", "zero", "zinc", "zone", "zoom",
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

type PassphraseConfig struct {
	Words      int
	Separator  string
	Capitalize bool
	AddNumber  bool
	Include    []string // user's own words to mix in
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

func generatePassphrase(cfg PassphraseConfig) (string, error) {
	// Start with user's included words
	words := make([]string, 0, cfg.Words)
	for _, w := range cfg.Include {
		w = strings.TrimSpace(w)
		if w == "" {
			continue
		}
		if cfg.Capitalize {
			w = strings.ToUpper(w[:1]) + w[1:]
		}
		words = append(words, w)
	}

	// Fill remaining slots with random words
	for len(words) < cfg.Words {
		idx, err := randInt(len(wordList))
		if err != nil {
			return "", err
		}
		w := wordList[idx]
		if cfg.Capitalize {
			w = strings.ToUpper(w[:1]) + w[1:]
		}
		words = append(words, w)
	}

	// Shuffle so user words aren't always at the start
	if err := shuffleStrings(words); err != nil {
		return "", err
	}

	passphrase := strings.Join(words, cfg.Separator)

	if cfg.AddNumber {
		n, err := randInt(1000)
		if err != nil {
			return "", err
		}
		passphrase += cfg.Separator + strconv.Itoa(n)
	}

	return passphrase, nil
}

func shuffleStrings(s []string) error {
	for i := len(s) - 1; i > 0; i-- {
		j, err := randInt(i + 1)
		if err != nil {
			return err
		}
		s[i], s[j] = s[j], s[i]
	}
	return nil
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
	fmt.Println("    1  Random     e.g. X7&kP2!qL9mR@wZ")
	fmt.Println("    2  Segmented  e.g. ab12-cd34-ef56")
	fmt.Println("    3  Passphrase e.g. Tiger-Maple-Cloud-97")
	fmt.Println()
	typeChoice := askChoice("  Choose type", []string{"1", "2", "3"}, "1")
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

	case "3":
		// ── Passphrase ──
		printDivider()
		fmt.Println("  Passphrase options")
		printDivider()
		includeRaw := askDefault("  Your words (comma separated, blank for all random)", "")
		var include []string
		if includeRaw != "" {
			for _, w := range strings.Split(includeRaw, ",") {
				w = strings.TrimSpace(w)
				if w != "" {
					include = append(include, w)
				}
			}
		}
		defWords := 4
		if len(include) > defWords {
			defWords = len(include) + 1
		}
		numWords   := askInt("  Total number of words", defWords)
		if numWords < len(include) {
			numWords = len(include)
		}
		separator  := askChoice("  Separator  (- or _)", []string{"-", "_"}, "-")
		capitalize := askYesNo("  Capitalize each word", true)
		addNumber  := askYesNo("  Add a random number at end", true)
		fmt.Println()

		cfg := PassphraseConfig{
			Words:      numWords,
			Separator:  separator,
			Capitalize: capitalize,
			AddNumber:  addNumber,
			Include:    include,
		}
		for i := 0; i < count; i++ {
			p, err := generatePassphrase(cfg)
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

	mode      := fs.String("type",      "random", "Password type: random, segment, or phrase")
	length    := fs.Int("length",       16,       "Password length (random mode)")
	count     := fs.Int("count",        1,        "Number of passwords to generate")
	noUpper   := fs.Bool("no-upper",    false,    "Exclude uppercase letters (A-Z)")
	noLower   := fs.Bool("no-lower",    false,    "Exclude lowercase letters (a-z)")
	noDigits  := fs.Bool("no-digits",   false,    "Exclude digits (0-9)")
	noSymbols := fs.Bool("no-symbols",  false,    "Exclude symbols (random mode only)")
	exclude   := fs.String("exclude",   "",       "Characters to exclude")
	segments  := fs.Int("segments",     3,        "Number of segments (segment mode)")
	segLen    := fs.Int("seg-length",   4,        "Characters per segment (segment mode)")
	separator := fs.String("separator", "-",      "Separator: - or _ (segment/phrase mode)")
	noCopy    := fs.Bool("no-copy",     false,    "Skip copying to clipboard")
	words     := fs.Int("words",        4,        "Number of words (phrase mode)")
	capitalize := fs.Bool("capitalize", true,     "Capitalize words (phrase mode)")
	addNum    := fs.Bool("add-number",  true,     "Add random number at end (phrase mode)")
	include   := fs.String("include",   "",       "Your words to mix in, comma separated (phrase mode)")

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
		fmt.Fprintln(os.Stderr, `  passgen -type phrase`)
		fmt.Fprintln(os.Stderr, `  passgen -type phrase -words 5 -separator _`)
		fmt.Fprintln(os.Stderr, `  passgen -type phrase -capitalize=false -add-number=false`)
		fmt.Fprintln(os.Stderr, `  passgen -type phrase -include "tiger,coffee"`)
		fmt.Fprintln(os.Stderr, `  passgen -type phrase -include "sun,moon" -words 5`)
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

	case "phrase", "passphrase":
		var inc []string
		if *include != "" {
			for _, w := range strings.Split(*include, ",") {
				w = strings.TrimSpace(w)
				if w != "" {
					inc = append(inc, w)
				}
			}
		}
		if *words < len(inc) {
			*words = len(inc)
		}
		if *words < 1 {
			fmt.Fprintln(os.Stderr, "error: -words must be >= 1")
			os.Exit(1)
		}
		cfg := PassphraseConfig{
			Words:      *words,
			Separator:  *separator,
			Capitalize: *capitalize,
			AddNumber:  *addNum,
			Include:    inc,
		}
		for i := 0; i < *count; i++ {
			p, err := generatePassphrase(cfg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			passwords = append(passwords, p)
		}

	default:
		fmt.Fprintf(os.Stderr, "error: unknown type %q — use random, segment, or phrase\n", *mode)
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
