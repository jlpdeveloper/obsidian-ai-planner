package main

import (
	"context"
	"fmt"
	_ "log"
	"obsidian-ai-planner/calendar"
	"obsidian-ai-planner/local_ai"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const gap = "\n\n"

var cal *calendar.GoogleCalendarIntegration
var calendarEvents []calendar.Event

func Today() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, now.Location())
}

type (
	errMsg error
)

type chatModel struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
	initialMsg  string
	spinner     spinner.Model
	loading     bool
	modelInfo   *local_ai.ModelInfo
}

func initialChatModel(initialMsg string) chatModel {
	ctx := context.Background()
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	cal = calendar.New(ctx)
	//replace _ with calendarEvents when ready to send to LLM
	var err error
	calendarEvents, err = cal.GetCalendarEvents(Today())

	if err != nil {
		return chatModel{err: err}
	}

	modelInfo := local_ai.NewOllamaModel(ctx)
	local_ai.DefinePlannerFlow(modelInfo)

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	initMsg := `Welcome to the obsidian planner!
Type a message and press Enter to send.`

	vp.SetContent(initMsg)
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return chatModel{
		textarea: ta,
		messages: []string{
			initMsg,
		},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
		initialMsg:  initialMsg,
		spinner:     s,
		loading:     false,
		modelInfo:   modelInfo,
	}
}

type cmdArgMsg string

func cmdWithStr(s string) tea.Cmd {
	return func() tea.Msg {
		if s == "" {
			return cmdArgMsg("Hello World!")
		}
		return cmdArgMsg(s)
	}
}

func (m chatModel) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, m.spinner.Tick, cmdWithStr(m.initialMsg))
}

func (m chatModel) runPlannerFlow(userPrompt string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		input := local_ai.PlannerInput{
			WeeklyGoals:  "", // TODO: Pull from Obsidian
			Calendar:     calendarEvents,
			JiraTickets:  []string{}, // TODO: Pull from Jira
			CurrentTasks: []string{}, // TODO: Pull from Daily Note
			UserPrompt:   userPrompt,
		}
		input.JiraTickets = append(input.JiraTickets, "Sample Jira Ticket")

		resp, err := m.modelInfo.GeneratePlan(ctx, input)
		if err != nil {
			return errMsg(err)
		}

		return cmdArgMsg(resp)
	}
}

func (m chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	var spCmd tea.Cmd
	m.spinner, spCmd = m.spinner.Update(msg)

	switch msg := msg.(type) {
	case cmdArgMsg:
		m.loading = false
		m.messages = append(m.messages, m.senderStyle.Render("Bot: ")+string(msg))
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		m.viewport.GotoBottom()
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)

		if len(m.messages) > 0 {
			// Wrap content before setting it.
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			if m.textarea.Value() != "" {
				userMsg := m.textarea.Value()
				m.messages = append(m.messages, m.senderStyle.Render("You: ")+userMsg)
				m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
				m.textarea.Reset()
				m.viewport.GotoBottom()
				m.loading = true
				return m, tea.Batch(
					tiCmd,
					vpCmd,
					spCmd,
					m.runPlannerFlow(userMsg),
				)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		m.loading = false
		m.messages = append(m.messages, m.senderStyle.Render("Error: ")+msg.Error())
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		m.viewport.GotoBottom()
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd, spCmd)
}

func (m chatModel) View() string {
	var s string
	if m.loading {
		s = m.spinner.View() + " Thinking..."
	} else {
		s = m.textarea.View()
	}
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		s,
	)
}
