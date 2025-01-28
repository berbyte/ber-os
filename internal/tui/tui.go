// Copyright 2025 BER - ber.run
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tui

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/berbyte/ber-os/internal/agent"
	llm "github.com/berbyte/ber-os/pkg/openai"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"

	_ "github.com/berbyte/ber-os/agents"
	"github.com/berbyte/ber-os/internal/services/memory_store"

	"github.com/berbyte/ber-os/internal/logger"
	"go.uber.org/zap"
)

var (
	agentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("63"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	chatStyle = lipgloss.NewStyle()

	userMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("63"))

	berMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))

	docStyle = lipgloss.NewStyle().Padding(1, 2)

	agentListStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 0)

	chatBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")).
			Padding(1, 0)

	logoStyle = lipgloss.NewStyle().
			Padding(0, 1)

	berLogo = `
██████╗ ███████╗██████╗    ╔══════════════════════════════╗
██╔══██╗██╔════╝██╔══██╗   ║    Your Agents, Your Rules   ║
██████╔╝█████╗  ██████╔╝   ╚══════════════════════════════╝
██╔══██╗██╔══╝  ██╔══██╗
██████╔╝███████╗██║  ██║
╚═════╝ ╚══════╝╚═╝  ╚═╝`
)

type state int

const (
	stateAgentSelect state = iota
	stateChat
)

type model struct {
	state         state
	agents        []agent.BERAgent
	selectedAgent agent.BERAgent
	textInput     textinput.Model
	messages      []llm.ChatMessage
	viewport      viewport.Model
	spinner       spinner.Model
	client        *llm.Client
	waiting       bool
	err           error
	width         int
	height        int
	activePanel   string // "agents" or "chat"
	memoryStore   *memory_store.MemoryStore
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Type your message..."
	ti.Focus()
	ti.CharLimit = 1000
	ti.Width = 80

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	vp := viewport.New(0, 0)
	vp.Style = chatStyle

	return model{
		state:       stateAgentSelect,
		agents:      agent.GetRegisteredAgents(),
		textInput:   ti,
		viewport:    vp,
		spinner:     s,
		client:      llm.NewClient(),
		messages:    []llm.ChatMessage{},
		activePanel: "agents",
		memoryStore: memory_store.NewMemoryStore(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "tab":
			// Toggle between panels
			if m.activePanel == "agents" {
				m.activePanel = "chat"
			} else {
				m.activePanel = "agents"
			}
			return m, nil

		case "pgup", "up":
			if m.activePanel == "chat" {
				m.viewport.LineUp(1)
				return m, nil
			}
			if m.activePanel == "agents" {
				currentIdx := -1
				// Find current selection
				for i, agent := range m.agents {
					if agent.Name == m.selectedAgent.Name {
						currentIdx = i
						break
					}
				}

				// Update selection
				if msg.String() == "up" {
					if currentIdx <= 0 {
						currentIdx = len(m.agents) - 1
					} else {
						currentIdx--
					}
				} else {
					if currentIdx >= len(m.agents)-1 {
						currentIdx = 0
					} else {
						currentIdx++
					}
				}

				m.selectedAgent = m.agents[currentIdx]
				return m, nil
			}

		case "pgdown", "down":
			if m.activePanel == "chat" {
				m.viewport.LineDown(1)
				return m, nil
			}
			if m.activePanel == "agents" {
				currentIdx := -1
				// Find current selection
				for i, agent := range m.agents {
					if agent.Name == m.selectedAgent.Name {
						currentIdx = i
						break
					}
				}

				// Update selection
				if msg.String() == "up" {
					if currentIdx <= 0 {
						currentIdx = len(m.agents) - 1
					} else {
						currentIdx--
					}
				} else {
					if currentIdx >= len(m.agents)-1 {
						currentIdx = 0
					} else {
						currentIdx++
					}
				}

				m.selectedAgent = m.agents[currentIdx]
				return m, nil
			}

		case "home":
			if m.activePanel == "chat" {
				m.viewport.GotoTop()
				return m, nil
			}

		case "end":
			if m.activePanel == "chat" {
				m.viewport.GotoBottom()
				return m, nil
			}

		case "enter":
			if m.activePanel == "agents" {
				m.state = stateChat
				m.activePanel = "chat"
				m.textInput.Reset()
				return m, nil
			}

			if m.activePanel == "chat" && m.state == stateChat {
				if m.waiting {
					return m, nil
				}

				// Handle chat input
				input := strings.TrimSpace(m.textInput.Value())
				if input == "" {
					return m, nil
				}

				m.messages = append(m.messages, llm.ChatMessage{
					Role:    llm.RoleUser,
					Content: input,
				})
				m.textInput.Reset()
				m.waiting = true

				// Update viewport with new message
				m.updateViewport()

				return m, tea.Batch(
					m.spinner.Tick,
					func() tea.Msg { return getBerResponse(m) },
				)
			}
		}

	case tea.MouseMsg:
		if m.activePanel == "chat" {
			switch {
			case msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonWheelUp:
				m.viewport.LineUp(3)
				return m, nil
			case msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonWheelDown:
				m.viewport.LineDown(3)
				return m, nil
			}
		}

	case berResponse:
		m.waiting = false
		if msg.result.Error != nil {
			m.err = msg.result.Error
			return m, nil
		}

		m.messages = append(m.messages, llm.ChatMessage{
			Role:    llm.RoleAssistant,
			Content: msg.result.Response,
		})
		m.updateViewport()
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update viewport height to fill available space
		logoHeight := 7   // Height of the ASCII art + padding
		headerHeight := 3 // Title + padding
		footerHeight := 3 // Input + padding
		m.viewport.Width = m.width/4*3 - 4
		m.viewport.Height = m.height - headerHeight - footerHeight - logoHeight - 2

		// Update input width
		m.textInput.Width = m.viewport.Width

		if m.viewport.Height < 0 {
			m.viewport.Height = 0
		}

		return m, nil
	}

	// Only update textinput if the active panel should receive input
	if (m.activePanel == "agents" && m.state == stateAgentSelect) ||
		(m.activePanel == "chat" && m.state == stateChat) {
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return m, cmd
}

func (m *model) updateViewport() {
	var b strings.Builder

	// Add scrolling instructions without any padding
	if len(m.messages) > 0 {
		b.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Render("(↑/↓ to scroll, Home/End to jump)") + "\n\n")
	}

	for _, msg := range m.messages {
		if msg.Role == llm.RoleUser {
			b.WriteString(userMsgStyle.Render("You: ") + msg.Content + "\n\n")
		} else {
			// Render markdown for assistant messages
			rendered, err := glamour.Render(msg.Content, "dark")
			content := msg.Content
			if err == nil {
				content = strings.TrimSpace(rendered)
			}
			b.WriteString(berMsgStyle.Render("BER: ") + content + "\n\n")
		}
	}

	m.viewport.SetContent(b.String())
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Charmbracelet's signature colors in reverse
	baseColors := []string{
		"#A076F2", // violet
		"#B77EE0", // purple
		"#D290E4", // lavender
		"#71BEF2", // azure
		"#66C2CD", // sky blue
		"#A8CC8C", // sage green
		"#DBAB79", // warm yellow
		"#DFA17F", // peach
		"#E88388", // coral
		"#F25D94", // soft pink
	}

	// Create transition colors between each base color
	var colors []string
	for i := 0; i < len(baseColors)-1; i++ {
		colors = append(colors, baseColors[i])
		// Add transition colors between base colors
		c1, _ := colorful.Hex(baseColors[i])
		c2, _ := colorful.Hex(baseColors[i+1])
		for j := 1; j < 8; j++ { // More transition steps for ultra-smooth gradient
			transitionalColor := c1.BlendLuv(c2, float64(j)/8.0)
			colors = append(colors, transitionalColor.Hex())
		}
	}
	colors = append(colors, baseColors[len(baseColors)-1])

	var coloredLogo strings.Builder
	lines := strings.Split(berLogo, "\n")

	// Calculate how many characters are in each line
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// Color each line with a gradient
	for _, line := range lines {
		if len(line) == 0 {
			coloredLogo.WriteString("\n")
			continue
		}

		// Create a gradient across each character in the line
		for i, char := range line {
			if char == ' ' {
				coloredLogo.WriteString(" ")
				continue
			}

			// Calculate color index based on character position
			colorIdx := (i * (len(colors) - 1)) / maxLen
			if colorIdx >= len(colors) {
				colorIdx = len(colors) - 1
			}

			coloredChar := lipgloss.NewStyle().
				Foreground(lipgloss.Color(colors[colorIdx])).
				Render(string(char))
			coloredLogo.WriteString(coloredChar)
		}
		coloredLogo.WriteString("\n")
	}

	// Create the logo section
	logo := logoStyle.Render(coloredLogo.String())

	// Calculate dimensions
	mainWidth := m.width/4*3 - 4
	sideWidth := m.width/4 - 4

	var leftView, rightView string

	// Update the styles based on active panel
	activeAgentStyle := agentListStyle.
		BorderForeground(lipgloss.Color("205"))
	activeChatStyle := chatBoxStyle.
		BorderForeground(lipgloss.Color("205"))
	inactiveStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 0)

	// Agent list (left side)
	{
		var b strings.Builder
		b.WriteString("Agents (↑↓ to select, Tab to switch)\n\n")

		for _, agent := range m.agents {
			prefix := "  "
			style := agentStyle
			if agent.Name == m.selectedAgent.Name {
				prefix = "▸ "
				style = selectedStyle
			}
			b.WriteString(fmt.Sprintf("%s%s\n", prefix, style.Render(agent.Name)))
		}

		leftView = lipgloss.NewStyle().Width(sideWidth).Render(b.String())
		if m.activePanel == "agents" {
			leftView = activeAgentStyle.Render(leftView)
		} else {
			leftView = inactiveStyle.Render(leftView)
		}
	}

	// Chat view (right side)
	{
		var b strings.Builder

		if m.state == stateChat {
			// Add viewport with scrollbar
			vp := m.viewport.View()
			scrollbar := ""
			totalLines := strings.Count(m.viewport.View(), "\n")
			if m.viewport.Height < totalLines {
				scrollPercent := float64(m.viewport.YOffset) / float64(totalLines-m.viewport.Height)
				scrollbarHeight := m.viewport.Height
				scrollPosition := int(float64(scrollbarHeight-1) * scrollPercent)

				scrollbar = lipgloss.NewStyle().
					Foreground(lipgloss.Color("240")).
					Render(strings.Repeat("│", scrollPosition) + "┃" +
						strings.Repeat("│", scrollbarHeight-scrollPosition-1))
			}

			if scrollbar != "" {
				// Join viewport and scrollbar side by side
				vp = lipgloss.JoinHorizontal(lipgloss.Top, vp, scrollbar)
			}
			b.WriteString(vp)

			if m.waiting {
				b.WriteString("\n" + m.spinner.View() + " Thinking...")
			}

			if m.err != nil {
				b.WriteString("\n" + errorStyle.Render("Error: "+m.err.Error()))
			}
		}

		b.WriteString("\n" + m.textInput.View())
		rightView = lipgloss.NewStyle().Width(mainWidth).Render(b.String())
		if m.activePanel == "chat" {
			rightView = activeChatStyle.Render(rightView)
		} else {
			rightView = inactiveStyle.Render(rightView)
		}
	}

	// Combine everything
	mainView := lipgloss.JoinVertical(
		lipgloss.Left,
		logo,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			leftView,
			rightView,
		),
	)

	return docStyle.Render(mainView)
}

// Update the berResponse type to match workflow result
type berResponse struct {
	result agent.WorkflowResult
}

// Replace the getBerResponse function with this new version
func getBerResponse(m model) berResponse {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute workflow
	result := agent.ExecuteWorkflow(ctx, agent.WorkflowInput{
		Message:     m.messages[len(m.messages)-1].Content,
		Agent:       &m.selectedAgent,
		ChatHistory: m.messages,
		LLMClient:   m.client,
		MemoryStore: m.memoryStore,
		SkillMatcher: func(msg string, a *agent.BERAgent) (agent.BaseSkill, error) {
			return agent.DefaultSkillMatcher(ctx, msg, a, m.client)
		},
	})

	return berResponse{result: result}
}

func StartTUI() {
	// Redirect all logging to /dev/null before starting the program
	log.SetOutput(io.Discard)
	logger.Log = zap.NewNop() // Replace the logger with a no-op logger

	// Create and start the program
	m := initialModel()
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
