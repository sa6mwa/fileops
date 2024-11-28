package fileops

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func ReplaceLineInFile(textfile, lineToReplace, replaceWithLine string, n int, matchFullStringNotJustPrefix, matchWithLeadingAndTrailingSpaces bool) error {
	f, err := os.OpenFile(textfile, os.O_RDWR, 0644)
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

	// Replace lineToReplace with replaceWithLine in lines slice, liens
	// slice will be modified
	if err := ReplaceLineInLines(&lines, lineToReplace, replaceWithLine, n, matchFullStringNotJustPrefix, matchWithLeadingAndTrailingSpaces); err != nil {
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

func ReplaceLineInLines(lines *[]string, lineToReplace string, replaceWithLine string, n int, matchFullStringNotJustPrefix, matchWithLeadingAndTrailingSpaces bool) error {
	if lines == nil {
		return errors.New("nil pointer")
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