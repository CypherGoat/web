package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const defaultAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

type huntEntry struct {
	Prize string `json:"prize"`
	Label string `json:"label"`
}

type countFlags map[string]int

func (f *countFlags) String() string {
	if f == nil || len(*f) == 0 {
		return ""
	}
	var parts []string
	for prize, count := range *f {
		parts = append(parts, fmt.Sprintf("%s=%d", prize, count))
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}

func (f *countFlags) Set(value string) error {
	prize, rawCount, ok := strings.Cut(value, "=")
	if !ok {
		return fmt.Errorf("invalid count %q: expected prize=count", value)
	}
	prize = normalisePrize(prize)
	if prize == "" {
		return fmt.Errorf("invalid count %q: missing prize", value)
	}
	count, err := strconv.Atoi(strings.TrimSpace(rawCount))
	if err != nil || count <= 0 {
		return fmt.Errorf("invalid count %q: count must be a positive integer", value)
	}
	if *f == nil {
		*f = make(countFlags)
	}
	(*f)[prize] += count
	return nil
}

type labelFlags map[string]string

func (f *labelFlags) String() string {
	if f == nil || len(*f) == 0 {
		return ""
	}
	var parts []string
	for prize, label := range *f {
		parts = append(parts, fmt.Sprintf("%s=%s", prize, label))
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}

func (f *labelFlags) Set(value string) error {
	prize, label, ok := strings.Cut(value, "=")
	if !ok {
		return fmt.Errorf("invalid label %q: expected prize=label", value)
	}
	prize = normalisePrize(prize)
	label = strings.TrimSpace(label)
	if prize == "" || label == "" {
		return fmt.Errorf("invalid label %q: prize and label are required", value)
	}
	if *f == nil {
		*f = make(labelFlags)
	}
	(*f)[prize] = label
	return nil
}

func main() {
	var counts countFlags
	var labels labelFlags

	outPath := flag.String("out", "-", "output file path, or - for stdout")
	existingPath := flag.String("existing", "", "existing hunt JSON to merge with and avoid collisions against")
	prefix := flag.String("prefix", "", "optional code prefix")
	length := flag.Int("length", 8, "random suffix length")
	alphabet := flag.String("alphabet", defaultAlphabet, "characters allowed in generated suffixes")
	flag.Var(&counts, "count", "prize=count pair; repeat for multiple prize types")
	flag.Var(&labels, "label", "prize=label override; repeat as needed")
	flag.Parse()

	if len(counts) == 0 {
		exitf("at least one -count prize=count flag is required")
	}
	if *length <= 0 {
		exitf("-length must be greater than zero")
	}

	cleanPrefix := strings.ToUpper(strings.TrimSpace(*prefix))

	cleanAlphabet := dedupeAlphabet(strings.ToUpper(strings.TrimSpace(*alphabet)))
	if len(cleanAlphabet) < 2 {
		exitf("-alphabet must contain at least two distinct characters")
	}

	entries, err := loadExisting(*existingPath)
	if err != nil {
		exitf("load existing entries: %v", err)
	}

	for prize, count := range counts {
		label := labels[prize]
		if label == "" {
			label = defaultLabel(prize)
		}

		for i := 0; i < count; i++ {
			code, err := nextUniqueCode(cleanPrefix, cleanAlphabet, *length, entries)
			if err != nil {
				exitf("generate code for %s: %v", prize, err)
			}
			entries[code] = huntEntry{
				Prize: prize,
				Label: label,
			}
		}
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		exitf("marshal output: %v", err)
	}
	data = append(data, '\n')

	if *outPath == "-" {
		if _, err := os.Stdout.Write(data); err != nil {
			exitf("write stdout: %v", err)
		}
		return
	}

	if err := os.WriteFile(*outPath, data, 0644); err != nil {
		exitf("write %s: %v", *outPath, err)
	}
}

func loadExisting(path string) (map[string]huntEntry, error) {
	entries := make(map[string]huntEntry)
	if strings.TrimSpace(path) == "" {
		return entries, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	for code, payload := range raw {
		var entry huntEntry
		if err := json.Unmarshal(payload, &entry); err != nil || entry.Prize == "" {
			continue
		}
		entries[strings.ToUpper(strings.TrimSpace(code))] = entry
	}

	return entries, nil
}

func nextUniqueCode(prefix, alphabet string, length int, existing map[string]huntEntry) (string, error) {
	for attempts := 0; attempts < 10000; attempts++ {
		suffix, err := randomString(alphabet, length)
		if err != nil {
			return "", err
		}
		code := suffix
		if prefix != "" {
			code = prefix + "-" + suffix
		}
		if _, taken := existing[code]; !taken {
			return code, nil
		}
	}
	return "", fmt.Errorf("failed to generate a unique code after repeated attempts")
}

func randomString(alphabet string, length int) (string, error) {
	buf := make([]byte, length)
	random := make([]byte, length)
	if _, err := rand.Read(random); err != nil {
		return "", err
	}
	for i := range buf {
		buf[i] = alphabet[int(random[i])%len(alphabet)]
	}
	return string(buf), nil
}

func dedupeAlphabet(alphabet string) string {
	var b strings.Builder
	seen := make(map[rune]struct{}, len(alphabet))
	for _, r := range alphabet {
		if _, ok := seen[r]; ok {
			continue
		}
		seen[r] = struct{}{}
		b.WriteRune(r)
	}
	return b.String()
}

func defaultLabel(prize string) string {
	if prize == "shirt" {
		return "CypherGoat T-Shirt"
	}

	if amount, ok := strings.CutSuffix(prize, "xmr"); ok {
		amount = strings.TrimSpace(amount)
		if amount != "" {
			return amount + " EUR in XMR"
		}
	}

	return strings.ToUpper(prize)
}

func normalisePrize(prize string) string {
	return strings.ToLower(strings.TrimSpace(prize))
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
