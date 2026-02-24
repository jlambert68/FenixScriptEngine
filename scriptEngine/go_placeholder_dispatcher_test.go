package scriptEngine

import (
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
