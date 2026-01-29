package main

import (
	"context"
	"fmt"
	_ "log"
	"obsidian-ai-planner/local_ai"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const gap = "\n\n"

type (
	errMsg error
)

type chatModel struct {
	viewport    viewport.Model
	messages    []string
	history     []local_ai.Message
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

	modelInfo, err := local_ai.NewOllamaModel(ctx)
	if err != nil {
		// Since initialChatModel is not designed to return an error,
		// we can't easily propagate it here.
		// However, NewOllamaModel currently only returns an error if
		// something fundamentally fails in GenKit/Ollama setup which
		// is unlikely given the current implementation (it always returns nil error).
		// If it did return an error, we'd probably want to handle it in the model.
		// For now, we'll just panic if it's a genuine failure.
		panic(err)
	}
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

func (m *chatModel) runChatFlow(userPrompt string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		input := local_ai.PlannerInput{
			UserPrompt: userPrompt,
			History:    m.history,
		}

		resp, err := m.modelInfo.Chat(ctx, input)
		if err != nil {
			return errMsg(err)
		}

		return cmdArgMsg(resp)
	}
}

func (m *chatModel) runGenerateFlow(userPrompt string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		input := local_ai.PlannerInput{
			UserPrompt: userPrompt,
			History:    m.history,
		}

		resp, err := m.modelInfo.GeneratePlan(ctx, input)
		if err != nil {
			return errMsg(err)
		}

		return cmdArgMsg(resp)
	}
}

type condenseMsg string

func (m *chatModel) runCondenseFlow() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		resp, err := m.modelInfo.Condense(ctx, m.history)
		if err != nil {
			return errMsg(err)
		}
		return condenseMsg(resp)
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
	case condenseMsg:
		m.loading = false
		summary := string(msg)
		// Reset history and append summary as the first message
		m.history = []local_ai.Message{{Role: "model", Content: "Summary of previous conversation: " + summary}}
		m.messages = []string{m.senderStyle.Render("Bot: ") + "Condensing conversation completed. Context cleared and summary added."}
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		m.viewport.GotoBottom()
	case cmdArgMsg:
		m.loading = false
		m.messages = append(m.messages, m.senderStyle.Render("Bot: ")+string(msg))
		m.history = append(m.history, local_ai.Message{Role: "model", Content: string(msg)})
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
				if strings.TrimSpace(strings.ToLower(userMsg)) == "/condense" {
					m.messages = append(m.messages, m.senderStyle.Render("You: ")+userMsg)
					m.messages = append(m.messages, m.senderStyle.Render("Bot: ")+"Condensing conversation...")
					m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
					m.textarea.Reset()
					m.viewport.GotoBottom()
					m.loading = true
					return m, tea.Batch(
						tiCmd,
						vpCmd,
						spCmd,
						m.runCondenseFlow(),
					)
				}
				m.messages = append(m.messages, m.senderStyle.Render("You: ")+userMsg)
				m.history = append(m.history, local_ai.Message{Role: "user", Content: userMsg})
				m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
				m.textarea.Reset()
				m.viewport.GotoBottom()
				m.loading = true
				var cmd tea.Cmd
				if strings.Contains(strings.ToLower(userMsg), "generate") {
					cmd = m.runGenerateFlow(userMsg)
				} else {
					cmd = m.runChatFlow(userMsg)
				}
				return m, tea.Batch(
					tiCmd,
					vpCmd,
					spCmd,
					cmd,
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
