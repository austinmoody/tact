# Tasks

## 1. Add Data Models

- [x] 1.1 Create `tui/model/project.go` with Project struct
- [x] 1.2 Create `tui/model/context.go` with ContextDocument struct

## 2. Add API Client Methods

- [x] 2.1 Add Project CRUD methods to `api/client.go`
- [x] 2.2 Add Project Context methods (`FetchProjectContext`, `CreateProjectContext`)
- [x] 2.3 Add Time Code Context methods (`FetchTimeCodeContext`, `CreateTimeCodeContext`)
- [x] 2.4 Add generic Context methods (`UpdateContext`, `DeleteContext`)

## 3. Add Projects Screen

- [x] 3.1 Create `ui/projects.go` with ProjectsScreen (list view)
- [x] 3.2 Create `ui/project_edit.go` with ProjectEditModal (add/edit form)
- [x] 3.3 Add ScreenProjects to app.go Screen enum
- [x] 3.4 Add project-related message types to app.go
- [x] 3.5 Wire up Projects screen navigation in app.go
- [x] 3.6 Add "Projects" to menu.go

## 4. Add Context Management

- [x] 4.1 Create `ui/context_list.go` with ContextListModal
- [x] 4.2 Create `ui/context_edit.go` with ContextEditModal (using bubbles/textarea)
- [x] 4.3 Implement textarea with Ctrl+S to save, Enter for newlines
- [x] 4.4 Add context-related modal enums to app.go
- [x] 4.5 Add context-related message types to app.go
- [x] 4.6 Wire up context modals in app.go

## 5. Integrate Context with Existing Screens

- [x] 5.1 Add `c` key handler to ProjectsScreen for opening context
- [x] 5.2 Add `c` key handler to TimeCodesScreen for opening context
- [x] 5.3 Update help text on both screens to show `[c] Context`

## Dependencies

- Task 2 depends on Task 1 (need models for API responses)
- Task 3 depends on Task 2 (need API methods for screen)
- Task 4 depends on Task 2 (need API methods for context)
- Task 5 depends on Tasks 3 and 4 (integrate context into screens)

## 6. Add Project Selection to Time Code Edit

- [x] 6.1 Update API client `CreateTimeCode` to accept `projectID` parameter
- [x] 6.2 Update API client `TimeCodeUpdate` struct to include `ProjectID` field
- [x] 6.3 Create project selector component for the edit modal
- [x] 6.4 Integrate project selector into `TimeCodeEditModal` (add/edit modes)
- [x] 6.5 Fetch available projects when opening the modal
- [x] 6.6 Test creating and editing time codes with project selection

## Notes

- Follow existing patterns from TimeCodesScreen and WorkTypesScreen
- Context edit uses `bubbles/textarea` for multi-line input
- Ctrl+S saves context (Enter adds newlines in textarea)
- All screens should support the existing key bindings (j/k, up/down, etc.)
- Project selector uses up/down to navigate, shows project name with ID
