package main

import (
	"analysis-migration-test/executor"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	"sap-no-debugger":                                "__test_files__/sap-no-debugger-sample.js",
	"sap-no-origin":                                  "__test_files__/sap-no-origin-sample.js",
	"sap-not-localized":                              "__test_files__/sap-not-localized-sample.js",
	"sap-concatenated-strings":                       "__test_files__/sap-concatenated-strings.js",
	"sap-no-hardcoded-color":                         "__test_files__/sap-hardcoded-color.js",
	"sap-no-window-alert":                            "__test_files__/sap-no-window-alert.js",
	"sap-no-console-log":                             "__test_files__/sap-no-console-log.js",
	"sap-no-eval":                                    "__test_files__/sap-no-eval.js",
	"sap-no-ui5base-prop":                            "__test_files__/sap-core-model-usage.js",
	"sap-unescaped-write":                            "__test_files__/sap-unescaped-write.js",
	"sap-controller-hook-missing-callback-signature": "__test_files__/sap-controller-hook-missing-callback-signature.js",
	"sap-controller-hook-bad-callback-signature":     "__test_files__/sap-controller-hook-bad-callback-signature.js",
	"sap-no-localstorage":                            "__test_files__/sap-no-localstorage.js",
	"sap-no-upload":                                  "__test_files__/sap-no-upload.js",
	"sap-controller-hook-name-convention":            "__test_files__/sap-controller-hook-name-convention.js",
}

var eslintRuleNameToAnalysisName = map[string]string{
	"sap-no-debugger":                                "JS_DEBUGGER_STATEMENT",
	"sap-no-origin":                                  "JS_ORIGIN_USED",
	"sap-not-localized":                              "JS_NOT_LOCALIZED",        // Analysis plugin relied on a different sample - to discuss
	"sap-concatenated-strings":                       "JS_CONCATENATED_STRINGS", // Eslint relied on different tests
	"sap-no-hardcoded-color":                         "JS_HARDCODED_COLOR",      // Wrong eslint rule name
	"sap-no-window-alert":                            "JS_WINDOW_ALERT",         // changed eslint rule name
	"sap-no-console-log":                             "JS_CONSOLE_LOG",
	"sap-no-eval":                                    "JS_EVAL_USED",
	"sap-no-ui5base-prop":                            "JS_CORE_MODEL_USAGE",
	"sap-unescaped-write":                            "JS_UNESCAPED_WRITE",
	"sap-controller-hook-missing-callback-signature": "JS_CONTROLLER_HOOK_NO_CALLBACK_SIGNATURE",
	"sap-controller-hook-bad-callback-signature":     "JS_CONTROLLER_HOOK_BAD_CALLBACK_SIGNATURE",
	"sap-no-localstorage":                            "JS_WEBSTORAGE", // Wrong eslint rule name
	"sap-no-upload":                                  "JS_UPLOAD_DECLARE",
	"sap-controller-hook-name-convention":            "JS_CONTROLLER_HOOK_NAME_CONVENTION",
}

var timestamp = time.Now().Format("20060102-150405")

var RULES_TO_BE_TEST = []string{
	"sap-no-debugger",                                // +
	"sap-no-origin",                                  // +
	"sap-not-localized",                              // +
	"sap-concatenated-strings",                       // +
	"sap-no-hardcoded-color",                         //+ elsint works when adding a color code like #cccccc, Analysis plugin doe not find the issue
	"sap-no-window-alert",                            //- ESLINT ok with no-alert, Analysis Plugin A issue ->only no-alert works, ticket needed
	"sap-no-console-log",                             // +
	"sap-no-eval",                                    //-  only    no-eval works with eslint -> same issue ticket
	"sap-no-ui5base-prop",                            //+
	"sap-unescaped-write",                            // +
	"sap-controller-hook-missing-callback-signature", // eslint works when removing annotation @callback sap.ca.scfld.md.controller.BaseDetailController~onDataReceived //Analiysi plugin fails
	"sap-controller-hook-bad-callback-signature",     // eslint passes, Analysis plugin fails
	"sap-no-localstorage",                            // +
	"sap-no-upload",                                  //- missing in configure.eslintrc //- did not pass with original file // passes with file from eslint plugin
	"sap-controller-hook-name-convention",            // + Eslint  ok, analysis plugin a fail
}

var dateTime = time.Now().Format("02 Jan 06 15:04 MST")
var dateDir = strings.ReplaceAll(dateTime, " ", "_")

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

		executionDir := filepath.Join("__executions__", dateDir, fmt.Sprintf("%s-%s", testName, timestamp))
		err := os.MkdirAll(executionDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create execution directory: %v", err)
		}

		filePath := ruleNameToTestFile[ruleName]

		if filePath == "" {
			t.Fatalf("File path for rule '%s' not found", ruleName)
		}
		t.Run("Test Rule: "+ruleName, func(t *testing.T) {

			/*t.Run("Single file", func(t *testing.T) {
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
					fmt.Println(">>>>> Rulename", ruleName, eslintRuleNameToAnalysisName[ruleName])
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

	executionDir := filepath.Join("__executions__", dateDir, fmt.Sprintf("%s-%s", testName, timestamp))
	err := os.MkdirAll(executionDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create execution directory: %v", err)
	}

	filePath := ruleNameToTestFile[ruleName]

	if filePath == "" {
		t.Fatalf("File path for rule '%s' not found", ruleName)
	}

	t.Run("Single file", func(t *testing.T) {
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
	})

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
