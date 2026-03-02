package scriptEngine

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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
	if len(input.Arguments) > 1 {
		return "", fmt.Errorf("Error - there must be exact 1 function argument. arguments: %v", input.Arguments)
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
