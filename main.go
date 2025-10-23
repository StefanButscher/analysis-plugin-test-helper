package main

import (
	"analysis-migration-test/executor"
	"analysis-migration-test/logsAnalyzer"
)

// RunChecks executes the provided executor and analyzes the resulting log file
// First error is error if smth went wrong, the last arg is linting errors
func RunChecks(basePath string, testName string, executors executor.Executor) (error, []error) {
	err, logPath := executors.RunCheck(basePath, testName)

	if err != nil {
		// If the executor returned an error, it could be either execution error or code check error
		// We still want to analyze the log file if it exists
		analyzerType := executors.GetLogAnalyzerType()
		analyzer := logsAnalyzer.GetAnalyzer(analyzerType)

		if analyzer != nil {
			codeErrors, analyzeErr := analyzer.AnalyzeLogFile(logPath)
			if analyzeErr == nil && len(codeErrors) > 0 {
				// Return both the execution error and code errors
				return err, codeErrors
			}
		}

		return err, nil
	}

	analyzerType := executors.GetLogAnalyzerType()
	analyzer := logsAnalyzer.GetAnalyzer(analyzerType)

	if analyzer == nil {
		return nil, nil
	}

	codeErrors, analyzeErr := analyzer.AnalyzeLogFile(logPath)
	if analyzeErr != nil {
		return analyzeErr, nil
	}

	return nil, codeErrors
}

// RunChecksAgainstFile executes the provided executor against a specific file and analyzes the resulting log file
// First error is error if smth went wrong, the last arg is linting errors
func RunChecksAgainstFile(basePath, testName string, executors executor.Executor, filePath string) (error, []error) {
	err, logPath := executors.RunCheckAgainstFile(basePath, testName, filePath)

	if err != nil {
		// If the executor returned an error, it could be either execution error or code check error
		// We still want to analyze the log file if it exists
		analyzerType := executors.GetLogAnalyzerType()
		analyzer := logsAnalyzer.GetAnalyzer(analyzerType)

		if analyzer != nil {
			codeErrors, analyzeErr := analyzer.AnalyzeLogFile(logPath)
			if analyzeErr == nil && len(codeErrors) > 0 {
				// Return both the execution error and code errors
				return err, codeErrors
			}
		}

		return err, nil
	}

	analyzerType := executors.GetLogAnalyzerType()
	analyzer := logsAnalyzer.GetAnalyzer(analyzerType)

	if analyzer == nil {
		return nil, nil
	}

	codeErrors, analyzeErr := analyzer.AnalyzeLogFile(logPath)
	if analyzeErr != nil {
		return analyzeErr, nil
	}

	return nil, codeErrors
}
