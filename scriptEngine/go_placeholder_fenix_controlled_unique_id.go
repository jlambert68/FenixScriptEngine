package scriptEngine

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	controlledUniqueRandomNumberPattern  = regexp.MustCompile(`%n\((\d+)\)%`)
	controlledUniqueRandomLowerPattern   = regexp.MustCompile(`%a\((\d+)\)%`)
	controlledUniqueRandomUpperPattern   = regexp.MustCompile(`%A\((\d+)\)%`)
	controlledUniqueRandomAaPattern      = regexp.MustCompile(`%aA\((\d+)\)%`)
	controlledUniqueRandomAnPattern      = regexp.MustCompile(`%an\((\d+)\)%`)
	controlledUniqueRandomAnUpperPattern = regexp.MustCompile(`%An\((\d+)\)%`)
	controlledUniqueRandomAaNPattern     = regexp.MustCompile(`%aAn\((\d+)\)%`)
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
	if len(input.Arguments) != 3 {
		return "", fmt.Errorf("Error - there must be exact 3 function arguments. arguments: %v", input.Arguments)
	}

	textToProcess := input.Arguments[0]
	useEntropyFromExecutionUUID, err := strconv.ParseBool(strings.TrimSpace(input.Arguments[1]))
	if err != nil {
		return "", fmt.Errorf("Error - second function argument must be a Boolean. arguments: %v", input.Arguments)
	}
	extraEntropy, err := strconv.ParseUint(strings.TrimSpace(input.Arguments[2]), 10, 64)
	if err != nil {
		return "", fmt.Errorf("Error - third function argument must be an Integer. arguments: %v", input.Arguments)
	}

	entropyToUse := extraEntropy
	if useEntropyFromExecutionUUID == true {
		entropyToUse = uint64(crc32.ChecksumIEEE([]byte(input.TestCaseExecutionUUID))) + extraEntropy
	}

	now := currentTimeProvider().In(time.Local)

	result := textToProcess
	result = strings.ReplaceAll(result, "%YYYY-MM-DD%", now.Format("2006-01-02"))
	result = strings.ReplaceAll(result, "%YYYYMMDD%", now.Format("20060102"))
	result = strings.ReplaceAll(result, "%YYMMDD%", now.Format("060102"))
	result = strings.ReplaceAll(result, "%hh:mm:ss%", now.Format("15:04:05"))
	result = strings.ReplaceAll(result, "%hh.mm.ss%", now.Format("15.04.05"))
	result = strings.ReplaceAll(result, "%hhmmss%", now.Format("150405"))
	result = strings.ReplaceAll(result, "%hhmm%", now.Format("1504"))
	result = strings.ReplaceAll(result, "%mmss%", now.Format("0405"))
	result = strings.ReplaceAll(result, "%ms%", fmt.Sprintf("%03d", now.Nanosecond()/1_000_000))
	result = strings.ReplaceAll(result, "%us%", fmt.Sprintf("%06d", now.Nanosecond()/1_000))
	result = strings.ReplaceAll(result, "%ns%", fmt.Sprintf("%09d", now.Nanosecond()))

	// Replace date/time components in mixed templates according to Jira token set.
	result = strings.ReplaceAll(result, "YYYY", now.Format("2006"))
	result = strings.ReplaceAll(result, "YY", now.Format("06"))
	result = strings.ReplaceAll(result, "MM", now.Format("01"))
	result = strings.ReplaceAll(result, "DD", now.Format("02"))
	result = strings.ReplaceAll(result, "hh", now.Format("15"))
	result = strings.ReplaceAll(result, "mm", now.Format("04"))
	result = strings.ReplaceAll(result, "ss", now.Format("05"))
	result = strings.ReplaceAll(result, "ms", fmt.Sprintf("%03d", now.Nanosecond()/1_000_000))
	result = strings.ReplaceAll(result, "us", fmt.Sprintf("%06d", now.Nanosecond()/1_000))
	result = strings.ReplaceAll(result, "ns", fmt.Sprintf("%09d", now.Nanosecond()))

	// Random value generation must stay deterministic for same array index + entropy.
	seed := int64(arrayPositionToUse) + int64(entropyToUse)
	result = replaceControlledUniqueRandomNumberPatterns(result, seed+11)
	result = replaceControlledUniqueRandomStringPatterns(result, controlledUniqueRandomLowerPattern, "abcdefghijklmnopqrstuvwxyz", seed+13)
	result = replaceControlledUniqueRandomStringPatterns(result, controlledUniqueRandomUpperPattern, "ABCDEFGHIJKLMNOPQRSTUVWXYZ", seed+17)
	result = replaceControlledUniqueRandomStringPatterns(result, controlledUniqueRandomAaPattern, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", seed+19)
	result = replaceControlledUniqueRandomStringPatterns(result, controlledUniqueRandomAnPattern, "abcdefghijklmnopqrstuvwxyz0123456789", seed+23)
	result = replaceControlledUniqueRandomStringPatterns(result, controlledUniqueRandomAnUpperPattern, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", seed+29)
	result = replaceControlledUniqueRandomStringPatterns(result, controlledUniqueRandomAaNPattern, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", seed+31)

	return result, nil
}

// replaceControlledUniqueRandomNumberPatterns replaces %n(length)% tokens.
func replaceControlledUniqueRandomNumberPatterns(input string, seed int64) string {
	rng := rand.New(rand.NewSource(seed))

	return controlledUniqueRandomNumberPattern.ReplaceAllStringFunc(input, func(match string) string {
		subMatches := controlledUniqueRandomNumberPattern.FindStringSubmatch(match)
		if len(subMatches) != 2 {
			return match
		}

		numberOfDigits, err := strconv.Atoi(subMatches[1])
		if err != nil || numberOfDigits <= 0 {
			return match
		}

		output := make([]byte, numberOfDigits)
		for i := 0; i < numberOfDigits; i++ {
			output[i] = byte('0' + rng.Intn(10))
		}

		return string(output)
	})
}

// replaceControlledUniqueRandomStringPatterns replaces Jira string token patterns like %a(n)%.
func replaceControlledUniqueRandomStringPatterns(input string, pattern *regexp.Regexp, characterSet string, seed int64) string {
	rng := rand.New(rand.NewSource(seed))

	return pattern.ReplaceAllStringFunc(input, func(match string) string {
		subMatches := pattern.FindStringSubmatch(match)
		if len(subMatches) != 2 {
			return match
		}

		lengthValue, err := strconv.Atoi(subMatches[1])
		if err != nil || lengthValue <= 0 {
			return match
		}

		return generateDeterministicStringFromCharset(lengthValue, characterSet, rng)
	})
}

// generateDeterministicStringFromCharset returns deterministic output for a fixed RNG stream.
func generateDeterministicStringFromCharset(length int, characterSet string, rng *rand.Rand) string {
	if length <= 0 {
		return ""
	}
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = characterSet[rng.Intn(len(characterSet))]
	}

	return string(result)
}
