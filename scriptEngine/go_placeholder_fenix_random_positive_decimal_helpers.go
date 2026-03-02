package scriptEngine

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// parseRandomPositiveFunctionArguments validates Jira-style parameters:
// (IntegerPrecision, FractionPrecision, IntegerFieldWidth, FractionFieldWidth, DecimalPointCharacter)
func parseRandomPositiveFunctionArguments(arguments []string) (functionArguments []int, decimalPointCharacter string, err error) {
	if len(arguments) != 5 {
		switch {
		case len(arguments) > 1:
			return nil, "", fmt.Errorf("Error - there must be exact 5 function parameter. '%v'", arguments)
		case len(arguments) == 1:
			return nil, "", fmt.Errorf("Error - there must be exact 5 function parameter. '[%s]'", arguments[0])
		default:
			return nil, "", fmt.Errorf("Error - there must be exact 5 function parameter but it is empty.")
		}
	}

	functionArguments = make([]int, 0, 4)
	for _, argument := range arguments[:4] {
		parsed, parseErr := strconv.Atoi(strings.TrimSpace(argument))
		if parseErr != nil {
			return nil, "", fmt.Errorf("Error - first four functions parameters must be of type Integer. '%v'", arguments)
		}
		functionArguments = append(functionArguments, parsed)
	}

	decimalPointCharacter = strings.TrimSpace(arguments[4])
	if decimalPointCharacter == "" {
		return nil, "", fmt.Errorf("Error - decimalPointCharacter must be provided. '%v'", arguments)
	}
	if utf8.RuneCountInString(decimalPointCharacter) != 1 {
		return nil, "", fmt.Errorf("Error - decimalPointCharacter must be a single character. '%s'", decimalPointCharacter)
	}

	return functionArguments, decimalPointCharacter, nil
}

// fenixRandomDecimalValueArrayValue generates one deterministic value for one array index.
func fenixRandomDecimalValueArrayValue(arrayPosition int, functionArguments []int, decimalPointCharacter string, entropy uint64) string {
	maxIntegerPartSize := functionArguments[0]
	numberOfDecimals := functionArguments[1]

	randomValue := randomizeDecimalValue(arrayPosition, maxIntegerPartSize, numberOfDecimals, entropy)
	formattedValue := formatDecimalValue(randomValue, numberOfDecimals)
	formattedValue = padValueWithZeros(formattedValue, functionArguments[2], functionArguments[3])

	return applyDecimalPointCharacter(formattedValue, decimalPointCharacter)
}

// fenixRandomDecimalValueSumArrayValue generates per-index values and sums/subtracts them.
func fenixRandomDecimalValueSumArrayValue(arrayPositions []int, functionArguments []int, decimalPointCharacter string, entropy uint64) string {
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
	formattedValue = padValueWithZeros(formattedValue, functionArguments[2], functionArguments[3])

	return applyDecimalPointCharacter(formattedValue, decimalPointCharacter)
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
	rounded := roundToDecimalPlaces(number, numberOfDecimals)
	if numberOfDecimals > 0 {
		return strconv.FormatFloat(rounded, 'f', numberOfDecimals, 64)
	}

	return strconv.FormatFloat(rounded, 'f', 0, 64)
}

// padValueWithZeros left-pads integer part and right-pads fraction part when requested.
func padValueWithZeros(valueAsString string, integerSpace int, fractionSpace int) string {
	sign := ""
	valueWithoutSign := valueAsString
	if strings.HasPrefix(valueAsString, "-") {
		sign = "-"
		valueWithoutSign = strings.TrimPrefix(valueAsString, "-")
	}

	integerPart := valueWithoutSign
	fractionPart := ""
	noFractions := true

	if matches := regexp.MustCompile(`^([0-9]+)\.([0-9]+)$`).FindStringSubmatch(valueWithoutSign); len(matches) == 3 {
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
		return sign + integerPart + "." + fractionPart
	}

	return sign + integerPart
}

// applyDecimalPointCharacter swaps '.' with user-selected decimal point character.
func applyDecimalPointCharacter(valueAsString string, decimalPointCharacter string) string {
	if decimalPointCharacter == "." {
		return valueAsString
	}

	return strings.Replace(valueAsString, ".", decimalPointCharacter, 1)
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
