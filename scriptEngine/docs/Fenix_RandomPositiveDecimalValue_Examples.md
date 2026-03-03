# Fenix.RandomPositiveDecimalValue - Examples

This file contains examples for `Fenix.RandomPositiveDecimalValue` based on current Go implementation and tests.

## Signature

```text
{{Fenix.RandomPositiveDecimalValue[arrayIndex](IntegerPrecision, FractionPrecision, IntegerFieldWidth, FractionFieldWidth, DecimalPointCharacter)}}
```

Rules:

- Array index is optional. Default is `1`.
- At most one array index is allowed.
- Exactly five function arguments are required.
- First four arguments must be integers.
- `DecimalPointCharacter` must be exactly one character.

## Valid Examples

```text
{{Fenix.RandomPositiveDecimalValue(2, 3, 2, 3, ".")}}
{{Fenix.RandomPositiveDecimalValue[2](2, 3, 2, 3, ".")}}
{{Fenix.RandomPositiveDecimalValue(1, 2, 3, 4, ".")}}
{{Fenix.RandomPositiveDecimalValue(3, 2, 3, 2, ",")}}
```

Output shape examples (from unit tests):

```text
[0](1,2,3,4,".") -> ^\d{3}\.\d{4}$ ; example: 004.5700
[2](2,2,2,2,".") -> ^\d{2}\.\d{2}$ ; example: 28.87
[3](0,2,1,2,".") -> ^0\.\d{2}$ ; example: 0.37
[4](3,0,3,0,".") -> ^\d{3}$ ; example: 293
[6](3,2,3,2,",") -> ^\d{3},\d{2}$ ; example: 489,03
```

## Invalid Examples

```text
{{Fenix.RandomPositiveDecimalValue(0)}}                       // wrong argument count
{{Fenix.RandomPositiveDecimalValue(1, 2, 3, 4)}}             // wrong argument count
{{Fenix.RandomPositiveDecimalValue(1, two, 3, 4, ".")}}      // non-integer argument
{{Fenix.RandomPositiveDecimalValue(1, 2, 3, 4, "")}}         // empty decimal point character
{{Fenix.RandomPositiveDecimalValue(1, 2, 3, 4, "..")}}       // decimal point must be one char
{{Fenix.RandomPositiveDecimalValue[1,2](2, 3, 2, 3, ".")}}   // too many array indexes
```
