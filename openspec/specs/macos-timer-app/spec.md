# macos-timer-app Specification

## Purpose
TBD - created by archiving change add-macos-timer-app. Update Purpose after archive.
## Requirements
### Requirement: Dock Menu Interface

The app SHALL provide a right-click dock menu as the primary interface for timer management.

#### Scenario: View timer status in dock menu

- Given: The app is running with 2 timers (1 running, 1 paused)
- When: User right-clicks the dock icon
- Then: The menu displays "2 timers (1 running)" at the top

#### Scenario: Start new timer from dock menu

- Given: The app is running
- When: User right-clicks the dock icon and selects "Start New Timer..."
- Then: A small popup window appears with a text field for entering the timer description

#### Scenario: Access timer list from dock menu

- Given: The app is running with at least one timer
- When: User right-clicks the dock icon and selects "View All Timers"
- Then: The timer list window opens showing all active timers

#### Scenario: Access preferences from dock menu

- Given: The app is running
- When: User right-clicks the dock icon and selects "Preferences..."
- Then: The preferences window opens

### Requirement: Timer Creation

The app SHALL allow creating new timers with a description via a popup window.

#### Scenario: Create timer with description

- Given: The "Start New Timer" popup is open
- When: User enters "Working on API refactor" and presses Enter
- Then: A new timer starts with that description
- And: The popup closes

#### Scenario: Create timer with Start button

- Given: The "Start New Timer" popup is open with description entered
- When: User clicks the Start button
- Then: A new timer starts with the entered description
- And: The popup closes

#### Scenario: Cancel timer creation

- Given: The "Start New Timer" popup is open
- When: User presses Escape or clicks Cancel
- Then: The popup closes without creating a timer

#### Scenario: Empty description validation

- Given: The "Start New Timer" popup is open
- When: User attempts to start with empty description
- Then: The timer is not created
- And: The text field is highlighted or shows validation feedback

### Requirement: Multiple Timer Management

The app SHALL support multiple concurrent timers with automatic pause behavior.

#### Scenario: Start timer pauses running timer

- Given: Timer A is currently running
- When: User starts a new Timer B
- Then: Timer A is paused
- And: Timer B becomes the running timer

#### Scenario: Only one timer running at a time

- Given: Multiple timers exist
- When: User resumes a paused timer
- Then: Any other running timer is paused
- And: The selected timer becomes the only running timer

### Requirement: Timer State Transitions

Each timer SHALL support running, paused, and stopped states with appropriate transitions.

#### Scenario: Pause running timer

- Given: A timer is in running state
- When: User clicks Pause
- Then: The timer enters paused state
- And: Elapsed time stops accumulating

#### Scenario: Resume paused timer

- Given: A timer is in paused state and no other timer is running
- When: User clicks Resume
- Then: The timer enters running state
- And: Elapsed time continues accumulating from where it left off

#### Scenario: Stop timer triggers API call

- Given: A timer has been running for 45 minutes with description "standup meeting"
- When: User clicks Stop
- Then: The app calls POST /entries with user_input "45m standup meeting"
- And: On success, the timer is removed from the list

### Requirement: Timer List Window

The app SHALL provide a window displaying all active timers with real-time updates.

#### Scenario: Display timer details

- Given: The timer list window is open with active timers
- Then: Each timer shows its description, elapsed time, and current state (running/paused)

#### Scenario: Real-time elapsed time update

- Given: The timer list window is open with a running timer
- When: One second passes
- Then: The elapsed time display updates to reflect the new time

#### Scenario: Timer controls in list

- Given: The timer list window is open
- Then: Each timer has Pause/Resume and Stop buttons appropriate to its state

#### Scenario: Empty state

- Given: The timer list window is open with no timers
- Then: A message is displayed indicating no active timers

### Requirement: Time Entry API Integration

The app SHALL create time entries in the Tact backend when timers are stopped.

#### Scenario: Format duration correctly

- Given: A timer has accumulated 90 minutes
- When: The timer is stopped
- Then: The duration is formatted as "1h30m" in the API request

#### Scenario: Format short duration

- Given: A timer has accumulated 25 minutes
- When: The timer is stopped
- Then: The duration is formatted as "25m" in the API request

#### Scenario: Handle API success

- Given: A timer is stopped and API call succeeds
- Then: The timer is removed from the active timer list

#### Scenario: Handle API failure

- Given: A timer is stopped but API call fails
- Then: An error alert is displayed
- And: The timer remains in the list for retry

### Requirement: Timer Persistence

The app SHALL persist timer state across app restarts.

#### Scenario: Save timers on state change

- Given: A timer state changes (start, pause, resume)
- Then: All timer data is persisted to storage

#### Scenario: Restore timers on launch

- Given: The app was quit with 2 active timers
- When: The app is launched again
- Then: Both timers are restored with their descriptions and states

#### Scenario: Calculate elapsed time after restore

- Given: A running timer was persisted 5 minutes ago
- When: The app is launched
- Then: The timer shows 5 additional minutes of elapsed time

### Requirement: Preferences

The app SHALL provide configurable settings for API connection.

#### Scenario: Default API URL

- Given: The app is freshly installed
- Then: The API URL defaults to "http://localhost:2100"

#### Scenario: Configure API URL

- Given: The preferences window is open
- When: User changes the API URL and closes the window
- Then: The new URL is persisted and used for subsequent API calls

#### Scenario: Keyboard shortcut for preferences

- Given: The app is running
- When: User presses Cmd+,
- Then: The preferences window opens

### Requirement: Quit Protection

The app SHALL warn users before quitting with active timers.

#### Scenario: Warn on quit with active timers

- Given: The app has active timers
- When: User attempts to quit the app
- Then: A confirmation dialog appears warning about losing timer data

#### Scenario: Allow quit without active timers

- Given: The app has no active timers
- When: User quits the app
- Then: The app quits immediately without warning

