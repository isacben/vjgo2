package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Sample JSON data
	jsonData := `{
        "user": {
            "name": "John",
            "list": [1, 2, "three", 4],
            "escaped": "{\"hello\": \"world\"}",
            "addresses": [
                {
                        "street": "123 Main St",
                        "zipcode": "12345"
                },
                {
                        "street": "456 Oak Ave",
                        "zipcode": "67890"
                }
            ]
        }
    }`

	// Parse JSON
	var data interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		panic(err)
	}

	// Build tree
	tree := BuildTree(data, "", nil)

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	p := tea.NewProgram(
		// model{input_str: strings.Split(json_str, "\n")}, tea.WithAltScreen())
		model{tree: tree}, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	tree             *JSONTree
	visibleLines     *VisibleLines
	firstVisibleLine int
	windowLines      int
	margin           int
	cursorY          int
	ready            bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		{
			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "up", "k":
				{
					m.cursorY--
					if m.cursorY < 0 {
						m.cursorY = 0
					}
					m.ScrollUp()
				}
			case "down", "j":
				{
					m.cursorY++
					if m.cursorY >= len(m.visibleLines.content) {
						m.cursorY = len(m.visibleLines.content) - 1
					}
					m.ScrollDown()
				}
			case "left", "h":
				{
					log.Println("trying to collapse")
					m.tree.Collapse("user.addresses")
					log.Println(m.tree.PrintAsJSONFromRoot())
					m.visibleLines.UpdateContent(m.tree.PrintAsJSONFromRoot())
					m.visibleLines.UpdateVisibleLines(m.visibleLines.firstLine, m.visibleLines.total)
				}
			}
		}

	case tea.WindowSizeMsg:
		{
			if !m.ready {
				//m.cursorY = 0
				//m.firstVisibleLine = 0
				m.margin = 3
				m.windowLines = msg.Height - 1 // for the status bar
				m.visibleLines = NewVisibleLines(
					m.firstVisibleLine, m.windowLines,
					m.tree.PrintAsJSONFromRoot(),
				)
				m.ready = true
			} else {
				m.windowLines = msg.Height - 1 // for the status bar
				m.firstVisibleLine = m.visibleLines.firstLine

				log.Printf("cursor: %d; firstVis: %d; vl_firstVis: %d",
					m.cursorY, m.firstVisibleLine, m.visibleLines.firstLine)

				if m.windowLines <= 2*m.margin+3 {
					m.margin = 0
				} else {
					m.margin = 3
				}

				// fix the cursor at the bottom
				// +3 because the cursor starts at 0, plus the status line
				// plus the firstVisibleLine is 0
				if m.cursorY+3 >= m.firstVisibleLine+m.windowLines {
					// +1 to composate for the status line
					m.firstVisibleLine = m.cursorY - m.windowLines + 1
				}

				m.visibleLines.UpdateVisibleLines(
					m.firstVisibleLine, m.windowLines)
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "loading"
	}
	s := m.Print()
	return s
}

func (m model) Print() string {
	s := ""
	for _, line := range m.visibleLines.linesOnScreen {
		cursor := " "
		if line.num == m.cursorY {
			cursor = ">"
		}
		s += fmt.Sprintf(
			"%d %s %s\n",
			line.num,
			cursor,
			line.content,
		)
	}

	return strings.TrimSuffix(s, "\n")
}

func (m model) ScrollDown() {
	if m.cursorY > m.visibleLines.firstLine+
		m.visibleLines.total-1-m.margin {
		m.firstVisibleLine = m.cursorY - m.windowLines + 1 +
			m.margin

		m.visibleLines.UpdateVisibleLines(m.firstVisibleLine,
			m.windowLines)
	}
}

func (m model) ScrollUp() {
	if m.cursorY < m.visibleLines.firstLine+m.margin {
		m.firstVisibleLine = max(0, m.cursorY-m.margin)

		m.visibleLines.UpdateVisibleLines(m.firstVisibleLine,
			m.windowLines)
	}
}
