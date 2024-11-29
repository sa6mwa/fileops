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

func ReplaceLineInFile(textfile, lineToReplace, replaceWithLine string, n int, matchFullStringNotJustPrefix, matchWithLeadingAndTrailingSpaces bool) error {
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
		fmt.Fprintf(os.Stderr, "ReplaceLineInFile(%q, %q, %q, %d, %t, %t)\n", textfile, lineToReplace, replaceWithLine, n, matchFullStringNotJustPrefix, matchWithLeadingAndTrailingSpaces)
	}

	// Replace lineToReplace with replaceWithLine in lines slice, liens
	// slice will be modified
	if err := ReplaceLineInLines(&lines, lineToReplace, replaceWithLine, n, matchFullStringNotJustPrefix, matchWithLeadingAndTrailingSpaces); err != nil {
		return orExit(err)
	}

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

func ReplaceLineInLines(lines *[]string, lineToReplace string, replaceWithLine string, n int, matchFullStringNotJustPrefix, matchWithLeadingAndTrailingSpaces bool) error {
	if lines == nil {
		return orExit(errors.New("nil pointer"))
	}

	// Deref lines pointer
	slice := *lines

	// Function to find indices
	findLine := func(target string, startIndex int) int {
		for i, l := range slice {
			if i < startIndex {
				continue
			}
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

	lineIndex := 0

	if n == -1 {
		n = len(slice)
	}
	for i := 0; i < n; i++ {
		// Find lineToReplace
		lineIndex = findLine(lineToReplace, lineIndex)
		if lineIndex == -1 {
			break // nothing (more) to replace
		}
		// Replace lineToReplace with replaceWithLine
		slice[lineIndex] = replaceWithLine
		if lineIndex < len(slice)-1 {
			lineIndex++
		} else {
			break
		}
	}

	*lines = slice
	return nil
}
