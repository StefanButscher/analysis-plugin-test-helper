package errors

type CodeChecksError struct {
	RowText  string
	RuleName string
}

func (c CodeChecksError) Error() string {
	return c.RowText
}

func NewCodeChecksError(rowText string) error {

	// Implementation for creating a new CodeChecksError
	return CodeChecksError{RowText: rowText}
}

type ExecutionError struct {
	rowText string
}

func (e ExecutionError) Error() string {
	return e.rowText
}

func NewExecutionError(rowText string) error {

	// Implementation for creating a new ExecutionError
	return ExecutionError{rowText: rowText}
}
