package executor

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNewAnalysisPluginExecutor(t *testing.T) {
	// Create temporary files for testing
	tempDir, err := ioutil.TempDir("", "test_target")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tempJar, err := ioutil.TempFile("", "test.jar")
	if err != nil {
		t.Fatalf("Failed to create temp jar: %v", err)
	}
	defer os.Remove(tempJar.Name())
	tempJar.Close()

	t.Run("ValidInputs", func(t *testing.T) {
		executor := NewAnalysisPluginExecutor(tempDir, "2.1.8-A", tempJar.Name())
		if executor == nil {
			t.Fatal("Expected executor to be created")
		}
		
		if executor.GetLogAnalyzerType() != "ANALYSIS_PLUGIN" {
			t.Errorf("Expected ANALYSIS_PLUGIN, got %s", executor.GetLogAnalyzerType())
		}
	})

	t.Run("InvalidTargetDir", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for invalid target directory")
			}
		}()
		NewAnalysisPluginExecutor("/nonexistent", "2.1.8-A", tempJar.Name())
	})

	t.Run("InvalidJarFile", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for invalid JAR file")
			}
		}()
		NewAnalysisPluginExecutor(tempDir, "2.1.8-A", "/nonexistent.jar")
	})
}


func TestAnalysisPluginExecutor_RunCheck(t *testing.T) {
	// This test would require actual Maven setup, so we'll create a minimal test
	// In a real scenario, you'd mock the command execution
	t.Skip("Skipping integration test - requires Maven setup")
}

func TestAnalysisPluginExecutor_RunCheckAgainstFile(t *testing.T) {
	// This test would require actual Java setup, so we'll create a minimal test
	// In a real scenario, you'd mock the command execution
	t.Skip("Skipping integration test - requires Java setup")
}