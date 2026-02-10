package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muster/cli/internal/standup"
)

// StandupSubmitModel represents the standup submission form
type StandupSubmitModel struct {
	standupService *standup.Service
	inputs         []textarea.Model
	focusedInput   int
	date           string
	errorMsg       string
	successMsg     string
	loading        bool
	shouldGoBack   bool
	onBack         func()
}

const (
	yesterdayInput = iota
	todayInput
	blockersInput
)

// NewStandupSubmitModel creates a new standup submission model
func NewStandupSubmitModel(standupService *standup.Service, onBack func()) StandupSubmitModel {
	m := StandupSubmitModel{
		standupService: standupService,
		inputs:         make([]textarea.Model, 3),
		focusedInput:   yesterdayInput,
		date:           time.Now().Format("2006-01-02"),
		onBack:         onBack,
	}

	// Yesterday textarea
	m.inputs[yesterdayInput] = textarea.New()
	m.inputs[yesterdayInput].Placeholder = "What did you work on yesterday?"
	m.inputs[yesterdayInput].Focus()
	m.inputs[yesterdayInput].CharLimit = 1000
	m.inputs[yesterdayInput].SetWidth(60)
	m.inputs[yesterdayInput].SetHeight(3)

	// Today textarea
	m.inputs[todayInput] = textarea.New()
	m.inputs[todayInput].Placeholder = "What will you work on today?"
	m.inputs[todayInput].CharLimit = 1000
	m.inputs[todayInput].SetWidth(60)
	m.inputs[todayInput].SetHeight(3)

	// Blockers textarea
	m.inputs[blockersInput] = textarea.New()
	m.inputs[blockersInput].Placeholder = "Any blockers or issues?"
	m.inputs[blockersInput].CharLimit = 1000
	m.inputs[blockersInput].SetWidth(60)
	m.inputs[blockersInput].SetHeight(3)

	return m
}

// Init initializes the standup submit model
func (m StandupSubmitModel) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles messages
func (m StandupSubmitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			m.shouldGoBack = true
			return m, nil

		case "tab", "shift+tab":
			// Cycle through inputs
			if msg.String() == "shift+tab" {
				m.focusedInput--
			} else {
				m.focusedInput++
			}

			if m.focusedInput < 0 {
				m.focusedInput = len(m.inputs) - 1
			} else if m.focusedInput >= len(m.inputs) {
				m.focusedInput = 0
			}

			// Update focus
			for i := range m.inputs {
				if i == m.focusedInput {
					m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}

			return m, nil

		case "ctrl+s":
			if m.loading {
				return m, nil
			}
			return m, m.handleSubmit()
		}

	case standupSubmitSuccessMsg:
		m.loading = false
		m.successMsg = "Standup submitted successfully!"
		return m, nil

	case standupSubmitErrorMsg:
		m.loading = false
		m.errorMsg = string(msg)
		return m, nil
	}

	// Update focused textarea
	var cmd tea.Cmd
	m.inputs[m.focusedInput], cmd = m.inputs[m.focusedInput].Update(msg)
	return m, cmd
}

// View renders the standup submission form
func (m StandupSubmitModel) View() string {

	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render(fmt.Sprintf("📝 Submit Standup - %s", m.date)))
	b.WriteString("\n\n")

	// Yesterday input
	b.WriteString(labelStyle.Render("What did you do yesterday?"))
	b.WriteString("\n")
	if m.focusedInput == yesterdayInput {
		b.WriteString(focusedInputStyle.Render(m.inputs[yesterdayInput].View()))
	} else {
		b.WriteString(blurredInputStyle.Render(m.inputs[yesterdayInput].View()))
	}
	b.WriteString("\n")

	// Today input
	b.WriteString(labelStyle.Render("What will you do today?"))
	b.WriteString("\n")
	if m.focusedInput == todayInput {
		b.WriteString(focusedInputStyle.Render(m.inputs[todayInput].View()))
	} else {
		b.WriteString(blurredInputStyle.Render(m.inputs[todayInput].View()))
	}
	b.WriteString("\n")

	// Blockers input
	b.WriteString(labelStyle.Render("Any blockers?"))
	b.WriteString("\n")
	if m.focusedInput == blockersInput {
		b.WriteString(focusedInputStyle.Render(m.inputs[blockersInput].View()))
	} else {
		b.WriteString(blurredInputStyle.Render(m.inputs[blockersInput].View()))
	}
	b.WriteString("\n\n")

	// Error message
	if m.errorMsg != "" {
		b.WriteString(errorStyle.Render("❌ " + m.errorMsg))
		b.WriteString("\n")
	}

	// Success message
	if m.successMsg != "" {
		b.WriteString(successStyle.Render("✓ " + m.successMsg))
		b.WriteString("\n")
	}

	// Loading indicator
	if m.loading {
		b.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("⏳ Submitting standup..."))
		b.WriteString("\n")
	}

	// Help text
	if !m.loading {
		b.WriteString(helpStyle.Render("tab: next field • ctrl+s: submit • esc: back"))
	}

	return baseStyle.Render(b.String())
}

// handleSubmit handles the standup form submission
func (m *StandupSubmitModel) handleSubmit() tea.Cmd {
	yesterday := strings.TrimSpace(m.inputs[yesterdayInput].Value())
	today := strings.TrimSpace(m.inputs[todayInput].Value())
	blockers := strings.TrimSpace(m.inputs[blockersInput].Value())

	// Clear previous messages
	m.errorMsg = ""
	m.successMsg = ""
	m.loading = true

	return func() tea.Msg {
		if _, err := m.standupService.Submit(m.date, yesterday, today, blockers); err != nil {
			return standupSubmitErrorMsg(err.Error())
		}
		return standupSubmitSuccessMsg{}
	}
}

// Message types
type standupSubmitSuccessMsg struct{}
type standupSubmitErrorMsg string
