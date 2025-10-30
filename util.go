package main

import (
	"analysis-migration-test/errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// CopyFile copies a file from srcFilePath to destFilePath.
// It returns an error if something goes wrong.
func CopyFile(srcFilePath, destFilePath string) error {
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destFilePath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Ensure data is flushed to disk
	err = destFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}

// Delete removes the specified file.
// It returns an error if something goes wrong.
func Delete(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true // file exists
	}
	if os.IsNotExist(err) {
		return false // file does not exist
	}
	return false // some other error (e.g., permission issue)
}

func FindAnalysisPluginIssue(checkErrors []error, ruleName string) error {
	for _, err := range checkErrors {
		codeErr, ok := err.(errors.CodeChecksError)
		if !ok {
			continue
		}

		if strings.Contains(codeErr.RowText, "Quality issue") && strings.Contains(codeErr.RowText, ruleName) {
			return err
		}
	}
	return nil
}

func FindEslintIssue(checkErrors []error, ruleName string) error {
	for _, err := range checkErrors {
		codeErr, ok := err.(errors.CodeChecksError)
		if !ok {
			continue
		}
		if strings.Contains(codeErr.RowText, "error") && strings.Contains(codeErr.RowText, ":") && strings.Contains(codeErr.RowText, ruleName) {
			return err
		}
	}
	return nil
}
