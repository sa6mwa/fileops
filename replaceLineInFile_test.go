package fileops

import "testing"

func TestReplaceLineInLines(t *testing.T) {
	lines := []string{
		"Hello world",
		"Replace me",
		"Replace me",
		"Next last line",
		"Last line",
	}

	linesCopy := make([]string, len(lines))
	copy(linesCopy, lines)

	expectedLines := []string{
		"Hello world",
		"Replaced was this line",
		"Replace me",
		"Next last line",
		"Last line",
	}

	if err := ReplaceLineInLines(&lines, "Replace", "Replaced was this line", 1, false, false); err != nil {
		t.Fatal(err)
	}
	compareLines(t, &lines, expectedLines)

	lines = make([]string, len(linesCopy))
	copy(lines, linesCopy)

	expectedLines = []string{
		"Hello world",
		"Replaced was this line",
		"Replaced was this line",
		"Next last line",
		"Last line",
	}

	if err := ReplaceLineInLines(&lines, "Replace", "Replaced was this line", -1, false, false); err != nil {
		t.Fatal(err)
	}
	compareLines(t, &lines, expectedLines)

	lines = make([]string, len(linesCopy))
	copy(lines, linesCopy)

	expectedLines = []string{
		"Hello world",
		"Replaced was this line",
		"Replaced was this line",
		"Next last line",
		"Last line",
	}

	if err := ReplaceLineInLines(&lines, "Replace", "Replaced was this line", 2, false, false); err != nil {
		t.Fatal(err)
	}
	compareLines(t, &lines, expectedLines)

	lines = make([]string, len(linesCopy))
	copy(lines, linesCopy)
	copy(expectedLines, linesCopy)

	if err := ReplaceLineInLines(&lines, "Does not exist", "Yep, it got replaced", -1, false, false); err != nil {
		t.Fatal(err)
	}
	compareLines(t, &lines, expectedLines)

	expectedLines = []string{
		"Hello world",
		"Replace me",
		"Replace me",
		"Next last line",
		"Yep, replaced",
	}

	if err := ReplaceLineInLines(&lines, "Last line", "Yep, replaced", -1, false, false); err != nil {
		t.Fatal(err)
	}
	compareLines(t, &lines, expectedLines)

	if err := ReplaceLineInLines(&lines, "Last line", "Yep, replaced", -1, false, false); err != nil {
		t.Fatal(err)
	}
	compareLines(t, &lines, expectedLines)

	expectedLines = []string{
		"Hello universe",
		"Replace me",
		"Replace me",
		"Next last line",
		"Yep, replaced",
	}

	if err := ReplaceLineInLines(&lines, "Hello world", "Hello universe", -1, false, false); err != nil {
		t.Fatal(err)
	}
	compareLines(t, &lines, expectedLines)

	lines = []string{
		"Replace me",
		"Every other line",
		"Replace me",
		"Another line",
		"Replace with something",
	}

	expectedLines = []string{
		"Replaced",
		"Every other line",
		"Replaced",
		"Another line",
		"Replaced",
	}

	if err := ReplaceLineInLines(&lines, "Replace", "Replaced", 3, false, false); err != nil {
		t.Fatal(err)
	}
	compareLines(t, &lines, expectedLines)

}
