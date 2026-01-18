import AppKit

class AppDelegate: NSObject, NSApplicationDelegate {
    
    private var timerManager: TimerManager!
    private var dockMenuController: DockMenuController!
    private var newTimerWindowController: NewTimerWindowController?
    private var timerListWindowController: TimerListWindowController?
    private var preferencesWindowController: PreferencesWindowController?
    
    func applicationDidFinishLaunching(_ notification: Notification) {
        timerManager = TimerManager()
        dockMenuController = DockMenuController(timerManager: timerManager)
        
        dockMenuController.onStartNewTimer = { [weak self] in
            self?.showNewTimerWindow()
        }
        
        dockMenuController.onViewAllTimers = { [weak self] in
            self?.showTimerListWindow()
        }
        
        dockMenuController.onOpenPreferences = { [weak self] in
            self?.showPreferencesWindow()
        }

        setupMainMenu()

        // Show the timer list window on launch
        showTimerListWindow()
    }
    
    func applicationDockMenu(_ sender: NSApplication) -> NSMenu? {
        return dockMenuController.buildMenu()
    }
    
    func applicationShouldTerminate(_ sender: NSApplication) -> NSApplication.TerminateReply {
        if timerManager.hasActiveTimers {
            let alert = NSAlert()
            alert.messageText = "Active Timers"
            alert.informativeText = "You have active timers. Quitting will preserve them, but any running timer will stop accumulating time. Are you sure you want to quit?"
            alert.alertStyle = .warning
            alert.addButton(withTitle: "Quit")
            alert.addButton(withTitle: "Cancel")
            
            let response = alert.runModal()
            if response == .alertSecondButtonReturn {
                return .terminateCancel
            }
        }
        return .terminateNow
    }
    
    func applicationWillTerminate(_ notification: Notification) {
        timerManager.pauseRunningTimer()
        timerManager.save()
    }
    
    private func setupMainMenu() {
        let mainMenu = NSMenu()
        
        // App menu
        let appMenuItem = NSMenuItem()
        mainMenu.addItem(appMenuItem)
        let appMenu = NSMenu()
        appMenuItem.submenu = appMenu
        
        appMenu.addItem(withTitle: "About Tact Timer", action: #selector(NSApplication.orderFrontStandardAboutPanel(_:)), keyEquivalent: "")
        appMenu.addItem(NSMenuItem.separator())
        
        let preferencesItem = NSMenuItem(title: "Preferences...", action: #selector(showPreferencesWindow), keyEquivalent: ",")
        preferencesItem.target = self
        appMenu.addItem(preferencesItem)
        
        appMenu.addItem(NSMenuItem.separator())
        appMenu.addItem(withTitle: "Quit Tact Timer", action: #selector(NSApplication.terminate(_:)), keyEquivalent: "q")
        
        // Window menu
        let windowMenuItem = NSMenuItem()
        mainMenu.addItem(windowMenuItem)
        let windowMenu = NSMenu(title: "Window")
        windowMenuItem.submenu = windowMenu
        
        let timerListItem = NSMenuItem(title: "Timer List", action: #selector(showTimerListWindow), keyEquivalent: "l")
        timerListItem.target = self
        windowMenu.addItem(timerListItem)
        
        NSApplication.shared.mainMenu = mainMenu
    }
    
    @objc private func showNewTimerWindow() {
        if newTimerWindowController == nil {
            newTimerWindowController = NewTimerWindowController()
        }
        newTimerWindowController?.onTimerCreated = { [weak self] description in
            self?.timerManager.startNewTimer(description: description)
            self?.timerListWindowController?.refresh()
        }
        newTimerWindowController?.showWindow(nil)
        NSApp.activate(ignoringOtherApps: true)
    }
    
    @objc private func showTimerListWindow() {
        if timerListWindowController == nil {
            timerListWindowController = TimerListWindowController(timerManager: timerManager)
            timerListWindowController?.onNewTimer = { [weak self] in
                self?.showNewTimerWindow()
            }
        }
        timerListWindowController?.showWindow(nil)
        NSApp.activate(ignoringOtherApps: true)
    }
    
    @objc private func showPreferencesWindow() {
        if preferencesWindowController == nil {
            preferencesWindowController = PreferencesWindowController()
        }
        preferencesWindowController?.showWindow(nil)
        NSApp.activate(ignoringOtherApps: true)
    }
}
