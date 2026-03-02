package scriptEngine

import "testing"

func logPlaceholderInputMatrix(t *testing.T, callLabel string, input GoPlaceholderInput) {
	t.Helper()

	t.Logf(
		"Input matrix [%s]\n  Placeholder: %q\n  FunctionName: %q\n  ArrayIndexes: %v\n  Arguments: %v\n  UseEntropyFromExecutionUUID: %t\n  ExtraEntropy: %d\n  Entropy: %d\n  TestCaseExecutionUUID: %q",
		callLabel,
		input.Placeholder,
		input.FunctionName,
		input.ArrayIndexes,
		input.Arguments,
		input.UseEntropyFromExecutionUUID,
		input.ExtraEntropy,
		input.Entropy,
		input.TestCaseExecutionUUID,
	)
}

func logPlaceholderExecutionResult(t *testing.T, callLabel string, result string, err error) {
	t.Helper()
	if err != nil {
		t.Logf("Execution result [%s]\n  Result: %q\n  Error: %v", callLabel, result, err)
		return
	}

	t.Logf("Execution result [%s]\n  Result: %q\n  Error: <nil>", callLabel, result)
}
