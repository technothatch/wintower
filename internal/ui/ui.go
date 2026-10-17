package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/technothatch/wintower/internal/evtlog"
	"github.com/technothatch/wintower/internal/procs"
	"github.com/technothatch/wintower/internal/services"
)

type tab int

const (
	tabServices tab = iota
	tabProcs
	tabEvents
	tabAbout
)

var tabNames = []string{"Services", "Processes", "Event Log", "About"}

type servicesMsg struct {
	items []services.Service
	err   error
}
type procsMsg struct {
	items []procs.Proc
	err   error
}
type eventsMsg struct {
	items []evtlog.Event
	err   error
}
type actionDoneMsg struct{ err error }

type model struct {
	version  string
	width    int
	height   int
	active   tab
	status   string
	confirm  string
	services table.Model
	procs    table.Model
	events   table.Model
	svcRows  []services.Service
	procRows []procs.Proc
	loading  bool
}

func newTable(header []table.Column) table.Model {
	t := table.New(
		table.WithColumns(header),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	st := table.DefaultStyles()
	st.Header = st.Header.BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).Bold(true)
	st.Selected = st.Selected.Foreground(lipgloss.Color("212")).Bold(true)
	t.SetStyles(st)
	return t
}

func New(version string) model {
	return model{
		version: version,
		services: newTable([]table.Column{
			{Title: "Name", Width: 30},
			{Title: "Display Name", Width: 42},
			{Title: "Status", Width: 12},
			{Title: "Start", Width: 12},
		}),
		procs: newTable([]table.Column{
			{Title: "PID", Width: 8},
			{Title: "Name", Width: 46},
			{Title: "Memory (MB)", Width: 14},
		}),
		events: newTable([]table.Column{
			{Title: "Time", Width: 20},
			{Title: "Level", Width: 12},
			{Title: "ID", Width: 7},
			{Title: "Provider", Width: 28},
			{Title: "Message", Width: 60},
		}),
	}
}

func Run(version string) error {
	p := tea.NewProgram(New(version), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.loadServices(), m.loadProcs())
}

func (m *model) loadServices() tea.Cmd {
	return func() tea.Msg {
		items, err := services.List()
		return servicesMsg{items: items, err: err}
	}
}

func (m *model) loadProcs() tea.Cmd {
	return func() tea.Msg {
		items, err := procs.List()
		return procsMsg{items: items, err: err}
	}
}

func (m *model) loadEvents() tea.Cmd {
	return func() tea.Msg {
		items, err := evtlog.Events()
		return eventsMsg{items: items, err: err}
	}
}

func formatMem(b uint64) string {
	return fmt.Sprintf("%.1f", float64(b)/(1024*1024))
}

func (m *model) setServiceRows(items []services.Service) {
	m.svcRows = items
	rows := make([]table.Row, len(items))
	for i, s := range items {
		rows[i] = table.Row{s.Name, s.DisplayName, s.Status, s.StartType}
	}
	m.services.SetRows(rows)
}

func (m *model) setProcRows(items []procs.Proc) {
	m.procRows = items
	rows := make([]table.Row, len(items))
	for i, p := range items {
		rows[i] = table.Row{fmt.Sprintf("%d", p.PID), p.Name, formatMem(p.MemBytes)}
	}
	m.procs.SetRows(rows)
}

func (m *model) setEventRows(items []evtlog.Event) {
	rows := make([]table.Row, len(items))
	for i, e := range items {
		ts := ""
		if !e.Time.IsZero() {
			ts = e.Time.Local().Format("15:04:05 01-02")
		}
		rows[i] = table.Row{ts, e.Level, fmt.Sprintf("%d", e.ID), e.Provider, e.Message}
	}
	m.events.SetRows(rows)
}

func (m *model) selectedService() (services.Service, bool) {
	if m.active != tabServices {
		return services.Service{}, false
	}
	row := m.services.SelectedRow()
	if row == nil || len(m.svcRows) == 0 {
		return services.Service{}, false
	}
	for _, s := range m.svcRows {
		if s.Name == row[0] {
			return s, true
		}
	}
	return services.Service{}, false
}

func (m *model) selectedProc() (procs.Proc, bool) {
	if m.active != tabProcs {
		return procs.Proc{}, false
	}
	row := m.procs.SelectedRow()
	if row == nil {
		return procs.Proc{}, false
	}
	var pid uint32
	fmt.Sscanf(row[0], "%d", &pid)
	for _, p := range m.procRows {
		if p.PID == pid {
			return p, true
		}
	}
	return procs.Proc{}, false
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		h := msg.Height - 7
		if h < 3 {
			h = 3
		}
		m.services.SetHeight(h)
		m.procs.SetHeight(h)
		m.events.SetHeight(h)
		return m, nil

	case servicesMsg:
		m.loading = false
		if msg.err != nil {
			m.status = "services: " + msg.err.Error()
			return m, nil
		}
		m.setServiceRows(msg.items)
		m.status = fmt.Sprintf("%d services", len(msg.items))
		return m, nil

	case procsMsg:
		m.loading = false
		if msg.err != nil {
			m.status = "processes: " + msg.err.Error()
			return m, nil
		}
		m.setProcRows(msg.items)
		m.status = fmt.Sprintf("%d processes", len(msg.items))
		return m, nil

	case eventsMsg:
		m.loading = false
		if msg.err != nil {
			m.status = "event log: " + msg.err.Error()
			return m, nil
		}
		m.setEventRows(msg.items)
		m.status = fmt.Sprintf("%d events (System, Critical/Error/Warning)", len(msg.items))
		return m, nil

	case actionDoneMsg:
		m.loading = false
		if msg.err != nil {
			m.status = "error: " + msg.err.Error()
		} else {
			m.status = "done"
		}
		switch m.active {
		case tabServices:
			return m, m.loadServices()
		case tabProcs:
			return m, m.loadProcs()
		}
		return m, nil

	case tea.KeyMsg:
		if m.confirm != "" {
			switch msg.String() {
			case "y", "Y":
				kind, name := splitConfirm(m.confirm)
				m.confirm = ""
				m.loading = true
				switch kind {
				case "stop-svc", "start-svc", "restart-svc":
					action := strings.TrimSuffix(kind, "-svc")
					return m, func() tea.Msg {
						return actionDoneMsg{err: services.Control(name, action)}
					}
				case "kill":
					var pid uint32
					fmt.Sscanf(name, "%d", &pid)
					return m, func() tea.Msg {
						return actionDoneMsg{err: procs.Kill(pid)}
					}
				}
			default:
				m.confirm = ""
				m.status = "cancelled"
			}
			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab", "l":
			m.active = (m.active + 1) % tab(len(tabNames))
			if m.active == tabEvents && len(m.events.Rows()) == 0 {
				m.loading = true
				return m, m.loadEvents()
			}
			return m, nil
		case "shift+tab", "h":
			m.active = (m.active + tab(len(tabNames)) - 1) % tab(len(tabNames))
			if m.active == tabEvents && len(m.events.Rows()) == 0 {
				m.loading = true
				return m, m.loadEvents()
			}
			return m, nil
		case "1":
			m.active = tabServices
			return m, nil
		case "2":
			m.active = tabProcs
			return m, nil
		case "3":
			m.active = tabEvents
			if len(m.events.Rows()) == 0 {
				m.loading = true
				return m, m.loadEvents()
			}
			return m, nil
		case "4":
			m.active = tabAbout
			return m, nil
		case "r":
			m.loading = true
			switch m.active {
			case tabServices:
				return m, m.loadServices()
			case tabProcs:
				return m, m.loadProcs()
			case tabEvents:
				return m, m.loadEvents()
			}
			return m, nil
		case "s", "S":
			if s, ok := m.selectedService(); ok && s.Status == "Stopped" {
				m.confirm = "start-svc:" + s.Name
				m.status = fmt.Sprintf("start %s? press y to confirm", s.Name)
			}
			return m, nil
		case "t", "T":
			if s, ok := m.selectedService(); ok {
				m.confirm = "restart-svc:" + s.Name
				m.status = fmt.Sprintf("restart %s? press y to confirm", s.Name)
			}
			return m, nil
		case "x", "X":
			if s, ok := m.selectedService(); ok && s.Status == "Running" {
				m.confirm = "stop-svc:" + s.Name
				m.status = fmt.Sprintf("stop %s? press y to confirm", s.Name)
				return m, nil
			}
			if p, ok := m.selectedProc(); ok {
				m.confirm = "kill:" + fmt.Sprintf("%d", p.PID)
				m.status = fmt.Sprintf("kill %s (pid %d)? press y to confirm", p.Name, p.PID)
			}
			return m, nil
		}

		var cmd tea.Cmd
		switch m.active {
		case tabServices:
			m.services, cmd = m.services.Update(msg)
		case tabProcs:
			m.procs, cmd = m.procs.Update(msg)
		case tabEvents:
			m.events, cmd = m.events.Update(msg)
		}
		return m, cmd
	}
	return m, nil
}

func splitConfirm(s string) (string, string) {
	i := strings.Index(s, ":")
	if i < 0 {
		return s, ""
	}
	return s[:i], s[i+1:]
}

func (m model) View() string {
	var b strings.Builder

	tabs := make([]string, len(tabNames))
	for i, name := range tabNames {
		style := lipgloss.NewStyle().Padding(0, 2)
		if tab(i) == m.active {
			style = style.Bold(true).Foreground(lipgloss.Color("212")).Underline(true)
		}
		tabs[i] = style.Render(name)
	}
	b.WriteString(lipgloss.NewStyle().Bold(true).Render("wintower "+m.version) + "  " + strings.Join(tabs, " ") + "\n\n")

	switch m.active {
	case tabServices:
		b.WriteString(m.services.View())
	case tabProcs:
		b.WriteString(m.procs.View())
	case tabEvents:
		b.WriteString(m.events.View())
	case tabAbout:
		b.WriteString(aboutText(m.version))
	}

	b.WriteString("\n")
	if m.loading {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("loading…") + "\n")
	}
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("210")).Render(m.status) + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(
		"tab/h/l switch · 1-4 jump · r refresh · s start · t restart · x stop/kill · y confirm · q quit",
	))
	return b.String()
}

func aboutText(version string) string {
	return fmt.Sprintf(`wintower %s

A terminal cockpit for Windows: services, processes and the System event log
in one keyboard-driven UI.

 · Services tab talks to the Windows Service Control Manager natively.
 · Processes tab enumerates processes via the Win32 API.
 · Event Log tab tails Critical/Error/Warning entries from the System log.

Project: https://github.com/technothatch/wintower
MIT License.
`, version)
}
