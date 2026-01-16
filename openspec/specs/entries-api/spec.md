# entries-api Specification

## Purpose
TBD - created by archiving change add-entries-api. Update Purpose after archive.
## Requirements
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

### Requirement: List Entries

The API SHALL return entries with optional filtering, pagination, and parse notes included.

#### Scenario: List all entries

- Given: Multiple entries exist
- When: GET `/entries` is called
- Then: All entries are returned

#### Scenario: List includes parse notes

- Given: Multiple parsed entries exist
- When: GET `/entries` is called
- Then: Each entry in the response includes `parse_notes` field

#### Scenario: Filter by status

- Given: Entries with various statuses exist
- When: GET `/entries?status=pending` is called
- Then: Only entries with `status="pending"` are returned

#### Scenario: Filter by time_code_id

- Given: Entries with various time codes exist
- When: GET `/entries?time_code_id=PROJ-001` is called
- Then: Only entries with that time code are returned

#### Scenario: Filter by work_type_id

- Given: Entries with various work types exist
- When: GET `/entries?work_type_id=development` is called
- Then: Only entries with that work type are returned

#### Scenario: Filter by date range

- Given: Entries with various dates exist
- When: GET `/entries?from_date=2026-01-01&to_date=2026-01-31` is called
- Then: Only entries within that date range (inclusive) are returned

#### Scenario: Pagination

- Given: Many entries exist
- When: GET `/entries?limit=10&offset=20` is called
- Then: Up to 10 entries are returned starting from offset 20

### Requirement: Get Single Entry

The API SHALL return a single entry by ID, including parse notes.

#### Scenario: Entry exists

- Given: An entry with a specific ID exists
- When: GET `/entries/{id}` is called
- Then: The entry is returned

#### Scenario: Entry with parse notes

- Given: An entry exists that has been parsed
- When: GET `/entries/{id}` is called
- Then: The response includes `parse_notes` field
- And: `parse_notes` contains the LLM reasoning and context info from parsing

#### Scenario: Entry not found

- Given: No entry with ID "unknown-uuid" exists
- When: GET `/entries/unknown-uuid` is called
- Then: HTTP 404 is returned

### Requirement: Update Entry

The API SHALL allow partial updates to an entry via PATCH request. Updates mark the entry as `manually_corrected=true`. When the `learn` query parameter is true (default) and the entry has a `time_code_id`, a context document is created for the time code to improve future parsing.

#### Scenario: Successful update with learning enabled

- Given: An entry exists with user_input "2h standup meeting"
- When: PATCH `/entries/{id}` is called with `time_code_id="PROJ-DEV"` and `work_type_id="meetings"`
- Then: The entry is updated with the specified fields
- And: `manually_corrected` is set to `true`
- And: A context document is created for time code "PROJ-DEV"
- And: HTTP 200 is returned with the updated entry

### Requirement: Delete Entry

The API SHALL hard-delete an entry.

#### Scenario: Successful delete

- Given: An entry exists
- When: DELETE `/entries/{id}` is called
- Then: The entry is permanently removed
- And: HTTP 204 is returned

#### Scenario: Delete non-existent entry

- Given: No entry with ID "unknown-uuid" exists
- When: DELETE `/entries/unknown-uuid` is called
- Then: HTTP 404 is returned

