package scriptEngine

import (
	"hash/crc32"
	"strings"
	"testing"
	"time"
)

func TestExecuteGoPlaceholderFunction_ShouldHandleRegisteredFunction(t *testing.T) {
	// Uses a registered function name and verifies Go dispatch happens before Lua fallback.
	input := []interface{}{
		"{{Fenix.TodayShiftDay()}}",
		"Fenix_TodayShiftDay",
		[]interface{}{},
		[]interface{}{""},
		true,
		uint64(0),
	}

	value, handled, err := executeGoPlaceholderFunction(input, "123e4567-e89b-12d3-a456-426614174000")
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

	_, handled, err := executeGoPlaceholderFunction(input, "execution-uuid")
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

	_, handled, err := executeGoPlaceholderFunction(input, "execution-uuid")
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
	_, err := parseGoPlaceholderInput(invalidUseEntropy, "execution-uuid")
	if err == nil {
		t.Fatalf("expected type validation error for input parameter 4")
	}
	if strings.Contains(err.Error(), "input parameter 4 ('useEntropyFromExecutionUuid') must be bool") == false {
		t.Fatalf("unexpected error for input parameter 4: %v", err)
	}

	invalidExtraEntropy := append([]interface{}{}, baseInput...)
	invalidExtraEntropy[5] = 0
	_, err = parseGoPlaceholderInput(invalidExtraEntropy, "execution-uuid")
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

	parsedInput, err := parseGoPlaceholderInput(inputWithoutUUIDEntropy, testCaseExecutionUUID)
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

	parsedInput, err = parseGoPlaceholderInput(inputWithUUIDEntropy, testCaseExecutionUUID)
	if err != nil {
		t.Fatalf("expected no parse error, got: %v", err)
	}
	expectedEntropy := uint64(crc32.ChecksumIEEE([]byte(testCaseExecutionUUID))) + uint64(7)
	if parsedInput.Entropy != expectedEntropy {
		t.Fatalf("expected entropy=%d, got %d", expectedEntropy, parsedInput.Entropy)
	}
}
