## ADDED Requirements

### Requirement: Entry List Display

The web UI SHALL display a list of recent time entries on the home page.

#### Scenario: Show recent entries

- **WHEN** the home page is loaded
- **THEN** the most recent 15 entries SHALL be displayed
- **AND** entries SHALL be grouped by date with date headers

#### Scenario: Entry row display

- **WHEN** an entry is displayed in the list
- **THEN** the following SHALL be visible: user input text, time code (ID and name), status
- **AND** status SHALL be color-coded (parsed=green, pending=yellow, failed=red)

#### Scenario: Empty state

- **WHEN** no entries exist
- **THEN** a message SHALL indicate no entries exist
- **AND** a prominent "Create Entry" button SHALL be displayed

### Requirement: Entry Creation

The web UI SHALL allow creating new time entries.

#### Scenario: Create entry form

- **WHEN** user clicks "New Entry" button
- **THEN** an entry form SHALL be displayed
- **AND** the form SHALL have a text input field for the entry

#### Scenario: Submit new entry

- **WHEN** user submits the entry form with text
- **THEN** the entry SHALL be created via the backend API
- **AND** the entries list SHALL refresh to show the new entry
- **AND** the form SHALL close

#### Scenario: Cancel entry creation

- **WHEN** user cancels the entry form
- **THEN** the form SHALL close without creating an entry

### Requirement: Entry Detail View

The web UI SHALL display detailed information for a selected entry.

#### Scenario: View entry details

- **WHEN** user clicks on an entry in the list
- **THEN** a detail view SHALL be displayed showing all entry fields
- **AND** fields SHALL include: raw_text, status, duration_minutes, time_code (with name), work_type, description, confidence, parse_notes, entry_date

#### Scenario: Close detail view

- **WHEN** user closes the detail view
- **THEN** focus SHALL return to the entries list

### Requirement: Entry Editing

The web UI SHALL allow editing existing entries.

#### Scenario: Edit entry

- **WHEN** user clicks "Edit" on an entry detail view
- **THEN** an edit form SHALL be displayed with current values
- **AND** editable fields SHALL include: user_input, duration_minutes, time_code_id, work_type_id, description, entry_date

#### Scenario: Save entry edits

- **WHEN** user saves the edit form
- **THEN** the entry SHALL be updated via the backend API
- **AND** the detail view SHALL refresh with updated data

### Requirement: Entry Reparse

The web UI SHALL allow triggering a reparse of an entry.

#### Scenario: Reparse entry

- **WHEN** user clicks "Reparse" on an entry
- **THEN** the entry SHALL be sent for reparsing via the backend API
- **AND** the status SHALL update to "pending"
- **AND** a success message SHALL be displayed

### Requirement: Entry Filtering

The web UI SHALL allow filtering the entries list.

#### Scenario: Filter by status

- **WHEN** user selects a status filter (all, parsed, pending, failed)
- **THEN** only entries matching that status SHALL be displayed

#### Scenario: Filter by date range

- **WHEN** user selects a date range
- **THEN** only entries within that range SHALL be displayed
