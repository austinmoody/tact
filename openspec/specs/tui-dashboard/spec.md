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
- Then: Options for "Projects", "Time Codes", and "Work Types" are displayed

#### Scenario: Select projects

- Given: The menu modal is open with "Projects" selected
- When: Enter is pressed
- Then: The projects management screen is displayed

### Requirement: Time Codes Management

The TUI SHALL allow viewing, adding, editing, and deactivating time codes.

#### Scenario: Full add time code form

- Given: The time codes screen is displayed
- When: a is pressed
- Then: An add modal opens with ID, name, description, keywords, and examples fields
- And: All fields are editable

#### Scenario: Submit full add form

- Given: The add modal is open with all fields populated
- When: Enter is pressed
- Then: The time code is created via API with all provided fields
- And: The list is refreshed

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

### Requirement: Text Input Handling

The TUI SHALL handle text input without inserting control characters.

#### Scenario: Paste text into input field

- Given: A text input field is focused
- When: Text is pasted (Ctrl+V or terminal paste)
- Then: The text appears in the input field
- And: No control characters are inserted

#### Scenario: Arrow key navigation in text field

- Given: A text input field is focused with text entered
- When: Left or right arrow key is pressed
- Then: The cursor moves within the text
- And: No control characters are inserted into the input

### Requirement: Responsive Form Layout

The TUI SHALL adapt form widths to the terminal size.

#### Scenario: Wide terminal

- Given: The terminal width is greater than 100 characters
- When: A modal form is displayed
- Then: Input fields use a width up to 80 characters

#### Scenario: Narrow terminal

- Given: The terminal width is less than 60 characters
- When: A modal form is displayed
- Then: Input fields use a minimum width of 30 characters

#### Scenario: Medium terminal

- Given: The terminal width is between 60 and 100 characters
- When: A modal form is displayed
- Then: Input fields scale proportionally to available space

### Requirement: Projects Management

The TUI SHALL allow viewing, adding, editing, and deactivating projects.

#### Scenario: List projects

- Given: The projects screen is displayed
- Then: All projects are listed with ID, name, and active status

#### Scenario: Add project

- Given: The projects screen is displayed
- When: a is pressed
- Then: An add modal opens with ID, name, and description fields

#### Scenario: Submit add project

- Given: The add modal is open with fields populated
- When: Enter is pressed
- Then: The project is created via API
- And: The list is refreshed

#### Scenario: Edit project

- Given: A project is selected
- When: e is pressed
- Then: An edit modal opens with name and description fields

#### Scenario: Save project edits

- Given: The edit modal is open with changes
- When: Enter is pressed
- Then: The project is updated via API
- And: The list is refreshed

#### Scenario: Deactivate project

- Given: An active project is selected
- When: d is pressed
- Then: The project is deactivated via API
- And: The list is refreshed showing inactive status

#### Scenario: Return to home from projects

- Given: The projects screen is displayed
- When: Esc is pressed
- Then: The home screen is displayed

### Requirement: Context Management

The TUI SHALL allow viewing and managing context documents for projects and time codes.

#### Scenario: Open context for project

- Given: The projects screen is displayed with a project selected
- When: c is pressed
- Then: A context list modal opens showing context documents for that project

#### Scenario: Open context for time code

- Given: The time codes screen is displayed with a time code selected
- When: c is pressed
- Then: A context list modal opens showing context documents for that time code

#### Scenario: Context list display

- Given: The context list modal is open
- Then: Each context document shows truncated content (first 50 characters)
- And: The full content is visible when editing

#### Scenario: Add context document

- Given: The context list modal is open
- When: a is pressed
- Then: A context edit modal opens with a content field

#### Scenario: Submit context document

- Given: The context edit modal is open with content entered
- When: Ctrl+S is pressed
- Then: The context document is created via API
- And: The context list is refreshed

#### Scenario: Multi-line context editing

- Given: The context edit modal is open
- When: Enter is pressed
- Then: A newline is inserted in the textarea
- And: The cursor moves to the new line

#### Scenario: Edit context document

- Given: The context list modal is open with a document selected
- When: e is pressed
- Then: A context edit modal opens with the content pre-filled

#### Scenario: Delete context document

- Given: The context list modal is open with a document selected
- When: d is pressed
- Then: The context document is deleted via API
- And: The context list is refreshed

#### Scenario: Close context list

- Given: The context list modal is open
- When: Esc is pressed
- Then: The modal closes and focus returns to the parent screen

