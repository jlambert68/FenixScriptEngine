package scriptEngine

import (
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestGoFenixTodayShiftDay_ExamplesFromMain(t *testing.T) {
	// Mirrors main.go examples for TodayShiftDay with stable date assertions.
	originalTimeProvider := currentTimeProvider
	currentTimeProvider = func() time.Time {
		return time.Date(2026, time.February, 24, 12, 30, 45, 0, time.Local)
	}
	defer func() {
		currentTimeProvider = originalTimeProvider
	}()

	cases := []struct {
		name     string
		args     []string
		expected string
	}{
		{name: "no args", args: []string{}, expected: "2026-02-24"},
		{name: "shift -1", args: []string{"-1"}, expected: "2026-02-23"},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			result, err := goFenixTodayShiftDay(GoPlaceholderInput{
				FunctionName: "Fenix_TodayShiftDay",
				ArrayIndexes: []int{},
				Arguments:    testCase.args,
			})
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if result != testCase.expected {
				t.Fatalf("expected %s, got %s", testCase.expected, result)
			}
		})
	}
}

func TestGoFenixControlledUniqueID_ExamplesFromMain(t *testing.T) {
	// Covers the two ControlledUniqueId examples wired in main.go.
	originalTimeProvider := currentTimeProvider
	currentTimeProvider = func() time.Time {
		return time.Date(2026, time.February, 24, 13, 7, 9, 0, time.Local)
	}
	defer func() {
		currentTimeProvider = originalTimeProvider
	}()

	cases := []struct {
		name        string
		arrayIndex  []int
		inputString string
	}{
		{
			name:        "simple date replacement",
			arrayIndex:  []int{},
			inputString: "Date: %YYYY-MM-DD%",
		},
		{
			name:       "extended token replacement",
			arrayIndex: []int{0},
			inputString: "Date: %YYYY-MM-DD%, Date: %YYYYMMDD%, Date: %YYMMDD%, Time: %hh:mm:ss%, Time: %hhmmss%, " +
				"Time: %hhmm%, Random Number: %nnnnn%, Random String: %a(5; 11)%, Random String Uppercase: %A(5; 10)%, " +
				"Time: %hh:mm:ss%, Time: %hh.mm.ss% ",
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			input := GoPlaceholderInput{
				FunctionName: "Fenix_ControlledUniqueId",
				ArrayIndexes: testCase.arrayIndex,
				Arguments:    []string{testCase.inputString},
				Entropy:      0,
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
				t.Fatalf("expected deterministic output, got '%s' and '%s'", resultOne, resultTwo)
			}
			if strings.Contains(resultOne, "%") {
				t.Fatalf("expected no placeholder tokens left in output, got: %s", resultOne)
			}
		})
	}
}

func TestGoFenixRandomPositiveDecimalValue_ExamplesFromMain_SuccessPatterns(t *testing.T) {
	// main.go has many value examples; assert format contracts instead of exact RNG digits.
	cases := []struct {
		name       string
		arrayIndex []int
		args       []string
		entropy    uint64
		pattern    string
	}{
		{name: "2,3", arrayIndex: []int{}, args: []string{"2", "3"}, entropy: 0, pattern: `^\d{1,2}\.\d{3}$`},
		{name: "[1] 2,3", arrayIndex: []int{1}, args: []string{"2", "3"}, entropy: 0, pattern: `^\d{1,2}\.\d{3}$`},
		{name: "[3] 2,3", arrayIndex: []int{3}, args: []string{"2", "3"}, entropy: 0, pattern: `^\d{1,2}\.\d{3}$`},
		{name: "[2] 2,3", arrayIndex: []int{2}, args: []string{"2", "3"}, entropy: 0, pattern: `^\d{1,2}\.\d{3}$`},
		{name: "1,2", arrayIndex: []int{}, args: []string{"1", "2"}, entropy: 0, pattern: `^\d\.\d{2}$`},
		{name: "[2] 1,2", arrayIndex: []int{2}, args: []string{"1", "2"}, entropy: 0, pattern: `^\d\.\d{2}$`},
		{name: "1,1", arrayIndex: []int{}, args: []string{"1", "1"}, entropy: 0, pattern: `^\d\.\d$`},
		{name: "1,1 entropy=1", arrayIndex: []int{}, args: []string{"1", "1"}, entropy: 1, pattern: `^\d\.\d$`},
		{name: "[1] 1,1 entropy=1", arrayIndex: []int{1}, args: []string{"1", "1"}, entropy: 1, pattern: `^\d\.\d$`},
		{name: "0,1", arrayIndex: []int{}, args: []string{"0", "1"}, entropy: 0, pattern: `^0\.\d$`},
		{name: "[1] 1,0", arrayIndex: []int{1}, args: []string{"1", "0"}, entropy: 0, pattern: `^\d$`},
		{name: "[1] 0,0", arrayIndex: []int{1}, args: []string{"0", "0"}, entropy: 0, pattern: `^0$`},
		{name: "[1] 0,0,2,3", arrayIndex: []int{1}, args: []string{"0", "0", "2", "3"}, entropy: 0, pattern: `^00$`},
		{name: "[1] 0,2,3,4", arrayIndex: []int{1}, args: []string{"0", "2", "3", "4"}, entropy: 0, pattern: `^000\.\d{4}$`},
		{name: "[1] 2,2,3,4", arrayIndex: []int{1}, args: []string{"2", "2", "3", "4"}, entropy: 0, pattern: `^\d{3}\.\d{4}$`},
		{name: "6,6", arrayIndex: []int{}, args: []string{"6", "6"}, entropy: 0, pattern: `^\d{1,6}\.\d{6}$`},
		{name: "6,10", arrayIndex: []int{}, args: []string{"6", "10"}, entropy: 0, pattern: `^\d{1,6}\.\d{10}$`},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			input := GoPlaceholderInput{
				FunctionName: "Fenix_RandomPositiveDecimalValue",
				ArrayIndexes: testCase.arrayIndex,
				Arguments:    testCase.args,
				Entropy:      testCase.entropy,
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

func TestGoFenixRandomPositiveDecimalValue_ExamplesFromMain_ErrorCases(t *testing.T) {
	// Ensure all error scenarios demonstrated in main.go stay validated in code.
	cases := []struct {
		name          string
		arrayIndex    []int
		args          []string
		expectedError string
	}{
		{
			name:          "one argument only",
			arrayIndex:    []int{1},
			args:          []string{"0"},
			expectedError: "exact 2 or 4 function parameter",
		},
		{
			name:          "empty arguments",
			arrayIndex:    []int{1},
			args:          []string{},
			expectedError: "exact 2 or 4 function parameter",
		},
		{
			name:          "three arguments",
			arrayIndex:    []int{1},
			args:          []string{"1", "2", "3"},
			expectedError: "exact 2 or 4 function parameter",
		},
		{
			name:          "too many array indexes",
			arrayIndex:    []int{1, 2},
			args:          []string{"2", "3"},
			expectedError: "maximum of one value",
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			_, err := goFenixRandomPositiveDecimalValue(GoPlaceholderInput{
				FunctionName: "Fenix_RandomPositiveDecimalValue",
				ArrayIndexes: testCase.arrayIndex,
				Arguments:    testCase.args,
				Entropy:      0,
			})
			if err == nil {
				t.Fatalf("expected error")
			}
			if strings.Contains(err.Error(), testCase.expectedError) == false {
				t.Fatalf("expected error containing '%s', got '%s'", testCase.expectedError, err.Error())
			}
		})
	}
}

func TestGoFenixRandomPositiveDecimalValueSum_ExamplesFromMain(t *testing.T) {
	// Covers sum examples including subtraction and padding variants.
	cases := []struct {
		name       string
		arrayIndex []int
		args       []string
		pattern    string
	}{
		{name: "[1] 2,3", arrayIndex: []int{1}, args: []string{"2", "3"}, pattern: `^\d{1,2}\.\d{3}$`},
		{name: "[-1,2] 2,3", arrayIndex: []int{-1, 2}, args: []string{"2", "3"}, pattern: `^-?\d{1,3}\.\d{3}$`},
		{name: "[1,-2] 2,3", arrayIndex: []int{1, -2}, args: []string{"2", "3"}, pattern: `^-?\d{1,3}\.\d{3}$`},
		{name: "[-1,-2] 2,3", arrayIndex: []int{-1, -2}, args: []string{"2", "3"}, pattern: `^-\d{1,3}\.\d{3}$`},
		{name: "[1,2,3] 2,3", arrayIndex: []int{1, 2, 3}, args: []string{"2", "3"}, pattern: `^\d{1,3}\.\d{3}$`},
		{name: "[1,2,3] 2,3,2,3", arrayIndex: []int{1, 2, 3}, args: []string{"2", "3", "2", "3"}, pattern: `^\d{2,3}\.\d{3}$`},
		{name: "[1,2,3] 2,3,4,4", arrayIndex: []int{1, 2, 3}, args: []string{"2", "3", "4", "4"}, pattern: `^\d{4}\.\d{4}$`},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			input := GoPlaceholderInput{
				FunctionName: "Fenix_RandomPositiveDecimalValue_Sum",
				ArrayIndexes: testCase.arrayIndex,
				Arguments:    testCase.args,
				Entropy:      0,
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
