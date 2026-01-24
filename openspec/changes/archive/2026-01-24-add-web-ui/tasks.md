## 1. Project Setup

- [x] 1.1 Create `webui/` directory structure (main.go, handlers/, templates/, api/, static/)
- [x] 1.2 Initialize Go module (`go mod init tact-webui`)
- [x] 1.3 Add dependencies: Templ, Chi router (or standard net/http)
- [x] 1.4 Download and add static files: htmx.min.js, pico.min.css
- [x] 1.5 Create Makefile with build, dev, and templ generate targets

## 2. Core Infrastructure

- [x] 2.1 Create main.go with HTTP server setup (port 2200, configurable)
- [x] 2.2 Implement API client for backend communication (GET, POST, PUT, DELETE helpers)
- [x] 2.3 Add API URL configuration (flag and environment variable support)
- [x] 2.4 Create base layout template with navigation header
- [x] 2.5 Add static file serving for CSS and JS
- [x] 2.6 Implement error page template and error handling middleware

## 3. Entry List Page (Home)

- [x] 3.1 Create entries handler with list endpoint
- [x] 3.2 Create home page template with entry list grouped by date
- [x] 3.3 Create entry row component with status color coding
- [x] 3.4 Add empty state display when no entries exist
- [x] 3.5 Implement HTMX partial refresh for entry list

## 4. Entry Creation

- [x] 4.1 Create new entry form component (modal or inline)
- [x] 4.2 Implement POST handler for entry creation
- [x] 4.3 Add HTMX form submission with list refresh
- [x] 4.4 Add cancel functionality

## 5. Entry Detail and Edit

- [x] 5.1 Create entry detail view template showing all fields
- [x] 5.2 Implement GET handler for single entry
- [x] 5.3 Create entry edit form with editable fields
- [x] 5.4 Implement PUT handler for entry updates
- [x] 5.5 Add reparse button with POST handler
- [x] 5.6 Implement HTMX detail view open/close

## 6. Entry Filtering

- [x] 6.1 Add status filter dropdown (all, parsed, pending, failed)
- [x] 6.2 Add date range filter inputs
- [x] 6.3 Implement filter query parameters in list handler
- [x] 6.4 Add HTMX filter application without page reload

## 7. Timer Core

- [x] 7.1 Create timer page template with display and controls
- [x] 7.2 Implement timer state management (start, pause, resume, stop)
- [x] 7.3 Create SSE endpoint `/timer/stream` for real-time updates
- [x] 7.4 Add HTMX SSE integration for timer display updates
- [x] 7.5 Implement timer tick updates (1 second interval)

## 8. Timer Controls

- [x] 8.1 Create start timer form with description input
- [x] 8.2 Implement pause/resume toggle button
- [x] 8.3 Implement stop button with confirmation dialog
- [x] 8.4 Add timer-to-entry conversion form on stop
- [x] 8.5 Implement discard timer option

## 9. Timer Navigation Indicator

- [x] 9.1 Add compact timer display to navigation header
- [x] 9.2 Show timer indicator only when timer is active
- [x] 9.3 Make indicator clickable to navigate to timer page
- [x] 9.4 Update indicator via SSE alongside main timer display

## 10. Projects Management

- [x] 10.1 Create projects list page template
- [x] 10.2 Implement projects handler with CRUD operations
- [x] 10.3 Create add project form (ID, name, description)
- [x] 10.4 Create edit project form (name, description)
- [x] 10.5 Add deactivate button with confirmation
- [x] 10.6 Implement search/filter functionality

## 11. Time Codes Management

- [x] 11.1 Create time codes list page template
- [x] 11.2 Implement time codes handler with CRUD operations
- [x] 11.3 Create add time code form with project dropdown
- [x] 11.4 Create edit time code form
- [x] 11.5 Add deactivate button with confirmation

## 12. Work Types Management

- [x] 12.1 Create work types list page template
- [x] 12.2 Implement work types handler with CRUD operations
- [x] 12.3 Create quick-add form (name only, auto-generate ID)
- [x] 12.4 Create edit form
- [x] 12.5 Add deactivate button with confirmation

## 13. Context Documents

- [x] 13.1 Create context list component (reusable for projects and time codes)
- [x] 13.2 Implement context handler with CRUD operations
- [x] 13.3 Create add context form with multi-line textarea
- [x] 13.4 Create edit context form
- [x] 13.5 Add delete with confirmation
- [x] 13.6 Integrate context button into projects and time codes pages

## 14. Polish and Testing

- [x] 14.1 Add loading indicators for all HTMX requests
- [x] 14.2 Verify dark mode styling across all pages
- [x] 14.3 Test responsive layout on mobile viewport sizes
- [x] 14.4 Add SSE reconnection logic for timer
- [x] 14.5 Manual end-to-end testing of all features
- [x] 14.6 Update project README with web UI setup instructions

## 15. DevOps and Infrastructure

- [x] 15.1 Add Makefile targets (webui-generate, webui-build, webui-run, webui-dev)
- [x] 15.2 Create Dockerfile for webui with multi-stage build
- [x] 15.3 Add webui service to docker-compose.yml
- [x] 15.4 Add GitHub Actions CI workflow for webui (build, test, vet, staticcheck)
- [x] 15.5 Add GitHub Actions Docker workflow for webui image publishing
