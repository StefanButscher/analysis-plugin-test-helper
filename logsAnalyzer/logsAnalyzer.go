package logsAnalyzer

// LogsAnalyzer interface for analyzing log files from different executors
type LogsAnalyzer interface {
	// AnalyzeLogFile takes a log file path and returns a slice of errors found
	AnalyzeLogFile(logFilePath string) ([]error, error)
	// GetAnalyzerType returns the type of analyzer
	GetAnalyzerType() string
}

// GetAnalyzer factory function to get the appropriate analyzer based on type
func GetAnalyzer(analyzerType string) LogsAnalyzer {
	switch analyzerType {
	case "ANALYSIS_PLUGIN":
		return &AnalysisPluginLogAnalyzer{}
	case "ESLINT":
		return &EsLintLogAnalyzer{}
	default:
		return nil
	}
}