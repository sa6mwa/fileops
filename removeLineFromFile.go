package fileops

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
)

// RemoveLineFromFile removes line n number of times (or all of them
// if n is -1) from textfile. If before and/or after are not nil, the
// line before and/or after line to be removed must contain the
// after/before string respectively. If both before and after are nil,
// line is removed from anywhere in the file.
func RemoveLineFromFile(textfile, line string, n int, before, after *string, matchFullStringNotJustPrefix, matchWithLeadingAndTrailingSpaces bool) error {
	f, err := os.OpenFile(textfile, os.O_RDWR, 0644)
	if err != nil {
		return orExit(err)
	}
	defer f.Close()

	var lines []string
	var originalLines []string

	// Read all lines from textfile
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return orExit(err)
	}

	if DryRun {
		originalLines = make([]string, len(lines))
		copy(originalLines, lines)
		fmt.Fprintf(os.Stderr, "RemoveLineFromFile(%q, %q, %d, %+v, %+v, %t, %t)\n", textfile, line, n, before, after, matchFullStringNotJustPrefix, matchWithLeadingAndTrailingSpaces)
	}

	// Remove the target line up to `n` times
	lineRemoved := false
	removalCount := 0
	filteredLines := []string{}

	for i := 0; i < len(lines); i++ {
		if n != -1 && removalCount >= n {
			// If we've removed `n` lines, append the rest unchanged
			filteredLines = append(filteredLines, lines[i:]...)
			break
		}

		currentLine := lines[i]
		trimmedLine := currentLine
		if !matchWithLeadingAndTrailingSpaces {
			trimmedLine = strings.TrimSpace(currentLine)
		}

		matches := false
		if matchFullStringNotJustPrefix {
			matches = trimmedLine == line
		} else {
			matches = strings.HasPrefix(trimmedLine, line)
		}

		if matches {
			// Check `before` and `after` conditions
			beforeMatch := before == nil || (i > 0 && strings.Contains(lines[i-1], *before))
			afterMatch := after == nil || (i < len(lines)-1 && strings.Contains(lines[i+1], *after))

			if beforeMatch && afterMatch {
				lineRemoved = true
				removalCount++
				continue // Skip this line
			}
		}

		// Keep this line
		filteredLines = append(filteredLines, currentLine)
	}

	// If no line is removed, exit early
	if !lineRemoved {
		return nil
	}

	lines = filteredLines

	if DryRun {
		// Show diff
		origStrings := strings.Join(originalLines, "\n")
		edits := myers.ComputeEdits(span.URIFromPath(path.Join("a", textfile)), origStrings, strings.Join(lines, "\n"))
		diff := fmt.Sprint(gotextdiff.ToUnified(path.Join("a", textfile), path.Join("b", textfile), origStrings, edits))
		if len(diff) > 0 {
			fmt.Fprintln(os.Stderr, diff)
		}
		return nil
	}

	// Write lines back to textfile
	if err := f.Truncate(0); err != nil {
		return orExit(err)
	}
	if _, err := f.Seek(0, 0); err != nil {
		return orExit(err)
	}
	writer := bufio.NewWriter(f)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return orExit(err)
		}
	}
	return orExit(writer.Flush())
}
