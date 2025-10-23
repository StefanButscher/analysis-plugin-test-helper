package logsAnalyzer

import (
	"analysis-migration-test/errors"
	"bufio"
	"log"
	"os"
	"strings"
)

type EsLintLogAnalyzer struct{}

func (e *EsLintLogAnalyzer) GetAnalyzerType() string {
	return "ESLINT"
}

func (e *EsLintLogAnalyzer) AnalyzeLogFile(logFilePath string) ([]error, error) {
	log.Printf("[EsLintLogAnalyzer] Starting analysis of log file: %s", logFilePath)
	
	file, err := os.Open(logFilePath)
	if err != nil {
		log.Printf("[EsLintLogAnalyzer] Failed to open log file: %v", err)
		return nil, err
	}
	defer file.Close()

	var codeErrors []error
	scanner := bufio.NewScanner(file)
	totalLines := 0
	errorLines := 0
	warningLines := 0
	fileHeaders := 0

	for scanner.Scan() {
		totalLines++
		line := scanner.Text()

		// Check for execution errors first
		if strings.Contains(line, "Oops! Something went wrong!") {
			log.Printf("[EsLintLogAnalyzer] Execution error detected in log")
			return []error{errors.NewExecutionError(line)}, nil
		}
		
		// Check for file not found errors
		lowerLine := strings.ToLower(line)
		if strings.Contains(lowerLine, "enoent") || 
		   strings.Contains(lowerLine, "no such file") ||
		   strings.Contains(lowerLine, "cannot open") ||
		   strings.Contains(lowerLine, "not found") {
			log.Printf("[EsLintLogAnalyzer] File not found error detected in log")
			return []error{errors.NewExecutionError(line)}, nil
		}

		// Check for linting errors - looking for lines that contain error/warning info with line numbers
		if (strings.Contains(line, "error") || strings.Contains(line, "warning")) && strings.Contains(line, ":") {
			// Create CodeChecksError for each error line found
			codeErrors = append(codeErrors, errors.NewCodeChecksError(line))
			if strings.Contains(line, "error") {
				errorLines++
			} else {
				warningLines++
			}
		}
		// Also check for file header lines that contain .js (these indicate files with issues)
		if strings.Contains(line, ".js") && !strings.Contains(line, ":") && strings.HasPrefix(strings.TrimSpace(line), "/") {
			// This is a file header line, indicating a file with issues
			codeErrors = append(codeErrors, errors.NewCodeChecksError(line))
			fileHeaders++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("[EsLintLogAnalyzer] Error reading log file: %v", err)
		return codeErrors, err
	}

	log.Printf("[EsLintLogAnalyzer] Analysis completed: %d total lines processed, %d errors found (%d error lines, %d warning lines, %d file headers)", 
		totalLines, len(codeErrors), errorLines, warningLines, fileHeaders)
	return codeErrors, nil
}
