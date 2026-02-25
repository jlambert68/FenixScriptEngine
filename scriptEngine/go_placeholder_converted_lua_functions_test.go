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

func TestGoFenixControlledUniqueID_DocExamples_SupportedReplacements(t *testing.T) {
	// Freeze clock so replacement outputs remain stable.
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
		entropy     uint64
		input       string
		expectedOut string
	}{
		{name: "yyyy-mm-dd", input: "Date=%YYYY-MM-DD%", expectedOut: "Date=2026-02-24"},
		{name: "yyyymmdd", input: "Date=%YYYYMMDD%", expectedOut: "Date=20260224"},
		{name: "yymmdd", input: "Date=%YYMMDD%", expectedOut: "Date=260224"},
		{name: "time-hh-mm-ss", input: "Time=%hh:mm:ss%", expectedOut: "Time=13:07:09"},
		{name: "time-hh-mm-ss-dot", input: "Time=%hh.mm.ss%", expectedOut: "Time=13.07.09"},
		{name: "time-hhmmss", input: "Time=%hhmmss%", expectedOut: "Time=130709"},
		{name: "time-hhmm", input: "Time=%hhmm%", expectedOut: "Time=1307"},
		{name: "random-number", input: "Rand=%nnnnn%", expectedOut: "Rand=79410"},
		{name: "random-lower", input: "Rand=%a(5; 11)%", expectedOut: "Rand=gbrma"},
		{name: "random-upper", input: "Rand=%A(5; 10)%", expectedOut: "Rand=IMPVG"},
		{
			name:       "full-mix",
			arrayIndex: []int{0},
			entropy:    0,
			input: "Date: %YYYY-MM-DD%, Date: %YYYYMMDD%, Date: %YYMMDD%, Time: %hh:mm:ss%, Time: %hhmmss%, " +
				"Time: %hhmm%, Random Number: %nnnnn%, Random String: %a(5; 11)%, Random String Uppercase: %A(5; 10)%, " +
				"Time: %hh:mm:ss%, Time: %hh.mm.ss%",
			expectedOut: "Date: 2026-02-24, Date: 20260224, Date: 260224, Time: 13:07:09, Time: 130709, " +
				"Time: 1307, Random Number: 65505, Random String: gbrma, Random String Uppercase: IMPVG, " +
				"Time: 13:07:09, Time: 13.07.09",
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			result, err := goFenixControlledUniqueID(GoPlaceholderInput{
				FunctionName: "Fenix_ControlledUniqueId",
				ArrayIndexes: testCase.arrayIndex,
				Arguments:    []string{testCase.input},
				Entropy:      testCase.entropy,
			})
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if result != testCase.expectedOut {
				t.Fatalf("expected '%s', got '%s'", testCase.expectedOut, result)
			}
		})
	}
}

func TestGoFenixControlledUniqueID_DocExamples_ProcessingBehavior(t *testing.T) {
	cases := []struct {
		name        string
		arrayIndex  []int
		entropy     uint64
		input       string
		expectedOut string
	}{
		{
			name:        "default array index maps to 1",
			arrayIndex:  []int{},
			entropy:     0,
			input:       "Rand=%nnn%",
			expectedOut: "Rand=410",
		},
		{
			name:        "array index changes numeric seed",
			arrayIndex:  []int{2},
			entropy:     0,
			input:       "Rand=%nnn%",
			expectedOut: "Rand=511",
		},
		{
			name:        "extra entropy changes numeric seed",
			arrayIndex:  []int{1},
			entropy:     7,
			input:       "Rand=%nnn%",
			expectedOut: "Rand=840",
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			result, err := goFenixControlledUniqueID(GoPlaceholderInput{
				FunctionName: "Fenix_ControlledUniqueId",
				ArrayIndexes: testCase.arrayIndex,
				Arguments:    []string{testCase.input},
				Entropy:      testCase.entropy,
			})
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if result != testCase.expectedOut {
				t.Fatalf("expected '%s', got '%s'", testCase.expectedOut, result)
			}
		})
	}

	_, err := goFenixControlledUniqueID(GoPlaceholderInput{
		FunctionName: "Fenix_ControlledUniqueId",
		ArrayIndexes: []int{1, 2},
		Arguments:    []string{"X"},
		Entropy:      0,
	})
	if err == nil {
		t.Fatalf("expected error for more than one array index")
	}
	if strings.Contains(err.Error(), "there cant be more than 1 value") == false {
		t.Fatalf("unexpected error for array index validation: %v", err)
	}

	_, err = goFenixControlledUniqueID(GoPlaceholderInput{
		FunctionName: "Fenix_ControlledUniqueId",
		ArrayIndexes: []int{},
		Arguments:    []string{},
		Entropy:      0,
	})
	if err == nil {
		t.Fatalf("expected error when text input is missing")
	}
	if strings.Contains(err.Error(), "textToProcess must be a string, got nil") == false {
		t.Fatalf("unexpected error for missing text input: %v", err)
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
