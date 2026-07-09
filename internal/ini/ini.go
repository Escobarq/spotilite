// Package ini provides a minimal INI-style parser used by both the spicetify
// config loader and the theme color.ini reader. Splits lines into a map of
// section -> key -> value. Comments starting with ; or # are ignored.
//
// The implementation is deliberately lenient: quoted values keep their
// surrounding quotes (matching spicetify-cli's behavior), no escape handling
// is performed, and unknown sections/keys round-trip unchanged.
package ini

import (
	"bufio"
	"fmt"
	"sort"
	"strings"
)

// Sections maps section name -> key/value pairs.
type Sections map[string]map[string]string

// Parse reads an INI document into out. Unknown sections are preserved.
func Parse(text string, out Sections) error {
	if out == nil {
		out = Sections{}
	}
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Buffer(make([]byte, 1024), 1024*1024)

	var current string
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		trim := strings.TrimSpace(raw)
		if trim == "" {
			continue
		}
		if trim[0] == ';' || trim[0] == '#' {
			continue
		}
		if trim[0] == '[' && trim[len(trim)-1] == ']' {
			name := strings.TrimSpace(trim[1 : len(trim)-1])
			if name == "" {
				return fmt.Errorf("line %d: empty section name", lineNum)
			}
			if out[name] == nil {
				out[name] = map[string]string{}
			}
			current = name
			continue
		}
		if current == "" {
			// key/value outside a section: silently dropped.
			continue
		}
		eq := strings.IndexByte(trim, '=')
		if eq < 0 {
			return fmt.Errorf("line %d: expected key=value, got %q", lineNum, trim)
		}
		key := strings.TrimSpace(trim[:eq])
		val := strings.TrimSpace(trim[eq+1:])
		if key == "" {
			return fmt.Errorf("line %d: empty key", lineNum)
		}
		if out[current] == nil {
			out[current] = map[string]string{}
		}
		out[current][key] = val
	}
	return scanner.Err()
}

// Serialize converts a Sections map back to INI text with stable ordering.
// Sections named "Setting", "Preprocesses", "AdditionalOptions", "Patch"
// appear first (mirroring spicetify-cli ordering), followed by any others
// alphabetically.
func Serialize(in Sections) string {
	var b strings.Builder
	canonical := []string{"Setting", "Preprocesses", "AdditionalOptions", "Patch"}
	seen := make(map[string]bool, len(in))
	for _, name := range canonical {
		if _, ok := in[name]; !ok {
			continue
		}
		writeSection(&b, name, in[name])
		seen[name] = true
		b.WriteString("\n")
	}
	var others []string
	for name := range in {
		if !seen[name] {
			others = append(others, name)
		}
	}
	sort.Strings(others)
	for _, name := range others {
		writeSection(&b, name, in[name])
		b.WriteString("\n")
	}
	return b.String()
}

func writeSection(b *strings.Builder, name string, kv map[string]string) {
	b.WriteString("[")
	b.WriteString(name)
	b.WriteString("]\n")
	keys := make([]string, 0, len(kv))
	for k := range kv {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		b.WriteString(k)
		b.WriteString(" = ")
		b.WriteString(kv[k])
		b.WriteString("\n")
	}
}

// Get is a shorthand for `out[section][key]` with empty fallback.
func Get(in Sections, section, key string) string {
	if sec, ok := in[section]; ok {
		return sec[key]
	}
	return ""
}

// GetBool reports whether value equals "1".
func GetBool(in Sections, section, key string) bool {
	return strings.TrimSpace(Get(in, section, key)) == "1"
}
