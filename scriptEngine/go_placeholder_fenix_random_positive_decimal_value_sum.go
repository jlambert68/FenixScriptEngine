package scriptEngine

// goFenixRandomPositiveDecimalValueSum mirrors Fenix_RandomPositiveDecimalValue_Sum in Lua.
func goFenixRandomPositiveDecimalValueSum(input GoPlaceholderInput) (string, error) {
	arrayIndexes := append([]int{}, input.ArrayIndexes...)
	if len(arrayIndexes) == 0 {
		arrayIndexes = []int{1}
	}

	functionArguments, decimalPointCharacter, err := parseRandomPositiveFunctionArguments(input.Arguments)
	if err != nil {
		return "", err
	}

	return fenixRandomDecimalValueSumArrayValue(arrayIndexes, functionArguments, decimalPointCharacter, input.Entropy), nil
}
