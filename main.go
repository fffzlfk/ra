package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	currentDir   string
	entrys       []os.DirEntry
	cursor       int
	selectedPath os.DirEntry
}

func (m *model) updateFiles() error {
	entrys, err := os.ReadDir(m.currentDir)
	if err != nil {
		log.Fatal(err)
		return err
	}
	m.entrys = entrys
	return nil
}

func initalModel() model {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	m := model{
		currentDir: currentDir,
	}
	err = m.updateFiles()
	if err != nil {
		panic(err)
	}
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.entrys)-1 {
				m.cursor++
			}
		case "enter", "l":
			if len(m.entrys) <= 0 {
				return m, nil
			}
			m.selectedPath = m.entrys[m.cursor]
			if m.selectedPath.IsDir() {
				m.currentDir += "/" + m.selectedPath.Name()
			} else {
				return m, nil
			}
			m.updateFiles()
			m.cursor = 0
			log.Printf("selected path: %s", m.selectedPath.Name())
			return m, nil
		}
	}
	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("current dir: %s\n", m.currentDir)
	for i, entry := range m.entrys {
		cursor := " "
		if m.cursor == i {
			cursor = ">" // cursor!
		}
		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, filepath.Base(entry.Name()))
	}
	return s
}

func main() {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	p := tea.NewProgram(initalModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
