## ADDED Requirements

### Requirement: Duration Rounding

The parser SHALL round parsed durations up to a configurable increment.

#### Scenario: No rounding (default)

- Given: `TACT_DURATION_ROUNDING` is not set or set to `none`
- And: Raw text "7m standup"
- When: The entry is parsed
- Then: duration_minutes is 7

#### Scenario: Round to 15 minutes

- Given: `TACT_DURATION_ROUNDING=15`
- And: Raw text "7m standup"
- When: The entry is parsed
- Then: duration_minutes is 15 (rounded up from 7)

#### Scenario: Round to 15 minutes at boundary

- Given: `TACT_DURATION_ROUNDING=15`
- And: Raw text "15m standup"
- When: The entry is parsed
- Then: duration_minutes is 15 (already at boundary)

#### Scenario: Round to 15 minutes just over boundary

- Given: `TACT_DURATION_ROUNDING=15`
- And: Raw text "16m standup"
- When: The entry is parsed
- Then: duration_minutes is 30 (rounded up from 16)

#### Scenario: Round to 30 minutes

- Given: `TACT_DURATION_ROUNDING=30`
- And: Raw text "7m standup"
- When: The entry is parsed
- Then: duration_minutes is 30 (rounded up from 7)

#### Scenario: Round to 30 minutes just over boundary

- Given: `TACT_DURATION_ROUNDING=30`
- And: Raw text "31m planning"
- When: The entry is parsed
- Then: duration_minutes is 60 (rounded up from 31)

#### Scenario: Preserve original user input

- Given: `TACT_DURATION_ROUNDING=15`
- And: Raw text "7m standup"
- When: The entry is parsed
- Then: user_input remains "7m standup" unchanged
- And: duration_minutes is 15

#### Scenario: Zero duration not rounded

- Given: `TACT_DURATION_ROUNDING=15`
- And: LLM could not extract a duration (returns null)
- When: The entry is parsed
- Then: duration_minutes remains null

## MODIFIED Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `TACT_DURATION_ROUNDING` | `none` | Duration rounding: `none`, `15`, or `30` minutes |
