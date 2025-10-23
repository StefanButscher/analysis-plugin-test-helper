package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestExecutorLogFileCreation is an integration test that verifies executors properly create log files
func TestExecutorLogFileCreation(t *testing.T) {
	// Set timeout to 5 minutes
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create execution directory for this test run
	timestamp := time.Now().Format("20060102-150405")
	executionDir := filepath.Join("__executions__", fmt.Sprintf("LogFileTest-%s", timestamp))

	err := os.MkdirAll(executionDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create execution directory: %v", err)
	}

	t.Logf("Integration test execution directory: %s", executionDir)

	// Test AnalysisPluginExecutor log file creation
	t.Run("AnalysisPluginLogFileCreation", func(t *testing.T) {

		testName := "LogFileTest"

		t.Logf("Testing AnalysisPluginExecutor log file creation")

		// Execute the analysis
		execErr, logPath := analysisPluginA.RunCheck(executionDir, testName)

		// Verify log file was created regardless of execution result
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			t.Errorf("Expected log file to be created at %s, but file does not exist", logPath)
			return
		}

		t.Logf("✓ Log file created successfully at: %s", logPath)

		// Verify log file is not empty
		fileInfo, err := os.Stat(logPath)
		if err != nil {
			t.Errorf("Failed to get log file info: %v", err)
			return
		}

		if fileInfo.Size() == 0 {
			t.Errorf("Log file is empty, expected some content")
			return
		}

		t.Logf("✓ Log file has content: %d bytes", fileInfo.Size())

		// Verify log file contains expected Maven output patterns
		content, err := os.ReadFile(logPath)
		if err != nil {
			t.Errorf("Failed to read log file: %v", err)
			return
		}

		contentStr := string(content)
		expectedPatterns := []string{
			"[INFO]", // Maven info logs
		}

		for _, pattern := range expectedPatterns {
			if !strings.Contains(contentStr, pattern) {
				t.Logf("WARNING: Expected pattern '%s' not found in log file", pattern)
			} else {
				t.Logf("✓ Found expected pattern '%s' in log file", pattern)
			}
		}

		// Log execution result
		if execErr != nil {
			t.Logf("AnalysisPlugin execution result: %v (this may be expected)", execErr)
		} else {
			t.Logf("✓ AnalysisPlugin execution completed successfully")
		}

		// Verify expected log file naming pattern
		expectedLogName := fmt.Sprintf("%s-analysisPlugin-2.1.8-A.log", testName)
		if !strings.HasSuffix(logPath, expectedLogName) {
			t.Errorf("Log file name doesn't match expected pattern. Expected to end with '%s', got '%s'", expectedLogName, logPath)
		} else {
			t.Logf("✓ Log file naming pattern is correct: %s", expectedLogName)
		}
	})

	// Test EsLintExecutor log file creation
	t.Run("EsLintLogFileCreation", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatalf("Test timed out before ESLint execution: %v", ctx.Err())
		default:
		}

		testName := "LogFileTest"

		t.Logf("Testing EsLintExecutor log file creation")

		// Execute the linting
		execErr, logPath := eslintPlugin.RunCheck(executionDir, testName)

		// Verify log file was created regardless of execution result
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			t.Errorf("Expected log file to be created at %s, but file does not exist", logPath)
			return
		}

		t.Logf("✓ Log file created successfully at: %s", logPath)

		// Verify log file is not empty
		fileInfo, err := os.Stat(logPath)
		if err != nil {
			t.Errorf("Failed to get log file info: %v", err)
			return
		}

		if fileInfo.Size() == 0 {
			t.Errorf("Log file is empty, expected some content")
			return
		}

		t.Logf("✓ Log file has content: %d bytes", fileInfo.Size())

		// Read and analyze log content
		content, err := os.ReadFile(logPath)
		if err != nil {
			t.Errorf("Failed to read log file: %v", err)
			return
		}

		contentStr := string(content)

		// Look for typical ESLint output patterns
		if strings.Contains(contentStr, ".js") ||
			strings.Contains(contentStr, "error") ||
			strings.Contains(contentStr, "warning") ||
			strings.Contains(contentStr, "problems") {
			t.Logf("✓ Found expected ESLint output patterns in log file")
		} else {
			t.Logf("WARNING: Expected ESLint output patterns not found, but log file was created")
		}

		// Log execution result
		if execErr != nil {
			t.Logf("ESLint execution result: %v (this may be expected for files with linting issues)", execErr)
		} else {
			t.Logf("✓ ESLint execution completed successfully")
		}

		// Verify expected log file naming pattern
		expectedLogName := fmt.Sprintf("%s-eslint.log", testName)
		if !strings.HasSuffix(logPath, expectedLogName) {
			t.Errorf("Log file name doesn't match expected pattern. Expected to end with '%s', got '%s'", expectedLogName, logPath)
		} else {
			t.Logf("✓ Log file naming pattern is correct: %s", expectedLogName)
		}
	})

	// Test RunCheckAgainstFile for both executors
	t.Run("FileSpecificLogCreation", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatalf("Test timed out before file-specific execution: %v", ctx.Err())
		default:
		}

		testName := "FileTest"

		// Test file path (assuming this exists in the target directory)
		testFilePath := filepath.Join(targetDir, "webapp", "Component.js")

		// Verify test file exists before running
		if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
			t.Skipf("Test file %s does not exist, skipping file-specific tests", testFilePath)
		}

		// Test AnalysisPlugin file-specific execution
		t.Run("AnalysisPluginFileSpecific", func(t *testing.T) {
			execErr, logPath := analysisPluginA.RunCheckAgainstFile(executionDir, testName, testFilePath)

			// Verify log file was created
			if _, err := os.Stat(logPath); os.IsNotExist(err) {
				t.Errorf("Expected log file to be created at %s", logPath)
				return
			}

			t.Logf("✓ AnalysisPlugin file-specific log created: %s", logPath)

			// Verify naming pattern for file-specific logs
			expectedLogName := fmt.Sprintf("%s-analysisPlugin-2.1.8-A-file.log", testName)
			if !strings.HasSuffix(logPath, expectedLogName) {
				t.Errorf("File-specific log name doesn't match pattern. Expected '%s', got '%s'", expectedLogName, filepath.Base(logPath))
			}

			if execErr != nil {
				t.Logf("AnalysisPlugin file execution result: %v", execErr)
			} else {
				t.Logf("✓ AnalysisPlugin file execution completed successfully")
			}
		})

		// Test ESLint file-specific execution
		t.Run("ESLintFileSpecific", func(t *testing.T) {
			execErr, logPath := eslintPlugin.RunCheckAgainstFile(executionDir, testName, testFilePath)

			// Verify log file was created
			if _, err := os.Stat(logPath); os.IsNotExist(err) {
				t.Errorf("Expected log file to be created at %s", logPath)
				return
			}

			t.Logf("✓ ESLint file-specific log created: %s", logPath)

			// Verify naming pattern for file-specific logs
			expectedLogName := fmt.Sprintf("%s-eslint-file.log", testName)
			if !strings.HasSuffix(logPath, expectedLogName) {
				t.Errorf("File-specific log name doesn't match pattern. Expected '%s', got '%s'", expectedLogName, filepath.Base(logPath))
			}

			if execErr != nil {
				t.Logf("ESLint file execution result: %v", execErr)
			} else {
				t.Logf("✓ ESLint file execution completed successfully")
			}
		})
	})

	// Test using the exported RunChecks and RunChecksAgainstFile functions
	t.Run("ExportedFunctionsTest", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Fatalf("Test timed out before exported functions test: %v", ctx.Err())
		default:
		}

		testName := "ExportedFunctionTest"
		
		t.Logf("Testing exported RunChecks function with AnalysisPluginA")
		
		// Test RunChecks function
		execErr, codeErrors := RunChecks(executionDir, testName, analysisPluginA)
		
		if execErr != nil {
			t.Logf("RunChecks returned execution error: %v", execErr)
		}
		
		t.Logf("RunChecks found %d code errors", len(codeErrors))
		for i, err := range codeErrors {
			t.Logf("  Code Error %d: %v", i+1, err)
		}
		
		// Verify log file was created
		expectedLogPath := filepath.Join(executionDir, fmt.Sprintf("%s-analysisPlugin-2.1.8-A.log", testName))
		if _, err := os.Stat(expectedLogPath); os.IsNotExist(err) {
			t.Errorf("Expected log file to be created by RunChecks at %s", expectedLogPath)
		} else {
			t.Logf("✓ RunChecks successfully created log file: %s", expectedLogPath)
		}
		
		// Test RunChecksAgainstFile function
		testFilePath := filepath.Join(targetDir, "webapp", "Component.js")
		if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
			t.Skipf("Test file %s does not exist, skipping RunChecksAgainstFile test", testFilePath)
		} else {
			t.Logf("Testing exported RunChecksAgainstFile function")
			
			fileTestName := "FileExportedTest"
			execErr, codeErrors := RunChecksAgainstFile(executionDir, fileTestName, analysisPluginA, testFilePath)
			
			if execErr != nil {
				t.Logf("RunChecksAgainstFile returned execution error: %v", execErr)
			}
			
			t.Logf("RunChecksAgainstFile found %d code errors", len(codeErrors))
			for i, err := range codeErrors {
				t.Logf("  File Code Error %d: %v", i+1, err)
			}
			
			// Verify log file was created
			expectedFileLogPath := filepath.Join(executionDir, fmt.Sprintf("%s-analysisPlugin-2.1.8-A-file.log", fileTestName))
			if _, err := os.Stat(expectedFileLogPath); os.IsNotExist(err) {
				t.Errorf("Expected log file to be created by RunChecksAgainstFile at %s", expectedFileLogPath)
			} else {
				t.Logf("✓ RunChecksAgainstFile successfully created log file: %s", expectedFileLogPath)
			}
		}
	})

	// Summary
	t.Logf("✓ Integration test completed. All log files saved in: %s", executionDir)

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
