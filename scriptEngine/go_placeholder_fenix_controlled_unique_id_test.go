package scriptEngine

import (
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestGoFenixControlledUniqueID_ShouldReplaceDateTimeAndRandomPatterns(t *testing.T) {
	originalTimeProvider := currentTimeProvider
	currentTimeProvider = func() time.Time {
		return time.Date(2026, time.February, 24, 13, 7, 9, 0, time.Local)
	}
	defer func() {
		currentTimeProvider = originalTimeProvider
	}()

	input := GoPlaceholderInput{
		Placeholder:  "{{Fenix.ControlledUniqueId[1](%YYYY-MM-DD%|%YYYYMMDD%|%YYMMDD%|%hh:mm:ss%|%hh.mm.ss%|%hhmmss%|%hhmm%|%nnn%|%a(4; 11)%|%A(4; 11)%)}}",
		FunctionName: "Fenix_ControlledUniqueId",
		ArrayIndexes: []int{1},
		Arguments: []string{
			"%YYYY-MM-DD%|%YYYYMMDD%|%YYMMDD%|%hh:mm:ss%|%hh.mm.ss%|%hhmmss%|%hhmm%|%nnn%|%a(4; 11)%|%A(4; 11)%",
		},
		Entropy: 17,
	}

	logPlaceholderInputMatrix(t, "replace-date-time-random", input)
	resultOne, err := goFenixControlledUniqueID(input)
	logPlaceholderExecutionResult(t, "replace-date-time-random", resultOne, err)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	logPlaceholderInputMatrix(t, "replace-date-time-random-repeat", input)
	resultTwo, err := goFenixControlledUniqueID(input)
	logPlaceholderExecutionResult(t, "replace-date-time-random-repeat", resultTwo, err)
	if err != nil {
		t.Fatalf("expected no error on second call, got: %v", err)
	}

	if resultOne != resultTwo {
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

func TestGoFenixControlledUniqueID_JiraPatternsAndEntropyBehavior(t *testing.T) {
	originalTimeProvider := currentTimeProvider
	currentTimeProvider = func() time.Time {
		return time.Date(2026, time.February, 26, 8, 23, 59, 0, time.Local)
	}
	defer func() {
		currentTimeProvider = originalTimeProvider
	}()

	cases := []struct {
		name        string
		arrayIndex  []int
		entropy     uint64
		input       string
		expectedOut string
	}{
		{name: "yyyy-mm-dd", input: "%YYYY-MM-DD%", expectedOut: "2026-02-26"},
		{name: "yyyymmdd", input: "%YYYYMMDD%", expectedOut: "20260226"},
		{name: "yymmdd", input: "%YYMMDD%", expectedOut: "260226"},
		{name: "hh:mm:ss", input: "%hh:mm:ss%", expectedOut: "08:23:59"},
		{name: "hh.mm.ss", input: "%hh.mm.ss%", expectedOut: "08.23.59"},
		{name: "hhmmss", input: "%hhmmss%", expectedOut: "082359"},
		{name: "hhmm", input: "%hhmm%", expectedOut: "0823"},
		{name: "mmss-is-not-replaced", input: "%mmss%", expectedOut: "%mmss%"},
		{name: "random-number-n-pattern", input: "%nnnnn%", expectedOut: "79410"},
		{name: "random-lower", input: "%a(5; 11)%", expectedOut: "gbrma"},
		{name: "random-upper", input: "%A(5; 10)%", expectedOut: "IMPVG"},
		{name: "legacy-jira-n-parenthesis-is-kept", input: "%n(4)%", expectedOut: "%n(4)%"},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			input := GoPlaceholderInput{
				Placeholder:  "{{Fenix.ControlledUniqueId(" + testCase.input + ")}}",
				FunctionName: "Fenix_ControlledUniqueId",
				ArrayIndexes: testCase.arrayIndex,
				Arguments:    []string{testCase.input},
				Entropy:      testCase.entropy,
			}
			logPlaceholderInputMatrix(t, testCase.name, input)

			result, err := goFenixControlledUniqueID(input)
			logPlaceholderExecutionResult(t, testCase.name, result, err)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if result != testCase.expectedOut {
				t.Fatalf("expected '%s', got '%s'", testCase.expectedOut, result)
			}
		})
	}
}

func TestGoFenixControlledUniqueID_ShouldValidateInput(t *testing.T) {
	input := GoPlaceholderInput{
		Placeholder:  "{{Fenix.ControlledUniqueId[1,2](X)}}",
		FunctionName: "Fenix_ControlledUniqueId",
		ArrayIndexes: []int{1, 2},
		Arguments:    []string{"X"},
		Entropy:      0,
	}
	logPlaceholderInputMatrix(t, "too-many-array-indexes", input)
	_, err := goFenixControlledUniqueID(input)
	logPlaceholderExecutionResult(t, "too-many-array-indexes", "", err)
	if err == nil {
		t.Fatalf("expected error for more than one array index")
	}
	if strings.Contains(err.Error(), "there cant be more than 1 value") == false {
		t.Fatalf("unexpected error for array index validation: %v", err)
	}

	input = GoPlaceholderInput{
		Placeholder:  "{{Fenix.ControlledUniqueId()}}",
		FunctionName: "Fenix_ControlledUniqueId",
		ArrayIndexes: []int{},
		Arguments:    []string{},
		Entropy:      0,
	}
	logPlaceholderInputMatrix(t, "missing-text-argument", input)
	_, err = goFenixControlledUniqueID(input)
	logPlaceholderExecutionResult(t, "missing-text-argument", "", err)
	if err == nil {
		t.Fatalf("expected error when text input is missing")
	}
	if strings.Contains(err.Error(), "textToProcess must be a string, got nil") == false {
		t.Fatalf("unexpected error for missing text input: %v", err)
	}

	input = GoPlaceholderInput{
		Placeholder:  "{{Fenix.ControlledUniqueId(A,B)}}",
		FunctionName: "Fenix_ControlledUniqueId",
		ArrayIndexes: []int{},
		Arguments:    []string{"A", "B"},
		Entropy:      0,
	}
	logPlaceholderInputMatrix(t, "too-many-text-arguments", input)
	_, err = goFenixControlledUniqueID(input)
	logPlaceholderExecutionResult(t, "too-many-text-arguments", "", err)
	if err == nil {
		t.Fatalf("expected error when more than one text argument is provided")
	}
	if strings.Contains(err.Error(), "exact 1 function argument") == false {
		t.Fatalf("unexpected error for too many text arguments: %v", err)
	}
}
