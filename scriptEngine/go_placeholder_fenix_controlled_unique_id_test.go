package scriptEngine

import (
	"fmt"
	"regexp"
	"strconv"
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
		Placeholder:  "{{Fenix.ControlledUniqueId(..., false, 17)}}",
		FunctionName: "Fenix_ControlledUniqueId",
		ArrayIndexes: []int{1},
		Arguments: []string{
			"%YYYY-MM-DD%|%YYYYMMDD%|%YYMMDD%|%hh:mm:ss%|%hh.mm.ss%|%hhmmss%|%hhmm%|%mmss%|%n(3)%|%a(4)%|%A(4)%|%aA(4)%|%an(4)%|%An(4)%|%aAn(4)%",
			"false",
			"17",
		},
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
		t.Fatalf("expected all supported placeholders to be replaced, got: %s", resultOne)
	}

	expectedPattern := regexp.MustCompile(`^2026-02-24\|20260224\|260224\|13:07:09\|13\.07\.09\|130709\|1307\|0709\|[0-9]{3}\|[a-z]{4}\|[A-Z]{4}\|[a-zA-Z]{4}\|[a-z0-9]{4}\|[A-Z0-9]{4}\|[a-zA-Z0-9]{4}$`)
	if expectedPattern.MatchString(resultOne) == false {
		t.Fatalf("output did not match expected Jira pattern set: %s", resultOne)
	}
}

func TestGoFenixControlledUniqueID_JiraPatternsAndEntropyBehavior(t *testing.T) {
	originalTimeProvider := currentTimeProvider
	currentTimeProvider = func() time.Time {
		return time.Date(2026, time.February, 26, 8, 23, 59, 49208485, time.Local)
	}
	defer func() {
		currentTimeProvider = originalTimeProvider
	}()

	testCaseExecutionUUID := "f8c06f7e-0a8a-4d75-9f25-5e5fb8d2a6d3"

	cases := []struct {
		name            string
		arrayIndex      []int
		useEntropy      bool
		extraEntropy    uint64
		input           string
		expectedExact   string
		expectedPattern string
	}{
		{name: "yyyy-mm-dd", input: "%YYYY-MM-DD%", expectedExact: "2026-02-26"},
		{name: "yyyymmdd", input: "%YYYYMMDD%", expectedExact: "20260226"},
		{name: "yymmdd", input: "%YYMMDD%", expectedExact: "260226"},
		{name: "hh:mm:ss", input: "%hh:mm:ss%", expectedExact: "08:23:59"},
		{name: "hh.mm.ss", input: "%hh.mm.ss%", expectedExact: "08.23.59"},
		{name: "hhmmss", input: "%hhmmss%", expectedExact: "082359"},
		{name: "hhmm", input: "%hhmm%", expectedExact: "0823"},
		{name: "mmss", input: "%mmss%", expectedExact: "2359"},
		{
			name:          "mixed-date-time-components",
			input:         "%Year: YYYY, Month: MM, Day: DD, Hour: hh, Minute: mm, Second: ss, Milliseconds: ms, Microseconds: us, Nanoseconds: ns%",
			expectedExact: "%Year: 2026, Month: 02, Day: 26, Hour: 08, Minute: 23, Second: 59, Milliseconds: 049, Microseconds: 049208, Nanoseconds: 049208485%",
		},
		{name: "random-number-jira", input: "%n(5)%", useEntropy: false, extraEntropy: 1, expectedPattern: `^[0-9]{5}$`},
		{name: "random-lower-jira", input: "%a(5)%", useEntropy: false, extraEntropy: 1, expectedPattern: `^[a-z]{5}$`},
		{name: "random-upper-jira", input: "%A(5)%", useEntropy: false, extraEntropy: 1, expectedPattern: `^[A-Z]{5}$`},
		{name: "random-aa-jira", input: "%aA(5)%", useEntropy: false, extraEntropy: 1, expectedPattern: `^[a-zA-Z]{5}$`},
		{name: "random-an-jira", input: "%an(5)%", useEntropy: false, extraEntropy: 1, expectedPattern: `^[a-z0-9]{5}$`},
		{name: "random-An-jira", input: "%An(5)%", useEntropy: false, extraEntropy: 1, expectedPattern: `^[A-Z0-9]{5}$`},
		{name: "random-aAn-jira", input: "%aAn(5)%", useEntropy: false, extraEntropy: 1, expectedPattern: `^[a-zA-Z0-9]{5}$`},
		{name: "legacy-number-pattern-not-supported", input: "%nnnnn%", expectedExact: "%nnnnn%"},
		{name: "legacy-lower-pattern-not-supported", input: "%a(5; 11)%", expectedExact: "%a(5; 11)%"},
		{name: "legacy-upper-pattern-not-supported", input: "%A(5; 10)%", expectedExact: "%A(5; 10)%"},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			input := GoPlaceholderInput{
				Placeholder:                 "{{Fenix.ControlledUniqueId(" + testCase.input + ", " + strconv.FormatBool(testCase.useEntropy) + ", " + fmt.Sprintf("%d", testCase.extraEntropy) + ")}}",
				FunctionName:                "Fenix_ControlledUniqueId",
				ArrayIndexes:                testCase.arrayIndex,
				Arguments:                   []string{testCase.input, strconv.FormatBool(testCase.useEntropy), fmt.Sprintf("%d", testCase.extraEntropy)},
				TestCaseExecutionUUID:       testCaseExecutionUUID,
				UseEntropyFromExecutionUUID: testCase.useEntropy,
				ExtraEntropy:                testCase.extraEntropy,
			}
			logPlaceholderInputMatrix(t, testCase.name, input)

			resultOne, err := goFenixControlledUniqueID(input)
			logPlaceholderExecutionResult(t, testCase.name, resultOne, err)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}

			resultTwo, err := goFenixControlledUniqueID(input)
			logPlaceholderExecutionResult(t, testCase.name+"-repeat", resultTwo, err)
			if err != nil {
				t.Fatalf("expected no error on second call, got: %v", err)
			}
			if resultOne != resultTwo {
				t.Fatalf("expected deterministic output, got '%s' and '%s'", resultOne, resultTwo)
			}

			if testCase.expectedExact != "" {
				if resultOne != testCase.expectedExact {
					t.Fatalf("expected '%s', got '%s'", testCase.expectedExact, resultOne)
				}
				return
			}

			if regexp.MustCompile(testCase.expectedPattern).MatchString(resultOne) == false {
				t.Fatalf("result '%s' did not match expected pattern '%s'", resultOne, testCase.expectedPattern)
			}
		})
	}
}

func TestGoFenixControlledUniqueID_ShouldUseEntropyArguments(t *testing.T) {
	baseInput := "%n(4)%"
	testCaseExecutionUUID := "f8c06f7e-0a8a-4d75-9f25-5e5fb8d2a6d3"

	inputUseUUID := GoPlaceholderInput{
		Placeholder:           "{{Fenix.ControlledUniqueId(%n(4)%, true, 1)}}",
		FunctionName:          "Fenix_ControlledUniqueId",
		Arguments:             []string{baseInput, "true", "1"},
		TestCaseExecutionUUID: testCaseExecutionUUID,
	}
	inputWithoutUUID := GoPlaceholderInput{
		Placeholder:           "{{Fenix.ControlledUniqueId(%n(4)%, false, 1)}}",
		FunctionName:          "Fenix_ControlledUniqueId",
		Arguments:             []string{baseInput, "false", "1"},
		TestCaseExecutionUUID: testCaseExecutionUUID,
	}
	inputExtraEntropy := GoPlaceholderInput{
		Placeholder:           "{{Fenix.ControlledUniqueId(%n(4)%, true, 2)}}",
		FunctionName:          "Fenix_ControlledUniqueId",
		Arguments:             []string{baseInput, "true", "2"},
		TestCaseExecutionUUID: testCaseExecutionUUID,
	}

	logPlaceholderInputMatrix(t, "entropy-use-uuid", inputUseUUID)
	useUUIDValue, err := goFenixControlledUniqueID(inputUseUUID)
	logPlaceholderExecutionResult(t, "entropy-use-uuid", useUUIDValue, err)
	if err != nil {
		t.Fatalf("expected no error for inputUseUUID, got: %v", err)
	}

	logPlaceholderInputMatrix(t, "entropy-without-uuid", inputWithoutUUID)
	withoutUUIDValue, err := goFenixControlledUniqueID(inputWithoutUUID)
	logPlaceholderExecutionResult(t, "entropy-without-uuid", withoutUUIDValue, err)
	if err != nil {
		t.Fatalf("expected no error for inputWithoutUUID, got: %v", err)
	}

	logPlaceholderInputMatrix(t, "entropy-extra-entropy", inputExtraEntropy)
	extraEntropyValue, err := goFenixControlledUniqueID(inputExtraEntropy)
	logPlaceholderExecutionResult(t, "entropy-extra-entropy", extraEntropyValue, err)
	if err != nil {
		t.Fatalf("expected no error for inputExtraEntropy, got: %v", err)
	}

	if useUUIDValue == withoutUUIDValue {
		t.Fatalf("expected different outputs when toggling useEntropyFromExecutionUUID, got '%s'", useUUIDValue)
	}
	if useUUIDValue == extraEntropyValue {
		t.Fatalf("expected different outputs when changing extraEntropy, got '%s'", useUUIDValue)
	}
}

func TestGoFenixControlledUniqueID_ShouldValidateInput(t *testing.T) {
	input := GoPlaceholderInput{
		Placeholder:  "{{Fenix.ControlledUniqueId[1,2](X, false, 0)}}",
		FunctionName: "Fenix_ControlledUniqueId",
		ArrayIndexes: []int{1, 2},
		Arguments:    []string{"X", "false", "0"},
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
		Placeholder:  "{{Fenix.ControlledUniqueId(X)}}",
		FunctionName: "Fenix_ControlledUniqueId",
		ArrayIndexes: []int{},
		Arguments:    []string{"X"},
	}
	logPlaceholderInputMatrix(t, "missing-arguments", input)
	_, err = goFenixControlledUniqueID(input)
	logPlaceholderExecutionResult(t, "missing-arguments", "", err)
	if err == nil {
		t.Fatalf("expected error when required arguments are missing")
	}
	if strings.Contains(err.Error(), "exact 3 function arguments") == false {
		t.Fatalf("unexpected error for missing arguments: %v", err)
	}

	input = GoPlaceholderInput{
		Placeholder:  "{{Fenix.ControlledUniqueId(A,B,C,D)}}",
		FunctionName: "Fenix_ControlledUniqueId",
		ArrayIndexes: []int{},
		Arguments:    []string{"A", "B", "C", "D"},
	}
	logPlaceholderInputMatrix(t, "too-many-arguments", input)
	_, err = goFenixControlledUniqueID(input)
	logPlaceholderExecutionResult(t, "too-many-arguments", "", err)
	if err == nil {
		t.Fatalf("expected error when too many arguments are provided")
	}
	if strings.Contains(err.Error(), "exact 3 function arguments") == false {
		t.Fatalf("unexpected error for too many arguments: %v", err)
	}

	input = GoPlaceholderInput{
		Placeholder:  "{{Fenix.ControlledUniqueId(X, maybe, 1)}}",
		FunctionName: "Fenix_ControlledUniqueId",
		ArrayIndexes: []int{},
		Arguments:    []string{"X", "maybe", "1"},
	}
	logPlaceholderInputMatrix(t, "invalid-boolean-argument", input)
	_, err = goFenixControlledUniqueID(input)
	logPlaceholderExecutionResult(t, "invalid-boolean-argument", "", err)
	if err == nil {
		t.Fatalf("expected error when second argument is not Boolean")
	}
	if strings.Contains(err.Error(), "second function argument must be a Boolean") == false {
		t.Fatalf("unexpected error for invalid boolean argument: %v", err)
	}

	input = GoPlaceholderInput{
		Placeholder:  "{{Fenix.ControlledUniqueId(X, true, entropy)}}",
		FunctionName: "Fenix_ControlledUniqueId",
		ArrayIndexes: []int{},
		Arguments:    []string{"X", "true", "entropy"},
	}
	logPlaceholderInputMatrix(t, "invalid-integer-argument", input)
	_, err = goFenixControlledUniqueID(input)
	logPlaceholderExecutionResult(t, "invalid-integer-argument", "", err)
	if err == nil {
		t.Fatalf("expected error when third argument is not Integer")
	}
	if strings.Contains(err.Error(), "third function argument must be an Integer") == false {
		t.Fatalf("unexpected error for invalid integer argument: %v", err)
	}
}
