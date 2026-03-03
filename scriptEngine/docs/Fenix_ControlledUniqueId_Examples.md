# Fenix.ControlledUniqueId - Runtime Behavior

This document describes the current Go implementation in `go_placeholder_fenix_controlled_unique_id.go`.

## Signature

`Fenix.ControlledUniqueId` requires exactly three function arguments:

1. `textToProcess`
2. `useEntropyFromExecutionUUID` (`true`/`false`)
3. `extraEntropy` (integer)

Optional array index:

- No index -> default index `1`
- One index is allowed
- More than one index returns validation error

## Deterministic Entropy

Entropy is calculated inside this function from argument 2 and 3:

```text
entropy = extraEntropy
if useEntropyFromExecutionUUID == true {
  entropy = crc32(testCaseExecutionUUID) + extraEntropy
}
seedBase = arrayIndex + entropy
```

## Supported Token Replacements

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

Component replacements inside mixed strings:

- `YYYY`, `YY`, `MM`, `DD`, `hh`, `mm`, `ss`, `ms`, `us`, `ns`

Random Jira tokens:

- `%n(length)%`
- `%a(length)%`
- `%A(length)%`
- `%aA(length)%`
- `%an(length)%`
- `%An(length)%`
- `%aAn(length)%`

Legacy random formats are not replaced and remain unchanged.

## Parser-Safe Examples

These examples avoid commas inside the first argument:

```text
{{Fenix.ControlledUniqueId(%YYYY-MM-DD%, true, 0)}}
{{Fenix.ControlledUniqueId(%n(5)%-%a(5)%-%A(5)%, true, 5)}}
{{Fenix.ControlledUniqueId(Year=YYYY-Month=MM-Day=DD, false, 1)}}
{{Fenix.ControlledUniqueId[2](ID-%aAn(4)%, true, 0)}}
```

## Validation Errors

Common failures:

- More than one array index.
- Missing arguments (`len != 3`).
- Argument 2 is not boolean.
- Argument 3 is not integer.

## Unit Test Coverage

`go_placeholder_fenix_controlled_unique_id_test.go` verifies:

- Date/time token replacement.
- Jira random token patterns.
- Determinism for repeated calls.
- Entropy effect differences (`true/false`, `extraEntropy` changes).
- Unsupported legacy token behavior.
- Input validation and error messages.

All test calls log input matrix and execution result via:

- `logPlaceholderInputMatrix(...)`
- `logPlaceholderExecutionResult(...)`
