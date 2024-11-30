package fileops

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// PutFile writes content into local file destination with mode
// filePerm. If destination is a path with directories they will be
// created using optional dirPerm or mode 0755 by default. Returns
// error if something failed.
func PutFile(destination, content string, filePerm os.FileMode, dirPerm ...os.FileMode) error {
	// Determine directory permissions
	var directoryPermission os.FileMode = 0755
	if len(dirPerm) > 0 {
		directoryPermission = dirPerm[0]
	}

	// Get the directory path from the destination
	dirPath := filepath.Dir(destination)

	// Make sure content ends with a new line
	newline := "\n"
	if os.PathSeparator == '\\' {
		// Windows
		newline = "\r\n"
	}
	if !strings.HasSuffix(content, newline) {
		content += newline
	}

	if DryRun {
		fmt.Fprintf(os.Stderr, "os.MkdirAll(%q, %v)\n", dirPath, directoryPermission)
	} else {
		// Create directories if they do not exist
		err := os.MkdirAll(dirPath, directoryPermission)
		if err != nil {
			return orExit(fmt.Errorf("failed to create directories: %w", err))
		}
	}

	if DryRun {
		fmt.Fprintf(os.Stderr, "os.WriteFile(%q, %q, %v)\n", destination, content, filePerm)
	} else {
		// Write the file
		if err := os.WriteFile(destination, []byte(content), filePerm); err != nil {
			return orExit(fmt.Errorf("failed to write file: %w", err))
		}
	}

	if DryRun {
		fmt.Fprintf(os.Stderr, "os.Chmod(%q, %v)\n", destination, filePerm)
	} else {
		os.Chmod(destination, filePerm)
	}

	return nil
}

// PutFileFromFS copies a file or recursively copies a directory from
// an fs.FS interface to a target path on the local
// filesystem. Returns error in case of failure.
func PutFileFromFS(fsys fs.FS, source string, destination string, filePerm os.FileMode, dirPerm ...os.FileMode) error {
	if DryRun {
		if len(dirPerm) > 0 {
			fmt.Fprintf(os.Stderr, "PutFileFromFS(<fs>, %q, %q, %v, %v)\n", source, destination, filePerm, dirPerm[0])
		} else {
			fmt.Fprintf(os.Stderr, "PutFileFromFS(<fs>, %q, %q, %v)\n", source, destination, filePerm)
		}

	}

	var directoryPermission os.FileMode = 0755
	if len(dirPerm) > 0 {
		directoryPermission = dirPerm[0]
	}

	// Get the file information from the source path.
	srcInfo, err := fs.Stat(fsys, source)
	if err != nil {
		return orExit(fmt.Errorf("failed to stat source path: %w", err))
	}

	// Handle directories recursively.
	if srcInfo.IsDir() {
		return orExit(copyDir(fsys, source, destination, filePerm, directoryPermission))
	}

	// Handle single file copy.
	return orExit(copyFile(fsys, source, destination, filePerm, directoryPermission))
}

// copyDir recursively copies a directory and its contents.
func copyDir(fsys fs.FS, srcDir string, destDir string, filePerm os.FileMode, dirPerm os.FileMode) error {
	err := fs.WalkDir(fsys, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return orExit(fmt.Errorf("failed to walk directory: %w", err))
		}

		// Compute the relative path and destination path.
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return orExit(fmt.Errorf("failed to compute relative path: %w", err))
		}
		destPath := filepath.Join(destDir, relPath)

		// Handle directories.
		if d.IsDir() {
			if DryRun {
				fmt.Fprintf(os.Stderr, "os.MkdirAll(%q, %v)\n", destPath, dirPerm)
				return nil
			}
			if err := os.MkdirAll(destPath, dirPerm); err != nil {
				return orExit(fmt.Errorf("failed to create directory: %w", err))
			}
			return nil
		}

		// Handle files.
		return orExit(copyFile(fsys, path, destPath, filePerm, dirPerm))
	})

	return orExit(err)
}

// copyFile copies a single file from fs.FS to the local filesystem.
func copyFile(fsys fs.FS, srcFile string, destFile string, filePerm os.FileMode, dirPerm os.FileMode) error {
	// Open the source file.
	src, err := fsys.Open(srcFile)
	if err != nil {
		return orExit(fmt.Errorf("failed to open source file: %w", err))
	}
	defer src.Close()

	// Create the destination file's directory.
	destDir := filepath.Dir(destFile)
	if DryRun {
		fmt.Fprintf(os.Stderr, "os.MkdirAll(%q, %v)\n", destDir, dirPerm)
		fmt.Fprintf(os.Stderr, "%q <- %q\n", destFile, srcFile)
	} else {
		if err := os.MkdirAll(destDir, dirPerm); err != nil {
			return orExit(fmt.Errorf("failed to create destination directory: %w", err))
		}
		// Create the destination file.
		dest, err := os.OpenFile(destFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filePerm)
		//dest, err := os.Create(destFile)
		if err != nil {
			return orExit(fmt.Errorf("failed to create destination file: %w", err))
		}
		defer dest.Close()
		// Copy the file content.
		_, err = io.Copy(dest, src)
		if err != nil {
			return orExit(fmt.Errorf("failed to copy file content: %w", err))
		}
	}

	return nil
}

// ListFiles recursively lists all files in the given fs.FS starting
// from the root directory.
func ListFiles(fsys fs.FS, root string) ([]string, error) {
	if DryRun {
		fmt.Fprintf(os.Stderr, "ListFiles(<fs>, %q)\n", root)
	}
	var files []string
	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, path) // Collect the file path
			if DryRun {
				fmt.Fprintln(os.Stderr, path)
			}
		}
		return nil
	})
	return files, err
}
