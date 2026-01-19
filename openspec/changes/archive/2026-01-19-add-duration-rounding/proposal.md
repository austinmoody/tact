# Add Duration Rounding

## Problem

Many companies require time entries to be logged in fixed increments (e.g., 15 or 30 minutes). Currently, Tact stores exact durations parsed from user input, requiring manual adjustment before submitting to external time tracking systems.

## Solution

Add configurable duration rounding in the backend API. After the LLM parses a duration from user input, the backend rounds UP to the configured increment before storing.

## Configuration

Environment variable: `TACT_DURATION_ROUNDING`

Values:
- `none` (default) - No rounding, store exact parsed minutes
- `15` - Round up to nearest 15 minutes (7m → 15m, 16m → 30m)
- `30` - Round up to nearest 30 minutes (7m → 30m, 31m → 60m)

## Scope

- Applies to ALL time entries processed by the parser
- Rounding happens after LLM extraction, before database storage
- Original `user_input` preserved unchanged
- Only `duration_minutes` field is affected

## Out of Scope

- Per-entry rounding overrides
- UI configuration (admin panel)
- macOS app displaying rounded preview (could be future enhancement)
