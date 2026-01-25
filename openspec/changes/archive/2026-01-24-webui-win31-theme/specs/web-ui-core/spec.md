## MODIFIED Requirements

### Requirement: Base Layout

The web UI SHALL provide a consistent layout across all pages using a Windows 3.1 visual theme.

#### Scenario: Navigation header

- **WHEN** any page is displayed
- **THEN** a navigation header SHALL be visible with Windows 3.1 styling
- **AND** the header SHALL contain links to: Home, Entries, Timer, Projects, Time Codes, Work Types
- **AND** the header SHALL use a blue (#000080) title bar with white text

#### Scenario: Active navigation state

- **WHEN** a page is displayed
- **THEN** the corresponding navigation link SHALL be visually highlighted

#### Scenario: Windows 3.1 color scheme

- **WHEN** any page is displayed
- **THEN** the background color SHALL be Windows gray (#C0C0C0)
- **AND** windows/cards SHALL have beveled 3D borders (white top/left, dark gray bottom/right)
- **AND** title bars SHALL be navy blue (#000080) with white text

#### Scenario: Windows 3.1 typography

- **WHEN** any page is displayed
- **THEN** the primary font SHALL be Arial, Helvetica, sans-serif (approximating MS Sans Serif)
- **AND** title bar text SHALL be bold

#### Scenario: Windows 3.1 buttons

- **WHEN** a button is displayed
- **THEN** the button SHALL have a raised 3D beveled appearance
- **AND** the button background SHALL be Windows gray (#C0C0C0)
- **AND** the button SHALL have white top/left border and dark gray (#808080) bottom/right border

#### Scenario: Windows 3.1 form inputs

- **WHEN** a text input or select field is displayed
- **THEN** the field SHALL have a sunken 3D inset border effect
- **AND** the field SHALL have dark gray top/left border and white bottom/right border

## REMOVED Requirements

### Requirement: Dark mode support

**Reason**: Windows 3.1 theme provides a consistent retro aesthetic that replaces the dark mode styling.
**Migration**: Remove dark mode CSS and replace with Windows 3.1 theme styles.
