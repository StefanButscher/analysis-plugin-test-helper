package executor

import (
	"analysis-migration-test/errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type AnalysisPluginExecutor struct {
	targetDir          string
	version            string
	standaloneFilePath string
}

// CLAUDE comments:
// targetDir - path to the target directory containing the files to be analyzed
// You should run all mvn commands against this directory
// stnadaloneFilePath - path to the standalone jar file - you should use this to RunCheckAgainstFile
// please ensure that this file actually exists
func NewAnalysisPluginExecutor(targetDir string, version string, standaloneFilePath string) Executor {
	// Validate that standalone file exists
	if _, err := os.Stat(standaloneFilePath); os.IsNotExist(err) {
		panic(fmt.Sprintf("Standalone JAR file does not exist: %s", standaloneFilePath))
	}

	// Validate target directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		panic(fmt.Sprintf("Target directory does not exist: %s", targetDir))
	}

	return &AnalysisPluginExecutor{
		targetDir:          targetDir,
		version:            version,
		standaloneFilePath: standaloneFilePath,
	}
}

// CLAUDE comments:
// To run check u can use this command, where -Dsap.ca.analysis-plugin.version
// is version of the plugin.
// During instantiating process please check that this version actually exist, using mvn command
// mvn clean verify -Panalysis.plugin -Dsap.ca.analysis-plugin.version=2.1.8-A
func (e *AnalysisPluginExecutor) RunCheck(basePath string, testName string) (error, string) {
	logFileName := fmt.Sprintf("%s-analysisPlugin-%s.log", testName, e.version)
	logFilePath := filepath.Join(basePath, logFileName)

	log.Printf("[AnalysisPluginExecutor] Starting analysis for test '%s' with version %s", testName, e.version)
	log.Printf("[AnalysisPluginExecutor] Target directory: %s", e.targetDir)
	log.Printf("[AnalysisPluginExecutor] Log file: %s", logFilePath)

	// Build the Maven command
	mvnCmd := fmt.Sprintf("mvn clean verify -Panalysis.plugin -Dsap.ca.analysis-plugin.version=%s", e.version)

	// Create the full command with tee using absolute path
	absLogFilePath, absErr := filepath.Abs(logFilePath)
	if absErr != nil {
		log.Printf("[AnalysisPluginExecutor] Failed to get absolute path for log file: %v", absErr)
		absLogFilePath = logFilePath
	}
	fullCmd := fmt.Sprintf("cd %s && %s | tee %s", e.targetDir, mvnCmd, absLogFilePath)

	log.Printf("[AnalysisPluginExecutor] Executing: %s", mvnCmd)

	// Execute the command
	cmd := exec.Command("sh", "-c", fullCmd)
	err := cmd.Run()

	if err != nil {
		log.Printf("[AnalysisPluginExecutor] Maven command failed: %v", err)
		return errors.NewExecutionError(fmt.Sprintf("Maven command failed: %v", err)), logFilePath
	}

	log.Printf("[AnalysisPluginExecutor] Analysis completed successfully")
	return nil, logFilePath
}

// CLAUDE comments:
// java -jar fiori-js-analysis-standalone.jar ./webapp/Component.js
func (e *AnalysisPluginExecutor) RunCheckAgainstFile(basePath string, testName string, filePath string) (error, string) {
	logFileName := fmt.Sprintf("%s-analysisPlugin-%s-file.log", testName, e.version)
	logFilePath := filepath.Join(basePath, logFileName)

	log.Printf("[AnalysisPluginExecutor] Starting file analysis for test '%s' with version %s", testName, e.version)
	log.Printf("[AnalysisPluginExecutor] Target file: %s", filePath)
	log.Printf("[AnalysisPluginExecutor] Standalone JAR: %s", e.standaloneFilePath)
	log.Printf("[AnalysisPluginExecutor] Log file: %s", logFilePath)

	// Build the Java command for standalone analysis
	javaCmd := fmt.Sprintf("java -jar %s %s", e.standaloneFilePath, filePath)

	// Create the full command with tee using absolute path
	absLogFilePath, absErr := filepath.Abs(logFilePath)
	if absErr != nil {
		log.Printf("[AnalysisPluginExecutor] Failed to get absolute path for log file: %v", absErr)
		absLogFilePath = logFilePath
	}
	fullCmd := fmt.Sprintf("%s | tee %s", javaCmd, absLogFilePath)

	log.Printf("[AnalysisPluginExecutor] Executing: %s", javaCmd)

	// Execute the command
	cmd := exec.Command("sh", "-c", fullCmd)
	err := cmd.Run()

	if err != nil {
		log.Printf("[AnalysisPluginExecutor] Java command failed: %v", err)
		return errors.NewExecutionError(fmt.Sprintf("Java command failed: %v", err)), logFilePath
	}

	log.Printf("[AnalysisPluginExecutor] File analysis completed successfully")
	return nil, logFilePath
}

func (e *AnalysisPluginExecutor) GetLogAnalyzerType() string {
	return "ANALYSIS_PLUGIN"
}
