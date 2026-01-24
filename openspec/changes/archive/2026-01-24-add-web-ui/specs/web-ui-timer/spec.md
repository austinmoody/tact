## ADDED Requirements

### Requirement: Timer Display

The web UI SHALL display the current timer state prominently.

#### Scenario: No active timer

- **WHEN** no timer is running
- **THEN** a "Start Timer" button SHALL be displayed
- **AND** the elapsed time display SHALL show "00:00:00"

#### Scenario: Timer running

- **WHEN** a timer is running
- **THEN** the elapsed time SHALL be displayed in HH:MM:SS format
- **AND** the timer description SHALL be visible
- **AND** "Pause" and "Stop" buttons SHALL be available

#### Scenario: Timer paused

- **WHEN** a timer is paused
- **THEN** the elapsed time SHALL be displayed but not advancing
- **AND** "Resume" and "Stop" buttons SHALL be available
- **AND** a visual indicator SHALL show the paused state

### Requirement: Real-time Timer Updates

The web UI SHALL update the timer display in real-time using Server-Sent Events.

#### Scenario: SSE connection

- **WHEN** the timer page is loaded
- **THEN** an SSE connection SHALL be established to `/timer/stream`
- **AND** timer updates SHALL be received without page refresh

#### Scenario: Timer tick updates

- **WHEN** a timer is running
- **THEN** the elapsed time SHALL update every second via SSE
- **AND** the update SHALL be smooth without flickering

#### Scenario: SSE reconnection

- **WHEN** the SSE connection is lost
- **THEN** the client SHALL automatically attempt to reconnect
- **AND** the timer state SHALL be restored upon reconnection

### Requirement: Timer Controls

The web UI SHALL provide controls to manage the timer.

#### Scenario: Start timer

- **WHEN** user clicks "Start Timer"
- **THEN** a form SHALL appear to enter timer description
- **AND** upon submission, the timer SHALL start
- **AND** the display SHALL update to show running state

#### Scenario: Pause timer

- **WHEN** user clicks "Pause" on a running timer
- **THEN** the timer SHALL pause
- **AND** elapsed time SHALL stop incrementing
- **AND** the pause time SHALL be recorded

#### Scenario: Resume timer

- **WHEN** user clicks "Resume" on a paused timer
- **THEN** the timer SHALL resume from where it was paused
- **AND** elapsed time SHALL continue incrementing

#### Scenario: Stop timer

- **WHEN** user clicks "Stop" on a running or paused timer
- **THEN** a confirmation dialog SHALL appear
- **AND** upon confirmation, the timer SHALL stop

### Requirement: Timer to Entry Conversion

The web UI SHALL allow converting a stopped timer into a time entry.

#### Scenario: Create entry from timer

- **WHEN** a timer is stopped
- **THEN** a form SHALL appear to create an entry from the timer
- **AND** the form SHALL be pre-filled with: description as user_input, calculated duration_minutes
- **AND** user SHALL be able to edit before submission

#### Scenario: Discard timer

- **WHEN** user chooses to discard instead of create entry
- **THEN** the timer data SHALL be cleared
- **AND** no entry SHALL be created

### Requirement: Timer Persistence

The web UI timer state SHALL persist across page loads.

#### Scenario: Page refresh with running timer

- **WHEN** the page is refreshed while a timer is running
- **THEN** the timer SHALL continue running
- **AND** the elapsed time SHALL be accurate (not reset)

#### Scenario: Timer state on navigation

- **WHEN** user navigates to another page
- **THEN** the timer SHALL continue running in the background
- **AND** a timer indicator SHALL be visible in the navigation

### Requirement: Timer in Navigation

The web UI SHALL show timer status in the navigation header when a timer is active.

#### Scenario: Timer indicator in header

- **WHEN** a timer is running or paused
- **THEN** a compact timer indicator SHALL appear in the navigation
- **AND** clicking it SHALL navigate to the full timer page
