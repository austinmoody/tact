package ui

import (
	tea "charm.land/bubbletea/v2"

	"tact-tui/api"
	"tact-tui/model"
)

// Screen represents the current view
type Screen int

const (
	ScreenHome Screen = iota
	ScreenTimeCodes
	ScreenWorkTypes
	ScreenProjects
)

// Modal represents the current modal overlay
type Modal int

const (
	ModalNone Modal = iota
	ModalNewEntry
	ModalEntryDetail
	ModalMenu
	ModalTimeCodeAdd
	ModalTimeCodeEdit
	ModalWorkTypeAdd
	ModalWorkTypeEdit
	ModalProjectAdd
	ModalProjectEdit
	ModalContextList
	ModalContextAdd
	ModalContextEdit
)

// App is the root model that manages screens and modals
type App struct {
	client *api.Client
	width  int
	height int

	// Current screen and modal
	screen Screen
	modal  Modal

	// Screen models
	home      *Home
	timeCodes *TimeCodesScreen
	workTypes *WorkTypesScreen
	projects  *ProjectsScreen

	// Modal models
	entryInput      *EntryInputModal
	entryDetail     *EntryDetailModal
	menu            *MenuModal
	timeCodeEdit    *TimeCodeEditModal
	workTypeEdit    *WorkTypeEditModal
	projectEdit     *ProjectEditModal
	contextList     *ContextListModal
	contextEdit     *ContextEditModal

	// Shared data for modals
	selectedEntry    *model.Entry
	selectedTimeCode *model.TimeCode
	selectedWorkType *model.WorkType
	selectedProject  *model.Project
	contextOwner     ContextOwner
}

func NewApp(client *api.Client) *App {
	return &App{
		client:    client,
		screen:    ScreenHome,
		modal:     ModalNone,
		home:      NewHome(client),
		timeCodes: NewTimeCodesScreen(client),
		workTypes: NewWorkTypesScreen(client),
		projects:  NewProjectsScreen(client),
		menu:      NewMenuModal(),
	}
}

func (a *App) Init() tea.Cmd {
	return a.home.Init()
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		// Propagate to current screen
		switch a.screen {
		case ScreenHome:
			a.home.SetSize(msg.Width, msg.Height)
		case ScreenTimeCodes:
			a.timeCodes.SetSize(msg.Width, msg.Height)
		case ScreenWorkTypes:
			a.workTypes.SetSize(msg.Width, msg.Height)
		case ScreenProjects:
			a.projects.SetSize(msg.Width, msg.Height)
		}
		return a, nil

	case tea.KeyPressMsg:
		// Handle modal input first
		if a.modal != ModalNone {
			return a.updateModal(msg)
		}
		// Handle screen input
		return a.updateScreen(msg)

	case tea.PasteMsg:
		// Route paste events to modal if one is open
		if a.modal != ModalNone {
			return a.updateModalPaste(msg)
		}

	// Handle modal result messages
	case EntryCreatedMsg:
		a.modal = ModalNone
		a.entryInput = nil
		return a, a.home.Refresh()

	case EntryReparseMsg:
		a.modal = ModalNone
		a.entryDetail = nil
		return a, a.home.Refresh()

	case ModalCloseMsg:
		a.modal = ModalNone
		a.entryInput = nil
		a.entryDetail = nil
		a.timeCodeEdit = nil
		a.workTypeEdit = nil
		a.projectEdit = nil
		a.contextList = nil
		a.contextEdit = nil
		return a, nil

	case MenuSelectMsg:
		a.modal = ModalNone
		switch msg.Selection {
		case "projects":
			a.screen = ScreenProjects
			a.projects.SetSize(a.width, a.height)
			return a, a.projects.Init()
		case "timecodes":
			a.screen = ScreenTimeCodes
			a.timeCodes.SetSize(a.width, a.height)
			return a, a.timeCodes.Init()
		case "worktypes":
			a.screen = ScreenWorkTypes
			a.workTypes.SetSize(a.width, a.height)
			return a, a.workTypes.Init()
		}
		return a, nil

	case NavigateHomeMsg:
		a.screen = ScreenHome
		return a, a.home.Refresh()

	case TimeCodeCreatedMsg, TimeCodeUpdatedMsg, TimeCodeDeletedMsg:
		a.modal = ModalNone
		a.timeCodeEdit = nil
		return a, a.timeCodes.Refresh()

	case WorkTypeCreatedMsg, WorkTypeUpdatedMsg, WorkTypeDeletedMsg:
		a.modal = ModalNone
		a.workTypeEdit = nil
		return a, a.workTypes.Refresh()

	case OpenTimeCodeAddMsg:
		a.timeCodeEdit = NewTimeCodeEditModal(a.client, nil, a.width)
		a.modal = ModalTimeCodeAdd
		return a, a.timeCodeEdit.Init()

	case OpenTimeCodeEditMsg:
		a.selectedTimeCode = msg.TimeCode
		a.timeCodeEdit = NewTimeCodeEditModal(a.client, msg.TimeCode, a.width)
		a.modal = ModalTimeCodeEdit
		return a, a.timeCodeEdit.Init()

	case OpenWorkTypeAddMsg:
		a.workTypeEdit = NewWorkTypeEditModal(a.client, nil, a.width)
		a.modal = ModalWorkTypeAdd
		return a, a.workTypeEdit.Init()

	case OpenWorkTypeEditMsg:
		a.selectedWorkType = msg.WorkType
		a.workTypeEdit = NewWorkTypeEditModal(a.client, msg.WorkType, a.width)
		a.modal = ModalWorkTypeEdit
		return a, a.workTypeEdit.Init()

	case ProjectCreatedMsg, ProjectUpdatedMsg, ProjectDeletedMsg:
		a.modal = ModalNone
		a.projectEdit = nil
		return a, a.projects.Refresh()

	case OpenProjectAddMsg:
		a.projectEdit = NewProjectEditModal(a.client, nil, a.width)
		a.modal = ModalProjectAdd
		return a, a.projectEdit.Init()

	case OpenProjectEditMsg:
		a.selectedProject = msg.Project
		a.projectEdit = NewProjectEditModal(a.client, msg.Project, a.width)
		a.modal = ModalProjectEdit
		return a, a.projectEdit.Init()

	case OpenProjectContextMsg:
		a.contextOwner = ContextOwner{ProjectID: &msg.Project.ID, Name: msg.Project.Name}
		a.contextList = NewContextListModal(a.client, a.contextOwner, a.width)
		a.modal = ModalContextList
		return a, a.contextList.Init()

	case OpenTimeCodeContextMsg:
		a.contextOwner = ContextOwner{TimeCodeID: &msg.TimeCode.ID, Name: msg.TimeCode.ID}
		a.contextList = NewContextListModal(a.client, a.contextOwner, a.width)
		a.modal = ModalContextList
		return a, a.contextList.Init()

	case ContextCreatedMsg, ContextUpdatedMsg:
		// After creating/updating context, go back to context list
		a.contextEdit = nil
		a.contextList = NewContextListModal(a.client, a.contextOwner, a.width)
		a.modal = ModalContextList
		return a, a.contextList.Init()

	case ContextDeletedMsg:
		// Refresh context list after delete
		if a.contextList != nil {
			return a, a.contextList.Refresh()
		}
		return a, nil

	case OpenContextAddMsg:
		a.contextOwner = msg.Owner
		a.contextEdit = NewContextEditModal(a.client, a.contextOwner, nil, a.width, a.height)
		a.modal = ModalContextAdd
		return a, a.contextEdit.Init()

	case OpenContextEditMsg:
		a.contextOwner = msg.Owner
		a.contextEdit = NewContextEditModal(a.client, a.contextOwner, msg.Context, a.width, a.height)
		a.modal = ModalContextEdit
		return a, a.contextEdit.Init()
	}

	// Propagate other messages to current screen
	return a.propagateToScreen(msg)
}

func (a *App) updateModalPaste(msg tea.PasteMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch a.modal {
	case ModalNewEntry:
		if a.entryInput != nil {
			_, cmd = a.entryInput.Update(msg)
		}
	case ModalTimeCodeAdd, ModalTimeCodeEdit:
		if a.timeCodeEdit != nil {
			_, cmd = a.timeCodeEdit.Update(msg)
		}
	case ModalWorkTypeAdd, ModalWorkTypeEdit:
		if a.workTypeEdit != nil {
			_, cmd = a.workTypeEdit.Update(msg)
		}
	case ModalProjectAdd, ModalProjectEdit:
		if a.projectEdit != nil {
			_, cmd = a.projectEdit.Update(msg)
		}
	case ModalContextAdd, ModalContextEdit:
		if a.contextEdit != nil {
			_, cmd = a.contextEdit.Update(msg)
		}
	}

	return a, cmd
}

func (a *App) updateModal(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch a.modal {
	case ModalNewEntry:
		if a.entryInput != nil {
			_, cmd = a.entryInput.Update(msg)
		}
	case ModalEntryDetail:
		if a.entryDetail != nil {
			_, cmd = a.entryDetail.Update(msg)
		}
	case ModalMenu:
		_, cmd = a.menu.Update(msg)
	case ModalTimeCodeAdd, ModalTimeCodeEdit:
		if a.timeCodeEdit != nil {
			_, cmd = a.timeCodeEdit.Update(msg)
		}
	case ModalWorkTypeAdd, ModalWorkTypeEdit:
		if a.workTypeEdit != nil {
			_, cmd = a.workTypeEdit.Update(msg)
		}
	case ModalProjectAdd, ModalProjectEdit:
		if a.projectEdit != nil {
			_, cmd = a.projectEdit.Update(msg)
		}
	case ModalContextList:
		if a.contextList != nil {
			_, cmd = a.contextList.Update(msg)
		}
	case ModalContextAdd, ModalContextEdit:
		if a.contextEdit != nil {
			_, cmd = a.contextEdit.Update(msg)
		}
	}

	return a, cmd
}

func (a *App) updateScreen(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	// Global quit
	if matchesKey(msg, keys.Quit) {
		return a, tea.Quit
	}

	switch a.screen {
	case ScreenHome:
		return a.updateHomeScreen(msg)
	case ScreenTimeCodes:
		return a.updateTimeCodesScreen(msg)
	case ScreenWorkTypes:
		return a.updateWorkTypesScreen(msg)
	case ScreenProjects:
		return a.updateProjectsScreen(msg)
	}

	return a, nil
}

func (a *App) updateHomeScreen(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case matchesKey(msg, keys.NewEntry):
		a.entryInput = NewEntryInputModal(a.client, a.width)
		a.modal = ModalNewEntry
		return a, a.entryInput.Init()

	case matchesKey(msg, keys.Menu):
		a.modal = ModalMenu
		return a, nil

	case matchesKey(msg, keys.Enter):
		if entry := a.home.SelectedEntry(); entry != nil {
			a.selectedEntry = entry
			a.entryDetail = NewEntryDetailModal(a.client, entry)
			a.modal = ModalEntryDetail
		}
		return a, nil

	default:
		_, cmd := a.home.Update(msg)
		return a, cmd
	}
}

func (a *App) updateTimeCodesScreen(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case matchesKey(msg, keys.Escape):
		a.screen = ScreenHome
		return a, a.home.Refresh()
	default:
		_, cmd := a.timeCodes.Update(msg)
		return a, cmd
	}
}

func (a *App) updateWorkTypesScreen(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case matchesKey(msg, keys.Escape):
		a.screen = ScreenHome
		return a, a.home.Refresh()
	default:
		_, cmd := a.workTypes.Update(msg)
		return a, cmd
	}
}

func (a *App) updateProjectsScreen(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case matchesKey(msg, keys.Escape):
		a.screen = ScreenHome
		return a, a.home.Refresh()
	default:
		_, cmd := a.projects.Update(msg)
		return a, cmd
	}
}

func (a *App) propagateToScreen(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch a.screen {
	case ScreenHome:
		_, cmd = a.home.Update(msg)
	case ScreenTimeCodes:
		_, cmd = a.timeCodes.Update(msg)
	case ScreenWorkTypes:
		_, cmd = a.workTypes.Update(msg)
	case ScreenProjects:
		_, cmd = a.projects.Update(msg)
	}

	return a, cmd
}

func (a *App) View() tea.View {
	// Render current screen
	var screenView string
	switch a.screen {
	case ScreenHome:
		screenView = a.home.View()
	case ScreenTimeCodes:
		screenView = a.timeCodes.View()
	case ScreenWorkTypes:
		screenView = a.workTypes.View()
	case ScreenProjects:
		screenView = a.projects.View()
	}

	// Build the view with alt screen enabled
	var v tea.View
	if a.modal != ModalNone {
		v = tea.NewView(a.renderWithModal(screenView))
	} else {
		v = tea.NewView(screenView)
	}
	v.AltScreen = true
	return v
}

func (a *App) renderWithModal(screenView string) string {
	var modalView string

	switch a.modal {
	case ModalNewEntry:
		if a.entryInput != nil {
			modalView = a.entryInput.View()
		}
	case ModalEntryDetail:
		if a.entryDetail != nil {
			modalView = a.entryDetail.View()
		}
	case ModalMenu:
		modalView = a.menu.View()
	case ModalTimeCodeAdd, ModalTimeCodeEdit:
		if a.timeCodeEdit != nil {
			modalView = a.timeCodeEdit.View()
		}
	case ModalWorkTypeAdd, ModalWorkTypeEdit:
		if a.workTypeEdit != nil {
			modalView = a.workTypeEdit.View()
		}
	case ModalProjectAdd, ModalProjectEdit:
		if a.projectEdit != nil {
			modalView = a.projectEdit.View()
		}
	case ModalContextList:
		if a.contextList != nil {
			modalView = a.contextList.View()
		}
	case ModalContextAdd, ModalContextEdit:
		if a.contextEdit != nil {
			modalView = a.contextEdit.View()
		}
	}

	return renderModalOverlay(screenView, modalView, a.width, a.height)
}

// Messages for communication between components

type EntryCreatedMsg struct{}
type EntryReparseMsg struct{}
type ModalCloseMsg struct{}
type NavigateHomeMsg struct{}

type MenuSelectMsg struct {
	Selection string
}

type TimeCodeCreatedMsg struct{}
type TimeCodeUpdatedMsg struct{}
type TimeCodeDeletedMsg struct{}

type WorkTypeCreatedMsg struct{}
type WorkTypeUpdatedMsg struct{}
type WorkTypeDeletedMsg struct{}

type OpenTimeCodeAddMsg struct{}
type OpenTimeCodeEditMsg struct {
	TimeCode *model.TimeCode
}

type OpenWorkTypeAddMsg struct{}
type OpenWorkTypeEditMsg struct {
	WorkType *model.WorkType
}

type ProjectCreatedMsg struct{}
type ProjectUpdatedMsg struct{}
type ProjectDeletedMsg struct{}

type OpenProjectAddMsg struct{}
type OpenProjectEditMsg struct {
	Project *model.Project
}

type OpenProjectContextMsg struct {
	Project *model.Project
}

type OpenTimeCodeContextMsg struct {
	TimeCode *model.TimeCode
}

type ContextCreatedMsg struct{}
type ContextUpdatedMsg struct{}
type ContextDeletedMsg struct{}

type OpenContextAddMsg struct {
	Owner ContextOwner
}

type OpenContextEditMsg struct {
	Context *model.ContextDocument
	Owner   ContextOwner
}
