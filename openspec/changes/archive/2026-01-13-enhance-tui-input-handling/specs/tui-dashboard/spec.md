# tui-dashboard Specification Delta

## MODIFIED Requirements

### Requirement: Time Codes Management

The TUI SHALL allow viewing, adding, editing, and deactivating time codes.

#### Scenario: Full add time code form

- Given: The time codes screen is displayed
- When: a is pressed
- Then: An add modal opens with ID, name, description, keywords, and examples fields
- And: All fields are editable

#### Scenario: Submit full add form

- Given: The add modal is open with all fields populated
- When: Enter is pressed
- Then: The time code is created via API with all provided fields
- And: The list is refreshed

## ADDED Requirements

### Requirement: Text Input Handling

The TUI SHALL handle text input without inserting control characters.

#### Scenario: Paste text into input field

- Given: A text input field is focused
- When: Text is pasted (Ctrl+V or terminal paste)
- Then: The text appears in the input field
- And: No control characters are inserted

#### Scenario: Arrow key navigation in text field

- Given: A text input field is focused with text entered
- When: Left or right arrow key is pressed
- Then: The cursor moves within the text
- And: No control characters are inserted into the input

### Requirement: Responsive Form Layout

The TUI SHALL adapt form widths to the terminal size.

#### Scenario: Wide terminal

- Given: The terminal width is greater than 100 characters
- When: A modal form is displayed
- Then: Input fields use a width up to 80 characters

#### Scenario: Narrow terminal

- Given: The terminal width is less than 60 characters
- When: A modal form is displayed
- Then: Input fields use a minimum width of 30 characters

#### Scenario: Medium terminal

- Given: The terminal width is between 60 and 100 characters
- When: A modal form is displayed
- Then: Input fields scale proportionally to available space
