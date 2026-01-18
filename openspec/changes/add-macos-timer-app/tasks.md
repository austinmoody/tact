## 1. Project Setup
- [ ] 1.1 Create Xcode project for "Tact Timer" in `/macos` directory
- [ ] 1.2 Configure as macOS app with AppKit lifecycle
- [ ] 1.3 Set up basic AppDelegate with dock menu placeholder
- [ ] 1.4 Verify app appears in Dock when running

## 2. Timer Model & Manager
- [ ] 2.1 Create `Timer` struct with id, description, state, startedAt, accumulatedSeconds
- [ ] 2.2 Create `TimerManager` class to manage timer collection
- [ ] 2.3 Implement timer state transitions (start, pause, resume, stop)
- [ ] 2.4 Implement "start new pauses current" logic
- [ ] 2.5 Add persistence to UserDefaults (save/load timers as JSON)
- [ ] 2.6 Restore timers on app launch with correct elapsed time

## 3. Dock Menu
- [ ] 3.1 Create dock menu showing timer count and running status
- [ ] 3.2 Add "Start New Timer..." menu item
- [ ] 3.3 Add "View All Timers" menu item
- [ ] 3.4 Add "Preferences..." menu item
- [ ] 3.5 Wire dock menu to AppDelegate applicationDockMenu

## 4. New Timer Popup
- [ ] 4.1 Create small floating window with text field
- [ ] 4.2 Add Start and Cancel buttons
- [ ] 4.3 Handle Enter key to start, Escape to cancel
- [ ] 4.4 Close window and start timer on confirmation

## 5. Timer List Window
- [ ] 5.1 Create window with table/list of timers
- [ ] 5.2 Display description, elapsed time, and state for each
- [ ] 5.3 Add Pause/Resume button per timer
- [ ] 5.4 Add Stop button per timer
- [ ] 5.5 Update elapsed time display every second
- [ ] 5.6 Handle empty state message

## 6. API Integration
- [ ] 6.1 Create APIClient with configurable base URL
- [ ] 6.2 Implement POST /entries with user_input body
- [ ] 6.3 Create TimeFormatter to format duration (45m, 1h30m)
- [ ] 6.4 On timer stop, format entry and call API
- [ ] 6.5 Handle API errors with alert, keep timer for retry
- [ ] 6.6 Remove timer from list on API success

## 7. Preferences Window
- [ ] 7.1 Create window with API URL text field
- [ ] 7.2 Load/save URL from UserDefaults
- [ ] 7.3 Default to http://localhost:2100
- [ ] 7.4 Wire Cmd+, shortcut to open preferences

## 8. Polish
- [ ] 8.1 Warn on quit if active timers exist
- [ ] 8.2 Add simple app icon
- [ ] 8.3 Test persistence across restart
- [ ] 8.4 Test API integration with running backend

## Verification

1. Build and run Tact Timer from Xcode
2. Right-click dock icon → "Start New Timer..." → enter description → Start
3. Verify timer count shows in dock menu
4. Open "View All Timers" → verify elapsed time updates live
5. Start second timer → verify first timer auto-pauses
6. Stop a timer → verify entry created in backend:
   ```bash
   curl http://localhost:2100/entries | jq '.[-1]'
   ```
7. Quit and relaunch → verify timers restored with correct elapsed time
8. Test Preferences → API URL persists across restart
