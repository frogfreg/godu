package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/frogfreg/du-test/fileinfo"
)

type model struct {
	currentDir    string
	loading       bool
	selectedIndex int
	err           error
	table         table.Model
}

type fileInfoResponse struct {
	data []fileinfo.FileInfo
	err  error
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m model) Init() tea.Cmd {
	return getFileInfoCmd(m.currentDir)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case fileInfoResponse:
		m.loading = false
		fir := fileInfoResponse(msg)
		if fir.err != nil {
			m.err = fir.err
			log.Fatalf("something went wrong: %v", m.err)
		}

		m.table.SetRows(fileinfo.FileInfosToRow(fir.data))
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
		case "enter", "right", "l":
			sr := m.table.SelectedRow()
			if len(sr) != 0 && sr[1] == "dir" {
				m = m.updateCurrentDir(m.table.SelectedRow()[0], false)
				return m, getFileInfoCmd(m.currentDir)
			}
		case "left", "h", "backspace":
			m = m.updateCurrentDir(filepath.Dir(m.currentDir), true)

			return m, getFileInfoCmd(m.currentDir)
		}

	}

	m.table, _ = m.table.Update(msg)

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	if m.loading {
		return "Reading files..."
	}

	viewString := fmt.Sprintf("Current directory: %q\n", m.currentDir)
	viewString += baseStyle.Render(m.table.View()) + "\n"

	return viewString
}

func (m model) updateCurrentDir(dir string, replace bool) model {
	m.currentDir = filepath.Join(m.currentDir, dir)
	if replace {
		m.currentDir = dir
	}
	m.loading = true
	return m
}

func getFileInfoCmd(dir string) tea.Cmd {
	f := func() tea.Msg {
		var res fileInfoResponse
		data, err := fileinfo.GetRootInfo(dir)
		if err != nil {
			res.err = err
			return res
		}

		res.data = data
		return res
	}

	return f
}

func getInitialTable() table.Model {
	t := table.New(table.WithColumns(
		[]table.Column{
			{Title: "Name", Width: 10},
			{Title: "Type", Width: 10},
			{Title: "Size", Width: 10},
		}),
		table.WithRows(
			[]table.Row{
				{"file1.txt", "file", "1000"},
				{"file2.txt", "file", "2000"},
			}),
		table.WithFocused(true),
		table.WithHeight(10))

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)

	return t
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("something went wrong: %v", err)
	}

	m := model{
		currentDir:    cwd,
		loading:       true,
		selectedIndex: 0,
		table:         getInitialTable(),
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("something went wrong %v\n", err)
	}
}
