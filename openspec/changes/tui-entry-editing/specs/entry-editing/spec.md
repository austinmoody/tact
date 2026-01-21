## ADDED Requirements

### Requirement: User can edit entry user input
The entry detail modal SHALL allow users to edit the `user_input` field of a time entry.

#### Scenario: Enter edit mode from detail view
- **WHEN** user presses `e` key while viewing entry details
- **THEN** the modal SHALL switch to edit mode with editable text fields

#### Scenario: Edit user input text
- **WHEN** user is in edit mode
- **THEN** the user_input field SHALL be editable with cursor support

#### Scenario: Cancel edit mode
- **WHEN** user presses `Esc` while in edit mode
- **THEN** the modal SHALL return to view mode without saving changes

### Requirement: User can edit entry date
The entry detail modal SHALL allow users to change the `entry_date` of a time entry.

#### Scenario: Edit entry date field
- **WHEN** user is in edit mode
- **THEN** the entry_date field SHALL be editable in YYYY-MM-DD format

#### Scenario: Tab between fields
- **WHEN** user presses `Tab` in edit mode
- **THEN** focus SHALL move between user_input and entry_date fields

#### Scenario: Invalid date format rejected
- **WHEN** user enters an invalid date format and attempts to save
- **THEN** the modal SHALL display an error message and remain in edit mode

### Requirement: Save entry changes with learn flag
The entry detail modal SHALL save changes via the API with the learn flag enabled.

#### Scenario: Save changes successfully
- **WHEN** user presses `Enter` in edit mode with valid changes
- **THEN** the system SHALL call PATCH /entries/{id}?learn=true with updated fields
- **AND** the modal SHALL show updated entry details
- **AND** the home screen entry list SHALL refresh

#### Scenario: API error displayed
- **WHEN** the API returns an error during save
- **THEN** the modal SHALL display the error message
- **AND** remain in edit mode so user can retry

#### Scenario: Date change moves entry to correct group
- **WHEN** user changes entry_date from "2026-01-20" to "2026-01-18" and saves
- **THEN** the entry SHALL appear under the "Saturday - Jan 18, 2026" date header on home screen
