## ADDED Requirements

### Requirement: Entry list displays time code names

The TUI entry list SHALL display time code names alongside IDs for parsed entries.

#### Scenario: Show time code in entry list

- **WHEN** the home screen displays entries with parsed time codes
- **THEN** each entry SHALL show the time code ID and name between the user input and status
- **AND** long time code names SHALL be truncated to fit available space

#### Scenario: Entry without time code

- **WHEN** an entry has no parsed time code
- **THEN** the time code column SHALL be empty or show a placeholder
- **AND** the layout SHALL remain consistent with other entries

### Requirement: Entry detail displays time code names

The TUI entry detail modal SHALL display time code names in the parsed fields section.

#### Scenario: Show time code name in parsed fields

- **WHEN** the entry detail modal displays a parsed time code
- **THEN** the time code line SHALL show: ID, name, and confidence percentage
- **AND** the format SHALL be "Time Code: ID - Name (confidence%)"

#### Scenario: Time code name truncation in detail

- **WHEN** the time code name is longer than available space
- **THEN** the name SHALL be truncated with ellipsis
- **AND** the ID and confidence SHALL always be fully visible

### Requirement: Time codes are cached for display

The TUI SHALL cache time codes for efficient name lookups.

#### Scenario: Fetch time codes on startup

- **WHEN** the TUI starts and loads the home screen
- **THEN** time codes SHALL be fetched from the API
- **AND** a lookup map from ID to name SHALL be created

#### Scenario: Refresh time codes

- **WHEN** the user presses the refresh key (r)
- **THEN** time codes SHALL be re-fetched along with entries
- **AND** the lookup map SHALL be updated
