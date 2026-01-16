# entries-api Specification Delta

## MODIFIED Requirements

### Requirement: Create Entry

The API SHALL allow creating a new time entry via POST request. Only `user_input` is required. Entry date defaults to today if not specified. Status is set to `pending`.

#### Scenario: Successful creation with user_input only

- Given: A payload with only `user_input`
- When: POST `/entries` is called
- Then: The entry is created with `status="pending"` and `entry_date` set to today
- And: HTTP 201 is returned with the created entry including generated UUID

#### Scenario: Successful creation with entry_date specified

- Given: A payload with `user_input` and `entry_date`
- When: POST `/entries` is called
- Then: The entry is created with the specified `entry_date`
- And: HTTP 201 is returned

#### Scenario: Missing user_input

- Given: A payload without `user_input`
- When: POST `/entries` is called
- Then: HTTP 422 is returned with validation error

### Requirement: Update Entry

The API SHALL allow partial updates to an entry via PATCH request. Updates mark the entry as `manually_corrected=true`. When the `learn` query parameter is true (default) and the entry has a `time_code_id`, a context document is created for the time code to improve future parsing.

#### Scenario: Successful update with learning enabled

- Given: An entry exists with user_input "2h standup meeting"
- When: PATCH `/entries/{id}` is called with `time_code_id="PROJ-DEV"` and `work_type_id="meetings"`
- Then: The entry is updated with the specified fields
- And: `manually_corrected` is set to `true`
- And: A context document is created for time code "PROJ-DEV"
- And: HTTP 200 is returned with the updated entry
