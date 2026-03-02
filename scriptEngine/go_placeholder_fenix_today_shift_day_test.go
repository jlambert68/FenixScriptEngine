package scriptEngine

import (
	"strings"
	"testing"
	"time"
)

func TestGoFenixTodayShiftDay_JiraBasicVariants(t *testing.T) {
	originalTimeProvider := currentTimeProvider
	currentTimeProvider = func() time.Time {
		return time.Date(2026, time.February, 26, 12, 30, 45, 0, time.Local)
	}
	defer func() {
		currentTimeProvider = originalTimeProvider
	}()

	testCases := []struct {
		name     string
		args     []string
		expected string
	}{
		{name: "shift 0", args: []string{"0"}, expected: "2026-02-26"},
		{name: "shift -1", args: []string{"-1"}, expected: "2026-02-25"},
		{name: "shift +1", args: []string{"1"}, expected: "2026-02-27"},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			input := GoPlaceholderInput{
				Placeholder:  "{{Fenix.TodayShiftDay(" + testCase.args[0] + ")}}",
				FunctionName: "Fenix_TodayShiftDay",
				Arguments:    testCase.args,
			}
			logPlaceholderInputMatrix(t, testCase.name, input)

			result, err := goFenixTodayShiftDay(input)
			logPlaceholderExecutionResult(t, testCase.name, result, err)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if result != testCase.expected {
				t.Fatalf("expected %s, got %s", testCase.expected, result)
			}
		})
	}
}

func TestGoFenixTodayShiftDay_ShouldValidateInput(t *testing.T) {
	input := GoPlaceholderInput{
		Placeholder:  "{{Fenix.TodayShiftDay[1](0)}}",
		FunctionName: "Fenix_TodayShiftDay",
		ArrayIndexes: []int{1},
		Arguments:    []string{"0"},
	}
	logPlaceholderInputMatrix(t, "array-index-not-supported", input)
	_, err := goFenixTodayShiftDay(input)
	logPlaceholderExecutionResult(t, "array-index-not-supported", "", err)
	if err == nil {
		t.Fatalf("expected error when array index is provided")
	}
	if strings.Contains(err.Error(), "array index is not supported") == false {
		t.Fatalf("unexpected error: %v", err)
	}

	input = GoPlaceholderInput{
		Placeholder:  "{{Fenix.TodayShiftDay(abc)}}",
		FunctionName: "Fenix_TodayShiftDay",
		Arguments:    []string{"abc"},
	}
	logPlaceholderInputMatrix(t, "non-integer-argument", input)
	_, err = goFenixTodayShiftDay(input)
	logPlaceholderExecutionResult(t, "non-integer-argument", "", err)
	if err == nil {
		t.Fatalf("expected error for non integer argument")
	}
	if strings.Contains(err.Error(), "not an Integer") == false {
		t.Fatalf("unexpected error: %v", err)
	}

	input = GoPlaceholderInput{
		Placeholder:  "{{Fenix.TodayShiftDay(1,2)}}",
		FunctionName: "Fenix_TodayShiftDay",
		Arguments:    []string{"1", "2"},
	}
	logPlaceholderInputMatrix(t, "too-many-arguments", input)
	_, err = goFenixTodayShiftDay(input)
	logPlaceholderExecutionResult(t, "too-many-arguments", "", err)
	if err == nil {
		t.Fatalf("expected error when more than one argument is provided")
	}
	if strings.Contains(err.Error(), "exact 1 parameter argument") == false {
		t.Fatalf("unexpected error: %v", err)
	}

	input = GoPlaceholderInput{
		Placeholder:  "{{Fenix.TodayShiftDay()}}",
		FunctionName: "Fenix_TodayShiftDay",
		Arguments:    []string{},
	}
	logPlaceholderInputMatrix(t, "missing-argument", input)
	_, err = goFenixTodayShiftDay(input)
	logPlaceholderExecutionResult(t, "missing-argument", "", err)
	if err == nil {
		t.Fatalf("expected error when argument is missing")
	}
	if strings.Contains(err.Error(), "exact 1 parameter argument") == false {
		t.Fatalf("unexpected error: %v", err)
	}
}
