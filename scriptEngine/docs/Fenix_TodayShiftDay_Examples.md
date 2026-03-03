# Fenix.TodayShiftDay - Examples

This file contains examples for `Fenix.TodayShiftDay` based on current Go implementation.

## Signature

```text
{{Fenix.TodayShiftDay(shiftDays)}}
```

Rules:

- Exactly one integer argument is required.
- Array indexes are not supported.
- Output format is `YYYY-MM-DD` (local time).

## Valid Examples

```text
{{Fenix.TodayShiftDay(0)}}
{{Fenix.TodayShiftDay(-1)}}
{{Fenix.TodayShiftDay(1)}}
{{Fenix.TodayShiftDay(10)}}
```

Example outputs when current date is `2026-02-26`:

```text
{{Fenix.TodayShiftDay(0)}}  -> 2026-02-26
{{Fenix.TodayShiftDay(-1)}} -> 2026-02-25
{{Fenix.TodayShiftDay(1)}}  -> 2026-02-27
{{Fenix.TodayShiftDay(10)}} -> 2026-03-08
```

## Invalid Examples

```text
{{Fenix.TodayShiftDay()}}         // missing argument
{{Fenix.TodayShiftDay(1,2)}}      // too many arguments
{{Fenix.TodayShiftDay(abc)}}      // non-integer argument
{{Fenix.TodayShiftDay[1](0)}}     // array index not supported
```

