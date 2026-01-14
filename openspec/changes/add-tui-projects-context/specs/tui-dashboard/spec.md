# tui-dashboard Spec Delta

## ADDED Requirements

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

## MODIFIED Requirements

### Requirement: Menu Navigation

The TUI SHALL provide a menu for accessing management screens.

#### Scenario: Menu options

- Given: The menu modal is open
- Then: Options for "Projects", "Time Codes", and "Work Types" are displayed

#### Scenario: Select projects

- Given: The menu modal is open with "Projects" selected
- When: Enter is pressed
- Then: The projects management screen is displayed
