# Design: Enhance TUI with Entry Management

## Screen Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│  Home Screen (Default)                                          │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  Tact - Time Tracking                      [n] New Entry  │  │
│  ├───────────────────────────────────────────────────────────┤  │
│  │                                                           │  │
│  │  Recent Entries                                           │  │
│  │  ─────────────────                                        │  │
│  │  > 2h coding on Project Alpha          parsed   Today     │  │
│  │    meeting with team                   pending  Today     │  │
│  │    1.5 hours on alpha project          parsed   Today     │  │
│  │    half an hour reviewing PRs          parsed   Today     │  │
│  │    stuff                               parsed   Today     │  │
│  │                                                           │  │
│  ├───────────────────────────────────────────────────────────┤  │
│  │  [n] New  [Enter] Details  [m] Menu  [r] Refresh  [q] Quit│  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## Modal Designs

### New Entry Modal

```
┌─────────────────────────────────────────┐
│  New Entry                              │
├─────────────────────────────────────────┤
│                                         │
│  Enter time entry:                      │
│  ┌───────────────────────────────────┐  │
│  │ 2h working on Project Alpha_      │  │
│  └───────────────────────────────────┘  │
│                                         │
│  [Enter] Submit    [Esc] Cancel         │
│                                         │
└─────────────────────────────────────────┘
```

### Entry Detail Modal

```
┌─────────────────────────────────────────────────┐
│  Entry Details                                  │
├─────────────────────────────────────────────────┤
│                                                 │
│  Raw Text:    "2h coding on Project Alpha"      │
│  Status:      parsed                            │
│  Duration:    120 minutes                       │
│  Time Code:   PROJ-001 (Project Alpha)          │
│  Work Type:   development                       │
│  Description: Coding                            │
│  Confidence:  0.85                              │
│  Entry Date:  2026-01-11                        │
│                                                 │
│  [p] Reparse    [Esc] Close                     │
│                                                 │
└─────────────────────────────────────────────────┘
```

### Menu Modal

```
┌─────────────────────────────────────────┐
│  Menu                                   │
├─────────────────────────────────────────┤
│                                         │
│  > Time Codes                           │
│    Work Types                           │
│                                         │
│  [Enter] Select    [Esc] Close          │
│                                         │
└─────────────────────────────────────────┘
```

## Management Screens

### Time Codes Screen

```
┌─────────────────────────────────────────────────────────────────┐
│  Time Codes                                          [a] Add    │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  > PROJ-001    Project Alpha         active                     │
│    ADMIN-01    Admin Tasks           active                     │
│    OLD-PROJ    Old Project           inactive                   │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│  [a] Add  [e] Edit  [Enter] Details  [Esc] Back  [d] Deactivate │
└─────────────────────────────────────────────────────────────────┘
```

### Time Code Edit Modal

```
┌─────────────────────────────────────────────────┐
│  Edit Time Code                                 │
├─────────────────────────────────────────────────┤
│                                                 │
│  ID:          PROJ-001 (readonly)               │
│                                                 │
│  Name:        ┌─────────────────────────────┐   │
│               │ Project Alpha               │   │
│               └─────────────────────────────┘   │
│                                                 │
│  Description: ┌─────────────────────────────┐   │
│               │ Main project work           │   │
│               └─────────────────────────────┘   │
│                                                 │
│  Keywords:    ┌─────────────────────────────┐   │
│               │ alpha, main                 │   │
│               └─────────────────────────────┘   │
│               (comma-separated)                 │
│                                                 │
│  [Enter] Save    [Esc] Cancel                   │
│                                                 │
└─────────────────────────────────────────────────┘
```

### Quick Add Time Code Modal

```
┌─────────────────────────────────────────────────┐
│  Add Time Code                                  │
├─────────────────────────────────────────────────┤
│                                                 │
│  ID:          ┌─────────────────────────────┐   │
│               │ PROJ-NEW                    │   │
│               └─────────────────────────────┘   │
│                                                 │
│  Name:        ┌─────────────────────────────┐   │
│               │ New Project                 │   │
│               └─────────────────────────────┘   │
│                                                 │
│  [Enter] Create    [Esc] Cancel                 │
│                                                 │
└─────────────────────────────────────────────────┘
```

### Work Types Screen

```
┌─────────────────────────────────────────────────────────────────┐
│  Work Types                                          [a] Add    │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  > development     Development              active              │
│    meeting         Meeting                  active              │
│    code-review     Code Review              active              │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│  [a] Add  [e] Edit  [Esc] Back  [d] Deactivate                  │
└─────────────────────────────────────────────────────────────────┘
```

### Quick Add Work Type Modal

```
┌─────────────────────────────────────────────────┐
│  Add Work Type                                  │
├─────────────────────────────────────────────────┤
│                                                 │
│  Name:        ┌─────────────────────────────┐   │
│               │ Research                    │   │
│               └─────────────────────────────┘   │
│               (ID auto-generated as slug)       │
│                                                 │
│  [Enter] Create    [Esc] Cancel                 │
│                                                 │
└─────────────────────────────────────────────────┘
```

## Component Structure

```
tui/
├── main.go
├── api/
│   └── client.go           # Extended with entries, mutations
├── ui/
│   ├── app.go              # NEW: Root app model, screen routing
│   ├── home.go             # NEW: Home screen (entries list)
│   ├── timecodes.go        # NEW: Time codes management screen
│   ├── worktypes.go        # NEW: Work types management screen
│   ├── modal/
│   │   ├── entry_input.go  # NEW: New entry modal
│   │   ├── entry_detail.go # NEW: Entry detail modal
│   │   ├── menu.go         # NEW: Menu modal
│   │   ├── timecode_edit.go# NEW: Time code edit/add modal
│   │   └── worktype_edit.go# NEW: Work type edit/add modal
│   ├── keys.go             # Extended keybindings
│   └── styles.go           # Extended styles
├── model/
│   ├── entry.go            # NEW: Entry struct
│   ├── timecode.go
│   ├── worktype.go
│   └── time.go
└── go.mod
```

## Navigation Flow

```
                    ┌─────────────┐
                    │   Home      │
                    │  (entries)  │
                    └─────────────┘
                          │
           ┌──────────────┼──────────────┐
           │              │              │
           ▼              ▼              ▼
    ┌─────────────┐ ┌───────────┐ ┌───────────┐
    │ New Entry   │ │  Entry    │ │   Menu    │
    │   Modal     │ │  Detail   │ │   Modal   │
    │   (n)       │ │  (Enter)  │ │   (m)     │
    └─────────────┘ └───────────┘ └───────────┘
                                       │
                          ┌────────────┴────────────┐
                          ▼                         ▼
                   ┌─────────────┐           ┌─────────────┐
                   │ Time Codes  │           │ Work Types  │
                   │   Screen    │           │   Screen    │
                   └─────────────┘           └─────────────┘
                          │                         │
              ┌───────────┼───────────┐    ┌───────┴───────┐
              ▼           ▼           ▼    ▼               ▼
         ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐
         │  Add   │ │  Edit  │ │ Detail │ │  Add   │ │  Edit  │
         │ Modal  │ │ Modal  │ │ Modal  │ │ Modal  │ │ Modal  │
         └────────┘ └────────┘ └────────┘ └────────┘ └────────┘
```

## API Client Extensions

```go
// Entries
func (c *Client) FetchEntries(limit int) ([]model.Entry, error)
func (c *Client) CreateEntry(rawText string) (*model.Entry, error)
func (c *Client) ReparseEntry(id string) (*model.Entry, error)

// Time Codes (mutations)
func (c *Client) CreateTimeCode(id, name string) (*model.TimeCode, error)
func (c *Client) UpdateTimeCode(id string, updates TimeCodeUpdate) (*model.TimeCode, error)
func (c *Client) DeleteTimeCode(id string) error  // soft-delete

// Work Types (mutations)
func (c *Client) CreateWorkType(name string) (*model.WorkType, error)
func (c *Client) UpdateWorkType(id string, updates WorkTypeUpdate) (*model.WorkType, error)
func (c *Client) DeleteWorkType(id string) error  // soft-delete
```

## Keybinding Summary

### Global
- `q` / `Ctrl+C` - Quit application
- `Esc` - Close modal / Go back

### Home Screen
- `n` - New entry modal
- `j` / `↓` - Move cursor down
- `k` / `↑` - Move cursor up
- `Enter` - Open entry detail modal
- `m` - Open menu
- `r` - Refresh entries

### Entry Detail Modal
- `p` - Reparse entry
- `Esc` - Close modal

### Menu Modal
- `j` / `↓` - Move cursor down
- `k` / `↑` - Move cursor up
- `Enter` - Select option
- `Esc` - Close menu

### Time Codes / Work Types Screen
- `a` - Quick add
- `e` - Edit selected
- `d` - Deactivate selected
- `Enter` - View details (time codes only)
- `j` / `↓` - Move cursor down
- `k` / `↑` - Move cursor up
- `Esc` - Back to home

### Edit/Add Modals
- `Tab` - Next field
- `Shift+Tab` - Previous field
- `Enter` - Save/Submit
- `Esc` - Cancel
