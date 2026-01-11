# tui-dashboard Specification

## Purpose
TBD - created by archiving change add-tui-dashboard. Update Purpose after archive.
## Requirements
### Requirement: Dashboard Display

The TUI SHALL display time codes and work types in a split-pane layout.

#### Scenario: Show time codes pane

- Given: The TUI is running and connected to the backend
- When: Time codes are fetched successfully
- Then: Time codes are displayed in the left pane with code and name

#### Scenario: Show work types pane

- Given: The TUI is running and connected to the backend
- When: Work types are fetched successfully
- Then: Work types are displayed in the right pane with id and name

#### Scenario: Show detail panel

- Given: An item is selected
- When: Enter is pressed
- Then: A detail panel shows additional information for the selected item

### Requirement: Keyboard Navigation

The TUI SHALL support vim-style keyboard navigation.

#### Scenario: Navigate within pane

- Given: The TUI is displaying data
- When: j or down arrow is pressed
- Then: The cursor moves down one item

#### Scenario: Navigate up

- Given: The cursor is not on the first item
- When: k or up arrow is pressed
- Then: The cursor moves up one item

#### Scenario: Switch panes

- Given: The left pane is focused
- When: l is pressed
- Then: Focus moves to the right pane

#### Scenario: Switch panes back

- Given: The right pane is focused
- When: h is pressed
- Then: Focus moves to the left pane

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

