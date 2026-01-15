# entries-api Delta

## MODIFIED Requirements

### Requirement: Update Entry

The API SHALL allow partial updates to an entry via PATCH request. Updates mark the entry as `manually_corrected=true`. When the `learn` query parameter is true (default) and the entry has a `time_code_id`, a context document is created for the time code to improve future parsing.

#### Scenario: Successful update with learning enabled

- Given: An entry exists with raw_text "2h standup meeting"
- When: PATCH `/entries/{id}` is called with `time_code_id="PROJ-DEV"` and `work_type_id="meetings"`
- Then: The entry is updated with the specified fields
- And: `manually_corrected` is set to `true`
- And: A context document is created for time code "PROJ-DEV"
- And: HTTP 200 is returned with the updated entry

#### Scenario: Update with learning disabled

- Given: An entry exists
- When: PATCH `/entries/{id}?learn=false` is called with updated fields
- Then: The entry is updated
- And: `manually_corrected` is set to `true`
- And: No context document is created
- And: HTTP 200 is returned

#### Scenario: Update without time_code_id does not learn

- Given: An entry exists without a time_code_id
- When: PATCH `/entries/{id}` is called with only `duration_minutes=60`
- Then: The entry is updated
- And: No context document is created (no time code to associate with)
- And: HTTP 200 is returned

#### Scenario: Update entry_date

- Given: An entry exists with entry_date "2026-01-01"
- When: PATCH `/entries/{id}` is called with `entry_date="2026-01-15"`
- Then: The entry_date is updated to "2026-01-15"

#### Scenario: Update non-existent entry

- Given: No entry with ID "unknown-uuid" exists
- When: PATCH `/entries/unknown-uuid` is called
- Then: HTTP 404 is returned
