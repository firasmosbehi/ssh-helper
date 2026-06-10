package tui

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/firasmosbehi/ssh-helper/internal/config"
	"github.com/firasmosbehi/ssh-helper/internal/core"
	"github.com/firasmosbehi/ssh-helper/internal/mcp"
	"github.com/firasmosbehi/ssh-helper/internal/ssh"
	"github.com/firasmosbehi/ssh-helper/internal/store"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	itemStyle  = lipgloss.NewStyle().PaddingLeft(2)
	selected   = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#7D56F4"))
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
)

// App is the main TUI application entrypoint.
func App() error {
	m := newModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

type state int

const (
	stateMenu state = iota
	stateHosts
	stateJobs
	stateKeys
	stateMCP
	stateMCPTools
	stateMCPArgs
	stateMCPResult
)

type mcpToolsMsg struct {
	tools []mcpgo.Tool
	err   error
}

type mcpCallMsg struct {
	result string
	err    error
}

type model struct {
	state  state
	width  int
	height int
	cursor int
	filter string
	msg    string

	hosts []core.Host
	jobs  []core.SyncJob
	keys  []ssh.KeyInfo

	cfgPath    string
	mcpServers map[string]config.MCPClientConfig
	mcpTools   []mcpgo.Tool
	mcpResult  string
	inputBuf   string
	selectedServer string
	selectedTool   string
}

func newModel() model {
	cfg, _ := config.NewManager()
	c, _ := cfg.Load()
	path := c.SSHConfigPath
	if path == "" {
		path = os.ExpandEnv("$HOME/.ssh/config")
	}
	return model{
		state:      stateMenu,
		cfgPath:    path,
		mcpServers: c.MCPClients,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case mcpToolsMsg:
		if msg.err != nil {
			m.msg = fmt.Sprintf("Error: %v", msg.err)
		} else {
			m.mcpTools = msg.tools
			m.state = stateMCPTools
			m.cursor = 0
			m.msg = ""
		}
		return m, nil
	case mcpCallMsg:
		if msg.err != nil {
			m.msg = fmt.Sprintf("Error: %v", msg.err)
			m.state = stateMCPResult
		} else {
			m.mcpResult = msg.result
			m.state = stateMCPResult
			m.msg = ""
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			switch m.state {
			case stateMenu:
				if m.cursor < 4 {
					m.cursor++
				}
			case stateHosts:
				if m.cursor < len(m.hosts)-1 {
					m.cursor++
				}
			case stateJobs:
				if m.cursor < len(m.jobs)-1 {
					m.cursor++
				}
			case stateKeys:
				if m.cursor < len(m.keys)-1 {
					m.cursor++
				}
			case stateMCP:
				if m.cursor < len(mcpServerNames(m.mcpServers))-1 {
					m.cursor++
				}
			case stateMCPTools:
				if m.cursor < len(m.mcpTools)-1 {
					m.cursor++
				}
			}
		case "enter":
			return m.handleEnter()
		case "esc":
			switch m.state {
			case stateMCPResult:
				m.state = stateMCPTools
				m.cursor = 0
				m.msg = ""
			case stateMCPArgs:
				m.state = stateMCPTools
				m.cursor = 0
				m.msg = ""
			case stateMCPTools:
				m.state = stateMCP
				m.cursor = 0
				m.msg = ""
			case stateMCP:
				m.state = stateMenu
				m.cursor = 0
				m.msg = ""
			case stateHosts, stateJobs, stateKeys:
				m.state = stateMenu
				m.cursor = 0
				m.msg = ""
			}
		case "backspace":
			if m.state == stateMCPArgs && len(m.inputBuf) > 0 {
				m.inputBuf = m.inputBuf[:len(m.inputBuf)-1]
			} else if m.state != stateMenu && len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.applyFilter()
			}
		default:
			if m.state == stateMCPArgs {
				m.inputBuf += msg.String()
			} else if m.state != stateMenu && len(msg.String()) == 1 {
				m.filter += msg.String()
				m.applyFilter()
			}
		}
	}
	return m, nil
}

func (m *model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case stateMenu:
		switch m.cursor {
		case 0:
			m.state = stateHosts
			m.loadHosts()
		case 1:
			m.state = stateJobs
			m.loadJobs()
		case 2:
			m.state = stateKeys
			m.loadKeys()
		case 3:
			m.state = stateMCP
		case 4:
			return m, tea.Quit
		}
		m.cursor = 0
		m.filter = ""
	case stateHosts:
		if m.cursor < len(m.hosts) {
			h := m.hosts[m.cursor]
			m.msg = fmt.Sprintf("Connect to %s (%s@%s:%d)", h.Name, h.User, h.Hostname, h.Port)
		}
	case stateJobs:
		if m.cursor < len(m.jobs) {
			j := m.jobs[m.cursor]
			m.msg = fmt.Sprintf("Run job %s (%s -> %s)", j.Name, j.Source, j.Dest)
		}
	case stateKeys:
		if m.cursor < len(m.keys) {
			k := m.keys[m.cursor]
			m.msg = fmt.Sprintf("Key %s (%s)", k.Name, k.Fingerprint)
		}
	case stateMCP:
		names := mcpServerNames(m.mcpServers)
		if m.cursor < len(names) {
			m.selectedServer = names[m.cursor]
			return m, m.listMCPToolsCmd(m.mcpServers[m.selectedServer])
		}
	case stateMCPTools:
		if m.cursor < len(m.mcpTools) {
			t := m.mcpTools[m.cursor]
			m.selectedTool = t.Name
			m.state = stateMCPArgs
			m.inputBuf = "{}"
			m.cursor = 0
			m.msg = fmt.Sprintf("Tool: %s - %s", t.Name, t.Description)
		}
	case stateMCPArgs:
		return m, m.callMCPToolCmd(m.mcpServers[m.selectedServer], m.selectedTool, m.inputBuf)
	case stateMCPResult:
		m.state = stateMCPTools
		m.cursor = 0
		m.msg = ""
	}
	return m, nil
}

func (m *model) applyFilter() {
	m.cursor = 0
}

func (m *model) loadHosts() {
	cfg, err := ssh.ParseConfig(m.cfgPath)
	if err == nil {
		m.hosts = ssh.HostsFromConfig(cfg)
	}
}

func (m *model) loadJobs() {
	dir, _ := os.UserConfigDir()
	s := store.NewJSONStore(dir + "/ssh-helper")
	jobs, _ := s.ListSyncJobs()
	m.jobs = jobs
}

func (m *model) loadKeys() {
	home, _ := os.UserHomeDir()
	keys, _ := ssh.ListKeys(home + "/.ssh")
	m.keys = keys
}

func (m *model) listMCPToolsCmd(srv config.MCPClientConfig) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		tools, err := mcp.ListTools(ctx, srv)
		return mcpToolsMsg{tools: tools, err: err}
	}
}

func (m *model) callMCPToolCmd(srv config.MCPClientConfig, toolName, args string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		res, err := mcp.CallToolString(ctx, srv, toolName, args)
		return mcpCallMsg{result: res, err: err}
	}
}

func mcpServerNames(servers map[string]config.MCPClientConfig) []string {
	names := make([]string, 0, len(servers))
	for name := range servers {
		names = append(names, name)
	}
	// Sort for stable ordering
	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			if names[j] < names[i] {
				names[i], names[j] = names[j], names[i]
			}
		}
	}
	return names
}

func (m model) View() string {
	switch m.state {
	case stateMenu:
		return m.viewMenu()
	case stateHosts:
		return m.viewHosts()
	case stateJobs:
		return m.viewJobs()
	case stateKeys:
		return m.viewKeys()
	case stateMCP:
		return m.viewMCP()
	case stateMCPTools:
		return m.viewMCPTools()
	case stateMCPArgs:
		return m.viewMCPArgs()
	case stateMCPResult:
		return m.viewMCPResult()
	}
	return ""
}

func (m model) viewMenu() string {
	items := []string{"Hosts", "Sync Jobs", "Keys", "MCP", "Quit"}
	var b strings.Builder
	b.WriteString(titleStyle.Render("ssh-helper") + "\n\n")
	for i, item := range items {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			b.WriteString(selected.Render(fmt.Sprintf("%s %s", cursor, item)) + "\n")
		} else {
			b.WriteString(itemStyle.Render(fmt.Sprintf("%s %s", cursor, item)) + "\n")
		}
	}
	b.WriteString("\n" + helpStyle.Render("↑/↓ navigate • enter select • q quit"))
	return b.String()
}

func (m model) viewHosts() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Hosts") + "\n")
	if m.filter != "" {
		b.WriteString(helpStyle.Render("filter: "+m.filter) + "\n")
	}
	b.WriteString("\n")
	filtered := filterHosts(m.hosts, m.filter)
	for i, h := range filtered {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			b.WriteString(selected.Render(fmt.Sprintf("%s %s (%s@%s:%d)", cursor, h.Name, h.User, h.Hostname, h.Port)) + "\n")
		} else {
			b.WriteString(itemStyle.Render(fmt.Sprintf("%s %s (%s@%s:%d)", cursor, h.Name, h.User, h.Hostname, h.Port)) + "\n")
		}
	}
	if m.msg != "" {
		b.WriteString("\n" + m.msg + "\n")
	}
	b.WriteString("\n" + helpStyle.Render("↑/↓ navigate • enter select • esc back • type to filter"))
	return b.String()
}

func (m model) viewJobs() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Sync Jobs") + "\n\n")
	for i, j := range m.jobs {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			b.WriteString(selected.Render(fmt.Sprintf("%s %s -> %s", cursor, j.Source, j.Dest)) + "\n")
		} else {
			b.WriteString(itemStyle.Render(fmt.Sprintf("%s %s -> %s", cursor, j.Source, j.Dest)) + "\n")
		}
	}
	if m.msg != "" {
		b.WriteString("\n" + m.msg + "\n")
	}
	b.WriteString("\n" + helpStyle.Render("↑/↓ navigate • enter run • esc back"))
	return b.String()
}

func (m model) viewKeys() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("SSH Keys") + "\n\n")
	for i, k := range m.keys {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			b.WriteString(selected.Render(fmt.Sprintf("%s %s (%s)", cursor, k.Name, k.Type)) + "\n")
		} else {
			b.WriteString(itemStyle.Render(fmt.Sprintf("%s %s (%s)", cursor, k.Name, k.Type)) + "\n")
		}
	}
	if m.msg != "" {
		b.WriteString("\n" + m.msg + "\n")
	}
	b.WriteString("\n" + helpStyle.Render("↑/↓ navigate • esc back"))
	return b.String()
}

func (m model) viewMCP() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("MCP Servers") + "\n\n")
	names := mcpServerNames(m.mcpServers)
	for i, name := range names {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			b.WriteString(selected.Render(fmt.Sprintf("%s %s", cursor, name)) + "\n")
		} else {
			b.WriteString(itemStyle.Render(fmt.Sprintf("%s %s", cursor, name)) + "\n")
		}
	}
	if m.msg != "" {
		b.WriteString("\n" + m.msg + "\n")
	}
	b.WriteString("\n" + helpStyle.Render("↑/↓ navigate • enter connect • esc back"))
	return b.String()
}

func (m model) viewMCPTools() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render(fmt.Sprintf("MCP Tools (%s)", m.selectedServer)) + "\n\n")
	for i, t := range m.mcpTools {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			b.WriteString(selected.Render(fmt.Sprintf("%s %s - %s", cursor, t.Name, t.Description)) + "\n")
		} else {
			b.WriteString(itemStyle.Render(fmt.Sprintf("%s %s", cursor, t.Name)) + "\n")
		}
	}
	if m.msg != "" {
		b.WriteString("\n" + m.msg + "\n")
	}
	b.WriteString("\n" + helpStyle.Render("↑/↓ navigate • enter select • esc back"))
	return b.String()
}

func (m model) viewMCPArgs() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render(fmt.Sprintf("Call %s", m.selectedTool)) + "\n\n")
	b.WriteString(helpStyle.Render("JSON arguments:") + "\n")
	b.WriteString(m.inputBuf + "\n")
	if m.msg != "" {
		b.WriteString("\n" + m.msg + "\n")
	}
	b.WriteString("\n" + helpStyle.Render("type to edit • enter call • esc back"))
	return b.String()
}

func (m model) viewMCPResult() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Result") + "\n\n")
	b.WriteString(m.mcpResult + "\n")
	if m.msg != "" {
		b.WriteString("\n" + m.msg + "\n")
	}
	b.WriteString("\n" + helpStyle.Render("enter/esc back"))
	return b.String()
}

func filterHosts(hosts []core.Host, filter string) []core.Host {
	if filter == "" {
		return hosts
	}
	var out []core.Host
	for _, h := range hosts {
		if strings.Contains(h.Name, filter) || strings.Contains(h.Hostname, filter) {
			out = append(out, h)
		}
	}
	return out
}
