package logsAnalyzer

import (
	"analysis-migration-test/errors"
	"bufio"
	"log"
	"os"
	"strings"
)

type AnalysisPluginLogAnalyzer struct{}

func (a *AnalysisPluginLogAnalyzer) GetAnalyzerType() string {
	return "ANALYSIS_PLUGIN"
}

func (a *AnalysisPluginLogAnalyzer) AnalyzeLogFile(logFilePath string) ([]error, error) {
	log.Printf("[AnalysisPluginLogAnalyzer] Starting analysis of log file: %s", logFilePath)

	file, err := os.Open(logFilePath)
	if err != nil {
		log.Printf("[AnalysisPluginLogAnalyzer] Failed to open log file: %v", err)
		return nil, err
	}
	defer file.Close()

	var codeErrors []error
	scanner := bufio.NewScanner(file)
	totalLines := 0

	for scanner.Scan() {
		totalLines++
		line := scanner.Text()
		if strings.Contains(line, "[ERROR]") || strings.Contains(line, "Quality issue") {
			// Create CodeChecksError for each error line found
			codeErrors = append(codeErrors, errors.NewCodeChecksError(line))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("[AnalysisPluginLogAnalyzer] Error reading log file: %v", err)
		return codeErrors, err
	}

	log.Printf("[AnalysisPluginLogAnalyzer] Analysis completed: %d total lines processed, %d errors found", totalLines, len(codeErrors))
	return codeErrors, nil
}
