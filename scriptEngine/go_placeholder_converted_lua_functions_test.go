package scriptEngine

import (
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestGoFenixControlledUniqueID_ShouldReplaceDateTimeAndRandomPatterns(t *testing.T) {
	// Freeze clock so date/time replacement assertions are deterministic.
	originalTimeProvider := currentTimeProvider
	currentTimeProvider = func() time.Time {
		return time.Date(2026, time.February, 24, 13, 7, 9, 0, time.Local)
	}
	defer func() {
		currentTimeProvider = originalTimeProvider
	}()

	input := GoPlaceholderInput{
		FunctionName: "Fenix_ControlledUniqueId",
		ArrayIndexes: []int{1},
		Arguments: []string{
			"%YYYY-MM-DD%|%YYYYMMDD%|%YYMMDD%|%hh:mm:ss%|%hh.mm.ss%|%hhmmss%|%hhmm%|%nnn%|%a(4; 11)%|%A(4; 11)%",
		},
		Entropy: 17,
	}

	resultOne, err := goFenixControlledUniqueID(input)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	resultTwo, err := goFenixControlledUniqueID(input)
	if err != nil {
		t.Fatalf("expected no error on second call, got: %v", err)
	}

	if resultOne != resultTwo {
		// Same seed inputs should always produce same output.
		t.Fatalf("expected deterministic output, got '%s' and '%s'", resultOne, resultTwo)
	}

	if strings.Contains(resultOne, "%") {
		t.Fatalf("expected all placeholders to be replaced, got: %s", resultOne)
	}

	if strings.Contains(resultOne, "2026-02-24") == false {
		t.Fatalf("expected formatted date in output, got: %s", resultOne)
	}

	if strings.Contains(resultOne, "20260224") == false || strings.Contains(resultOne, "260224") == false {
		t.Fatalf("expected compact dates in output, got: %s", resultOne)
	}

	if strings.Contains(resultOne, "13:07:09") == false ||
		strings.Contains(resultOne, "13.07.09") == false ||
		strings.Contains(resultOne, "130709") == false ||
		strings.Contains(resultOne, "1307") == false {
		t.Fatalf("expected formatted times in output, got: %s", resultOne)
	}

	if regexp.MustCompile(`\|[0-9]{3}\|`).MatchString(resultOne) == false {
		t.Fatalf("expected replaced random numeric token in output, got: %s", resultOne)
	}
}

func TestGoFenixRandomPositiveDecimalValue_ShouldBeDeterministicAndPadded(t *testing.T) {
	// 4-argument form enables integer/fraction zero-padding.
	input := GoPlaceholderInput{
		FunctionName: "Fenix_RandomPositiveDecimalValue",
		ArrayIndexes: []int{2},
		Arguments:    []string{"1", "2", "3", "4"},
		Entropy:      12345,
	}

	resultOne, err := goFenixRandomPositiveDecimalValue(input)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	resultTwo, err := goFenixRandomPositiveDecimalValue(input)
	if err != nil {
		t.Fatalf("expected no error on second call, got: %v", err)
	}

	if resultOne != resultTwo {
		t.Fatalf("expected deterministic output, got '%s' and '%s'", resultOne, resultTwo)
	}

	if regexp.MustCompile(`^[0-9]{3}\.[0-9]{4}$`).MatchString(resultOne) == false {
		t.Fatalf("expected zero-padded decimal value, got: %s", resultOne)
	}
}

func TestGoFenixRandomPositiveDecimalValue_ShouldFailOnInvalidArguments(t *testing.T) {
	// Lua contract requires either 2 or 4 arguments.
	input := GoPlaceholderInput{
		FunctionName: "Fenix_RandomPositiveDecimalValue",
		ArrayIndexes: []int{1},
		Arguments:    []string{"1", "2", "3"},
		Entropy:      1,
	}

	_, err := goFenixRandomPositiveDecimalValue(input)
	if err == nil {
		t.Fatalf("expected validation error for invalid number of arguments")
	}
}

func TestGoFenixRandomPositiveDecimalValueSum_ShouldBeDeterministic(t *testing.T) {
	// Sum variant supports positive/negative array indexes with deterministic seeding.
	input := GoPlaceholderInput{
		FunctionName: "Fenix_RandomPositiveDecimalValue_Sum",
		ArrayIndexes: []int{1, -2, 3},
		Arguments:    []string{"2", "2"},
		Entropy:      999,
	}

	resultOne, err := goFenixRandomPositiveDecimalValueSum(input)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	resultTwo, err := goFenixRandomPositiveDecimalValueSum(input)
	if err != nil {
		t.Fatalf("expected no error on second call, got: %v", err)
	}

	if resultOne != resultTwo {
		t.Fatalf("expected deterministic output, got '%s' and '%s'", resultOne, resultTwo)
	}

	if regexp.MustCompile(`^-?[0-9]+(\.[0-9]+)?$`).MatchString(resultOne) == false {
		t.Fatalf("expected numeric output, got: %s", resultOne)
	}
}
