## ADDED Requirements

### Requirement: Projects Management

The web UI SHALL provide a page to manage projects.

#### Scenario: List projects

- **WHEN** the projects page is loaded
- **THEN** all projects SHALL be displayed
- **AND** each project SHALL show: ID, name, description, active status

#### Scenario: Add project

- **WHEN** user clicks "Add Project"
- **THEN** a form SHALL appear with fields: ID, name, description
- **AND** upon submission, the project SHALL be created
- **AND** the list SHALL refresh

#### Scenario: Edit project

- **WHEN** user clicks "Edit" on a project
- **THEN** a form SHALL appear with current values
- **AND** editable fields SHALL be: name, description

#### Scenario: Deactivate project

- **WHEN** user clicks "Deactivate" on an active project
- **THEN** the project SHALL be deactivated via API
- **AND** the list SHALL update to show inactive status

### Requirement: Time Codes Management

The web UI SHALL provide a page to manage time codes.

#### Scenario: List time codes

- **WHEN** the time codes page is loaded
- **THEN** all time codes SHALL be displayed
- **AND** each time code SHALL show: ID, name, project name, active status

#### Scenario: Add time code

- **WHEN** user clicks "Add Time Code"
- **THEN** a form SHALL appear with fields: ID, name, project (dropdown)
- **AND** upon submission, the time code SHALL be created
- **AND** the list SHALL refresh

#### Scenario: Edit time code

- **WHEN** user clicks "Edit" on a time code
- **THEN** a form SHALL appear with current values
- **AND** editable fields SHALL be: name, project

#### Scenario: Deactivate time code

- **WHEN** user clicks "Deactivate" on an active time code
- **THEN** the time code SHALL be deactivated via API
- **AND** the list SHALL update to show inactive status

### Requirement: Work Types Management

The web UI SHALL provide a page to manage work types.

#### Scenario: List work types

- **WHEN** the work types page is loaded
- **THEN** all work types SHALL be displayed
- **AND** each work type SHALL show: ID, name, active status

#### Scenario: Quick add work type

- **WHEN** user clicks "Add Work Type"
- **THEN** a form SHALL appear with name field only
- **AND** ID SHALL be auto-generated
- **AND** upon submission, the work type SHALL be created

#### Scenario: Edit work type

- **WHEN** user clicks "Edit" on a work type
- **THEN** a form SHALL appear with current name
- **AND** the name SHALL be editable

#### Scenario: Deactivate work type

- **WHEN** user clicks "Deactivate" on an active work type
- **THEN** the work type SHALL be deactivated via API
- **AND** the list SHALL update to show inactive status

### Requirement: Context Documents Management

The web UI SHALL allow managing context documents for projects and time codes.

#### Scenario: View context for project

- **WHEN** user clicks "Context" on a project
- **THEN** a list of context documents for that project SHALL be displayed
- **AND** each document SHALL show truncated content (first 50 characters)

#### Scenario: View context for time code

- **WHEN** user clicks "Context" on a time code
- **THEN** a list of context documents for that time code SHALL be displayed

#### Scenario: Add context document

- **WHEN** user clicks "Add Context" in the context list
- **THEN** a form SHALL appear with a multi-line text area
- **AND** upon submission, the context document SHALL be created

#### Scenario: Edit context document

- **WHEN** user clicks "Edit" on a context document
- **THEN** a form SHALL appear with the full content
- **AND** the content SHALL be editable in a multi-line text area

#### Scenario: Delete context document

- **WHEN** user clicks "Delete" on a context document
- **THEN** a confirmation SHALL be requested
- **AND** upon confirmation, the document SHALL be deleted
- **AND** the list SHALL refresh

### Requirement: Search and Filter

The web UI SHALL provide search and filter capabilities on management pages.

#### Scenario: Search projects

- **WHEN** user types in the search box on projects page
- **THEN** the list SHALL filter to show only matching projects
- **AND** search SHALL match against ID and name

#### Scenario: Filter by active status

- **WHEN** user toggles the "Show inactive" checkbox
- **THEN** inactive items SHALL be shown or hidden accordingly
