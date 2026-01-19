## 1. Backend Implementation

- [x] 1.1 Add `round_duration()` utility function in `backend/src/tact/utils/duration.py`
- [x] 1.2 Read `TACT_DURATION_ROUNDING` env var with validation (none/15/30)
- [x] 1.3 Apply rounding in parser after LLM extracts duration_minutes
- [x] 1.4 Add unit tests for rounding function (edge cases: 0, 1, 14, 15, 16, 29, 30, 31, etc.)

## 2. Documentation

- [x] 2.1 Update backend README with new environment variable
- [x] 2.2 Add example to docker-compose.yml (commented out)

## 3. Testing & Verification

- [x] 3.1 Test with TACT_DURATION_ROUNDING=none (default behavior unchanged)
- [x] 3.2 Test with TACT_DURATION_ROUNDING=15 (7m → 15m, 16m → 30m)
- [x] 3.3 Test with TACT_DURATION_ROUNDING=30 (7m → 30m, 31m → 60m)
- [x] 3.4 Verify original user_input is preserved unchanged

## Verification

1. Set `TACT_DURATION_ROUNDING=15`
2. Create entry with "7m standup"
3. After parsing, verify `duration_minutes` is 15 (not 7)
4. Verify `user_input` still shows "7m standup"
