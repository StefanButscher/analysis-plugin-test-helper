package executor

type Executor interface {
	// Should return error and logFilePath
	// basePath - path to dir under __executions__ for current run
	// testName - current testName
	RunCheck(basePath string, testName string) (error, string)
	// Should return error and logFilePath
	// basePath - path to dir under __executions__ for current run
	// testName - current testName
	// filePath - path to the file to be checked
	RunCheckAgainstFile(basePath string, testName string, filePath string) (error, string)
	// Get Log Analyzer Type
	GetLogAnalyzerType() string
}
