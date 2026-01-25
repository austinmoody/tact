## ADDED Requirements

### Requirement: Web Server

The web UI SHALL run as a standalone Go HTTP server on a configurable port.

#### Scenario: Start web server

- **WHEN** the web UI application starts
- **THEN** an HTTP server SHALL listen on port 2200 by default
- **AND** the port SHALL be configurable via `--port` flag or `TACT_WEBUI_PORT` environment variable

#### Scenario: Serve static assets

- **WHEN** a request is made for static files (CSS, JS)
- **THEN** the server SHALL serve files from the static directory
- **AND** appropriate cache headers SHALL be set

### Requirement: Backend API Configuration

The web UI SHALL connect to the TACT backend API using configurable settings.

#### Scenario: Default API URL

- **WHEN** no API URL is configured
- **THEN** the web UI SHALL connect to `http://localhost:2100`

#### Scenario: Custom API URL via flag

- **WHEN** the web UI is started with `--api http://custom:2100`
- **THEN** all API requests SHALL be sent to the specified URL

#### Scenario: Custom API URL via environment

- **WHEN** `TACT_API_URL` environment variable is set
- **THEN** the web UI SHALL use that URL for API requests

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

### Requirement: HTMX Integration

The web UI SHALL use HTMX for dynamic page updates without full page reloads.

#### Scenario: Partial page updates

- **WHEN** a user action requires updating part of the page
- **THEN** only the affected section SHALL be replaced via HTMX
- **AND** browser history SHALL be updated appropriately

#### Scenario: Loading indicators

- **WHEN** an HTMX request is in progress
- **THEN** a loading indicator SHALL be displayed

### Requirement: Error Handling

The web UI SHALL display errors gracefully.

#### Scenario: API connection failure

- **WHEN** the backend API is unreachable
- **THEN** an error message SHALL be displayed to the user
- **AND** a retry option SHALL be available

#### Scenario: API error response

- **WHEN** the backend API returns an error
- **THEN** the error message SHALL be displayed in a user-friendly format
