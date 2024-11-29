package fileops

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
)

// EnsureLineInFile ensures line is in textfile, optionally before
// and/or after string(s) unless before/after is/are nil. Will trim
// leading and trailing spaces from line read from textfile before
// matching unless matchWithLeadingAndTrailingSpaces is true. Will
// treat after and before as prefix unless
// matchFullStringNotJustPrefix. If optional filePerm is specified,
// the first item in the slice is used as file mode if textfile does
// not exist. Returns error on failure.
func EnsureLineInFile(textfile, line string, before, after *string, matchFullStringNotJustPrefix, matchWithLeadingAndTrailingSpaces bool, filePerm ...os.FileMode) error {
	var fileMode os.FileMode = 0644
	if len(filePerm) > 0 {
		fileMode = filePerm[0]
	}
	f, err := os.OpenFile(textfile, os.O_RDWR|os.O_CREATE, fileMode)
	if err != nil {
		return err
	}
	defer f.Close()

	var lines []string
	var originalLines []string

	// Read all lines from textfile
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// If after and before is nil, avoid re-writing the file if the
		// exact line already exists in the file.
		if before == nil && after == nil && scanner.Text() == line {
			return nil
		}
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if DryRun {
		originalLines = make([]string, len(lines))
		copy(originalLines, lines)
		fmt.Fprintf(os.Stderr, "EnsureLineInFile(%q, %q, %v, %v, %t, %t)\n", textfile, line, before, after, matchFullStringNotJustPrefix, matchWithLeadingAndTrailingSpaces)
	}

	// Ensure line is in lines slice, lines slice will be modified
	if err := EnsureLineInLines(&lines, line, before, after, matchFullStringNotJustPrefix, matchWithLeadingAndTrailingSpaces); err != nil {
		return err
	}

	if DryRun {
		// Show diff
		origStrings := strings.Join(originalLines, "\n")
		edits := myers.ComputeEdits(span.URIFromPath(path.Join("a", textfile)), origStrings, strings.Join(lines, "\n"))
		diff := fmt.Sprint(gotextdiff.ToUnified(path.Join("a", textfile), path.Join("b", textfile), origStrings, edits))
		fmt.Fprintln(os.Stderr, diff)
		return nil
	}

	// Write lines back to textfile
	if err := f.Truncate(0); err != nil {
		return err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}
	writer := bufio.NewWriter(f)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return writer.Flush()
}

// EnsureLineInLines ensures line is in lines string pointer slice,
// optionally before and/or after string(s) unless before/after is/are
// nil. Will trim leading and trailing spaces from line read from
// lines before matching unless matchWithLeadingAndTrailingSpaces is
// true. Will treat after and before as prefix unless
// matchFullStringNotJustPrefix. Returns error on failure.
func EnsureLineInLines(lines *[]string, line string, before, after *string, matchFullStringNotJustPrefix, matchWithLeadingAndTrailingSpaces bool) error {
	if lines == nil {
		return errors.New("nil pointer")
	}

	// Deref lines pointer
	slice := *lines

	// Function to find indices
	findLine := func(target string) int {
		for i, l := range slice {
			if matchWithLeadingAndTrailingSpaces {
				if matchFullStringNotJustPrefix {
					if l == target {
						return i
					}
				} else {
					if strings.HasPrefix(l, target) {
						return i
					}
				}
			} else {
				if matchFullStringNotJustPrefix {
					if strings.TrimSpace(l) == target {
						return i
					}
				} else {
					if strings.HasPrefix(strings.TrimSpace(l), target) {
						return i
					}
				}
			}
		}
		return -1
	}

	// Remove line if it already exists
	lineIndex := findLine(line)
	if lineIndex != -1 {
		slice = append(slice[:lineIndex], slice[lineIndex+1:]...)
	}

	// Determine where to insert the line
	insertIndex := len(slice) // default: at the end
	if after != nil {
		afterIndex := findLine(*after)
		if afterIndex != -1 {
			insertIndex = afterIndex + 1
		}
	}
	if before != nil {
		beforeIndex := findLine(*before)
		if beforeIndex != -1 {
			insertIndex = beforeIndex
		}
	}

	// Insert the line at the correct position
	if insertIndex > len(slice) {
		slice = append(slice, line)
	} else {
		slice = append(slice[:insertIndex], append([]string{line}, slice[insertIndex:]...)...)
	}

	*lines = slice
	return nil
}
