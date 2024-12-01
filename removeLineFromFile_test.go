package fileops

import (
	"os"
	"strings"
	"testing"
)

func TestRemoveLineFromFile(t *testing.T) {
	// Helper function to create a temporary file
	createTempFile := func(t *testing.T, content string) *os.File {
		t.Helper()
		tmpFile, err := os.CreateTemp("", "testfile_*.txt")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		_, err = tmpFile.WriteString(content)
		if err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		if err := tmpFile.Close(); err != nil {
			t.Fatalf("Failed to close temp file: %v", err)
		}
		return tmpFile
	}

	tests := []struct {
		name                         string
		initialContent               string
		line                         string
		n                            int
		before, after                *string
		matchFullString, matchSpaces bool
		expectedContent              string
		expectError                  bool
	}{
		// Adjusted test cases
		{
			name:            "Remove single line without conditions",
			initialContent:  "line1\nline2\nline3\n",
			line:            "line2",
			n:               1,
			before:          nil,
			after:           nil,
			matchFullString: true,
			matchSpaces:     true,
			expectedContent: "line1\nline3\n",
			expectError:     false,
		},
		{
			name:            "Remove multiple lines with n=-1",
			initialContent:  "line1\nline2\nline2\nline3\n",
			line:            "line2",
			n:               -1,
			before:          nil,
			after:           nil,
			matchFullString: true,
			matchSpaces:     true,
			expectedContent: "line1\nline3\n",
			expectError:     false,
		},
		{
			name:            "Remove line with before/after conditions",
			initialContent:  "line1\nbefore\nline2\nafter\nline3\n",
			line:            "line2",
			n:               1,
			before:          ptr("before"),
			after:           ptr("after"),
			matchFullString: true,
			matchSpaces:     true,
			expectedContent: "line1\nbefore\nafter\nline3\n",
			expectError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile := createTempFile(t, tt.initialContent)
			defer os.Remove(tmpFile.Name())

			err := RemoveLineFromFile(tmpFile.Name(), tt.line, tt.n, tt.before, tt.after, tt.matchFullString, tt.matchSpaces)
			if (err != nil) != tt.expectError {
				t.Fatalf("Expected error: %v, got: %v", tt.expectError, err)
			}

			updatedContent, err := os.ReadFile(tmpFile.Name())
			if err != nil {
				t.Fatalf("Failed to read updated file: %v", err)
			}

			if strings.TrimSpace(string(updatedContent)) != strings.TrimSpace(tt.expectedContent) {
				t.Errorf("Expected content:\n%q\nGot:\n%q\n", tt.expectedContent, updatedContent)
			}
		})
	}
}

// Helper to create a pointer to a string
func ptr(s string) *string {
	return &s
}
