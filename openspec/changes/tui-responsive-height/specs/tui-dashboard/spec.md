## MODIFIED Requirements

### Requirement: Dashboard Display

The TUI SHALL display recent entries on the home screen with quick access to entry creation, adapting the number of visible entries to fit the terminal height.

#### Scenario: Show recent entries

- Given: The TUI is running and connected to the backend
- When: Entries are fetched successfully
- Then: Entries are displayed grouped by date
- And: The number of entries shown SHALL adapt to terminal height
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

#### Scenario: Responsive height - tall terminal

- Given: The terminal has sufficient height for all entries
- When: The home screen is displayed
- Then: All fetched entries are displayed
- And: The help bar is visible at the bottom

#### Scenario: Responsive height - short terminal

- Given: The terminal height is too short to display all entries
- When: The home screen is displayed
- Then: Only entries that fit are displayed
- And: The help bar SHALL always remain visible
- And: A scroll indicator shows how many entries are hidden

#### Scenario: Responsive height - resize

- Given: The terminal is resized while displaying entries
- When: The new dimensions are received
- Then: The visible entry count SHALL update immediately
- And: The cursor SHALL remain on a visible entry

#### Scenario: Responsive height - minimum entries

- Given: The terminal is extremely short
- When: The home screen is displayed
- Then: At least one entry SHALL be shown if entries exist
- And: The help bar may be hidden to prioritize content
