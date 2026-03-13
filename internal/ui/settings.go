package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muster/cli/internal/api"
	"github.com/muster/cli/internal/organization"
)

const (
	settingsOrgName = iota
)

// SettingsModel represents the organization settings form
type SettingsModel struct {
	orgService   *organization.Service
	inputs       []textinput.Model
	focusIndex   int
	errorMsg     string
	successMsg   string
	loading      bool
	shouldGoBack bool
}

// NewSettingsModel creates a new settings model
func NewSettingsModel(orgService *organization.Service) SettingsModel {
	orgName := textinput.New()
	orgName.Placeholder = "Organization name"
	orgName.CharLimit = 255
	orgName.Width = 40

	return SettingsModel{
		orgService: orgService,
		inputs: []textinput.Model{
			orgName,
		},
		focusIndex: settingsOrgName,
		loading:    true,
	}
}

// Init loads current settings
func (m SettingsModel) Init() tea.Cmd {
	return m.fetchSettings()
}

// Update handles messages
func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			m.shouldGoBack = true
			return m, nil

		case "ctrl+s", "enter":
			if m.loading {
				return m, nil
			}
			return m, m.handleSubmit()
		}

	case settingsLoadedMsg:
		m.loading = false
		m.populateFromSettings(msg.settings)
		m.inputs[settingsOrgName].Focus()
		return m, nil

	case settingsUpdatedMsg:
		m.loading = false
		m.successMsg = "Settings updated successfully!"
		m.errorMsg = ""
		return m, nil

	case settingsErrorMsg:
		m.loading = false
		m.errorMsg = string(msg)
		return m, nil
	}

	// Update focused text input
	var cmd tea.Cmd
	m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
	return m, cmd
}

func (m *SettingsModel) populateFromSettings(s *api.OrgSettingsResponse) {
	m.inputs[settingsOrgName].SetValue(s.Name)
}

// View renders the settings form
func (m SettingsModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Organization Settings"))
	b.WriteString("\n\n")

	if m.loading && m.successMsg == "" && m.errorMsg == "" {
		b.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Loading settings..."))
		b.WriteString("\n")
		return baseStyle.Render(b.String())
	}

	// Organization name
	b.WriteString(labelStyle.Render("Organization name:"))
	b.WriteString("\n")
	if m.focusIndex == settingsOrgName {
		b.WriteString(focusedInputStyle.Render(m.inputs[settingsOrgName].View()))
	} else {
		b.WriteString(blurredInputStyle.Render(m.inputs[settingsOrgName].View()))
	}
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("Note: Summary and webhook settings have moved to project settings."))
	b.WriteString("\n")
	b.WriteString(labelStyle.Render("Use /project settings to configure per-project summaries."))
	b.WriteString("\n\n")

	if m.errorMsg != "" {
		b.WriteString(errorStyle.Render("Error: " + m.errorMsg))
		b.WriteString("\n")
	}

	if m.successMsg != "" {
		b.WriteString(successStyle.Render(m.successMsg))
		b.WriteString("\n")
	}

	if m.loading {
		b.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Saving settings..."))
		b.WriteString("\n")
	}

	if !m.loading {
		b.WriteString(helpStyle.Render("[Enter/Ctrl+S] Save  [Esc] Back"))
	}

	return baseStyle.Render(b.String())
}

func (m *SettingsModel) fetchSettings() tea.Cmd {
	return func() tea.Msg {
		settings, err := m.orgService.GetSettings()
		if err != nil {
			return settingsErrorMsg(err.Error())
		}
		return settingsLoadedMsg{settings: settings}
	}
}

func (m *SettingsModel) handleSubmit() tea.Cmd {
	orgName := strings.TrimSpace(m.inputs[settingsOrgName].Value())

	m.errorMsg = ""
	m.successMsg = ""
	m.loading = true

	return func() tea.Msg {
		req := &api.UpdateOrgSettingsRequest{}

		if orgName != "" {
			req.Name = &orgName
		}

		_, err := m.orgService.UpdateSettings(req)
		if err != nil {
			return settingsErrorMsg(err.Error())
		}
		return settingsUpdatedMsg{}
	}
}

// Message types
type settingsLoadedMsg struct {
	settings *api.OrgSettingsResponse
}
type settingsUpdatedMsg struct{}
type settingsErrorMsg string
