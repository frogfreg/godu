package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/frogfreg/godu/fileinfo"
)

type model struct {
	currentDir    string
	fileMap       map[string]fileinfo.FileInfo
	loading       bool
	deleting      bool
	deleteDir     string
	selectedIndex int
	err           error
	table         table.Model
}

type fileInfoResponse struct {
	data    []fileinfo.FileInfo
	fileMap map[string]fileinfo.FileInfo
	err     error
}

type deleteResponse struct {
	err error
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m model) Init() tea.Cmd {
	return getFileInfoCmd(nil, m.currentDir)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case fileInfoResponse:
		m.loading = false
		fir := fileInfoResponse(msg)
		m.err = fir.err
		if m.err != nil {
			return m, nil
		}

		m.fileMap = fir.fileMap

		m.table.SetRows(fileinfo.FileInfosToRow(fir.data))
	case deleteResponse:
		m.deleting = false
		m.err = deleteResponse(msg).err
		if m.err != nil {
			return m, nil
		}

		m.loading = true
		return m, getFileInfoCmd(m.fileMap, m.currentDir)

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
		case "enter", "right", "l":
			if !m.loading && m.err != nil {
				m.err = nil
				return m, nil
			}
			sr := m.table.SelectedRow()

			if !m.loading && !m.deleting && len(sr) != 0 && sr[1] == "dir" {
				m = m.updateCurrentDir(m.table.SelectedRow()[0], false)
				return m, getFileInfoCmd(m.fileMap, m.currentDir)
			}
		case "left", "h", "backspace":
			if !m.loading && !m.deleting {
				m = m.updateCurrentDir(filepath.Dir(m.currentDir), true)
				return m, getFileInfoCmd(m.fileMap, m.currentDir)
			}
		case "d":
			sr := m.table.SelectedRow()

			if !m.loading && !m.deleting && len(sr) != 0 {
				m.deleting = true
				m.deleteDir = filepath.Join(m.currentDir, sr[0])
				return m, deleteCmd(m.fileMap, m.deleteDir)
			}
		}

	}

	m.table, _ = m.table.Update(msg)

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("%v\n\nPress enter to continue", m.err.Error())
	}

	if m.deleting {
		return fmt.Sprintf("deleting %v...", m.deleteDir)
	}

	if m.loading {
		return "Reading files..."
	}

	viewString := "Controls:\n"
	viewString += "d: delete dir/file\nh,left arrow: go up a directory\nl, right arrow, enter: go down selected dir\n"
	viewString += "j,down arrow: move down\nk,up arrow: move up\n"
	viewString += fmt.Sprintf("\nCurrent directory: %q\n", m.currentDir)
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

func deleteCmd(m map[string]fileinfo.FileInfo, path string) tea.Cmd {
	f := func() tea.Msg {
		var res deleteResponse

		err := os.RemoveAll(path)
		if err != nil {
			res.err = err
		}

		fileinfo.CleanChildren(m, path)

		return res
	}

	return f
}

func getFileInfoCmd(m map[string]fileinfo.FileInfo, dir string) tea.Cmd {
	return func() tea.Msg {
		var res fileInfoResponse
		fm, err := fileinfo.GenerateFileMap(m, dir)
		if err != nil {
			res.err = err
			return res
		}
		data := fileinfo.GetSortedDirs(fm, dir)
		res.fileMap = fm

		res.data = data
		return res
	}
}

func getInitialTable() table.Model {
	t := table.New(table.WithColumns(
		[]table.Column{
			{Title: "Name", Width: 30},
			{Title: "Type", Width: 7},
			{Title: "Size", Width: 15},
		}),
		table.WithRows(
			[]table.Row{}),
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

	// If I slog anywhere in the model without first enabling the following code, the view gets messed up for some reason

	// f, err := tea.LogToFile("debug.log", "debug")
	// if err != nil {
	// 	log.Fatalf("fatal: %v", err)
	// }
	// defer f.Close()

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("something went wrong %v\n", err)
	}
}
