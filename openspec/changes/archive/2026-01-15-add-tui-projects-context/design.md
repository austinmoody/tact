# Design: TUI Projects and Context Management

## Overview

This design adds Projects and Context Document management to the TUI, following existing patterns for Time Codes and Work Types.

## Navigation Structure

```
Home
  └── Menu (m)
        ├── Projects      ← NEW
        ├── Time Codes
        └── Work Types

Projects Screen
  └── Context for Project (c) ← NEW

Time Codes Screen
  └── Context for Time Code (c) ← NEW
```

## Screen Layouts

### Projects Screen

```
Projects                                    [Esc] Back

Manage Projects
────────────────────────────────────────────
> izg         IZG Hub
  acme        ACME Corp Project
  default     Default Project [inactive]

[a] Add  [e] Edit  [d] Delete  [c] Context  [r] Refresh  [Esc] Back
```

### Projects Add/Edit Modal

```
┌─────────────────────────────────────┐
│ Add Project                         │
│                                     │
│ ID:          [________________]     │
│ Name:        [________________]     │
│ Description: [________________]     │
│                                     │
│ [Enter] Save  [Esc] Cancel          │
└─────────────────────────────────────┘
```

### Context List Modal (for Project or Time Code)

```
┌─────────────────────────────────────────────────────┐
│ Context for: FEDS-165                               │
│                                                     │
│ > Time code FEDS-165 should be used for: Incide... │
│   Do NOT use FEDS-165 for security incidents...    │
│                                                     │
│ [a] Add  [e] Edit  [d] Delete  [Esc] Close          │
└─────────────────────────────────────────────────────┘
```

### Context Add/Edit Modal

```
┌─────────────────────────────────────────────────────┐
│ Add Context                                         │
│                                                     │
│ Content:                                            │
│ ┌─────────────────────────────────────────────────┐ │
│ │ Time code FEDS-165 should be used for:         │ │
│ │ - Incident response activities                 │ │
│ │ - Help Desk activities and meetings            │ │
│ │ - IZG Hub updates                              │ │
│ │ _                                              │ │
│ └─────────────────────────────────────────────────┘ │
│                                                     │
│ [Ctrl+S] Save  [Esc] Cancel                         │
└─────────────────────────────────────────────────────┘
```

Uses `bubbles/textarea` for multi-line editing. Enter adds newlines, Ctrl+S saves.

## Data Models (Go)

### Project Model

```go
type Project struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Active      bool      `json:"active"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### ContextDocument Model

```go
type ContextDocument struct {
    ID         string    `json:"id"`
    ProjectID  *string   `json:"project_id"`
    TimeCodeID *string   `json:"time_code_id"`
    Content    string    `json:"content"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}
```

## API Client Methods

New methods needed:

```go
// Projects
FetchProjects() ([]Project, error)
CreateProject(id, name, description string) (*Project, error)
UpdateProject(id string, updates ProjectUpdate) (*Project, error)
DeleteProject(id string) error

// Context - Project
FetchProjectContext(projectID string) ([]ContextDocument, error)
CreateProjectContext(projectID, content string) (*ContextDocument, error)

// Context - Time Code
FetchTimeCodeContext(timeCodeID string) ([]ContextDocument, error)
CreateTimeCodeContext(timeCodeID, content string) (*ContextDocument, error)

// Context - Generic
UpdateContext(contextID, content string) (*ContextDocument, error)
DeleteContext(contextID string) error
```

## Component Structure

### New Files

- `tui/model/project.go` - Project struct
- `tui/model/context.go` - ContextDocument struct
- `tui/ui/projects.go` - ProjectsScreen
- `tui/ui/project_edit.go` - ProjectEditModal
- `tui/ui/context_list.go` - ContextListModal
- `tui/ui/context_edit.go` - ContextEditModal

### Modified Files

- `tui/api/client.go` - Add Project and Context API methods
- `tui/ui/app.go` - Add ScreenProjects, new modals, message types
- `tui/ui/menu.go` - Add "Projects" menu item
- `tui/ui/timecodes.go` - Add `c` key handler for context
- `tui/ui/keys.go` - Add Context key binding if needed

## Key Bindings

| Key | Context | Action |
|-----|---------|--------|
| `c` | Projects screen, item selected | Open context list for project |
| `c` | Time Codes screen, item selected | Open context list for time code |
| `a` | Context list modal | Add new context document |
| `e` | Context list modal | Edit selected context document |
| `d` | Context list modal | Delete selected context document |

## Multi-line Text Input

Context documents are typically multi-line. Use the `bubbles/textarea` component for proper multi-line editing.

Features:
- Multi-line text entry with natural line breaks (Enter key)
- Cursor navigation (arrow keys, home/end)
- Line wrapping within the textarea bounds
- Scrolling for content that exceeds visible area
- Submit with Ctrl+S or a dedicated "Save" action (since Enter adds newlines)

The context edit modal will use a textarea instead of a text input, with adjusted key bindings:
- **Enter**: Add newline
- **Ctrl+S** or **Tab then Enter**: Save and close
- **Esc**: Cancel and close
