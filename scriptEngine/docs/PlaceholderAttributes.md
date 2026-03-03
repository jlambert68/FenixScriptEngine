# Placeholder Attributes

This file lists accepted user parameters and validation behavior for each supported placeholder.

## Function Contracts

| Function | Array Index Part `[ ... ]` | Function Arguments `( ... )` | Validation Summary |
|---|---|---|---|
| `Fenix.TodayShiftDay` | Not allowed | Exactly one integer: `(shiftDays)` | Fails when array index exists, when arg count != 1, or argument is not integer |
| `Fenix.ControlledUniqueId` | Optional single integer index; default `1` | Exactly three args: `(textToProcess, useEntropyFromExecutionUUID, extraEntropy)` | Fails when more than one array index, wrong arg count, invalid boolean, invalid integer |
| `Fenix.RandomPositiveDecimalValue` | Optional single integer index; default `1` | Exactly five args: `(IntegerPrecision, FractionPrecision, IntegerFieldWidth, FractionFieldWidth, DecimalPointCharacter)` | Fails when more than one index, arg count != 5, non-integer among first four args, empty/multi-char decimal point |
| `Fenix.RandomPositiveDecimalValue.Sum` | One or more integers; negatives subtract; default `[1]` | Exactly same five args as value variant | Fails on invalid argument count/type/decimal-point character |
| `HappyLuaTime` | Not allowed | No arguments: `()` | Fails when array index or arguments are provided |

## ControlledUniqueId Token Set

Date/time tokens:

- `%YYYY-MM-DD%`
- `%YYYYMMDD%`
- `%YYMMDD%`
- `%hh:mm:ss%`
- `%hh.mm.ss%`
- `%hhmmss%`
- `%hhmm%`
- `%mmss%`
- `%ms%`
- `%us%`
- `%ns%`

Random Jira tokens:

- `%n(length)%`
- `%a(length)%`
- `%A(length)%`
- `%aA(length)%`
- `%an(length)%`
- `%An(length)%`
- `%aAn(length)%`

Legacy random token forms are treated as unsupported and remain unchanged.

## Entropy Inputs

Two entropy paths exist in current implementation:

1. Dispatcher entropy (`GoPlaceholderInput.Entropy`)
- Used by random decimal placeholders.
- Derived from optional trailing tail `(useEntropy, extraEntropy)`.

2. ControlledUniqueId internal entropy
- ControlledUniqueId computes entropy from its own function args (`useEntropyFromExecutionUUID`, `extraEntropy`) plus execution UUID.

Dispatcher entropy formula:

```text
entropy = extraEntropy
if useEntropyFromExecutionUUID == true {
  entropy = crc32(testCaseExecutionUUID) + extraEntropy
}
```

## Parser Constraints

From `placeholderReplacementEngine.match(...)`:

- Function arguments are split on commas.
- No escape/quote handling for commas inside a single argument.
- Template examples that include commas inside one argument are not parser-safe.

## Example Calls

```text
{{Fenix.TodayShiftDay(0)}}
{{Fenix.ControlledUniqueId(%YYYYMMDD%, true, 1)}}
{{Fenix.ControlledUniqueId[2](ID-%n(4)%-%A(4)%, false, 1)}}
{{Fenix.RandomPositiveDecimalValue(2, 3, 2, 3, ".")}}
{{Fenix.RandomPositiveDecimalValue.Sum[-1,2](2, 3, 3, 3, ".")}}
{{HappyLuaTime()}}
```

## TestData Placeholder

`ParseAndFormatPlaceholders(...)` recognizes:

- `TestData.Context.Column`
- `Context.TestData.Column` (legacy)

Malformed `TestData` references are returned as a readable error string in output text.
