package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muster/cli/internal/api"
	"github.com/muster/cli/internal/attendance"
	"github.com/muster/cli/internal/auth"
	"github.com/muster/cli/internal/config"
	"github.com/muster/cli/internal/invitation"
	"github.com/muster/cli/internal/leave"
	"github.com/muster/cli/internal/organization"
	"github.com/muster/cli/internal/project"
	"github.com/muster/cli/internal/standup"
	"github.com/muster/cli/internal/user"
	"github.com/muster/cli/internal/ui"
)

var (
	version        = "dev"
	defaultBaseURL = "http://localhost:3000"
)

func main() {
	// Set build-time base URL before loading config
	if defaultBaseURL != "" {
		config.DefaultBaseURL = defaultBaseURL
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Create API client
	apiClient := api.NewClientFromConfig(cfg)

	// Create services
	authService := auth.NewService(apiClient, cfg)
	standupService := standup.NewService(apiClient)
	attendanceService := attendance.NewService(apiClient)
	invitationService := invitation.NewService(apiClient)
	leaveService := leave.NewService(apiClient)
	userService := user.NewService(apiClient)
	orgService := organization.NewService(apiClient)
	projectService := project.NewService(apiClient)

	// Start the shell
	shellModel := ui.NewShellModel(authService, standupService, attendanceService, invitationService, leaveService, userService, orgService, projectService, cfg)
	p := tea.NewProgram(shellModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
