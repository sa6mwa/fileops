package fileops

import (
	"testing"
)

func TestEnsureLineInLines(t *testing.T) {
	lines := []string{
		"# Start",
		"ConfigOption = true",
		"AnotherOption = false",
		"# CommentedOut = false",
		"# End",
	}

	linesCopy := make([]string, len(lines))
	copy(linesCopy, lines)

	expectedLines := []string{
		"# Start",
		"ConfigOption = true",
		"AnotherOption = false",
		"CommentedOut = false",
		"# CommentedOut = false",
		"# End",
	}

	before := "# CommentedOut"
	if err := EnsureLineInLines(&lines, "CommentedOut = false", &before, nil, false, false); err != nil {
		t.Fatal(err)
	}
	compareLines(t, &lines, expectedLines)

	lines = make([]string, len(linesCopy))
	copy(lines, linesCopy)

	after := "AnotherOption"
	if err := EnsureLineInLines(&lines, "CommentedOut = false", nil, &after, false, false); err != nil {
		t.Fatal(err)
	}
	compareLines(t, &lines, expectedLines)

	lines = make([]string, len(linesCopy))
	copy(lines, linesCopy)

	before = "# CommentedOut"
	after = "AnotherOption = false"
	if err := EnsureLineInLines(&lines, "CommentedOut = false", &before, &after, false, false); err != nil {
		t.Fatal(err)
	}
	compareLines(t, &lines, expectedLines)

	lines = make([]string, len(linesCopy))
	copy(lines, linesCopy)

	before = "# CommentedOut = false"
	after = "AnotherOption = false"
	if err := EnsureLineInLines(&lines, "CommentedOut = false", &before, &after, true, true); err != nil {
		t.Fatal(err)
	}
	compareLines(t, &lines, expectedLines)

	lines = []string{}
	expectedLines = []string{
		"Hello = world",
	}

	if err := EnsureLineInLines(&lines, "Hello = world", nil, &[]string{"Hello"}[0], true, false); err != nil {
		t.Fatal(err)
	}
	compareLines(t, &lines, expectedLines)

	lines = []string{}
	if err := EnsureLineInLines(&lines, "Hello = world", nil, nil, true, true); err != nil {
		t.Fatal(err)
	}
	compareLines(t, &lines, expectedLines)
}

func compareLines(t *testing.T, lines *[]string, expectedLines []string) {
	if expected, got := len(expectedLines), len(*lines); expected != got {
		t.Fatalf("Expected %d lines, got %d", expected, got)
	}
	for i := range expectedLines {
		if (*lines)[i] != expectedLines[i] {
			t.Errorf("Expected line %d to be %x, but got %x", i+1, expectedLines[i], (*lines)[i])
		}
	}
}
