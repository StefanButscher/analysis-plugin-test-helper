package executor

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNewEsLintExecutor(t *testing.T) {
	// Create temporary files for testing
	tempDir, err := ioutil.TempDir("", "test_target")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tempConfig, err := ioutil.TempFile("", "eslintrc")
	if err != nil {
		t.Fatalf("Failed to create temp config: %v", err)
	}
	defer os.Remove(tempConfig.Name())
	tempConfig.Close()

	tempRulesDir, err := ioutil.TempDir("", "rules")
	if err != nil {
		t.Fatalf("Failed to create temp rules dir: %v", err)
	}
	defer os.RemoveAll(tempRulesDir)

	// Create temp binary file
	tempBinary, err := ioutil.TempFile("", "eslint.js")
	if err != nil {
		t.Fatalf("Failed to create temp binary: %v", err)
	}
	defer os.Remove(tempBinary.Name())
	tempBinary.Close()

	t.Run("ValidInputs", func(t *testing.T) {
		executor := NewEsLintExecutor(tempDir, tempConfig.Name(), tempRulesDir, tempBinary.Name())
		if executor == nil {
			t.Fatal("Expected executor to be created")
		}
		
		if executor.GetLogAnalyzerType() != "ESLINT" {
			t.Errorf("Expected ESLINT, got %s", executor.GetLogAnalyzerType())
		}
	})

	t.Run("InvalidTargetDir", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for invalid target directory")
			}
		}()
		NewEsLintExecutor("/nonexistent", tempConfig.Name(), tempRulesDir, tempBinary.Name())
	})

	t.Run("InvalidConfigFile", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for invalid config file")
			}
		}()
		NewEsLintExecutor(tempDir, "/nonexistent.json", tempRulesDir, tempBinary.Name())
	})

	t.Run("InvalidRulesDir", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for invalid rules directory")
			}
		}()
		NewEsLintExecutor(tempDir, tempConfig.Name(), "/nonexistent", tempBinary.Name())
	})

	t.Run("InvalidBinaryPath", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for invalid binary path")
			}
		}()
		NewEsLintExecutor(tempDir, tempConfig.Name(), tempRulesDir, "/nonexistent/eslint.js")
	})
}


func TestEsLintExecutor_RunCheck(t *testing.T) {
	// This test would require actual Node.js/ESLint setup, so we'll create a minimal test
	// In a real scenario, you'd mock the command execution
	t.Skip("Skipping integration test - requires Node.js/ESLint setup")
}

func TestEsLintExecutor_RunCheckAgainstFile(t *testing.T) {
	// This test would require actual Node.js/ESLint setup, so we'll create a minimal test
	// In a real scenario, you'd mock the command execution
	t.Skip("Skipping integration test - requires Node.js/ESLint setup")
}