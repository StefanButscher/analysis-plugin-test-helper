package logsAnalyzer

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestGetAnalyzer(t *testing.T) {
	t.Run("AnalysisPluginAnalyzer", func(t *testing.T) {
		analyzer := GetAnalyzer("ANALYSIS_PLUGIN")
		if analyzer == nil {
			t.Fatal("Expected analyzer, got nil")
		}
		if analyzer.GetAnalyzerType() != "ANALYSIS_PLUGIN" {
			t.Errorf("Expected ANALYSIS_PLUGIN, got %s", analyzer.GetAnalyzerType())
		}
	})

	t.Run("ESLintAnalyzer", func(t *testing.T) {
		analyzer := GetAnalyzer("ESLINT")
		if analyzer == nil {
			t.Fatal("Expected analyzer, got nil")
		}
		if analyzer.GetAnalyzerType() != "ESLINT" {
			t.Errorf("Expected ESLINT, got %s", analyzer.GetAnalyzerType())
		}
	})

	t.Run("InvalidAnalyzer", func(t *testing.T) {
		analyzer := GetAnalyzer("INVALID")
		if analyzer != nil {
			t.Errorf("Expected nil for invalid analyzer, got %v", analyzer)
		}
	})
}

func TestAnalysisPluginLogAnalyzer(t *testing.T) {
	analyzer := &AnalysisPluginLogAnalyzer{}

	t.Run("NoErrors", func(t *testing.T) {
		// Create temp file with no errors
		tempFile, err := ioutil.TempFile("", "analysis.log")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())

		content := `[INFO] Scanning for projects...
[INFO] Building project
[WARNING] Some warning
[INFO] BUILD SUCCESS`
		
		tempFile.WriteString(content)
		tempFile.Close()

		errors, err := analyzer.AnalyzeLogFile(tempFile.Name())
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(errors) != 0 {
			t.Errorf("Expected 0 errors, got %d", len(errors))
		}
	})

	t.Run("WithErrors", func(t *testing.T) {
		// Create temp file with errors from sample data
		tempFile, err := ioutil.TempFile("", "analysis.log")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())

		content := `[INFO] Scanning for projects...
[ERROR] Quality issue: JS_DEBUGGER_STATEMENT:Very High!
[ERROR] In file src/model/formatter.js (38:8)
[WARNING] Some warning
[ERROR] Quality issue: JS_CONSOLE_LOG:Very High!
[INFO] BUILD FAILURE`
		
		tempFile.WriteString(content)
		tempFile.Close()

		errors, err := analyzer.AnalyzeLogFile(tempFile.Name())
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(errors) != 3 {
			t.Errorf("Expected 3 errors, got %d", len(errors))
		}
		
		// Check first error content
		if !strings.Contains(errors[0].Error(), "JS_DEBUGGER_STATEMENT") {
			t.Errorf("Expected first error to contain JS_DEBUGGER_STATEMENT, got %s", errors[0].Error())
		}
	})

	t.Run("NonexistentFile", func(t *testing.T) {
		errors, err := analyzer.AnalyzeLogFile("/nonexistent.log")
		if err == nil {
			t.Error("Expected error for nonexistent file")
		}
		if errors != nil {
			t.Errorf("Expected nil errors for nonexistent file, got %v", errors)
		}
	})
}

func TestEsLintLogAnalyzer(t *testing.T) {
	analyzer := &EsLintLogAnalyzer{}

	t.Run("NoErrors", func(t *testing.T) {
		// Create temp file with no errors
		tempFile, err := ioutil.TempFile("", "eslint.log")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())

		content := `✨  Done in 2.45s.`
		
		tempFile.WriteString(content)
		tempFile.Close()

		errors, err := analyzer.AnalyzeLogFile(tempFile.Name())
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(errors) != 0 {
			t.Errorf("Expected 0 errors, got %d", len(errors))
		}
	})

	t.Run("WithExecutionError", func(t *testing.T) {
		// Create temp file with execution error
		tempFile, err := ioutil.TempFile("", "eslint.log")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())

		content := `Oops! Something went wrong! :(
Some additional error details`
		
		tempFile.WriteString(content)
		tempFile.Close()

		errors, err := analyzer.AnalyzeLogFile(tempFile.Name())
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(errors) != 1 {
			t.Errorf("Expected 1 error, got %d", len(errors))
		}
		
		// Check that it's an execution error
		if !strings.Contains(errors[0].Error(), "Oops! Something went wrong!") {
			t.Errorf("Expected execution error, got %s", errors[0].Error())
		}
	})

	t.Run("WithLintingErrors", func(t *testing.T) {
		// Create temp file with linting errors from sample data
		tempFile, err := ioutil.TempFile("", "eslint.log")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())

		content := `/Users/alexis/Projects/FioriPipelines/TestApps/ca.infra.testapp/webapp/Component.js
  1:1  warning  Legacy jQuery.sap usage is not allowed due to strict Content Security Policy  sap-ui5-legacy-jquerysap-usage
  3:1  warning  Legacy jQuery.sap usage is not allowed due to strict Content Security Policy  sap-ui5-legacy-jquerysap-usage

/Users/alexis/Projects/FioriPipelines/TestApps/ca.infra.testapp/webapp/model/formatter.js
  38:5  error    Debugger statement should not be part of the code that is committed to GIT!  sap-no-debugger
  38:5  error    Unexpected 'debugger' statement                                              no-debugger

✖ 185 problems (2 errors, 183 warnings)`
		
		tempFile.WriteString(content)
		tempFile.Close()

		errors, err := analyzer.AnalyzeLogFile(tempFile.Name())
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(errors) == 0 {
			t.Error("Expected at least some linting errors, got none")
		} else {
			t.Logf("Found %d errors in ESLint log", len(errors))
			// Check that we found some errors with file paths
			foundFileError := false
			for _, err := range errors {
				if strings.Contains(err.Error(), ".js") {
					foundFileError = true
					break
				}
			}
			if !foundFileError {
				t.Error("Expected at least one error with file path (.js)")
			}
		}
	})

	t.Run("NonexistentFile", func(t *testing.T) {
		errors, err := analyzer.AnalyzeLogFile("/nonexistent.log")
		if err == nil {
			t.Error("Expected error for nonexistent file")
		}
		if errors != nil {
			t.Errorf("Expected nil errors for nonexistent file, got %v", errors)
		}
	})
}