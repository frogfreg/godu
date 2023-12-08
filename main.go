package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/frogfreg/du-test/fileinfo"
)

type model struct {
	currentDir    string
	infoList      []fileinfo.FileInfo
	loading       bool
	selectedIndex int
	err           error
	table         table.Model
}

type fileInfoResponse struct {
	data []fileinfo.FileInfo
	err  error
}

func (m model) Init() tea.Cmd {
	f := func() tea.Msg {
		var res fileInfoResponse
		data, err := fileinfo.GetRootInfo(m.currentDir)
		if err != nil {
			res.err = err
			return res
		}

		res.data = data
		return res
	}

	return f
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case fileInfoResponse:
		m.loading = false
		fir := fileInfoResponse(msg)
		if fir.err != nil {
			m.err = fir.err
			return m, tea.Quit
		}
		// m.infoList = fir.data

		m.table.SetRows(fileinfo.FileInfosToRow(fir.data))
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
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

	viewString := fmt.Sprintf("Current directory: %q\n\n", m.currentDir)
	viewString += m.table.View() + "\n"

	return viewString
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	m := model{
		currentDir:    cwd,
		infoList:      []fileinfo.FileInfo{},
		loading:       true,
		selectedIndex: 0,
		table: table.New(table.WithColumns(
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
			table.WithHeight(10)),
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("something went wrong %v\n", err)
	}
}
