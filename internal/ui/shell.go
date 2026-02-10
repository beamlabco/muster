package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muster/cli/internal/attendance"
	"github.com/muster/cli/internal/auth"
	"github.com/muster/cli/internal/config"
	"github.com/muster/cli/internal/invitation"
	"github.com/muster/cli/internal/leave"
	"github.com/muster/cli/internal/standup"
	"github.com/muster/cli/internal/user"
)

// View states
type viewState int

const (
	viewShell viewState = iota
	viewLogin
	viewRegister
	viewStandupSubmit
	viewStandupToday
	viewStandupHistory
	viewAttendanceMark
	viewAttendanceCheckIn
	viewAttendanceCheckOut
	viewAttendanceToday
	viewAttendanceHistory
	viewInvitationCreate
	viewJoin
	viewLeaveRequest
	viewLeaveList
	viewLeaveReview
	viewLeaveCancel
	viewTeamList
	viewRoleUpdate
)

// ShellModel is the main REPL shell
type ShellModel struct {
	// Services
	authService       *auth.Service
	standupService    *standup.Service
	attendanceService *attendance.Service
	invitationService *invitation.Service
	leaveService      *leave.Service
	userService       *user.Service
	config            *config.Config

	// Shell components
	commandInput *CommandInput
	registry     *CommandRegistry

	// Current view state
	currentView viewState

	// Sub-models
	loginModel             LoginModel
	registerModel          RegisterModel
	standupSubmitModel     StandupSubmitModel
	standupTodayModel      StandupTodayModel
	standupHistoryModel    StandupHistoryModel
	attendanceMarkModel      AttendanceMarkModel
	attendanceCheckInModel   AttendanceCheckInModel
	attendanceCheckOutModel  AttendanceCheckOutModel
	attendanceTodayModel     AttendanceTodayModel
	attendanceHistoryModel   AttendanceHistoryModel
	invitationCreateModel    InvitationCreateModel
	joinModel                JoinModel
	leaveRequestModel        LeaveRequestModel
	leaveListModel           LeaveListModel
	leaveReviewModel         LeaveReviewModel
	leaveCancelModel         LeaveCancelModel
	teamListModel            TeamListModel
	roleUpdateModel          RoleUpdateModel

	// Output area
	output     string
	outputType string // "success", "error", "info"

	// State
	quitting bool
	width    int
	height   int
}

// NewShellModel creates a new shell model
func NewShellModel(authService *auth.Service, standupService *standup.Service, attendanceService *attendance.Service, invitationService *invitation.Service, leaveService *leave.Service, userService *user.Service, cfg *config.Config) ShellModel {
	registry := NewCommandRegistry()
	isAuth := authService.IsAuthenticated()

	return ShellModel{
		authService:       authService,
		standupService:    standupService,
		attendanceService: attendanceService,
		invitationService: invitationService,
		leaveService:      leaveService,
		userService:       userService,
		config:            cfg,
		commandInput:      NewCommandInput(registry, isAuth),
		registry:          registry,
		currentView:       viewShell,
		width:             80,
		height:            24,
	}
}

// Init initializes the shell
func (m ShellModel) Init() tea.Cmd {
	return m.commandInput.Focus()
}

// Update handles messages
func (m ShellModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle window size
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = msg.Width
		m.height = msg.Height
	}

	// Route to sub-model if active
	switch m.currentView {
	case viewLogin:
		return m.updateLogin(msg)
	case viewRegister:
		return m.updateRegister(msg)
	case viewStandupSubmit:
		return m.updateStandupSubmit(msg)
	case viewStandupToday:
		return m.updateStandupToday(msg)
	case viewStandupHistory:
		return m.updateStandupHistory(msg)
	case viewAttendanceMark:
		return m.updateAttendanceMark(msg)
	case viewAttendanceCheckIn:
		return m.updateAttendanceCheckIn(msg)
	case viewAttendanceCheckOut:
		return m.updateAttendanceCheckOut(msg)
	case viewAttendanceToday:
		return m.updateAttendanceToday(msg)
	case viewAttendanceHistory:
		return m.updateAttendanceHistory(msg)
	case viewInvitationCreate:
		return m.updateInvitationCreate(msg)
	case viewJoin:
		return m.updateJoin(msg)
	case viewLeaveRequest:
		return m.updateLeaveRequest(msg)
	case viewLeaveList:
		return m.updateLeaveList(msg)
	case viewLeaveReview:
		return m.updateLeaveReview(msg)
	case viewLeaveCancel:
		return m.updateLeaveCancel(msg)
	case viewTeamList:
		return m.updateTeamList(msg)
	case viewRoleUpdate:
		return m.updateRoleUpdate(msg)
	}

	// Handle shell input
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "ctrl+l":
			// Clear screen
			m.output = ""
			m.outputType = ""
			return m, nil

		case "enter":
			// If dropdown is visible and has a selection, execute that command directly
			if selected := m.commandInput.SelectedCommand(); selected != nil {
				m.commandInput.Reset()
				// Check auth requirement
				if selected.RequiresAuth && !m.authService.IsAuthenticated() {
					m.output = fmt.Sprintf("You must be logged in to use /%s. Use /login first.", selected.Name)
					m.outputType = "error"
					return m, m.commandInput.Focus()
				}
				return m.handleCommand(selected.Name)
			}
			// Otherwise execute the typed command
			return m.executeCommand()
		}
	}

	// Update command input
	var cmd tea.Cmd
	m.commandInput, cmd = m.commandInput.Update(msg)
	return m, cmd
}

// executeCommand executes the current command
func (m ShellModel) executeCommand() (tea.Model, tea.Cmd) {
	input := strings.TrimSpace(m.commandInput.Value())
	m.commandInput.Reset()

	if input == "" {
		return m, m.commandInput.Focus()
	}

	// Must start with /
	if !strings.HasPrefix(input, "/") {
		m.output = "Commands must start with /"
		m.outputType = "error"
		return m, m.commandInput.Focus()
	}

	// Check for empty command (just "/")
	if input == "/" {
		m.output = "Type a command after /. Use /help for available commands."
		m.outputType = "error"
		return m, m.commandInput.Focus()
	}

	// Find command
	cmd := m.registry.FindExact(input)
	if cmd == nil {
		m.output = fmt.Sprintf("Unknown command: %s. Type /help for available commands.", input)
		m.outputType = "error"
		return m, m.commandInput.Focus()
	}

	// Check auth requirement
	if cmd.RequiresAuth && !m.authService.IsAuthenticated() {
		m.output = fmt.Sprintf("You must be logged in to use /%s. Use /login first.", cmd.Name)
		m.outputType = "error"
		return m, m.commandInput.Focus()
	}

	// Execute command
	return m.handleCommand(cmd.Name)
}

// handleCommand handles command execution
func (m ShellModel) handleCommand(cmdName string) (tea.Model, tea.Cmd) {
	switch cmdName {
	// Auth commands
	case "login":
		m.currentView = viewLogin
		m.loginModel = NewLoginModel(m.authService, nil)
		return m, m.loginModel.Init()

	case "signup":
		m.currentView = viewRegister
		m.registerModel = NewRegisterModel(m.authService, nil)
		return m, m.registerModel.Init()

	case "join":
		m.currentView = viewJoin
		m.joinModel = NewJoinModel(m.authService, nil)
		return m, m.joinModel.Init()

	case "invite":
		m.currentView = viewInvitationCreate
		m.invitationCreateModel = NewInvitationCreateModel(m.invitationService)
		return m, m.invitationCreateModel.Init()

	case "logout":
		if err := m.authService.Logout(); err != nil {
			m.output = fmt.Sprintf("Logout failed: %s", err.Error())
			m.outputType = "error"
		} else {
			m.output = "Logged out successfully"
			m.outputType = "success"
			m.commandInput.SetAuth(false)
		}
		return m, m.commandInput.Focus()

	// Standup commands
	case "standup":
		m.currentView = viewStandupSubmit
		m.standupSubmitModel = NewStandupSubmitModel(m.standupService, nil)
		return m, m.standupSubmitModel.Init()

	case "standup today":
		m.currentView = viewStandupToday
		m.standupTodayModel = NewStandupTodayModel(m.standupService, nil)
		return m, m.standupTodayModel.Init()

	case "standup history":
		m.currentView = viewStandupHistory
		m.standupHistoryModel = NewStandupHistoryModel(m.standupService, nil)
		return m, m.standupHistoryModel.Init()

	// Attendance commands
	case "checkin":
		m.currentView = viewAttendanceCheckIn
		m.attendanceCheckInModel = NewAttendanceCheckInModel(m.attendanceService)
		return m, m.attendanceCheckInModel.Init()

	case "checkout":
		m.currentView = viewAttendanceCheckOut
		m.attendanceCheckOutModel = NewAttendanceCheckOutModel(m.attendanceService)
		return m, m.attendanceCheckOutModel.Init()

	case "attendance":
		m.currentView = viewAttendanceMark
		m.attendanceMarkModel = NewAttendanceMarkModel(m.attendanceService)
		return m, m.attendanceMarkModel.Init()

	case "attendance today":
		m.currentView = viewAttendanceToday
		m.attendanceTodayModel = NewAttendanceTodayModel(m.attendanceService)
		return m, m.attendanceTodayModel.Init()

	case "attendance history":
		m.currentView = viewAttendanceHistory
		m.attendanceHistoryModel = NewAttendanceHistoryModel(m.attendanceService)
		return m, m.attendanceHistoryModel.Init()

	// Leave commands
	case "leave":
		m.currentView = viewLeaveRequest
		m.leaveRequestModel = NewLeaveRequestModel(m.leaveService)
		return m, m.leaveRequestModel.Init()

	case "leave list":
		m.currentView = viewLeaveList
		m.leaveListModel = NewLeaveListModel(m.leaveService)
		return m, m.leaveListModel.Init()

	case "leave review":
		m.currentView = viewLeaveReview
		m.leaveReviewModel = NewLeaveReviewModel(m.leaveService)
		return m, m.leaveReviewModel.Init()

	case "leave cancel":
		userID := 0
		if m.config.User != nil {
			userID = m.config.User.ID
		}
		m.currentView = viewLeaveCancel
		m.leaveCancelModel = NewLeaveCancelModel(m.leaveService, userID)
		return m, m.leaveCancelModel.Init()

	// User management commands
	case "team":
		m.currentView = viewTeamList
		m.teamListModel = NewTeamListModel(m.userService)
		return m, m.teamListModel.Init()

	case "role":
		m.currentView = viewRoleUpdate
		m.roleUpdateModel = NewRoleUpdateModel(m.userService)
		return m, m.roleUpdateModel.Init()

	// Utility commands
	case "help":
		m.output = m.renderHelp()
		m.outputType = "info"
		return m, m.commandInput.Focus()

	case "whoami":
		m.output = m.renderWhoami()
		m.outputType = "info"
		return m, m.commandInput.Focus()

	case "clear":
		m.output = ""
		m.outputType = ""
		return m, m.commandInput.Focus()

	case "quit":
		m.quitting = true
		return m, tea.Quit

	default:
		m.output = fmt.Sprintf("Command not implemented: %s", cmdName)
		m.outputType = "error"
		return m, m.commandInput.Focus()
	}
}

// Sub-model update handlers

func (m ShellModel) updateLogin(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Check for escape or success
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			m.currentView = viewShell
			return m, m.commandInput.Focus()
		}
	case loginSuccessMsg:
		m.currentView = viewShell
		m.output = "Login successful! Welcome back."
		m.outputType = "success"
		m.commandInput.SetAuth(true)
		return m, m.commandInput.Focus()
	case loginErrorMsg:
		// Let the login model handle displaying the error
	}

	// Update login model
	newModel, cmd := m.loginModel.Update(msg)
	m.loginModel = newModel.(LoginModel)

	// Check if quitting (user pressed ctrl+c)
	if m.loginModel.quitting {
		m.currentView = viewShell
		return m, m.commandInput.Focus()
	}

	return m, cmd
}

func (m ShellModel) updateRegister(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			m.currentView = viewShell
			return m, m.commandInput.Focus()
		}
	case registerSuccessMsg:
		m.currentView = viewShell
		m.output = fmt.Sprintf("Welcome to Muster, %s! Your account has been created.", string(msg))
		m.outputType = "success"
		m.commandInput.SetAuth(true)
		return m, m.commandInput.Focus()
	}

	newModel, cmd := m.registerModel.Update(msg)
	m.registerModel = newModel.(RegisterModel)

	if m.registerModel.quitting {
		m.currentView = viewShell
		return m, m.commandInput.Focus()
	}

	return m, cmd
}

func (m ShellModel) updateStandupSubmit(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			m.currentView = viewShell
			return m, m.commandInput.Focus()
		}
	case standupSubmitSuccessMsg:
		m.currentView = viewShell
		m.output = "Standup submitted successfully!"
		m.outputType = "success"
		return m, m.commandInput.Focus()
	}

	newModel, cmd := m.standupSubmitModel.Update(msg)
	m.standupSubmitModel = newModel.(StandupSubmitModel)

	if m.standupSubmitModel.shouldGoBack {
		m.currentView = viewShell
		return m, m.commandInput.Focus()
	}

	return m, cmd
}

func (m ShellModel) updateStandupToday(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" || msg.String() == "q" {
			m.currentView = viewShell
			return m, m.commandInput.Focus()
		}
	}

	newModel, cmd := m.standupTodayModel.Update(msg)
	m.standupTodayModel = newModel.(StandupTodayModel)

	if m.standupTodayModel.shouldGoBack {
		m.currentView = viewShell
		return m, m.commandInput.Focus()
	}

	return m, cmd
}

func (m ShellModel) updateStandupHistory(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" || msg.String() == "q" {
			m.currentView = viewShell
			return m, m.commandInput.Focus()
		}
	}

	newModel, cmd := m.standupHistoryModel.Update(msg)
	m.standupHistoryModel = newModel.(StandupHistoryModel)

	if m.standupHistoryModel.shouldGoBack {
		m.currentView = viewShell
		return m, m.commandInput.Focus()
	}

	return m, cmd
}

func (m ShellModel) updateAttendanceMark(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			m.currentView = viewShell
			return m, m.commandInput.Focus()
		}
	case attendanceMarkSuccessMsg:
		m.currentView = viewShell
		m.output = fmt.Sprintf("Attendance marked as %s!", msg.status)
		m.outputType = "success"
		return m, m.commandInput.Focus()
	}

	newModel, cmd := m.attendanceMarkModel.Update(msg)
	m.attendanceMarkModel = newModel.(AttendanceMarkModel)

	if m.attendanceMarkModel.shouldGoBack {
		m.currentView = viewShell
		return m, m.commandInput.Focus()
	}

	return m, cmd
}

func (m ShellModel) updateAttendanceCheckIn(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			m.currentView = viewShell
			return m, m.commandInput.Focus()
		}
	case checkInSuccessMsg:
		m.currentView = viewShell
		m.output = fmt.Sprintf("Checked in as %s at %s", msg.status, msg.time)
		m.outputType = "success"
		return m, m.commandInput.Focus()
	}

	newModel, cmd := m.attendanceCheckInModel.Update(msg)
	m.attendanceCheckInModel = newModel.(AttendanceCheckInModel)

	if m.attendanceCheckInModel.shouldGoBack {
		m.currentView = viewShell
		return m, m.commandInput.Focus()
	}

	return m, cmd
}

func (m ShellModel) updateAttendanceCheckOut(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			m.currentView = viewShell
			return m, m.commandInput.Focus()
		}
	case checkOutSuccessMsg:
		m.currentView = viewShell
		m.output = fmt.Sprintf("Checked out at %s", msg.time)
		m.outputType = "success"
		return m, m.commandInput.Focus()
	}

	newModel, cmd := m.attendanceCheckOutModel.Update(msg)
	m.attendanceCheckOutModel = newModel.(AttendanceCheckOutModel)

	if m.attendanceCheckOutModel.shouldGoBack {
		m.currentView = viewShell
		return m, m.commandInput.Focus()
	}

	return m, cmd
}

func (m ShellModel) updateAttendanceToday(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" || msg.String() == "q" {
			m.currentView = viewShell
			return m, m.commandInput.Focus()
		}
	}

	newModel, cmd := m.attendanceTodayModel.Update(msg)
	m.attendanceTodayModel = newModel.(AttendanceTodayModel)

	if m.attendanceTodayModel.shouldGoBack {
		m.currentView = viewShell
		return m, m.commandInput.Focus()
	}

	return m, cmd
}

func (m ShellModel) updateAttendanceHistory(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" || msg.String() == "q" {
			m.currentView = viewShell
			return m, m.commandInput.Focus()
		}
	}

	newModel, cmd := m.attendanceHistoryModel.Update(msg)
	m.attendanceHistoryModel = newModel.(AttendanceHistoryModel)

	if m.attendanceHistoryModel.shouldGoBack {
		m.currentView = viewShell
		return m, m.commandInput.Focus()
	}

	return m, cmd
}

func (m ShellModel) updateInvitationCreate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			m.currentView = viewShell
			return m, m.commandInput.Focus()
		}
	case invitationCreateSuccessMsg:
		m.currentView = viewShell
		m.output = fmt.Sprintf("Invitation sent to %s (token: %s)", msg.email, msg.token)
		m.outputType = "success"
		return m, m.commandInput.Focus()
	}

	newModel, cmd := m.invitationCreateModel.Update(msg)
	m.invitationCreateModel = newModel.(InvitationCreateModel)

	if m.invitationCreateModel.shouldGoBack {
		m.currentView = viewShell
		return m, m.commandInput.Focus()
	}

	return m, cmd
}

func (m ShellModel) updateJoin(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			m.currentView = viewShell
			return m, m.commandInput.Focus()
		}
	case joinSuccessMsg:
		m.currentView = viewShell
		m.output = fmt.Sprintf("Welcome to %s, %s! You've joined the organization.", msg.org, msg.name)
		m.outputType = "success"
		m.commandInput.SetAuth(true)
		return m, m.commandInput.Focus()
	}

	newModel, cmd := m.joinModel.Update(msg)
	m.joinModel = newModel.(JoinModel)

	if m.joinModel.quitting {
		m.currentView = viewShell
		return m, m.commandInput.Focus()
	}

	return m, cmd
}

func (m ShellModel) updateLeaveRequest(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			m.currentView = viewShell
			return m, m.commandInput.Focus()
		}
	case leaveRequestSuccessMsg:
		m.currentView = viewShell
		m.output = fmt.Sprintf("Leave request created! (%s: %s to %s)", msg.leaveType, msg.startDate, msg.endDate)
		m.outputType = "success"
		return m, m.commandInput.Focus()
	}

	newModel, cmd := m.leaveRequestModel.Update(msg)
	m.leaveRequestModel = newModel.(LeaveRequestModel)

	if m.leaveRequestModel.shouldGoBack {
		m.currentView = viewShell
		return m, m.commandInput.Focus()
	}

	return m, cmd
}

func (m ShellModel) updateLeaveList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" || msg.String() == "q" {
			m.currentView = viewShell
			return m, m.commandInput.Focus()
		}
	}

	newModel, cmd := m.leaveListModel.Update(msg)
	m.leaveListModel = newModel.(LeaveListModel)

	if m.leaveListModel.shouldGoBack {
		m.currentView = viewShell
		return m, m.commandInput.Focus()
	}

	return m, cmd
}

func (m ShellModel) updateLeaveReview(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" || msg.String() == "q" {
			m.currentView = viewShell
			return m, m.commandInput.Focus()
		}
	}

	newModel, cmd := m.leaveReviewModel.Update(msg)
	m.leaveReviewModel = newModel.(LeaveReviewModel)

	if m.leaveReviewModel.shouldGoBack {
		m.currentView = viewShell
		return m, m.commandInput.Focus()
	}

	return m, cmd
}

func (m ShellModel) updateLeaveCancel(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" || msg.String() == "q" {
			m.currentView = viewShell
			return m, m.commandInput.Focus()
		}
	}

	newModel, cmd := m.leaveCancelModel.Update(msg)
	m.leaveCancelModel = newModel.(LeaveCancelModel)

	if m.leaveCancelModel.shouldGoBack {
		m.currentView = viewShell
		return m, m.commandInput.Focus()
	}

	return m, cmd
}

func (m ShellModel) updateTeamList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" || msg.String() == "q" {
			m.currentView = viewShell
			return m, m.commandInput.Focus()
		}
	}

	newModel, cmd := m.teamListModel.Update(msg)
	m.teamListModel = newModel.(TeamListModel)

	if m.teamListModel.shouldGoBack {
		m.currentView = viewShell
		return m, m.commandInput.Focus()
	}

	return m, cmd
}

func (m ShellModel) updateRoleUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			m.currentView = viewShell
			return m, m.commandInput.Focus()
		}
	case roleUpdateSuccessMsg:
		m.currentView = viewShell
		m.output = fmt.Sprintf("Updated %s's role to %s", msg.name, msg.role)
		m.outputType = "success"
		return m, m.commandInput.Focus()
	}

	newModel, cmd := m.roleUpdateModel.Update(msg)
	m.roleUpdateModel = newModel.(RoleUpdateModel)

	if m.roleUpdateModel.shouldGoBack {
		m.currentView = viewShell
		return m, m.commandInput.Focus()
	}

	return m, cmd
}

// View renders the shell
func (m ShellModel) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	// Render sub-model if active
	switch m.currentView {
	case viewLogin:
		return m.loginModel.View()
	case viewRegister:
		return m.registerModel.View()
	case viewStandupSubmit:
		return m.standupSubmitModel.View()
	case viewStandupToday:
		return m.standupTodayModel.View()
	case viewStandupHistory:
		return m.standupHistoryModel.View()
	case viewAttendanceMark:
		return m.attendanceMarkModel.View()
	case viewAttendanceCheckIn:
		return m.attendanceCheckInModel.View()
	case viewAttendanceCheckOut:
		return m.attendanceCheckOutModel.View()
	case viewAttendanceToday:
		return m.attendanceTodayModel.View()
	case viewAttendanceHistory:
		return m.attendanceHistoryModel.View()
	case viewInvitationCreate:
		return m.invitationCreateModel.View()
	case viewJoin:
		return m.joinModel.View()
	case viewLeaveRequest:
		return m.leaveRequestModel.View()
	case viewLeaveList:
		return m.leaveListModel.View()
	case viewLeaveReview:
		return m.leaveReviewModel.View()
	case viewLeaveCancel:
		return m.leaveCancelModel.View()
	case viewTeamList:
		return m.teamListModel.View()
	case viewRoleUpdate:
		return m.roleUpdateModel.View()
	}

	// Render shell view
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(mutedColor).
		Width(m.width - 4).
		Padding(0, 1)

	headerText := "Muster CLI"
	if m.authService.IsAuthenticated() {
		user := m.authService.GetUser()
		org := m.authService.GetOrganization()
		if user != nil && org != nil {
			headerText = fmt.Sprintf("Muster CLI  •  %s @ %s", user.Name, org.Name)
		}
	}
	b.WriteString(headerStyle.Render(headerText))
	b.WriteString("\n\n")

	// Output area
	if m.output != "" {
		var outputStyle lipgloss.Style
		switch m.outputType {
		case "success":
			outputStyle = successStyle
		case "error":
			outputStyle = errorStyle
		default:
			outputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
		}
		b.WriteString(outputStyle.Render(m.output))
		b.WriteString("\n\n")
	}

	// Command input
	b.WriteString(m.commandInput.View())
	b.WriteString("\n\n")

	// Help footer
	footerStyle := lipgloss.NewStyle().Foreground(mutedColor)
	b.WriteString(footerStyle.Render("[Tab] Accept  [↑↓] Navigate  [Enter] Execute  [Ctrl+C] Quit"))

	return baseStyle.Render(b.String())
}

// Helper methods

func (m ShellModel) renderHelp() string {
	var b strings.Builder

	b.WriteString("Available Commands:\n\n")

	categories := m.registry.GetByCategory(m.authService.IsAuthenticated())
	categoryOrder := []string{"auth", "standup", "attendance", "leave", "admin", "utility"}
	categoryNames := map[string]string{
		"auth":       "Authentication",
		"standup":    "Standups",
		"attendance": "Attendance",
		"leave":      "Leaves",
		"admin":      "Admin",
		"utility":    "Utility",
	}

	for _, cat := range categoryOrder {
		cmds, ok := categories[cat]
		if !ok || len(cmds) == 0 {
			continue
		}

		b.WriteString(fmt.Sprintf("  %s\n", categoryNames[cat]))
		for _, cmd := range cmds {
			b.WriteString(fmt.Sprintf("    /%s - %s\n", cmd.Name, cmd.Description))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (m ShellModel) renderWhoami() string {
	if !m.authService.IsAuthenticated() {
		return "Not logged in"
	}

	user := m.authService.GetUser()
	org := m.authService.GetOrganization()

	if user == nil || org == nil {
		return "User info not available"
	}

	return fmt.Sprintf("User: %s (%s)\nEmail: %s\nRole: %s\nOrganization: %s",
		user.Name, user.Role, user.Email, user.Role, org.Name)
}

