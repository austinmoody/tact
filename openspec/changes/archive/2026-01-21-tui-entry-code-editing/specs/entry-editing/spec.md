## ADDED Requirements

### Requirement: User can edit entry time code
The entry detail modal SHALL allow users to select a different time_code for a time entry.

#### Scenario: Time code dropdown in edit mode
- **WHEN** user is in edit mode
- **THEN** the time_code field SHALL display a dropdown of available time codes
- **AND** each option SHALL show the code ID and description

#### Scenario: Navigate time code options
- **WHEN** time_code field is focused in edit mode
- **THEN** user can press `j` or `↓` to select next option
- **AND** user can press `k` or `↑` to select previous option

#### Scenario: Time code options loaded on edit
- **WHEN** user enters edit mode
- **THEN** available time codes SHALL be fetched from the API
- **AND** a loading indicator SHALL display until fetch completes

### Requirement: User can edit entry work type
The entry detail modal SHALL allow users to select a different work_type for a time entry.

#### Scenario: Work type dropdown in edit mode
- **WHEN** user is in edit mode
- **THEN** the work_type field SHALL display a dropdown of available work types
- **AND** each option SHALL show the work type name

#### Scenario: Navigate work type options
- **WHEN** work_type field is focused in edit mode
- **THEN** user can press `j` or `↓` to select next option
- **AND** user can press `k` or `↑` to select previous option

### Requirement: Learn flag triggers on code corrections
The system SHALL use the learn flag to improve AI parsing when users correct time_code or work_type assignments.

#### Scenario: Code change triggers learn
- **WHEN** user changes time_code or work_type and saves
- **THEN** the system SHALL call PATCH /entries/{id}?learn=true

#### Scenario: Text-only changes do not trigger learn
- **WHEN** user changes only user_input or entry_date and saves
- **THEN** the system SHALL call PATCH /entries/{id}?learn=false

#### Scenario: Mixed changes trigger learn
- **WHEN** user changes user_input AND time_code and saves
- **THEN** the system SHALL call PATCH /entries/{id}?learn=true

## MODIFIED Requirements

### Requirement: Save entry changes
The entry detail modal SHALL save changes via the API.

#### Scenario: Save changes successfully
- **WHEN** user presses `Enter` in edit mode with valid changes
- **THEN** the system SHALL call PATCH /entries/{id} with updated fields
- **AND** the learn flag SHALL be set based on which fields changed
- **AND** the modal SHALL show updated entry details
- **AND** the home screen entry list SHALL refresh

Note: learn=true when time_code or work_type changed (AI parsing corrections). learn=false when only user_input or entry_date changed.

#### Scenario: API error displayed
- **WHEN** the API returns an error during save
- **THEN** the modal SHALL display the error message
- **AND** remain in edit mode so user can retry

#### Scenario: Date change moves entry to correct group
- **WHEN** user changes entry_date from "2026-01-20" to "2026-01-18" and saves
- **THEN** the entry SHALL appear under the "Saturday - Jan 18, 2026" date header on home screen

### Requirement: User can edit entry user input
The entry detail modal SHALL allow users to edit the `user_input` field of a time entry.

#### Scenario: Enter edit mode from detail view
- **WHEN** user presses `e` key while viewing entry details
- **THEN** the modal SHALL switch to edit mode with editable fields
- **AND** available time codes and work types SHALL be fetched

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
- **THEN** focus SHALL cycle through user_input, entry_date, time_code, and work_type fields

#### Scenario: Invalid date format rejected
- **WHEN** user enters an invalid date format and attempts to save
- **THEN** the modal SHALL display an error message and remain in edit mode
