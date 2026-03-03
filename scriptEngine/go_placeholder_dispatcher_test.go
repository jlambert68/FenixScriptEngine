package scriptEngine

import (
	"hash/crc32"
	"strings"
	"testing"
	"time"
)

func logDispatcherInputMatrix(t *testing.T, callLabel string, input []interface{}, testCaseExecutionUUID string) {
	t.Helper()
	t.Logf(
		"Dispatcher input [%s]\n  Input: %v\n  TestCaseExecutionUUID: %q",
		callLabel,
		input,
		testCaseExecutionUUID,
	)
}

func logDispatcherExecutionResult(t *testing.T, callLabel string, value string, handled bool, err error) {
	t.Helper()
	t.Logf(
		"Dispatcher execution [%s]\n  Value: %q\n  Handled: %t\n  Error: %v",
		callLabel,
		value,
		handled,
		err,
	)
}

func logDispatcherParseResult(t *testing.T, callLabel string, parsedInput GoPlaceholderInput, err error) {
	t.Helper()
	if err != nil {
		t.Logf("Dispatcher parse [%s]\n  ParsedInput: <invalid>\n  Error: %v", callLabel, err)
		return
	}

	logPlaceholderInputMatrix(t, callLabel+"-parsed", parsedInput)
	t.Logf("Dispatcher parse [%s]\n  Error: <nil>", callLabel)
}

func TestExecuteGoPlaceholderFunction_ShouldHandleRegisteredFunction(t *testing.T) {
	// Uses a registered function name and verifies Go dispatch happens before Lua fallback.
	input := []interface{}{
		"{{Fenix.TodayShiftDay(0)}}",
		"Fenix_TodayShiftDay",
		[]interface{}{},
		[]interface{}{"0"},
		true,
		uint64(0),
	}

	testCaseExecutionUUID := "123e4567-e89b-12d3-a456-426614174000"
	logDispatcherInputMatrix(t, "registered-function", input, testCaseExecutionUUID)
	value, handled, err := executeGoPlaceholderFunction(input, testCaseExecutionUUID)
	logDispatcherExecutionResult(t, "registered-function", value, handled, err)
	if handled == false {
		t.Fatalf("expected Go handler to process function")
	}
	if err != nil {
		t.Fatalf("did not expect error, got: %v", err)
	}
	if _, parseErr := time.Parse("2006-01-02", value); parseErr != nil {
		t.Fatalf("expected date in format YYYY-MM-DD, got: %s", value)
	}
}

func TestExecuteGoPlaceholderFunction_ShouldReturnErrorForInvalidArgument(t *testing.T) {
	// Invalid argument should still be "handled by Go", but return a validation error.
	input := []interface{}{
		"{{Fenix.TodayShiftDay(notAnInt)}}",
		"Fenix_TodayShiftDay",
		[]interface{}{},
		[]interface{}{"notAnInt"},
		true,
		uint64(0),
	}

	testCaseExecutionUUID := "execution-uuid"
	logDispatcherInputMatrix(t, "invalid-argument", input, testCaseExecutionUUID)
	value, handled, err := executeGoPlaceholderFunction(input, testCaseExecutionUUID)
	logDispatcherExecutionResult(t, "invalid-argument", value, handled, err)
	if handled == false {
		t.Fatalf("expected Go handler to process function")
	}
	if err == nil {
		t.Fatalf("expected error when argument is not an integer")
	}
}

func TestExecuteGoPlaceholderFunction_ShouldIgnoreUnknownFunction(t *testing.T) {
	// Unknown function names must not fail at dispatcher level; Lua path should decide.
	input := []interface{}{
		"{{Fenix.Unknown()}}",
		"Fenix_Unknown",
		[]interface{}{},
		[]interface{}{""},
		true,
		uint64(0),
	}

	testCaseExecutionUUID := "execution-uuid"
	logDispatcherInputMatrix(t, "unknown-function", input, testCaseExecutionUUID)
	value, handled, err := executeGoPlaceholderFunction(input, testCaseExecutionUUID)
	logDispatcherExecutionResult(t, "unknown-function", value, handled, err)
	if handled == true {
		t.Fatalf("expected unknown function to be handled by Lua fallback")
	}
	if err != nil {
		t.Fatalf("did not expect error for unknown function, got: %v", err)
	}
}

func TestParseGoPlaceholderInput_ShouldValidateEntropyTypes(t *testing.T) {
	baseInput := []interface{}{
		"{{Fenix.ControlledUniqueId(X)}}",
		"Fenix_ControlledUniqueId",
		[]interface{}{},
		[]interface{}{"X"},
		true,
		uint64(0),
	}

	invalidUseEntropy := append([]interface{}{}, baseInput...)
	invalidUseEntropy[4] = "false"
	testCaseExecutionUUID := "execution-uuid"
	logDispatcherInputMatrix(t, "parse-invalid-use-entropy", invalidUseEntropy, testCaseExecutionUUID)
	parsedInput, err := parseGoPlaceholderInput(invalidUseEntropy, testCaseExecutionUUID)
	logDispatcherParseResult(t, "parse-invalid-use-entropy", parsedInput, err)
	if err == nil {
		t.Fatalf("expected type validation error for input parameter 4")
	}
	if strings.Contains(err.Error(), "input parameter 4 ('useEntropyFromExecutionUuid') must be bool") == false {
		t.Fatalf("unexpected error for input parameter 4: %v", err)
	}

	invalidExtraEntropy := append([]interface{}{}, baseInput...)
	invalidExtraEntropy[5] = 0
	logDispatcherInputMatrix(t, "parse-invalid-extra-entropy", invalidExtraEntropy, testCaseExecutionUUID)
	parsedInput, err = parseGoPlaceholderInput(invalidExtraEntropy, testCaseExecutionUUID)
	logDispatcherParseResult(t, "parse-invalid-extra-entropy", parsedInput, err)
	if err == nil {
		t.Fatalf("expected type validation error for input parameter 5")
	}
	if strings.Contains(err.Error(), "input parameter 5 ('extraEntropy') must be uint64") == false {
		t.Fatalf("unexpected error for input parameter 5: %v", err)
	}
}

func TestParseGoPlaceholderInput_ShouldCalculateEntropyFromTail(t *testing.T) {
	testCaseExecutionUUID := "f8c06f7e-0a8a-4d75-9f25-5e5fb8d2a6d3"

	inputWithoutUUIDEntropy := []interface{}{
		"{{Fenix.ControlledUniqueId(X)}(false, 7)}",
		"Fenix_ControlledUniqueId",
		[]interface{}{},
		[]interface{}{"X"},
		false,
		uint64(7),
	}

	logDispatcherInputMatrix(t, "parse-without-uuid-entropy", inputWithoutUUIDEntropy, testCaseExecutionUUID)
	parsedInput, err := parseGoPlaceholderInput(inputWithoutUUIDEntropy, testCaseExecutionUUID)
	logDispatcherParseResult(t, "parse-without-uuid-entropy", parsedInput, err)
	if err != nil {
		t.Fatalf("expected no parse error, got: %v", err)
	}
	if parsedInput.UseEntropyFromExecutionUUID != false {
		t.Fatalf("expected useEntropyFromExecutionUUID=false")
	}
	if parsedInput.Entropy != 7 {
		t.Fatalf("expected entropy=7 when UUID entropy is disabled, got %d", parsedInput.Entropy)
	}

	inputWithUUIDEntropy := []interface{}{
		"{{Fenix.ControlledUniqueId(X)}(true, 7)}",
		"Fenix_ControlledUniqueId",
		[]interface{}{},
		[]interface{}{"X"},
		true,
		uint64(7),
	}

	logDispatcherInputMatrix(t, "parse-with-uuid-entropy", inputWithUUIDEntropy, testCaseExecutionUUID)
	parsedInput, err = parseGoPlaceholderInput(inputWithUUIDEntropy, testCaseExecutionUUID)
	logDispatcherParseResult(t, "parse-with-uuid-entropy", parsedInput, err)
	if err != nil {
		t.Fatalf("expected no parse error, got: %v", err)
	}
	expectedEntropy := uint64(crc32.ChecksumIEEE([]byte(testCaseExecutionUUID))) + uint64(7)
	if parsedInput.Entropy != expectedEntropy {
		t.Fatalf("expected entropy=%d, got %d", expectedEntropy, parsedInput.Entropy)
	}
}
