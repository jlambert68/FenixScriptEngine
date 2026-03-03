# Fenix_ControlledUniqueId - Supported Replacements and Processing

This file summarizes current Jira-style token support for `Fenix.ControlledUniqueId`.

## Function Signature

`Fenix.ControlledUniqueId` expects exactly three function arguments:

1. `textToProcess` (string)
2. `useEntropyFromExecutionUUID` (`true`/`false`)
3. `extraEntropy` (integer)

Examples:

```text
{{Fenix.ControlledUniqueId(%YYYY-MM-DD%, true, 0)}}
{{Fenix.ControlledUniqueId(%n(5)%-%a(5)%-%A(5)%, true, 5)}}
{{Fenix.ControlledUniqueId(%Year: YYYY, Month: MM, Day: DD%, false, 1)}}
```

## Supported Replacements

Date/time token replacements:

| Token | Replacement format |
|---|---|
| `%YYYY-MM-DD%` | Current date (`YYYY-MM-DD`) |
| `%YYYYMMDD%` | Current date (`YYYYMMDD`) |
| `%YYMMDD%` | Current date (`YYMMDD`) |
| `%hh:mm:ss%` | Current time (`HH:MM:SS`) |
| `%hh.mm.ss%` | Current time (`HH.MM.SS`) |
| `%hhmmss%` | Current time (`HHMMSS`) |
| `%hhmm%` | Current time (`HHMM`) |
| `%mmss%` | Current time (`MMSS`) |
| `%ms%` | Milliseconds (`000-999`) |
| `%us%` | Microseconds (`000000-999999`) |
| `%ns%` | Nanoseconds (`000000000-999999999`) |

Random Jira token replacements:

| Token | Character set |
|---|---|
| `%n(length)%` | digits (`0-9`) |
| `%a(length)%` | lowercase letters (`a-z`) |
| `%A(length)%` | uppercase letters (`A-Z`) |
| `%aA(length)%` | mixed letters (`a-zA-Z`) |
| `%an(length)%` | lowercase alphanumeric (`a-z0-9`) |
| `%An(length)%` | uppercase alphanumeric (`A-Z0-9`) |
| `%aAn(length)%` | mixed alphanumeric (`a-zA-Z0-9`) |

## Processing Behavior

- If array index is omitted, default index `1` is used.
- At most one array index is allowed.
- Output is deterministic for the same `textToProcess`, array index, execution UUID, and entropy values.

## Full Example

Input:

```text
{{Fenix.ControlledUniqueId(Date=%YYYY-MM-DD%, Compact=%hhmmss%, Rand=%n(5)%, Mix=%aAn(4)%, true, 0)}}
```

Possible output shape:

```text
Date=2026-03-02, Compact=153045, Rand=79410, Mix=g7Qx
```

## Notes

- Jira token formats are the supported random formats.
- Legacy non-Jira random formats are not supported.
