package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muster/cli/internal/api"
	"github.com/muster/cli/internal/attendance"
)

// AttendanceTodayModel represents the today's attendance view
type AttendanceTodayModel struct {
	attendanceService *attendance.Service
	attendance        []*api.AttendanceResponse
	date              string
	statusCounts      map[string]int
	totalMembers      int
	errorMsg          string
	loading           bool
	shouldGoBack      bool
}

// NewAttendanceTodayModel creates a new attendance today model
func NewAttendanceTodayModel(attendanceService *attendance.Service) AttendanceTodayModel {
	return AttendanceTodayModel{
		attendanceService: attendanceService,
		loading:           true,
	}
}

// Init initializes the attendance today model
func (m AttendanceTodayModel) Init() tea.Cmd {
	return m.fetchAttendance()
}

// Update handles messages
func (m AttendanceTodayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc", "q":
			m.shouldGoBack = true
			return m, nil

		case "r":
			m.loading = true
			m.errorMsg = ""
			return m, m.fetchAttendance()
		}

	case attendanceTodaySuccessMsg:
		m.loading = false
		m.attendance = msg.attendance
		m.date = msg.date
		m.statusCounts = msg.statusCounts
		m.totalMembers = msg.totalMembers
		return m, nil

	case attendanceTodayErrorMsg:
		m.loading = false
		m.errorMsg = string(msg)
		return m, nil
	}

	return m, nil
}

// View renders the today's attendance view
func (m AttendanceTodayModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render(fmt.Sprintf("Team Attendance - %s", m.date)))
	b.WriteString("\n\n")

	// Loading indicator
	if m.loading {
		b.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Loading attendance..."))
		b.WriteString("\n")
	} else if m.errorMsg != "" {
		b.WriteString(errorStyle.Render("Error: " + m.errorMsg))
		b.WriteString("\n")
	} else {
		// Summary
		summaryStyle := lipgloss.NewStyle().Foreground(mutedColor)
		b.WriteString(summaryStyle.Render(fmt.Sprintf("Marked: %d/%d", len(m.attendance), m.totalMembers)))
		b.WriteString("\n")

		// Status counts
		if m.statusCounts != nil {
			countStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
			counts := []string{}
			if c, ok := m.statusCounts["present"]; ok && c > 0 {
				counts = append(counts, fmt.Sprintf("✓ Present: %d", c))
			}
			if c, ok := m.statusCounts["remote"]; ok && c > 0 {
				counts = append(counts, fmt.Sprintf("🏠 Remote: %d", c))
			}
			if c, ok := m.statusCounts["half-day"]; ok && c > 0 {
				counts = append(counts, fmt.Sprintf("½ Half-day: %d", c))
			}
			if c, ok := m.statusCounts["absent"]; ok && c > 0 {
				counts = append(counts, fmt.Sprintf("✗ Absent: %d", c))
			}
			if len(counts) > 0 {
				b.WriteString(countStyle.Render(strings.Join(counts, "  ")))
				b.WriteString("\n")
			}
		}

		b.WriteString("\n")

		if len(m.attendance) == 0 {
			b.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("No attendance marked yet today"))
			b.WriteString("\n")
		} else {
			// Display attendance
			for i, record := range m.attendance {
				if i > 0 {
					b.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render(strings.Repeat("─", 50)))
					b.WriteString("\n")
				}

				// User name
				userName := "Unknown"
				if record.User != nil {
					userName = record.User.Name
				}

				nameStyle := lipgloss.NewStyle().Bold(true).Foreground(primaryColor)
				b.WriteString(nameStyle.Render(userName))
				b.WriteString("  ")

				// Status with icon
				statusStyle := getStatusStyle(record.Status)
				b.WriteString(statusStyle.Render(fmt.Sprintf("%s %s", getStatusIcon(record.Status), record.Status)))
				b.WriteString("\n")

				// Checkin/checkout times
				timeStyle := lipgloss.NewStyle().Foreground(mutedColor)
				if record.CheckinTime != nil {
					b.WriteString(timeStyle.Render(fmt.Sprintf("  In: %s", utcTimeToLocal(*record.CheckinTime))))
					if record.CheckoutTime != nil {
						b.WriteString(timeStyle.Render(fmt.Sprintf("  Out: %s", utcTimeToLocal(*record.CheckoutTime))))
					}
					b.WriteString("\n")
				}

				// Notes if present
				if record.Notes != nil && *record.Notes != "" {
					notesStyle := lipgloss.NewStyle().Foreground(mutedColor).Italic(true)
					b.WriteString(notesStyle.Render("  " + *record.Notes))
					b.WriteString("\n")
				}
			}
		}
	}

	// Help text
	b.WriteString("\n")
	if !m.loading {
		b.WriteString(helpStyle.Render("[r] Refresh  [Esc] Back"))
	}

	return baseStyle.Render(b.String())
}

// fetchAttendance fetches today's attendance
func (m *AttendanceTodayModel) fetchAttendance() tea.Cmd {
	return func() tea.Msg {
		resp, err := m.attendanceService.GetToday()
		if err != nil {
			return attendanceTodayErrorMsg(err.Error())
		}
		return attendanceTodaySuccessMsg{
			attendance:   resp.Attendance,
			date:         resp.Date,
			statusCounts: resp.StatusCounts,
			totalMembers: resp.TotalMembers,
		}
	}
}

func getStatusStyle(status string) lipgloss.Style {
	switch status {
	case "present":
		return lipgloss.NewStyle().Foreground(successColor)
	case "remote":
		return lipgloss.NewStyle().Foreground(primaryColor)
	case "half-day":
		return lipgloss.NewStyle().Foreground(secondaryColor)
	case "absent":
		return lipgloss.NewStyle().Foreground(errorColor)
	default:
		return lipgloss.NewStyle().Foreground(mutedColor)
	}
}

// Message types
type attendanceTodaySuccessMsg struct {
	attendance   []*api.AttendanceResponse
	date         string
	statusCounts map[string]int
	totalMembers int
}
type attendanceTodayErrorMsg string
