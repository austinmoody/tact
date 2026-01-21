### Requirement: Entries display under correct date header
The TUI home screen SHALL display time entries grouped under the date header matching their `entry_date` field.

#### Scenario: Entry with past entry_date displays under correct header
- **WHEN** an entry has `entry_date` of "2026-01-18" and `created_at` of "2026-01-20"
- **THEN** the entry SHALL appear under the "Saturday - Jan 18, 2026" date header

#### Scenario: Today's entries appear first
- **WHEN** entries exist for today and previous days
- **THEN** today's entries SHALL appear at the top, followed by older dates in descending order

#### Scenario: Entries within same date sorted by creation time
- **WHEN** multiple entries have the same `entry_date`
- **THEN** they SHALL be sorted by `created_at` descending (newest first) within that date group

### Requirement: Timer Completed Today uses local timezone
The timer panel "Completed Today" section SHALL show only timers that were stopped today in the user's local timezone.

#### Scenario: Timer stopped yesterday evening not shown as completed today
- **WHEN** a timer was stopped at 9:37 PM local time yesterday
- **THEN** it SHALL NOT appear in "Completed Today" section today

#### Scenario: Timer stopped after local midnight shown as completed today
- **WHEN** a timer is stopped at 12:01 AM local time today
- **THEN** it SHALL appear in "Completed Today" section

### Requirement: Old completed timers cleaned up using local timezone
The timer manager SHALL clean up completed timers from previous days based on local timezone midnight.

#### Scenario: Yesterday's completed timer removed on app start
- **WHEN** the TUI starts and a completed timer has `stopped_at` before local midnight today
- **THEN** that timer SHALL be removed from persistence
