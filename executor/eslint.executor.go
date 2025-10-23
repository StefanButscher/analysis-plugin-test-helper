package executor

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type EsLintExecutor struct {
	targetDir      string
	configFilePath string
	rulesDir       string
	binaryPath     string
}

// CLAUDE comments:
// This is example from VSCode lanch,json configuration, basically this is just eslint execution
// "name": "Debug eslint plugin",
// "request": "launch",
// "skipFiles": ["<node_internals>/**"],
// "type": "node",
// "args": [
//   "${workspaceFolder}/testArea/ca.infra.testapp/node_modules/eslint/bin/eslint.js",
//   "${workspaceFolder}/testArea/ca.infra.testapp/webapp/Component.js",
//   "-c",
//   "/Users/C5310597/Documents/GitHub/eslint-plugin-fiori-custom/configure.eslintrc",
//   "--no-eslintrc",
//   "--rulesdir",
//   "/Users/C5310597/Documents/GitHub/eslint-plugin-fiori-custom/lib/rules",
//   "--ignore-pattern",
//   "test/**",
//   "--ignore-pattern",
//   "src/test/**",
//   "--ignore-pattern",
//   "target/**",
//   "--ignore-pattern",
//   "webapp/test/**",
//   "--ignore-pattern",
//   "src/main/webapp/test/**",
//   "--ignore-pattern",
//   "webapp/localservice/**",
//   "--ignore-pattern",
//   "/src/main/webapp/localService/**",
//   "--ignore-pattern",
//   "backup/**",
//   "--ignore-pattern",
//   "Gruntfile.js",
//   "--ignore-pattern",
//   "changes_preview.js",
//   "--ignore-pattern",
//   "gulpfile.js"
// ]
// },
// eslint should be already install in the target project
// but check that it actually installed and log it version using npm commands

func NewEsLintExecutor(targetDir string, configFilePath string, rulesDir string, binaryPath string) Executor {
	// Validate that target directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		panic(fmt.Sprintf("Target directory does not exist: %s", targetDir))
	}

	// Validate that config file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		panic(fmt.Sprintf("Config file does not exist: %s", configFilePath))
	}

	// Validate that rules directory exists
	if _, err := os.Stat(rulesDir); os.IsNotExist(err) {
		panic(fmt.Sprintf("Rules directory does not exist: %s", rulesDir))
	}

	// Validate that binary path exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("ESLint binary does not exist: %s", binaryPath))
	}

	return &EsLintExecutor{
		targetDir:      targetDir,
		configFilePath: configFilePath,
		rulesDir:       rulesDir,
		binaryPath:     binaryPath,
	}
}

func (e *EsLintExecutor) RunCheck(basePath string, testName string) (error, string) {
	logFileName := fmt.Sprintf("%s-eslint.log", testName)
	logFilePath := filepath.Join(basePath, logFileName)

	log.Printf("[EsLintExecutor] Starting ESLint analysis for test '%s'", testName)
	log.Printf("[EsLintExecutor] Target directory: %s", e.targetDir)
	log.Printf("[EsLintExecutor] Config file: %s", e.configFilePath)
	log.Printf("[EsLintExecutor] Rules directory: %s", e.rulesDir)
	log.Printf("[EsLintExecutor] Binary path: %s", e.binaryPath)
	log.Printf("[EsLintExecutor] Log file: %s", logFilePath)

	// Build the ESLint command for entire project using node and binary path
	eslintCmd := fmt.Sprintf("cd %s && node %s webapp/ -c %s --no-eslintrc --rulesdir %s --ignore-pattern \"test/**\" --ignore-pattern \"src/test/**\" --ignore-pattern \"target/**\" --ignore-pattern \"webapp/test/**\" --ignore-pattern \"src/main/webapp/test/**\" --ignore-pattern \"webapp/localservice/**\" --ignore-pattern \"/src/main/webapp/localService/**\" --ignore-pattern \"backup/**\" --ignore-pattern \"Gruntfile.js\" --ignore-pattern \"changes_preview.js\" --ignore-pattern \"gulpfile.js\"",
		e.targetDir, e.binaryPath, e.configFilePath, e.rulesDir)

	// Create the full command with tee using absolute path
	absLogFilePath, absErr := filepath.Abs(logFilePath)
	if absErr != nil {
		log.Printf("[EsLintExecutor] Failed to get absolute path for log file: %v", absErr)
		absLogFilePath = logFilePath
	}
	fullCmd := fmt.Sprintf("%s 2>&1 | tee %s", eslintCmd, absLogFilePath)

	log.Printf("[EsLintExecutor] Executing: %s [with ignore patterns]", fullCmd)

	// Execute the command
	cmd := exec.Command("sh", "-c", fullCmd)
	err := cmd.Run()

	if err != nil {
		log.Printf("[EsLintExecutor] ESLint command failed: %v", err)
		// Don't return error immediately - let log analyzer determine if it's a real error or just code issues
		// ESLint exits with code 1 for linting errors and code 2 for system errors
		log.Printf("[EsLintExecutor] ESLint completed with non-zero exit code (may be linting issues)")
		return nil, logFilePath
	}

	log.Printf("[EsLintExecutor] ESLint analysis completed successfully")
	return nil, logFilePath
}

func (e *EsLintExecutor) RunCheckAgainstFile(basePath string, testName string, filePath string) (error, string) {
	logFileName := fmt.Sprintf("%s-eslint-file.log", testName)
	logFilePath := filepath.Join(basePath, logFileName)

	log.Printf("[EsLintExecutor] Starting ESLint file analysis for test '%s'", testName)
	log.Printf("[EsLintExecutor] Target file: %s", filePath)
	log.Printf("[EsLintExecutor] Config file: %s", e.configFilePath)
	log.Printf("[EsLintExecutor] Rules directory: %s", e.rulesDir)
	log.Printf("[EsLintExecutor] Binary path: %s", e.binaryPath)
	log.Printf("[EsLintExecutor] Log file: %s", logFilePath)

	// Build the ESLint command for specific file using node and binary path
	eslintCmd := fmt.Sprintf("node %s %s -c %s --no-eslintrc --rulesdir %s", e.binaryPath, filePath, e.configFilePath, e.rulesDir)

	// Create the full command with tee using absolute path
	absLogFilePath, absErr := filepath.Abs(logFilePath)
	if absErr != nil {
		log.Printf("[EsLintExecutor] Failed to get absolute path for log file: %v", absErr)
		absLogFilePath = logFilePath
	}
	fullCmd := fmt.Sprintf("%s 2>&1 | tee %s", eslintCmd, absLogFilePath)

	log.Printf("[EsLintExecutor] Executing: %s [with config and rules]", fullCmd)

	// Execute the command
	cmd := exec.Command("sh", "-c", fullCmd)
	err := cmd.Run()

	if err != nil {
		log.Printf("[EsLintExecutor] ESLint command failed: %v", err)
		// Don't return error immediately - let log analyzer determine if it's a real error or just code issues
		// ESLint exits with code 1 for linting errors and code 2 for system errors
		log.Printf("[EsLintExecutor] ESLint file analysis completed with non-zero exit code (may be linting issues)")
		return nil, logFilePath
	}

	log.Printf("[EsLintExecutor] ESLint file analysis completed successfully")
	return nil, logFilePath
}

func (e *EsLintExecutor) GetLogAnalyzerType() string {
	return "ESLINT"
}
