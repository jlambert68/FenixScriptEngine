package scriptEngine

import "fmt"

// goFenixRandomPositiveDecimalValue mirrors Fenix_RandomPositiveDecimalValue in Lua.
func goFenixRandomPositiveDecimalValue(input GoPlaceholderInput) (string, error) {
	if len(input.ArrayIndexes) > 1 {
		return "", fmt.Errorf("Error - array index array can only have a maximum of one value. %v", input.ArrayIndexes)
	}

	arrayIndexToUse := 1
	if len(input.ArrayIndexes) == 1 {
		arrayIndexToUse = input.ArrayIndexes[0]
	}

	functionArguments, decimalPointCharacter, err := parseRandomPositiveFunctionArguments(input.Arguments)
	if err != nil {
		return "", err
	}

	return fenixRandomDecimalValueArrayValue(arrayIndexToUse, functionArguments, decimalPointCharacter, input.Entropy), nil
}
