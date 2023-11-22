package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/frogfreg/du-test/fileinfo"
)

type model struct {
	currentDir    string
	infoList      []fileinfo.FileInfo
	loading       bool
	selectedIndex int
	err           error
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
		m.infoList = fir.data
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "down", "j":
			if m.selectedIndex < len(m.infoList) {
				m.selectedIndex++
			}
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
		}

	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	if m.loading {
		return "Reading files..."
	}

	viewString := fmt.Sprintf("current directory %q\n", m.currentDir)

	for i, fi := range m.infoList {

		pad := "   "
		fiString := fmt.Sprintf("%v | %v | %v bytes\n", filepath.Base(fi.Name), fi.FileType, fi.Size)

		if i == m.selectedIndex {
			pad = "-> "
		}

		fiString = pad + fiString

		viewString += fiString
	}

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
		selectedIndex: 0}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("something went wrong %v\n", err)
	}
}
