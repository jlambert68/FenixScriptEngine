package scriptEngine

import (
	"hash/crc32"
	"regexp"
	"strings"
	"testing"
)

func TestGoFenixRandomPositiveDecimalValueSum_ShouldBeDeterministic(t *testing.T) {
	input := GoPlaceholderInput{
		Placeholder:  "{{Fenix.RandomPositiveDecimalValue.Sum[1,-2,3](2, 2, 3, 3, \".\")}}",
		FunctionName: "Fenix_RandomPositiveDecimalValue_Sum",
		ArrayIndexes: []int{1, -2, 3},
		Arguments:    []string{"2", "2", "3", "3", "."},
		Entropy:      999,
	}

	logPlaceholderInputMatrix(t, "sum-deterministic", input)
	resultOne, err := goFenixRandomPositiveDecimalValueSum(input)
	logPlaceholderExecutionResult(t, "sum-deterministic", resultOne, err)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	logPlaceholderInputMatrix(t, "sum-deterministic-repeat", input)
	resultTwo, err := goFenixRandomPositiveDecimalValueSum(input)
	logPlaceholderExecutionResult(t, "sum-deterministic-repeat", resultTwo, err)
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

func TestGoFenixRandomPositiveDecimalValueSum_JiraExamplesPatterns(t *testing.T) {
	cases := []struct {
		name       string
		arrayIndex []int
		args       []string
		pattern    string
	}{
		{name: "[1] 2,3,2,3,'.'", arrayIndex: []int{1}, args: []string{"2", "3", "2", "3", "."}, pattern: `^\d{2}\.\d{3}$`},
		{name: "[-1,2] 2,3,3,3,'.'", arrayIndex: []int{-1, 2}, args: []string{"2", "3", "3", "3", "."}, pattern: `^-?\d{1,3}\.\d{3}$`},
		{name: "[1,-2] 2,3,3,3,'.'", arrayIndex: []int{1, -2}, args: []string{"2", "3", "3", "3", "."}, pattern: `^-?\d{1,3}\.\d{3}$`},
		{name: "[1,2,3] 2,3,4,4,'.'", arrayIndex: []int{1, 2, 3}, args: []string{"2", "3", "4", "4", "."}, pattern: `^\d{4}\.\d{4}$`},
		{name: "[1,2] 2,3,4,4,','", arrayIndex: []int{1, 2}, args: []string{"2", "3", "4", "4", ","}, pattern: `^\d{4},\d{4}$`},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			input := GoPlaceholderInput{
				Placeholder:  "{{Fenix.RandomPositiveDecimalValue.Sum(" + strings.Join(testCase.args, ", ") + ")}}",
				FunctionName: "Fenix_RandomPositiveDecimalValue_Sum",
				ArrayIndexes: testCase.arrayIndex,
				Arguments:    testCase.args,
				Entropy:      0,
			}
			logPlaceholderInputMatrix(t, testCase.name, input)

			resultOne, err := goFenixRandomPositiveDecimalValueSum(input)
			logPlaceholderExecutionResult(t, testCase.name, resultOne, err)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			logPlaceholderInputMatrix(t, testCase.name+"-repeat", input)
			resultTwo, err := goFenixRandomPositiveDecimalValueSum(input)
			logPlaceholderExecutionResult(t, testCase.name+"-repeat", resultTwo, err)
			if err != nil {
				t.Fatalf("expected no error on second call, got: %v", err)
			}

			if resultOne != resultTwo {
				t.Fatalf("expected deterministic output, got '%s' and '%s'", resultOne, resultTwo)
			}

			matches, compileErr := regexp.MatchString(testCase.pattern, resultOne)
			if compileErr != nil {
				t.Fatalf("invalid regex in test: %v", compileErr)
			}
			if matches == false {
				t.Fatalf("output '%s' did not match expected pattern '%s'", resultOne, testCase.pattern)
			}
		})
	}
}

func TestGoFenixRandomPositiveDecimalValueSum_ShouldFormatNegativeWithPadding(t *testing.T) {
	testCaseExecutionUUID := "f8c06f7e-0a8a-4d75-9f25-5e5fb8d2a6d3"
	input := GoPlaceholderInput{
		Placeholder:  "{{Fenix.RandomPositiveDecimalValue.Sum[-1,2](2, 3, 3, 3, \".\")}}",
		FunctionName: "Fenix_RandomPositiveDecimalValue_Sum",
		ArrayIndexes: []int{-1, 2},
		Arguments:    []string{"2", "3", "3", "3", "."},
		Entropy:      uint64(crc32.ChecksumIEEE([]byte(testCaseExecutionUUID))),
	}
	logPlaceholderInputMatrix(t, "negative-padding-format", input)

	result, err := goFenixRandomPositiveDecimalValueSum(input)
	logPlaceholderExecutionResult(t, "negative-padding-format", result, err)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result != "-044.613" {
		t.Fatalf("expected '-044.613', got '%s'", result)
	}
}

func TestGoFenixRandomPositiveDecimalValueSum_ShouldFailOnInvalidArguments(t *testing.T) {
	cases := []struct {
		name          string
		arrayIndex    []int
		args          []string
		expectedError string
	}{
		{name: "four arguments", arrayIndex: []int{1}, args: []string{"2", "3", "4", "4"}, expectedError: "exact 5 function parameter"},
		{name: "six arguments", arrayIndex: []int{1}, args: []string{"2", "3", "4", "4", ".", "x"}, expectedError: "exact 5 function parameter"},
		{name: "non-integer among first four", arrayIndex: []int{1}, args: []string{"2", "three", "4", "4", "."}, expectedError: "first four functions parameters must be of type Integer"},
		{name: "empty decimal point", arrayIndex: []int{1}, args: []string{"2", "3", "4", "4", ""}, expectedError: "decimalPointCharacter must be provided"},
		{name: "multi-char decimal point", arrayIndex: []int{1}, args: []string{"2", "3", "4", "4", ".."}, expectedError: "decimalPointCharacter must be a single character"},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			input := GoPlaceholderInput{
				Placeholder:  "{{Fenix.RandomPositiveDecimalValue.Sum(" + strings.Join(testCase.args, ", ") + ")}}",
				FunctionName: "Fenix_RandomPositiveDecimalValue_Sum",
				ArrayIndexes: testCase.arrayIndex,
				Arguments:    testCase.args,
				Entropy:      0,
			}
			logPlaceholderInputMatrix(t, testCase.name, input)
			_, err := goFenixRandomPositiveDecimalValueSum(input)
			logPlaceholderExecutionResult(t, testCase.name, "", err)
			if err == nil {
				t.Fatalf("expected error")
			}
			if strings.Contains(err.Error(), testCase.expectedError) == false {
				t.Fatalf("expected error containing '%s', got '%s'", testCase.expectedError, err.Error())
			}
		})
	}
}
