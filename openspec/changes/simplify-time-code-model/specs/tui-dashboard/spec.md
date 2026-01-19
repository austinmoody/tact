## MODIFIED Requirements

### Requirement: Time Codes Management

The TUI SHALL allow viewing, adding, editing, and deactivating time codes.

#### Scenario: Simplified add time code form

- Given: The time codes screen is displayed
- When: a is pressed
- Then: An add modal opens with ID, name, and project fields
- And: All fields are editable

#### Scenario: Submit simplified add form

- Given: The add modal is open with ID, name, and project populated
- When: Enter is pressed
- Then: The time code is created via API with the provided fields
- And: The list is refreshed
