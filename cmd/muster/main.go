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
	"github.com/muster/cli/internal/standup"
	"github.com/muster/cli/internal/user"
	"github.com/muster/cli/internal/ui"
)

func main() {
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

	// Start the shell
	shellModel := ui.NewShellModel(authService, standupService, attendanceService, invitationService, leaveService, userService, cfg)
	p := tea.NewProgram(shellModel)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
