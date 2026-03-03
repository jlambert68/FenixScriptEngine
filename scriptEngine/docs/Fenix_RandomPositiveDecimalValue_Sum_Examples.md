# Fenix.RandomPositiveDecimalValue.Sum - Examples

This file contains examples for `Fenix.RandomPositiveDecimalValue.Sum` based on current Go implementation and tests.

## Signature

```text
{{Fenix.RandomPositiveDecimalValue.Sum[arrayIndexes](IntegerPrecision, FractionPrecision, IntegerFieldWidth, FractionFieldWidth, DecimalPointCharacter)}}
```

Rules:

- One or more array indexes are supported.
- Positive index adds value, negative index subtracts value.
- If array index list is omitted, default is `[1]`.
- Exactly five function arguments are required.
- First four arguments must be integers.
- `DecimalPointCharacter` must be exactly one character.

## Valid Examples

```text
{{Fenix.RandomPositiveDecimalValue.Sum[1](2, 3, 2, 3, ".")}}
{{Fenix.RandomPositiveDecimalValue.Sum[-1,2](2, 3, 3, 3, ".")}}
{{Fenix.RandomPositiveDecimalValue.Sum[1,2,3](2, 3, 4, 4, ".")}}
{{Fenix.RandomPositiveDecimalValue.Sum[1,2](2, 3, 4, 4, ",")}}
```

Output shape examples (from unit tests):

```text
[1](2,3,2,3,".")     -> ^\d{2}\.\d{3}$ ; example: 44.560
[-1,2](2,3,3,3,".")  -> ^-?\d{1,3}\.\d{3}$ ; example: -044.613
[1,-2](2,3,3,3,".")  -> ^-?\d{1,3}\.\d{3}$ ; example: 044.613
[1,2,3](2,3,4,4,".") -> ^\d{4}\.\d{4}$ ; example: 0140.2380
[1,2](2,3,4,4,",")   -> ^\d{4},\d{4}$ ; example: 0087,5400
```

Deterministic regression example in tests:

```text
Indexes: [-1,2]
Args: (2,3,3,3,".")
Entropy: crc32(testCaseExecutionUUID)
Expected output: -044.613
```

## Invalid Examples

```text
{{Fenix.RandomPositiveDecimalValue.Sum[1](2, 3, 4, 4)}}            // wrong argument count
{{Fenix.RandomPositiveDecimalValue.Sum[1](2, three, 4, 4, ".")}}   // non-integer argument
{{Fenix.RandomPositiveDecimalValue.Sum[1](2, 3, 4, 4, "")}}        // empty decimal point character
{{Fenix.RandomPositiveDecimalValue.Sum[1](2, 3, 4, 4, "..")}}      // decimal point must be one char
```
