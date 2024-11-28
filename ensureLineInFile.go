package fileops

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

// EnsureLineInFile ensures line is in textfile, optionally before
// and/or after string(s) unless before/after is/are nil. Will trim
// leading and trailing spaces from line read from textfile before
// matching unless matchWithLeadingAndTrailingSpaces is true. Will
// treat after and before as prefix unless
// matchFullStringNotJustPrefix. Returns error on failure.
func EnsureLineInFile(textfile, line string, before, after *string, matchFullStringNotJustPrefix, matchWithLeadingAndTrailingSpaces bool) error {
	f, err := os.OpenFile(textfile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	var lines []string

	// Read all lines from textfile
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// Ensure line is in lines slice, lines slice will be modified
	if err := EnsureLineInLines(&lines, line, before, after, matchFullStringNotJustPrefix, matchWithLeadingAndTrailingSpaces); err != nil {
		return err
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
