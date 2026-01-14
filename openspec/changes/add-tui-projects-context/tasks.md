# Tasks

## 1. Add Data Models

- [ ] 1.1 Create `tui/model/project.go` with Project struct
- [ ] 1.2 Create `tui/model/context.go` with ContextDocument struct

## 2. Add API Client Methods

- [ ] 2.1 Add Project CRUD methods to `api/client.go`
- [ ] 2.2 Add Project Context methods (`FetchProjectContext`, `CreateProjectContext`)
- [ ] 2.3 Add Time Code Context methods (`FetchTimeCodeContext`, `CreateTimeCodeContext`)
- [ ] 2.4 Add generic Context methods (`UpdateContext`, `DeleteContext`)

## 3. Add Projects Screen

- [ ] 3.1 Create `ui/projects.go` with ProjectsScreen (list view)
- [ ] 3.2 Create `ui/project_edit.go` with ProjectEditModal (add/edit form)
- [ ] 3.3 Add ScreenProjects to app.go Screen enum
- [ ] 3.4 Add project-related message types to app.go
- [ ] 3.5 Wire up Projects screen navigation in app.go
- [ ] 3.6 Add "Projects" to menu.go

## 4. Add Context Management

- [ ] 4.1 Create `ui/context_list.go` with ContextListModal
- [ ] 4.2 Create `ui/context_edit.go` with ContextEditModal (using bubbles/textarea)
- [ ] 4.3 Implement textarea with Ctrl+S to save, Enter for newlines
- [ ] 4.4 Add context-related modal enums to app.go
- [ ] 4.5 Add context-related message types to app.go
- [ ] 4.6 Wire up context modals in app.go

## 5. Integrate Context with Existing Screens

- [ ] 5.1 Add `c` key handler to ProjectsScreen for opening context
- [ ] 5.2 Add `c` key handler to TimeCodesScreen for opening context
- [ ] 5.3 Update help text on both screens to show `[c] Context`

## Dependencies

- Task 2 depends on Task 1 (need models for API responses)
- Task 3 depends on Task 2 (need API methods for screen)
- Task 4 depends on Task 2 (need API methods for context)
- Task 5 depends on Tasks 3 and 4 (integrate context into screens)

## Notes

- Follow existing patterns from TimeCodesScreen and WorkTypesScreen
- Context edit uses `bubbles/textarea` for multi-line input
- Ctrl+S saves context (Enter adds newlines in textarea)
- All screens should support the existing key bindings (j/k, up/down, etc.)
