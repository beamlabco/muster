package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muster/cli/internal/api"
	"github.com/muster/cli/internal/project"
)

const (
	projSettingsName = iota
	projSettingsSummaryTime
	projSettingsTimezone
	projSettingsDiscordURL
	projSettingsEnabled
)

// ProjectSettingsModel represents the project settings form
type ProjectSettingsModel struct {
	projectService *project.Service
	projects       []*api.ProjectResponse
	selectedIdx    int
	inputs         []textinput.Model
	enabledToggle  bool
	focusIndex     int
	phase          string // "select" or "edit"
	errorMsg       string
	successMsg     string
	loading        bool
	shouldGoBack   bool
}

// NewProjectSettingsModel creates a new project settings model
func NewProjectSettingsModel(projectService *project.Service) ProjectSettingsModel {
	name := textinput.New()
	name.Placeholder = "Project name"
	name.CharLimit = 255
	name.Width = 40

	summaryTime := textinput.New()
	summaryTime.Placeholder = "09:00"
	summaryTime.CharLimit = 5
	summaryTime.Width = 20

	timezone := textinput.New()
	timezone.Placeholder = "America/New_York"
	timezone.CharLimit = 50
	timezone.Width = 30

	discordURL := textinput.New()
	discordURL.Placeholder = "https://discord.com/api/webhooks/..."
	discordURL.CharLimit = 500
	discordURL.Width = 55

	return ProjectSettingsModel{
		projectService: projectService,
		inputs: []textinput.Model{
			name,
			summaryTime,
			timezone,
			discordURL,
		},
		enabledToggle: true,
		focusIndex:    projSettingsName,
		phase:         "select",
		loading:       true,
	}
}

// Init loads projects
func (m ProjectSettingsModel) Init() tea.Cmd {
	return m.fetchProjects()
}

// Update handles messages
func (m ProjectSettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			if m.phase == "edit" {
				m.phase = "select"
				m.errorMsg = ""
				m.successMsg = ""
				return m, nil
			}
			m.shouldGoBack = true
			return m, nil

		case "up", "k":
			if m.phase == "select" && m.selectedIdx > 0 {
				m.selectedIdx--
			}
			return m, nil

		case "down", "j":
			if m.phase == "select" && m.selectedIdx < len(m.projects)-1 {
				m.selectedIdx++
			}
			return m, nil

		case "enter":
			if m.phase == "select" && len(m.projects) > 0 {
				m.phase = "edit"
				m.populateFromProject(m.projects[m.selectedIdx])
				m.inputs[projSettingsName].Focus()
				return m, nil
			}

		case "tab":
			if m.phase == "edit" {
				m.blurAll()
				m.focusIndex++
				if m.focusIndex > projSettingsEnabled {
					m.focusIndex = projSettingsName
				}
				m.focusCurrent()
				return m, nil
			}

		case "shift+tab":
			if m.phase == "edit" {
				m.blurAll()
				m.focusIndex--
				if m.focusIndex < projSettingsName {
					m.focusIndex = projSettingsEnabled
				}
				m.focusCurrent()
				return m, nil
			}

		case " ":
			if m.phase == "edit" && m.focusIndex == projSettingsEnabled {
				m.enabledToggle = !m.enabledToggle
				return m, nil
			}

		case "ctrl+s":
			if m.phase == "edit" && !m.loading {
				return m, m.handleSubmit()
			}
		}

	case projectSettingsProjectsMsg:
		m.loading = false
		m.projects = msg.projects
		return m, nil

	case projectSettingsUpdatedMsg:
		m.loading = false
		m.successMsg = "Project settings updated!"
		m.errorMsg = ""
		return m, nil

	case projectSettingsErrorMsg:
		m.loading = false
		m.errorMsg = string(msg)
		return m, nil
	}

	// Update focused text input in edit phase
	if m.phase == "edit" && m.focusIndex < projSettingsEnabled {
		var cmd tea.Cmd
		m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *ProjectSettingsModel) blurAll() {
	for i := range m.inputs {
		m.inputs[i].Blur()
	}
}

func (m *ProjectSettingsModel) focusCurrent() {
	if m.focusIndex < projSettingsEnabled {
		m.inputs[m.focusIndex].Focus()
	}
}

func (m *ProjectSettingsModel) populateFromProject(p *api.ProjectResponse) {
	m.inputs[projSettingsName].SetValue(p.Name)
	if p.SummaryTime != nil {
		m.inputs[projSettingsSummaryTime].SetValue(*p.SummaryTime)
	}
	if p.Timezone != nil {
		m.inputs[projSettingsTimezone].SetValue(*p.Timezone)
	}
	if p.DiscordWebhookURL != nil {
		m.inputs[projSettingsDiscordURL].SetValue(*p.DiscordWebhookURL)
	}
	m.enabledToggle = p.SummaryEnabled
}

// View renders the form
func (m ProjectSettingsModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Project Settings"))
	b.WriteString("\n\n")

	if m.loading && m.phase == "select" {
		b.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Loading projects..."))
		return baseStyle.Render(b.String())
	}

	if m.phase == "select" {
		b.WriteString(labelStyle.Render("Select a project to edit:"))
		b.WriteString("\n\n")

		for i, p := range m.projects {
			if i == m.selectedIdx {
				b.WriteString(fmt.Sprintf("  > %s\n", lipgloss.NewStyle().Bold(true).Foreground(primaryColor).Render(p.Name)))
			} else {
				b.WriteString(fmt.Sprintf("    %s\n", p.Name))
			}
		}

		b.WriteString("\n")
		b.WriteString(helpStyle.Render("[↑↓] Navigate  [Enter] Select  [Esc] Back"))
		return baseStyle.Render(b.String())
	}

	// Edit phase
	proj := m.projects[m.selectedIdx]
	b.WriteString(labelStyle.Render(fmt.Sprintf("Editing: %s", proj.Name)))
	b.WriteString("\n\n")

	// Name
	b.WriteString(labelStyle.Render("Project name:"))
	b.WriteString("\n")
	if m.focusIndex == projSettingsName {
		b.WriteString(focusedInputStyle.Render(m.inputs[projSettingsName].View()))
	} else {
		b.WriteString(blurredInputStyle.Render(m.inputs[projSettingsName].View()))
	}
	b.WriteString("\n\n")

	// Summary enabled toggle
	b.WriteString(labelStyle.Render("Daily summary enabled:"))
	b.WriteString("\n")
	if m.focusIndex == projSettingsEnabled {
		if m.enabledToggle {
			b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(successColor).Render("[x] Enabled"))
		} else {
			b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(errorColor).Render("[ ] Disabled"))
		}
		b.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("  (space to toggle)"))
	} else {
		if m.enabledToggle {
			b.WriteString(lipgloss.NewStyle().Foreground(successColor).Render("[x] Enabled"))
		} else {
			b.WriteString(lipgloss.NewStyle().Foreground(errorColor).Render("[ ] Disabled"))
		}
	}
	b.WriteString("\n\n")

	// Summary time
	b.WriteString(labelStyle.Render("Summary time (HH:MM):"))
	b.WriteString("\n")
	if m.focusIndex == projSettingsSummaryTime {
		b.WriteString(focusedInputStyle.Render(m.inputs[projSettingsSummaryTime].View()))
	} else {
		b.WriteString(blurredInputStyle.Render(m.inputs[projSettingsSummaryTime].View()))
	}
	b.WriteString("\n\n")

	// Timezone
	b.WriteString(labelStyle.Render("Timezone (IANA):"))
	b.WriteString("\n")
	if m.focusIndex == projSettingsTimezone {
		b.WriteString(focusedInputStyle.Render(m.inputs[projSettingsTimezone].View()))
	} else {
		b.WriteString(blurredInputStyle.Render(m.inputs[projSettingsTimezone].View()))
	}
	b.WriteString("\n\n")

	// Discord webhook URL
	b.WriteString(labelStyle.Render("Discord webhook URL:"))
	b.WriteString("\n")
	if m.focusIndex == projSettingsDiscordURL {
		b.WriteString(focusedInputStyle.Render(m.inputs[projSettingsDiscordURL].View()))
	} else {
		b.WriteString(blurredInputStyle.Render(m.inputs[projSettingsDiscordURL].View()))
	}
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
		b.WriteString(helpStyle.Render("[Tab] Next field  [Space] Toggle  [Ctrl+S] Save  [Esc] Back to list"))
	}

	return baseStyle.Render(b.String())
}

func (m *ProjectSettingsModel) fetchProjects() tea.Cmd {
	return func() tea.Msg {
		projects, err := m.projectService.GetAll()
		if err != nil {
			return projectSettingsErrorMsg(err.Error())
		}
		return projectSettingsProjectsMsg{projects: projects}
	}
}

func (m *ProjectSettingsModel) handleSubmit() tea.Cmd {
	name := strings.TrimSpace(m.inputs[projSettingsName].Value())
	summaryTime := strings.TrimSpace(m.inputs[projSettingsSummaryTime].Value())
	timezone := strings.TrimSpace(m.inputs[projSettingsTimezone].Value())
	discordURL := strings.TrimSpace(m.inputs[projSettingsDiscordURL].Value())
	enabled := m.enabledToggle

	m.errorMsg = ""
	m.successMsg = ""
	m.loading = true

	projectID := m.projects[m.selectedIdx].ID

	return func() tea.Msg {
		req := &api.UpdateProjectRequest{
			SummaryEnabled: &enabled,
		}

		if name != "" {
			req.Name = &name
		}
		if summaryTime != "" {
			req.SummaryTime = &summaryTime
		}
		if timezone != "" {
			req.Timezone = &timezone
		}
		if discordURL != "" {
			req.DiscordWebhookURL = &discordURL
		}

		_, err := m.projectService.Update(projectID, req)
		if err != nil {
			return projectSettingsErrorMsg(err.Error())
		}
		return projectSettingsUpdatedMsg{}
	}
}

// Message types
type projectSettingsProjectsMsg struct {
	projects []*api.ProjectResponse
}
type projectSettingsUpdatedMsg struct{}
type projectSettingsErrorMsg string
