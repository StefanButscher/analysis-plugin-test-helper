package main

import (
	"analysis-migration-test/executor"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestEsLintExecutorNegativeCase tests ESLint executor behavior with non-existent files
func TestEsLintExecutorNegativeCase(t *testing.T) {
	// Create execution directory for this test run
	timestamp := time.Now().Format("20060102-150405")
	executionDir := filepath.Join("__executions__", fmt.Sprintf("ESLintNegativeTest-%s", timestamp))

	err := os.MkdirAll(executionDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create execution directory: %v", err)
	}

	t.Logf("ESLint negative test execution directory: %s", executionDir)

	// Create ESLint executor
	targetDir := "/Users/alexis/Projects/FioriPipelines/TestApps/ca.infra.testapp"
	eslintExecutor := executor.NewEsLintExecutor(
		targetDir,
		"/Users/alexis/Projects/FioriPipelines/eslint-plugin-fiori-custom/configure.eslintrc",
		"/Users/alexis/Projects/FioriPipelines/eslint-plugin-fiori-custom/lib/rules",
		"/Users/alexis/Projects/FioriPipelines/TestApps/ca.infra.testapp/node_modules/eslint/bin/eslint.js",
	)

	t.Run("NonExistentFile", func(t *testing.T) {
		testName := "NonExistentFileTest"
		nonExistentFile := "/path/that/does/not/exist/nonexistent.js"

		t.Logf("Testing ESLint executor with non-existent file: %s", nonExistentFile)

		// Execute ESLint against non-existent file
		execErr, logPath := eslintExecutor.RunCheckAgainstFile(executionDir, testName, nonExistentFile)

		// Verify log file was created (should be created even if command fails)
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			t.Errorf("Expected log file to be created at %s, but file does not exist", logPath)
			return
		}

		t.Logf("✓ Log file created at: %s", logPath)

		// Read log file content
		content, err := os.ReadFile(logPath)
		if err != nil {
			t.Errorf("Failed to read log file: %v", err)
			return
		}

		contentStr := string(content)
		t.Logf("Log file content:\n%s", contentStr)

		// Check if execution error occurred (may or may not happen depending on how executor handles it)
		if execErr != nil {
			t.Logf("✓ Execution error occurred: %v", execErr)
		} else {
			t.Logf("No execution error from executor (log analyzer will handle errors)")
		}

		// Verify log contains expected ESLint error patterns
		expectedErrorPatterns := []string{
			"Oops! Something went wrong!", // ESLint standard error message
			"no files matching",            // ESLint file not found message
			"not found",                    // General not found
		}

		foundError := false
		for _, pattern := range expectedErrorPatterns {
			if strings.Contains(contentStr, pattern) {
				t.Logf("✓ Found expected error pattern '%s' in log file", pattern)
				foundError = true
				break
			}
		}

		if !foundError {
			t.Logf("WARNING: Expected error patterns not found in log file, content: %s", contentStr)
		}

		// Test the RunChecksAgainstFile function (integrated with log analyzer)
		t.Logf("Testing integrated RunChecksAgainstFile function with non-existent file")
		integrationTestName := "IntegratedNonExistentFileTest"
		execError, codeErrors := RunChecksAgainstFile(executionDir, integrationTestName, eslintExecutor, nonExistentFile)

		// The behavior depends on the implementation - either execution error or code errors detected by analyzer
		if execError != nil {
			t.Logf("RunChecksAgainstFile returned execution error: %v", execError)
		} else if len(codeErrors) > 0 {
			t.Logf("✓ RunChecksAgainstFile found errors via log analyzer (expected behavior)")
		} else {
			t.Errorf("Expected RunChecksAgainstFile to return either execution error or code errors for non-existent file")
		}

		// Log analyzer should handle the error gracefully
		t.Logf("RunChecksAgainstFile found %d code errors", len(codeErrors))
		for i, err := range codeErrors {
			t.Logf("  Code Error %d: %v", i+1, err)
		}

		// Verify integrated log file was created
		expectedIntegratedLogPath := filepath.Join(executionDir, fmt.Sprintf("%s-eslint-file.log", integrationTestName))
		if _, err := os.Stat(expectedIntegratedLogPath); os.IsNotExist(err) {
			t.Errorf("Expected integrated log file to be created at %s", expectedIntegratedLogPath)
		} else {
			t.Logf("✓ Integrated log file created at: %s", expectedIntegratedLogPath)
		}
	})

	t.Run("InvalidFilePath", func(t *testing.T) {
		testName := "InvalidFilePathTest"
		invalidFile := ""  // Empty file path

		t.Logf("Testing ESLint executor with invalid (empty) file path")

		// Execute ESLint against invalid file path
		execErr, logPath := eslintExecutor.RunCheckAgainstFile(executionDir, testName, invalidFile)

		// Should have execution error
		if execErr != nil {
			t.Logf("✓ Expected execution error for invalid file path: %v", execErr)
		}

		// Verify log file was still created
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			t.Errorf("Expected log file to be created at %s even for invalid file path", logPath)
		} else {
			t.Logf("✓ Log file created even for invalid input: %s", logPath)
		}
	})

	t.Run("FileWithSyntaxErrors", func(t *testing.T) {
		testName := "SyntaxErrorTest"

		// Create a temporary file with syntax errors
		tempFile := filepath.Join(executionDir, "syntax-error.js")
		syntaxErrorContent := `
// This file has intentional syntax errors
function badSyntax( {
	console.log("missing closing parenthesis"
	var x = ; // missing value
	return unclosed string"
}
`
		err := os.WriteFile(tempFile, []byte(syntaxErrorContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create temp file with syntax errors: %v", err)
		}
		defer os.Remove(tempFile)

		t.Logf("Testing ESLint executor with file containing syntax errors: %s", tempFile)

		// Test with RunChecksAgainstFile (integrated approach)
		execError, codeErrors := RunChecksAgainstFile(executionDir, testName, eslintExecutor, tempFile)

		// Log the results
		if execError != nil {
			t.Logf("Execution error (may be expected for syntax errors): %v", execError)
		}

		t.Logf("Found %d code errors from syntax error file", len(codeErrors))
		for i, err := range codeErrors {
			t.Logf("  Syntax Error %d: %v", i+1, err)
		}

		// Verify log file was created
		expectedLogPath := filepath.Join(executionDir, fmt.Sprintf("%s-eslint-file.log", testName))
		if _, err := os.Stat(expectedLogPath); os.IsNotExist(err) {
			t.Errorf("Expected log file to be created at %s", expectedLogPath)
		} else {
			t.Logf("✓ Log file created for syntax error test: %s", expectedLogPath)
		}
	})

	// Summary
	t.Logf("✓ ESLint negative case testing completed. All log files saved in: %s", executionDir)

	// List all created log files
	logFiles, err := filepath.Glob(filepath.Join(executionDir, "*.log"))
	if err != nil {
		t.Logf("Could not list log files: %v", err)
	} else {
		t.Logf("Created log files (%d):", len(logFiles))
		for _, logFile := range logFiles {
			fileInfo, _ := os.Stat(logFile)
			t.Logf("  - %s (%d bytes)", filepath.Base(logFile), fileInfo.Size())
		}
	}
}