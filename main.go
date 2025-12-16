package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

    "github.com/mattn/go-isatty"
)

var repeatBuffer string

func main() {
    var args []string
	for _, arg := range os.Args[1:] {
        switch arg {
		case "-h", "--help":
			fmt.Println(usage())
			return
		case "-v", "-V", "--version":
			fmt.Println("vj", version)
			return
        default:
			args = append(args, arg)
		}
	}

    fd := os.Stdin.Fd()
	stdinIsTty := isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)

	var src io.Reader

	if stdinIsTty {
		if len(args) == 0 {
			// $ vj
			fmt.Println(usage())
			return
		} else {
			// $ vj file.json
			filePath := args[0]
            file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
            if err != nil {
                fmt.Printf("Error reading file: %v\n", err)
                os.Exit(1)
            }
            defer file.Close()
            src = file
		}
	} else {
		// $ cat file.json | vj
		src = os.Stdin
	}

    jsonInputBytes, err := io.ReadAll(src)
    if err != nil {
        panic(err)
    }

	// Parse JSON
	var data interface{}
	if err := json.Unmarshal(jsonInputBytes, &data); err != nil {
		panic(err)
	}

	// Build tree
	tree := BuildTree(data, "", nil)
	SetCurrentTheme("dark")

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	p := tea.NewProgram(
		model{tree: tree}, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	tree             *JSONTree
	visibleLines2     *VisibleLines2
    VirtualToRealLines []int
	firstVisibleLine int
    currentPath      string
	windowLines      int
    width            int
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
            case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
				{
					repeatBuffer += fmt.Sprintf("%c", msg.Runes[0])
				}
            case "g":
                {
                    // Move the cursos to the top
                    m.cursorY = 0
					m.ScrollUp()
                }
            case "G":
                {
                    // Move the cursos to the end of the file
                    if len(m.tree.VirtualToRealLines) > 0 {
                        m.cursorY = len(m.tree.VirtualToRealLines) - 1
                        m.ScrollDown()
                    }
                }
			case "up", "k":
				{
                    steps := 1
                    if repeatBuffer != "" {
                        steps = timesToRepeat()
                    }
                    m.cursorY -= steps

					if m.cursorY < 0 {
						m.cursorY = 0
					}

                    physicalLine := m.tree.VirtualToRealLines[m.cursorY]
                    node, exists := m.tree.GetNodeAtLine(physicalLine)
                    if exists {
                        m.currentPath = node.Path
                    }

					m.ScrollUp()
				}
			case "down", "j":
				{
                    steps := 1
                    if repeatBuffer != "" {
                        steps = timesToRepeat()
                    }
                    m.cursorY += steps

					if m.cursorY >= len(m.visibleLines2.content) {
						m.cursorY = len(m.visibleLines2.content) - 1
					}
                    physicalLine := m.tree.VirtualToRealLines[m.cursorY]
                    node, exists := m.tree.GetNodeAtLine(physicalLine)
                    if exists {
                        m.currentPath = node.Path
                    }
					m.ScrollDown()
				}

			case "left", "h":
				{
					physicalLine := m.tree.VirtualToRealLines[m.cursorY]
                    for i,n := range m.tree.VirtualToRealLines {
                        log.Println(i, n)
                    }

					log.Println("collapse physical line:", physicalLine)
					node, exists := m.tree.GetNodeAtLine(physicalLine)
					if exists {
						m.tree.Collapse(node.Path)
						m.visibleLines2.UpdateContent2(m.tree.PrintAsJSON2())
						m.visibleLines2.UpdateVisibleLines2(m.visibleLines2.firstLine,
							m.visibleLines2.total)
					}
				}

			case "right", "l":
				{
					physicalLine := m.tree.VirtualToRealLines[m.cursorY]
					log.Println("expand physical line:", physicalLine)
					node, exists := m.tree.GetNodeAtLine(physicalLine)
					if exists {
						m.tree.Expand(node.Path)
						m.visibleLines2.UpdateContent2(m.tree.PrintAsJSON2())
						m.visibleLines2.UpdateVisibleLines2(m.visibleLines2.firstLine,
							m.visibleLines2.total)
					}
				}
			}
		}

	case tea.WindowSizeMsg:
		{
			if !m.ready {
				m.margin = 3
				m.windowLines = msg.Height - 1 // for the status bar
                m.width = msg.Width
                
				m.tree.PrintAsJSONFromRoot()

				 m.visibleLines2 = NewVisibleLines2(
					m.firstVisibleLine, m.windowLines,
			    	m.tree.PrintAsJSON2(),
				 )


				m.ready = true
			} else {
				m.windowLines = msg.Height - 1 // for the status bar
                m.width = msg.Width
				m.firstVisibleLine = m.visibleLines2.firstLine

				log.Printf("cursor: %d; firstVis: %d; vl_firstVis: %d",
					m.cursorY, m.firstVisibleLine, m.visibleLines2.firstLine)

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

				m.visibleLines2.UpdateVisibleLines2(
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
	s := m.Render()

    // Print ~ on blank lines
    blankLines := m.windowLines - 
        len(m.visibleLines2.content) +
        m.visibleLines2.firstLine

    for range blankLines {
        s += "\n" + blankChar.Render("~")
    }

    s += "\n" + m.UpdateStatusBar()
	return s
}

func (m model) UpdateStatusBar() string {
    s := m.currentPath
    return s
}

func (m model) Render() string {
	s := ""

	for i, line := range m.visibleLines2.linesOnScreen {
		// Print line at cursor
		if i+m.visibleLines2.firstLine == m.cursorY {
			num := m.tree.VirtualToRealLines[m.cursorY] + 1

			s += fmt.Sprintf(
				"%s %s \n",
				lineNumbersCol.Render(strconv.Itoa(num)+" "),
                RenderLine(line, true),
			)
		}

		// Print lines before cursor
		if i+m.visibleLines2.firstLine < m.cursorY {
			num := (m.cursorY - m.visibleLines2.firstLine) - i
			s += fmt.Sprintf(
				"%s %s \n",
				lineNumbersCol.Render(strconv.Itoa(num)),
				RenderLine(line, false),
			)
		}

		// Print lines after cursor
		if i+m.visibleLines2.firstLine > m.cursorY {
			num := i - (m.cursorY - m.visibleLines2.firstLine)
			s += fmt.Sprintf(
				"%s %s \n",
				lineNumbersCol.Render(strconv.Itoa(num)),
				RenderLine(line, false),
			)
		}
	}

	return strings.TrimSuffix(s, "\n")
}

func RenderLine(line LineMetadata, hasCursor bool) string {
    isSelected := false
	indent := strings.Repeat("  ", line.Indent)

	switch line.LineType {
	case ContentWithBrace:
        if line.IsArrayElement {
            if line.IsCollapsed {
                comma := ""
                if !line.IsLastChild {
                    comma = ","
                }

                return RenderIndent(indent, isSelected) +
                    RenderSyntax("{", hasCursor, isSelected) +
                    RenderSyntax("...}" + comma, false, isSelected)
            }

            return RenderIndent(indent, isSelected) + 
                RenderSyntax(line.BracketChar, hasCursor, isSelected)

        } else if line.Key != "" && !line.IsCollapsed {
			// Key with opening bracket: "user": {
			return RenderIndent(indent, isSelected) +
                RenderSyntax(`"`, hasCursor, isSelected) +
                RenderKey(line.Key, isSelected) +
                RenderSyntax(`": ` + line.BracketChar, false, isSelected)

		} else if line.IsCollapsed {
			// Collapsed: "user": {...} or "items": [...]
			//keyPart := keyStyle.Render(`"` + line.Key + `"`)
			collapsedContent := ""
            comma := ""

            if !line.IsLastChild {
                comma = ","
            }

			if line.BracketChar == "{" {
				collapsedContent = "{...}" + comma
			} else {
				collapsedContent = "[...]" + comma
			}

			return RenderIndent(indent, isSelected) + 
                RenderSyntax(`"`, hasCursor, isSelected) +
                RenderKey(line.Key, isSelected) + 
                RenderSyntax(`": ` + collapsedContent, false, isSelected)
		} //else {
			// Just opening bracket
		//	return indent + line.BracketChar
		//}

    case OpenBracket:
        if line.Key == "" {
            if  !line.IsCollapsed {
                return RenderIndent(indent, isSelected) + 
                    RenderSyntax(line.BracketChar, hasCursor, isSelected)
            } else {
                collapsedContent := ""
                if line.BracketChar == "{" {
                    collapsedContent = "...}"
                } else {
                    collapsedContent = "...]"
                }
                return RenderIndent(indent, isSelected) + 
                    RenderSyntax(line.BracketChar, hasCursor, isSelected) +
                    RenderSyntax(collapsedContent, false, isSelected)
            }
        }

	case CloseBracket:
		comma := ""
		if !line.IsLastChild {
			comma = ","
		}
		return RenderIndent(indent, isSelected) + 
            RenderSyntax(line.BracketChar, hasCursor, isSelected) + 
            RenderSyntax(comma, false, isSelected)

	case ContentLine:
		comma := ""
		if !line.IsLastChild {
			comma = ","
		}

		if line.IsArrayElement {
			// Array element: just the value
			switch line.NodeType {
			case StringType:
				return RenderIndent(indent, isSelected) + 
                    RenderString(`"` +
                        line.Content +`"`, 
                        hasCursor, isSelected) +
                    RenderSyntax(comma, false, isSelected)
			case NumberType:
				return RenderIndent(indent, isSelected) + 
                    RenderNumber(line.Content, hasCursor, isSelected) + 
                    RenderSyntax(comma, false, isSelected)
			case BoolType:
				return RenderIndent(indent, isSelected) + 
                    RenderBoolean(line.Content, hasCursor, isSelected) +
                    RenderSyntax(comma, false, isSelected)
			case NullType:
				return RenderIndent(indent, isSelected) + 
                    RenderNull("null", hasCursor, isSelected) + 
                    RenderSyntax(comma, false, isSelected)
			}

		} else {
			// Object property: "key": value
			valuePart := ""

			switch line.NodeType {
			case StringType:
				valuePart = stringStyle.Render(`"` + line.Content + `"`)
			case NumberType:
				valuePart = numberStyle.Render(line.Content)
			case BoolType:
				valuePart = booleanStyle.Render(line.Content)
			case NullType:
				valuePart = nullStyle.Render("null")
			}

			return RenderIndent(indent, isSelected) +
                RenderSyntax(`"`, hasCursor, isSelected) +
                RenderKey(line.Key, isSelected) +
                RenderSyntax(`": ` + valuePart + comma, false, isSelected)
		}
	}

	return line.Content
}

func (m model) ScrollDown() {
	if m.cursorY > m.visibleLines2.firstLine+
		m.visibleLines2.total-1-m.margin {
		m.firstVisibleLine = m.cursorY - m.windowLines + 1 +
			m.margin

		m.visibleLines2.UpdateVisibleLines2(m.firstVisibleLine,
			m.windowLines)
	}
}

func (m model) ScrollUp() {
	if m.cursorY < m.visibleLines2.firstLine+m.margin {
		m.firstVisibleLine = max(0, m.cursorY-m.margin)

		m.visibleLines2.UpdateVisibleLines2(m.firstVisibleLine,
			m.windowLines)
	}
}

func timesToRepeat() int {
	number, err := strconv.Atoi(repeatBuffer)

	if err != nil {
		log.Fatal("Error converting string to int:", err)
	}

	repeatBuffer = ""
	return number
}
