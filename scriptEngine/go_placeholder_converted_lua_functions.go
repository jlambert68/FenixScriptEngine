package scriptEngine

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// This file contains Go ports of previously Lua-only placeholder functions.
// The goal is to preserve deterministic behavior and validation contracts.
var (
	// Supported token patterns from Fenix_ControlledUniqueId.lua.
	controlledUniqueRandomNumberPattern = regexp.MustCompile(`%(n+)%`)
	controlledUniqueRandomLowerPattern  = regexp.MustCompile(`%a\((\d+);\s*(\d+)\)%`)
	controlledUniqueRandomUpperPattern  = regexp.MustCompile(`%A\((\d+);\s*(\d+)\)%`)
)

// goFenixControlledUniqueID mirrors token replacement behavior from Fenix_ControlledUniqueId.lua.
func goFenixControlledUniqueID(input GoPlaceholderInput) (string, error) {
	if len(input.ArrayIndexes) > 1 {
		return "", fmt.Errorf("Error - there cant be more than 1 value in 'arrayPositionTable'. %v", input.ArrayIndexes)
	}

	arrayPositionToUse := 1
	if len(input.ArrayIndexes) == 1 {
		arrayPositionToUse = input.ArrayIndexes[0]
	}

	if len(input.Arguments) == 0 {
		return "", fmt.Errorf("textToProcess must be a string, got nil")
	}

	textToProcess := input.Arguments[0]
	now := currentTimeProvider().In(time.Local)

	result := textToProcess
	result = strings.ReplaceAll(result, "%YYYY-MM-DD%", now.Format("2006-01-02"))
	result = strings.ReplaceAll(result, "%YYYYMMDD%", now.Format("20060102"))
	result = strings.ReplaceAll(result, "%YYMMDD%", now.Format("060102"))
	result = strings.ReplaceAll(result, "%hh:mm:ss%", now.Format("15:04:05"))
	result = strings.ReplaceAll(result, "%hh.mm.ss%", now.Format("15.04.05"))
	result = strings.ReplaceAll(result, "%hhmmss%", now.Format("150405"))
	result = strings.ReplaceAll(result, "%hhmm%", now.Format("1504"))

	// Random value generation must stay deterministic for same array index + entropy.
	seed := int64(arrayPositionToUse) + int64(input.Entropy)
	result = replaceControlledUniqueRandomNumberPatterns(result, seed)
	result = replaceControlledUniqueRandomStringPatterns(result, controlledUniqueRandomLowerPattern, false)
	result = replaceControlledUniqueRandomStringPatterns(result, controlledUniqueRandomUpperPattern, true)

	return result, nil
}

// replaceControlledUniqueRandomNumberPatterns replaces %nnn% style numeric tokens.
func replaceControlledUniqueRandomNumberPatterns(input string, seed int64) string {
	rng := rand.New(rand.NewSource(seed))

	return controlledUniqueRandomNumberPattern.ReplaceAllStringFunc(input, func(match string) string {
		numberOfDigits := len(match) - 2 // `%nnn%` => `nnn`
		if numberOfDigits <= 0 {
			return match
		}

		maxValue, ok := pow10Int64(numberOfDigits)
		if ok == false || maxValue <= 0 {
			// Keep token unchanged if width exceeds supported int64-safe bounds.
			return match
		}

		randomValue := rng.Int63n(maxValue)
		return fmt.Sprintf("%0*d", numberOfDigits, randomValue)
	})
}

// replaceControlledUniqueRandomStringPatterns replaces %a(len;seed)% and %A(len;seed)% tokens.
func replaceControlledUniqueRandomStringPatterns(input string, pattern *regexp.Regexp, upper bool) string {
	matchIndexes := pattern.FindAllStringSubmatchIndex(input, -1)
	if len(matchIndexes) == 0 {
		return input
	}

	var builder strings.Builder
	lastEnd := 0

	for _, indexes := range matchIndexes {
		fullStart := indexes[0]
		fullEnd := indexes[1]
		lengthStart := indexes[2]
		lengthEnd := indexes[3]
		seedStart := indexes[4]
		seedEnd := indexes[5]

		builder.WriteString(input[lastEnd:fullStart])

		lengthAsString := input[lengthStart:lengthEnd]
		seedAsString := input[seedStart:seedEnd]
		token := input[fullStart:fullEnd]

		lengthValue, lengthErr := strconv.Atoi(lengthAsString)
		seedValue, seedErr := strconv.ParseInt(seedAsString, 10, 64)
		if lengthErr != nil || seedErr != nil || lengthValue < 0 {
			builder.WriteString(token)
		} else {
			builder.WriteString(generateSeededAlphabetString(lengthValue, upper, seedValue))
		}

		lastEnd = fullEnd
	}

	builder.WriteString(input[lastEnd:])
	return builder.String()
}

// generateSeededAlphabetString returns deterministic alpha output for a fixed seed.
func generateSeededAlphabetString(length int, upper bool, seed int64) string {
	if length <= 0 {
		return ""
	}

	characters := "abcdefghijklmnopqrstuvwxyz"
	if upper {
		characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	rng := rand.New(rand.NewSource(seed))
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = characters[rng.Intn(len(characters))]
	}

	return string(result)
}

// goFenixRandomPositiveDecimalValue mirrors Fenix_RandomPositiveDecimalValue in Lua.
func goFenixRandomPositiveDecimalValue(input GoPlaceholderInput) (string, error) {
	if len(input.ArrayIndexes) > 1 {
		return "", fmt.Errorf("Error - array index array can only have a maximum of one value. %v", input.ArrayIndexes)
	}

	arrayIndexToUse := 1
	if len(input.ArrayIndexes) == 1 {
		arrayIndexToUse = input.ArrayIndexes[0]
	}

	functionArguments, err := parseRandomPositiveFunctionArguments(input.Arguments)
	if err != nil {
		return "", err
	}

	return fenixRandomDecimalValueArrayValue(arrayIndexToUse, functionArguments, input.Entropy), nil
}

// goFenixRandomPositiveDecimalValueSum mirrors Fenix_RandomPositiveDecimalValue_Sum in Lua.
func goFenixRandomPositiveDecimalValueSum(input GoPlaceholderInput) (string, error) {
	arrayIndexes := append([]int{}, input.ArrayIndexes...)
	if len(arrayIndexes) == 0 {
		arrayIndexes = []int{1}
	}

	functionArguments, err := parseRandomPositiveFunctionArguments(input.Arguments)
	if err != nil {
		return "", err
	}

	return fenixRandomDecimalValueSumArrayValue(arrayIndexes, functionArguments, input.Entropy), nil
}

// parseRandomPositiveFunctionArguments validates the 2-or-4 integer argument contract.
func parseRandomPositiveFunctionArguments(arguments []string) (functionArguments []int, err error) {
	if len(arguments) != 2 && len(arguments) != 4 {
		switch {
		case len(arguments) > 2:
			return nil, fmt.Errorf("Error - there must be exact 2 or 4 function parameter. '%v'", arguments)
		case len(arguments) == 1:
			return nil, fmt.Errorf("Error - there must be exact 2 or 4 function parameter. '[%s]'", arguments[0])
		default:
			return nil, fmt.Errorf("Error - there must be exact 2 or 4 function parameter but it is empty.")
		}
	}

	functionArguments = make([]int, 0, len(arguments))
	for _, argument := range arguments {
		parsed, parseErr := strconv.Atoi(strings.TrimSpace(argument))
		if parseErr != nil {
			return nil, fmt.Errorf("Error - functions parameters must be of type Integer. '%v'", arguments)
		}
		functionArguments = append(functionArguments, parsed)
	}

	return functionArguments, nil
}

// fenixRandomDecimalValueArrayValue generates one deterministic value for one array index.
func fenixRandomDecimalValueArrayValue(arrayPosition int, functionArguments []int, entropy uint64) string {
	maxIntegerPartSize := functionArguments[0]
	numberOfDecimals := functionArguments[1]

	randomValue := randomizeDecimalValue(arrayPosition, maxIntegerPartSize, numberOfDecimals, entropy)
	formattedValue := formatDecimalValue(randomValue, numberOfDecimals)

	if len(functionArguments) == 2 {
		return formattedValue
	}

	return padValueWithZeros(formattedValue, functionArguments[2], functionArguments[3])
}

// fenixRandomDecimalValueSumArrayValue generates per-index values and sums/subtracts them.
func fenixRandomDecimalValueSumArrayValue(arrayPositions []int, functionArguments []int, entropy uint64) string {
	maxIntegerPartSize := functionArguments[0]
	numberOfDecimals := functionArguments[1]

	sumOfValues := 0.0
	for _, arrayPositionValue := range arrayPositions {
		arrayPositionToUse := int(math.Abs(float64(arrayPositionValue)))
		randomValue := randomizeDecimalValue(arrayPositionToUse, maxIntegerPartSize, numberOfDecimals, entropy)

		if arrayPositionValue >= 0 {
			sumOfValues += randomValue
		} else {
			sumOfValues -= randomValue
		}
	}

	formattedValue := formatDecimalValue(sumOfValues, numberOfDecimals)
	if len(functionArguments) == 2 {
		return formattedValue
	}

	return padValueWithZeros(formattedValue, functionArguments[2], functionArguments[3])
}

// randomizeDecimalValue reproduces Lua randomization shape using deterministic seeding.
func randomizeDecimalValue(arrayIndex int, maxIntegerPartSize int, numberOfDecimals int, entropy uint64) float64 {
	seed := int64(arrayIndex) + int64(entropy)
	rng := rand.New(rand.NewSource(seed))

	integerMultiplier := math.Pow10(maxIntegerPartSize)
	randomIntegerPart := rng.Float64()
	integerPart := math.Floor(integerMultiplier * randomIntegerPart)

	decimalPart := 0.0
	if numberOfDecimals > 0 {
		randomDecimalPart := rng.Float64()
		decimalMultiplier := math.Pow10(numberOfDecimals)
		decimalPart = math.Floor(decimalMultiplier * randomDecimalPart)
	}

	randomNumber := integerPart
	if numberOfDecimals > 0 {
		randomNumber += math.Pow10(-numberOfDecimals) * decimalPart
	}

	return roundToDecimalPlaces(randomNumber, numberOfDecimals)
}

// roundToDecimalPlaces applies Lua-style rounding logic.
func roundToDecimalPlaces(value float64, decimalPlaces int) float64 {
	shift := math.Pow10(decimalPlaces)
	return math.Floor(value*shift+0.5) / shift
}

// formatDecimalValue enforces an exact decimal width (or integer-only string when width is 0).
func formatDecimalValue(number float64, numberOfDecimals int) string {
	valueAsString := strconv.FormatFloat(number, 'f', -1, 64)
	dotIndex := strings.Index(valueAsString, ".")

	if numberOfDecimals > 0 && dotIndex == -1 {
		dotIndex = len(valueAsString)
		valueAsString += "."
	}

	if numberOfDecimals > 0 {
		currentDecimals := 0
		if dotIndex != -1 {
			currentDecimals = len(valueAsString) - dotIndex - 1
		}

		if currentDecimals < numberOfDecimals {
			valueAsString += strings.Repeat("0", numberOfDecimals-currentDecimals)
		}

		return valueAsString
	}

	if dotIndex != -1 {
		return valueAsString[:dotIndex]
	}

	return valueAsString
}

// padValueWithZeros left-pads integer part and right-pads fraction part when requested.
func padValueWithZeros(valueAsString string, integerSpace int, fractionSpace int) string {
	integerPart := valueAsString
	fractionPart := ""
	noFractions := true

	if matches := regexp.MustCompile(`^([0-9]+)\.([0-9]+)$`).FindStringSubmatch(valueAsString); len(matches) == 3 {
		integerPart = matches[1]
		fractionPart = matches[2]
		noFractions = false
	}

	if len(integerPart) < integerSpace {
		integerPart = strings.Repeat("0", integerSpace-len(integerPart)) + integerPart
	}

	if noFractions == false && len(fractionPart) < fractionSpace {
		fractionPart += strings.Repeat("0", fractionSpace-len(fractionPart))
	}

	if noFractions == false {
		return integerPart + "." + fractionPart
	}

	return integerPart
}

// pow10Int64 returns 10^n where n is int64-safe for this code path.
func pow10Int64(numberOfDigits int) (int64, bool) {
	if numberOfDigits <= 0 || numberOfDigits > 18 {
		return 0, false
	}

	result := int64(1)
	for i := 0; i < numberOfDigits; i++ {
		result *= 10
	}

	return result, true
}
