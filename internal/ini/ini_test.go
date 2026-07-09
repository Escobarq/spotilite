package ini

import (
	"strings"
	"testing"
)

func TestParseBasic(t *testing.T) {
	in := `
; comment
[a]
k = v
[b]
n = 1
`
	out := Sections{}
	if err := Parse(in, out); err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if Get(out, "a", "k") != "v" {
		t.Errorf("got %q want v", Get(out, "a", "k"))
	}
	if Get(out, "b", "n") != "1" {
		t.Errorf("b.n wrong: %q", Get(out, "b", "n"))
	}
}

func TestParseEmptyAndComments(t *testing.T) {
	in := `# only comments here
; second line`
	out := Sections{}
	if err := Parse(in, out); err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty sections, got %v", out)
	}
}

func TestParseValueWithEquals(t *testing.T) {
	out := Sections{}
	if err := Parse("[s]\nk=a=b=c\n", out); err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if Get(out, "s", "k") != "a=b=c" {
		t.Errorf("got %q want a=b=c", Get(out, "s", "k"))
	}
}

func TestParseOutsideSectionDropped(t *testing.T) {
	out := Sections{}
	if err := Parse("k=v\n[s]\nx=y\n", out); err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if _, ok := out[""]; ok {
		t.Errorf("empty section should not exist")
	}
}

func TestParseInvalidMissingEquals(t *testing.T) {
	out := Sections{}
	err := Parse("[s]\nnoequals here\n", out)
	if err == nil {
		t.Fatalf("expected error on missing = inside section")
	}
}

func TestSerializeRoundTrip(t *testing.T) {
	original := `[Setting]
a = 1
b = 2

[AdditionalOptions]
x = y
`
	out := Sections{}
	if err := Parse(original, out); err != nil {
		t.Fatalf("Parse: %v", err)
	}
	serialized := Serialize(out)
	if !strings.Contains(serialized, "[Setting]") {
		t.Errorf("missing Setting section in:\n%s", serialized)
	}
	if !strings.Contains(serialized, "a = 1") {
		t.Errorf("missing a=1 in:\n%s", serialized)
	}
	// Round-trip should be idempotent at the data level.
	out2 := Sections{}
	if err := Parse(serialized, out2); err != nil {
		t.Fatalf("reparse: %v", err)
	}
	if Get(out2, "AdditionalOptions", "x") != "y" {
		t.Errorf("lost AdditionalOptions.x after round-trip")
	}
}

func TestGetBool(t *testing.T) {
	out := Sections{}
	Parse("[s]\nenabled = 1\ndisabled = 0\n", out)
	if !GetBool(out, "s", "enabled") {
		t.Error("enabled should be true")
	}
	if GetBool(out, "s", "disabled") {
		t.Error("disabled should be false")
	}
	if GetBool(out, "s", "missing") {
		t.Error("missing should be false")
	}
}
