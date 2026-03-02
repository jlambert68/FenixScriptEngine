package scriptEngine

import (
	"regexp"
	"strings"
	"testing"
)

func TestGoFenixRandomPositiveDecimalValue_ShouldBeDeterministicAndPadded(t *testing.T) {
	input := GoPlaceholderInput{
		Placeholder:  "{{Fenix.RandomPositiveDecimalValue[2](1, 2, 3, 4, \".\")}}",
		FunctionName: "Fenix_RandomPositiveDecimalValue",
		ArrayIndexes: []int{2},
		Arguments:    []string{"1", "2", "3", "4", "."},
		Entropy:      12345,
	}

	logPlaceholderInputMatrix(t, "deterministic-and-padded", input)
	resultOne, err := goFenixRandomPositiveDecimalValue(input)
	logPlaceholderExecutionResult(t, "deterministic-and-padded", resultOne, err)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	logPlaceholderInputMatrix(t, "deterministic-and-padded-repeat", input)
	resultTwo, err := goFenixRandomPositiveDecimalValue(input)
	logPlaceholderExecutionResult(t, "deterministic-and-padded-repeat", resultTwo, err)
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

func TestGoFenixRandomPositiveDecimalValue_JiraExamplesPatterns(t *testing.T) {
	cases := []struct {
		name       string
		arrayIndex []int
		args       []string
		entropy    uint64
		pattern    string
	}{
		{name: "[0] 1,2,3,4,'.'", arrayIndex: []int{0}, args: []string{"1", "2", "3", "4", "."}, entropy: 0, pattern: `^\d{3}\.\d{4}$`},
		{name: "[2] 2,2,2,2,'.'", arrayIndex: []int{2}, args: []string{"2", "2", "2", "2", "."}, entropy: 0, pattern: `^\d{2}\.\d{2}$`},
		{name: "[3] 0,2,1,2,'.'", arrayIndex: []int{3}, args: []string{"0", "2", "1", "2", "."}, entropy: 0, pattern: `^0\.\d{2}$`},
		{name: "[4] 3,0,3,0,'.'", arrayIndex: []int{4}, args: []string{"3", "0", "3", "0", "."}, entropy: 0, pattern: `^\d{3}$`},
		{name: "[6] 3,2,3,2,','", arrayIndex: []int{6}, args: []string{"3", "2", "3", "2", ","}, entropy: 0, pattern: `^\d{3},\d{2}$`},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			input := GoPlaceholderInput{
				Placeholder:  "{{Fenix.RandomPositiveDecimalValue(" + strings.Join(testCase.args, ", ") + ")}}",
				FunctionName: "Fenix_RandomPositiveDecimalValue",
				ArrayIndexes: testCase.arrayIndex,
				Arguments:    testCase.args,
				Entropy:      testCase.entropy,
			}
			logPlaceholderInputMatrix(t, testCase.name, input)

			resultOne, err := goFenixRandomPositiveDecimalValue(input)
			logPlaceholderExecutionResult(t, testCase.name, resultOne, err)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			logPlaceholderInputMatrix(t, testCase.name+"-repeat", input)
			resultTwo, err := goFenixRandomPositiveDecimalValue(input)
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

func TestGoFenixRandomPositiveDecimalValue_ShouldFailOnInvalidArguments(t *testing.T) {
	cases := []struct {
		name          string
		arrayIndex    []int
		args          []string
		expectedError string
	}{
		{name: "one argument", arrayIndex: []int{1}, args: []string{"0"}, expectedError: "exact 5 function parameter"},
		{name: "four arguments", arrayIndex: []int{1}, args: []string{"1", "2", "3", "4"}, expectedError: "exact 5 function parameter"},
		{name: "six arguments", arrayIndex: []int{1}, args: []string{"1", "2", "3", "4", ".", "x"}, expectedError: "exact 5 function parameter"},
		{name: "non-integer among first four", arrayIndex: []int{1}, args: []string{"1", "two", "3", "4", "."}, expectedError: "first four functions parameters must be of type Integer"},
		{name: "empty decimal point", arrayIndex: []int{1}, args: []string{"1", "2", "3", "4", ""}, expectedError: "decimalPointCharacter must be provided"},
		{name: "multi-char decimal point", arrayIndex: []int{1}, args: []string{"1", "2", "3", "4", ".."}, expectedError: "decimalPointCharacter must be a single character"},
		{name: "too many array indexes", arrayIndex: []int{1, 2}, args: []string{"2", "3", "2", "3", "."}, expectedError: "maximum of one value"},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			input := GoPlaceholderInput{
				Placeholder:  "{{Fenix.RandomPositiveDecimalValue(" + strings.Join(testCase.args, ", ") + ")}}",
				FunctionName: "Fenix_RandomPositiveDecimalValue",
				ArrayIndexes: testCase.arrayIndex,
				Arguments:    testCase.args,
				Entropy:      0,
			}
			logPlaceholderInputMatrix(t, testCase.name, input)
			_, err := goFenixRandomPositiveDecimalValue(input)
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
