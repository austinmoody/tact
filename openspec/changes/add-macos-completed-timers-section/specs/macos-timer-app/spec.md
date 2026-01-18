## ADDED Requirements

### Requirement: Completed Timers Section

The app SHALL display completed timers from today in a separate section of the Timer List window.

#### Scenario: Display completed timers section

- Given: The timer list window is open with 1 active timer and 2 completed timers from today
- Then: The active timer appears in the top section
- And: A "Completed Today" section header appears below
- And: The 2 completed timers appear in the completed section

#### Scenario: Completed timer shows final duration

- Given: A timer was stopped after 45 minutes of work
- When: The timer list window displays the completed timer
- Then: The completed timer shows description and "45m" duration
- And: The elapsed time does not continue counting

#### Scenario: Start new timer from completed timer

- Given: The timer list window shows a completed timer with description "API refactoring"
- When: User clicks "Start New" on that completed timer
- Then: A new timer is created with description "API refactoring"
- And: The new timer starts in running state
- And: The completed timer remains in the completed section unchanged

#### Scenario: Hide completed section when empty

- Given: The timer list window is open with no completed timers from today
- Then: The "Completed Today" section header is not displayed

### Requirement: Completed Timer Cleanup

The app SHALL automatically remove completed timers from previous days.

#### Scenario: Cleanup old completed timers on launch

- Given: The app has completed timers from yesterday and today in storage
- When: The app is launched
- Then: Yesterday's completed timers are removed from storage
- And: Today's completed timers are retained

#### Scenario: Completed timers persist across restart same day

- Given: A timer was completed earlier today
- When: The app is quit and relaunched on the same day
- Then: The completed timer still appears in the "Completed Today" section

## MODIFIED Requirements

### Requirement: Timer State Transitions

Each timer SHALL support running, paused, and stopped states with appropriate transitions.

#### Scenario: Stop timer triggers API call

- Given: A timer has been running for 45 minutes with description "standup meeting"
- When: User clicks Stop
- Then: The app calls POST /entries with user_input "45m standup meeting"
- And: On success, the timer moves to stopped state with stoppedAt timestamp
- And: The timer appears in the "Completed Today" section
