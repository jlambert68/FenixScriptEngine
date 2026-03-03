# ScriptEngine Docs - Overview

This folder contains placeholder documentation for the Go-first ScriptEngine implementation.

## Documentation Types

1. `Jira Epics, Stories and Subtasks.txt`
- Source backlog text from Jira epics/stories/subtasks.
- Used as requirement input, not runtime behavior source.

2. `PLACEHOLDERS.md`
- Runtime behavior documentation.
- Describes current placeholder syntax, execution path, and parser constraints.

3. `PlaceholderAttributes.md`
- Function-by-function argument reference.
- Focuses on validation rules and accepted parameter shapes.

4. `Fenix_ControlledUniqueId_Examples.md`
- Deep-dive for `Fenix.ControlledUniqueId` token support and deterministic behavior.

5. `Fenix_TodayShiftDay_Examples.md`
- Examples for `Fenix.TodayShiftDay`.

6. `Fenix_RandomPositiveDecimalValue_Examples.md`
- Examples for `Fenix.RandomPositiveDecimalValue`.

7. `Fenix_RandomPositiveDecimalValue_Sum_Examples.md`
- Examples for `Fenix.RandomPositiveDecimalValue.Sum`.

8. `UNIT_TESTS_AND_LOGGING.md`
- Test coverage map for all Go unit tests.
- Documents the input/output logging format used in tests.

## Runtime Source Of Truth

When docs and Jira text differ, current runtime behavior is defined by code and tests in:

- `scriptEngine/*.go`
- `scriptEngine/*_test.go`
- `placeholderReplacementEngine/*.go`
- `placeholderReplacementEngine/*_test.go`

## Important Notes

- Placeholders are parsed by `placeholderReplacementEngine.match(...)`.
- Function arguments are split on commas with no escaping/quoting support.
- Go handlers are executed before Lua fallback (`executeGoPlaceholderFunction(...)`).
- Function names in templates use dots (`Fenix.X`) and are normalized to underscores (`Fenix_X`) internally.
