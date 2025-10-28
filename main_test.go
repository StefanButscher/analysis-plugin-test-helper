package main

import (
	"analysis-migration-test/executor"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var targetDir = "/home/d050449/Documents/ca.infra.testapp"

// Pure Analysis Plugin
var analysisPluginA = executor.NewAnalysisPluginExecutor(
	targetDir,
	"2.1.9",
	"",
)

// Analysis Plugin with Disabled JS Checks
var analysisPluginB = executor.NewAnalysisPluginExecutor(
	targetDir,
	"2.1.9-b",
	"",
)

var eslintPlugin = executor.NewEsLintExecutor(
	targetDir,
	"/home/d050449/Documents/eslint-plugin-fiori-custom/configure.eslintrc",
	"/home/d050449/Documents/eslint-plugin-fiori-custom/lib/rules",
	"/home/d050449/Documents/eslint-plugin-fiori-custom/node_modules/eslint/bin/eslint.js",
)

var ruleNameToTestFile = map[string]string{
	"sap-no-debugger":                            "__test_files__/sap-no-debugger-sample.js",
	"sap-no-origin":                              "__test_files__/sap-no-origin-sample.js",
	"sap-not-localized":                          "__test_files__/sap-not-localized-sample.js",
	"sap-concatenated-strings":                   "__test_files__/sap-concatenated-strings.js",
	"sap-hardcoded-color":                        "__test_files__/sap-hardcoded-color.js",
	"sap-window-alert":                           "__test_files__/sap-window-alert.js",
	"sap-console-log":                            "__test_files__/sap-console-log.js",
	"sap-eval-used":                              "__test_files__/sap-eval-used.js",
	"sap-core-model-usage":                       "__test_files__/sap-core-model-usage.js",
	"sap-unescape-write":                         "__test_files__/sap-unescape-write.js",
	"sap-controller-hook-no-callback-signature":  "__test_files__/sap-controller-hook-no-callback-signature.js",
	"sap-controller-hook-bad-callback-signature": "__test_files__/sap-controller-hook-bad-callback-signature.js",
}

var eslintRuleNameToAnalysisName = map[string]string{
	"sap-no-debugger":                            "JS_DEBUGGER_STATEMENT",
	"sap-no-origin":                              "JS_ORIGIN_USED",
	"sap-not-localized":                          "JS_NOT_LOCALIZED",
	"sap-concatenated-strings":                   "JS_CONCATENATED_STRINGS",
	"sap-hardcoded-color":                        "JS_HARDCODED_COLOR",
	"sap-window-alert":                           "JS_WINDOW_ALERT",
	"sap-console-log":                            "JS_CONSOLE_LOG",
	"sap-eval-used":                              "JS_EVAL_USED",
	"sap-core-model-usage":                       "JS_CORE_MODEL_USAGE",
	"sap-unescape-write":                         "JS_UNESCAPED_WRITE",
	"sap-controller-hook-no-callback-signature":  "JS_CONTROLLER_HOOK_NO_CALLBACK_SIGNATURE",
	"sap-controller-hook-bad-callback-signature": "JS_CONTROLLER_HOOK_BAD_CALLBACK_SIGNATURE",
}

var timestamp = time.Now().Format("20060102-150405")

var RULES_TO_BE_TEST = []string{
	//"sap-no-debugger",
	//"sap-no-origin",
	//"sap-not-localized",
	//"sap-concatenated-strings",
	//"sap-hardcoded-color",
	//"sap-window-alert",
	//"sap-console-log",
	//"sap-eval-used",
	//"sap-core-model-usage",
	//"sap-unescape-write",
	//"sap-controller-hook-no-callback-signature",
	"sap-controller-hook-bad-callback-signature",
}

func TestRules(t *testing.T) {

	for _, ruleName_toTest := range RULES_TO_BE_TEST {
		assert.NotEqual(t, ruleNameToTestFile[ruleName_toTest], "", "Test file path should be defined for rule: %s", ruleName_toTest)
		assert.True(t, IsFileExists(ruleNameToTestFile[ruleName_toTest]), "Test file should exist for rule: %s", ruleName_toTest)
		assert.NotEqual(t, eslintRuleNameToAnalysisName[ruleName_toTest], "", "Analysis name should be defined for rule: %s", ruleName_toTest)
	}

	for _, ruleName := range RULES_TO_BE_TEST {

		testName := ruleName

		assert.NotEqual(t, ruleNameToTestFile[ruleName], "", "Test file path should be defined for rule: %s", ruleName)
		assert.NotEqual(t, eslintRuleNameToAnalysisName[ruleName], "", "Analysis name should be defined for rule: %s", ruleName)

		executionDir := filepath.Join("__executions__", fmt.Sprintf("%s-%s", testName, timestamp))
		err := os.MkdirAll(executionDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create execution directory: %v", err)
		}

		filePath := ruleNameToTestFile[ruleName]

		if filePath == "" {
			t.Fatalf("File path for rule '%s' not found", ruleName)
		}
		t.Run("Test Rule: "+ruleName, func(t *testing.T) {

			/*	t.Run("Single file", func(t *testing.T) {
				t.Run("AnalysisPluginA - single file", func(t *testing.T) {
					execError, codeCheckErrors := RunChecksAgainstFile(executionDir, testName, analysisPluginA, filePath)
					if execError != nil {
						t.Logf("Analysis Plugin A returned execution error: %v", execError)
						t.Fatalf("Test failed due to execution error: %v", execError)
					}

					assert.Greater(t, len(codeCheckErrors), 0, "Expected at least one code check error from Analysis Plugin A")

					targetError := FindAnalysisPluginIssue(codeCheckErrors, eslintRuleNameToAnalysisName[ruleName])
					assert.NotNilf(t, targetError, "Expected to find specific error for rule '%s'", ruleName)
				})

				t.Run("AnalysisPluginB - single file", func(t *testing.T) {
					execError, codeCheckErrors := RunChecksAgainstFile(executionDir, testName, analysisPluginB, filePath)
					if execError != nil {
						t.Logf("Analysis Plugin B returned execution error: %v", execError)
						t.Fatalf("Test failed due to execution error: %v", execError)
					}
					assert.Equal(t, 0, len(codeCheckErrors), "Expected no code check errors from Analysis Plugin B")
				})

				t.Run("ESLintPlugin - single file", func(t *testing.T) {
					execError, codeCheckErrors := RunChecksAgainstFile(executionDir, testName, eslintPlugin, filePath)
					if execError != nil {
						t.Logf("ESLint Plugin returned execution error: %v", execError)
						t.Fatalf("Test failed due to execution error: %v", execError)
					}
					assert.Greater(t, len(codeCheckErrors), 0, "Expected at least one code check error from ESLint Plugin")

					targetError := FindEslintIssue(codeCheckErrors, ruleName)
					assert.NotNilf(t, targetError, "Expected to find specific error for rule '%s'", ruleName)
				})
			})*/

			t.Run("Project code checks", func(t *testing.T) {
				// Copy test file
				destFilePath := filepath.Join(targetDir, "webapp/migration-test", fmt.Sprintf("testfile-%s.js", ruleName))
				err := CopyFile(filePath, destFilePath)
				if err != nil {
					t.Fatalf("Failed to copy test file: %v", err)
				}

				t.Run("AnalysisPluginA", func(t *testing.T) {
					// Use the exported RunChecks function from main.go
					execErr, codeErrors := RunChecks(executionDir, testName, analysisPluginA)

					if execErr != nil {
						// This could be either an execution error or code check error
						t.Fatalf("Test failed due to execution error: %v", execErr)
					}

					assert.Greater(t, len(codeErrors), 0, "Expected at least one code check error from Analysis Plugin A")
					// Check for code analysis errors
					targetError := FindAnalysisPluginIssue(codeErrors, eslintRuleNameToAnalysisName[ruleName])
					assert.NotNilf(t, targetError, "Expected to find specific error for rule '%s'", ruleName)
				})

				t.Run("AnalysisPluginB", func(t *testing.T) {
					// Use the exported RunChecks function from main.go
					execErr, codeErrors := RunChecks(executionDir, testName, analysisPluginB)

					if execErr != nil {
						// This could be either an execution error or code check error
						t.Fatalf("Test failed due to execution error: %v", execErr)
					}

					if len(codeErrors) > 0 {
						targetError := FindAnalysisPluginIssue(codeErrors, eslintRuleNameToAnalysisName[ruleName])
						assert.Nilf(t, targetError, "Did not expect to find specific error for rule '%s' in Analysis Plugin B", ruleName)
					}
				})

				t.Run("ESLintPlugin", func(t *testing.T) {
					// Use the exported RunChecks function from main.go
					execErr, codeErrors := RunChecks(executionDir, testName, eslintPlugin)
					if execErr != nil {
						// This could be either an execution error or code check error
						t.Fatalf("Test failed due to execution error: %v", execErr)
					}
					assert.Greater(t, len(codeErrors), 0, "Expected at least one code check error from ESLint Plugin")

					targetError := FindEslintIssue(codeErrors, ruleName)
					assert.NotNilf(t, targetError, "Expected to find specific error for rule '%s'", ruleName)
				})

				// Delete test file
				Delete(destFilePath)
			})

		})
	}
}

func TestOneRule(t *testing.T) {
	ruleName := "sap-no-origin"
	testName := ruleName

	assert.NotEqual(t, ruleNameToTestFile[ruleName], "", "Test file path should be defined for rule: %s", ruleName)
	assert.NotEqual(t, eslintRuleNameToAnalysisName[ruleName], "", "Analysis name should be defined for rule: %s", ruleName)

	executionDir := filepath.Join("__executions__", fmt.Sprintf("%s-%s", testName, timestamp))
	err := os.MkdirAll(executionDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create execution directory: %v", err)
	}

	filePath := ruleNameToTestFile[ruleName]

	if filePath == "" {
		t.Fatalf("File path for rule '%s' not found", ruleName)
	}

	/*	t.Run("Single file", func(t *testing.T) {
		t.Run("AnalysisPluginA - single file", func(t *testing.T) {
			execError, codeCheckErrors := RunChecksAgainstFile(executionDir, testName, analysisPluginA, filePath)
			if execError != nil {
				t.Logf("Analysis Plugin A returned execution error: %v", execError)
				t.Fatalf("Test failed due to execution error: %v", execError)
			}

			assert.Greater(t, len(codeCheckErrors), 0, "Expected at least one code check error from Analysis Plugin A")

			targetError := FindAnalysisPluginIssue(codeCheckErrors, eslintRuleNameToAnalysisName[ruleName])
			assert.NotNilf(t, targetError, "Expected to find specific error for rule '%s'", ruleName)
		})

		t.Run("AnalysisPluginB - single file", func(t *testing.T) {
			execError, codeCheckErrors := RunChecksAgainstFile(executionDir, testName, analysisPluginB, filePath)
			if execError != nil {
				t.Logf("Analysis Plugin B returned execution error: %v", execError)
				t.Fatalf("Test failed due to execution error: %v", execError)
			}
			assert.Equal(t, 0, len(codeCheckErrors), "Expected no code check errors from Analysis Plugin B")
		})

		t.Run("ESLintPlugin - single file", func(t *testing.T) {
			execError, codeCheckErrors := RunChecksAgainstFile(executionDir, testName, eslintPlugin, filePath)
			if execError != nil {
				t.Logf("ESLint Plugin returned execution error: %v", execError)
				t.Fatalf("Test failed due to execution error: %v", execError)
			}
			assert.Greater(t, len(codeCheckErrors), 0, "Expected at least one code check error from ESLint Plugin")

			targetError := FindEslintIssue(codeCheckErrors, ruleName)
			assert.NotNilf(t, targetError, "Expected to find specific error for rule '%s'", ruleName)
		})
	})*/

	t.Run("Project code checks", func(t *testing.T) {
		// Copy test file
		destFilePath := filepath.Join(targetDir, "webapp/migration-test", fmt.Sprintf("testfile-%s.js", ruleName))
		err := CopyFile(filePath, destFilePath)
		if err != nil {
			t.Fatalf("Failed to copy test file: %v", err)
		}

		t.Run("AnalysisPluginA", func(t *testing.T) {
			// Use the exported RunChecks function from main.go
			execErr, codeErrorsPluginA := RunChecks(executionDir, testName, analysisPluginA)

			if execErr != nil {
				// This could be either an execution error or code check error
				t.Fatalf("Test failed due to execution error: %v", execErr)
			}

			assert.Greater(t, len(codeErrorsPluginA), 0, "Expected at least one code check error from Analysis Plugin A")
			// Check for code analysis errors
			targetError := FindAnalysisPluginIssue(codeErrorsPluginA, eslintRuleNameToAnalysisName[ruleName])
			assert.NotNilf(t, targetError, "Expected to find specific error for rule '%s'", ruleName)
		})

		t.Run("AnalysisPluginB", func(t *testing.T) {
			// Use the exported RunChecks function from main.go
			execErr, codeErrors := RunChecks(executionDir, testName, analysisPluginB)

			if execErr != nil {
				// This could be either an execution error or code check error
				t.Fatalf("Test failed due to execution error: %v", execErr)
			}

			if len(codeErrors) > 0 {
				targetError := FindAnalysisPluginIssue(codeErrors, eslintRuleNameToAnalysisName[ruleName])
				assert.Nilf(t, targetError, "Did not expect to find specific error for rule '%s' in Analysis Plugin B", ruleName)
			}
		})

		t.Run("ESLintPlugin", func(t *testing.T) {
			// Use the exported RunChecks function from main.go
			execErr, codeErrors := RunChecks(executionDir, testName, eslintPlugin)
			if execErr != nil {
				// This could be either an execution error or code check error
				t.Fatalf("Test failed due to execution error: %v", execErr)
			}
			assert.Greater(t, len(codeErrors), 0, "Expected at least one code check error from ESLint Plugin")
			targetError := FindEslintIssue(codeErrors, ruleName)
			assert.NotNilf(t, targetError, "Expected to find specific error for rule '%s'", ruleName)
		})

		// Delete test file
		Delete(destFilePath)
	})

}
