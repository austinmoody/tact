# tui-dashboard Specification

## Purpose
TBD - created by archiving change add-tui-dashboard. Update Purpose after archive.
## Requirements
### Requirement: Dashboard Display

The TUI SHALL display recent entries on the home screen with quick access to entry creation.

#### Scenario: Show recent entries

- Given: The TUI is running and connected to the backend
- When: Entries are fetched successfully
- Then: The last 15 entries are displayed grouped by date
- And: Each entry shows raw text and status
- And: Date headers separate entries by day

#### Scenario: Empty entries list

- Given: The TUI is running and no entries exist
- When: The home screen is displayed
- Then: A message indicates no entries exist
- And: The new entry keybinding is still available

#### Scenario: Entry status display

- Given: Entries with various statuses exist
- When: The home screen is displayed
- Then: Each entry shows its status (pending, parsed, failed)
- And: Status is color-coded for visibility

### Requirement: Keyboard Navigation

The TUI SHALL support vim-style keyboard navigation across all screens.

#### Scenario: Navigate entries list

- Given: The home screen is displayed with entries
- When: j or down arrow is pressed
- Then: The cursor moves down one entry

#### Scenario: Navigate up in list

- Given: The cursor is not on the first item
- When: k or up arrow is pressed
- Then: The cursor moves up one item

#### Scenario: Open menu

- Given: The home screen is displayed
- When: m is pressed
- Then: The menu modal opens

#### Scenario: Navigate menu

- Given: The menu modal is open
- When: j/k or arrow keys are pressed
- Then: The cursor moves between menu options

#### Scenario: Close modal

- Given: A modal is open
- When: Esc is pressed
- Then: The modal closes and focus returns to previous screen

#### Scenario: Quit application

- Given: The TUI is running
- When: q or Ctrl+C is pressed
- Then: The application exits

### Requirement: API Configuration

The TUI SHALL support configurable backend URL.

#### Scenario: Use command line flag

- Given: The TUI is started with `--api http://localhost:2100`
- When: API requests are made
- Then: Requests are sent to the specified URL

#### Scenario: Use environment variable

- Given: TACT_API_URL is set to `http://localhost:2100`
- When: The TUI is started without --api flag
- Then: Requests are sent to the URL from environment variable

#### Scenario: Use default URL

- Given: No --api flag and no TACT_API_URL set
- When: The TUI is started
- Then: Requests are sent to `http://localhost:2100`

### Requirement: Data Refresh

The TUI SHALL support refreshing data from the backend.

#### Scenario: Manual refresh

- Given: The TUI is displaying data
- When: r is pressed
- Then: Data is re-fetched from the backend and display is updated

### Requirement: Error Handling

The TUI SHALL handle connection errors gracefully.

#### Scenario: Connection failure

- Given: The backend is not reachable
- When: Data fetch is attempted
- Then: An error message is displayed in the status bar
- And: The application does not crash

#### Scenario: Retry after error

- Given: An error is being displayed
- When: r is pressed
- Then: A new connection attempt is made

### Requirement: Entry Creation

The TUI SHALL allow creating new entries via a modal input.

#### Scenario: Open new entry modal

- Given: The home screen is displayed
- When: n is pressed
- Then: The new entry modal opens with a text input field

#### Scenario: Submit new entry

- Given: The new entry modal is open
- And: Text has been entered
- When: Enter is pressed
- Then: The entry is created via API
- And: The modal closes
- And: The entries list is refreshed

#### Scenario: Cancel new entry

- Given: The new entry modal is open
- When: Esc is pressed
- Then: The modal closes without creating an entry

### Requirement: Entry Details

The TUI SHALL display entry details in a modal view.

#### Scenario: Open entry detail modal

- Given: An entry is selected on the home screen
- When: Enter is pressed
- Then: A detail modal shows all entry fields

#### Scenario: Entry detail contents

- Given: The entry detail modal is open
- Then: The following fields are displayed: raw_text, status, duration_minutes, time_code, work_type, description, confidence, entry_date

#### Scenario: Reparse entry

- Given: The entry detail modal is open
- When: p is pressed
- Then: The entry is reparsed via API
- And: A success message is shown
- And: The entry status updates to pending

### Requirement: Menu Navigation

The TUI SHALL provide a menu for accessing management screens.

#### Scenario: Menu options

- Given: The menu modal is open
- Then: Options for "Time Codes" and "Work Types" are displayed

#### Scenario: Select time codes

- Given: The menu modal is open with "Time Codes" selected
- When: Enter is pressed
- Then: The time codes management screen is displayed

#### Scenario: Select work types

- Given: The menu modal is open with "Work Types" selected
- When: Enter is pressed
- Then: The work types management screen is displayed

### Requirement: Time Codes Management

The TUI SHALL allow viewing, adding, editing, and deactivating time codes.

#### Scenario: List time codes

- Given: The time codes screen is displayed
- Then: All time codes are listed with ID, name, and active status

#### Scenario: Quick add time code

- Given: The time codes screen is displayed
- When: a is pressed
- Then: A quick-add modal opens with ID and name fields

#### Scenario: Submit quick add

- Given: The quick-add modal is open with ID and name entered
- When: Enter is pressed
- Then: The time code is created via API
- And: The list is refreshed

#### Scenario: Edit time code

- Given: A time code is selected
- When: e is pressed
- Then: An edit modal opens with name, description, keywords, and examples fields

#### Scenario: Save time code edits

- Given: The edit modal is open with changes
- When: Enter is pressed
- Then: The time code is updated via API
- And: The list is refreshed

#### Scenario: Deactivate time code

- Given: An active time code is selected
- When: d is pressed
- Then: The time code is deactivated via API
- And: The list is refreshed showing inactive status

#### Scenario: Return to home

- Given: The time codes screen is displayed
- When: Esc is pressed
- Then: The home screen is displayed

### Requirement: Work Types Management

The TUI SHALL allow viewing, adding, editing, and deactivating work types.

#### Scenario: List work types

- Given: The work types screen is displayed
- Then: All work types are listed with ID, name, and active status

#### Scenario: Quick add work type

- Given: The work types screen is displayed
- When: a is pressed
- Then: A quick-add modal opens with name field only

#### Scenario: Submit quick add work type

- Given: The quick-add modal is open with name entered
- When: Enter is pressed
- Then: The work type is created via API with auto-generated ID
- And: The list is refreshed

#### Scenario: Edit work type

- Given: A work type is selected
- When: e is pressed
- Then: An edit modal opens with name field

#### Scenario: Save work type edits

- Given: The edit modal is open with changes
- When: Enter is pressed
- Then: The work type is updated via API
- And: The list is refreshed

#### Scenario: Deactivate work type

- Given: An active work type is selected
- When: d is pressed
- Then: The work type is deactivated via API
- And: The list is refreshed showing inactive status

#### Scenario: Return to home from work types

- Given: The work types screen is displayed
- When: Esc is pressed
- Then: The home screen is displayed

