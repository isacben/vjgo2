package main

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	lineNumbersCol lipgloss.Style
    blankChar      lipgloss.Style
	cursorStyle    lipgloss.Style
	keyStyle       lipgloss.Style
	stringStyle    lipgloss.Style
	nullStyle      lipgloss.Style
	booleanStyle   lipgloss.Style
	numberStyle    lipgloss.Style
	syntaxStyle    lipgloss.Style
    statusBarStyle lipgloss.Style
    errorStyle     lipgloss.Style
)

type Color string

type Theme struct {
	Cursor     Color
	StatusBar  Color
	Key        Color
	String     Color
	Null       Color
	Boolean    Color
	Number     Color
	LineNumber Color
    Syntax     Color
    Error      Color
}

var (
	currentTheme = themes["nocolor"]

	defaultCursor     = Color("#bb9af7")
	defaultStatusBar  = Color("#414868")
	defaultKey        = Color("#7dcfff")
	defaultString     = Color("#9ece6a")
	defaultNull       = Color("#565f89")
	defaultBoolean    = Color("#ff9e64")
	defaultNumber     = Color("#ff9e64")
	defaultLineNumber = Color("#565f89")
	defaultSyntax     = Color("")
    defaultError      = Color("9")
)

var themes = map[string]Theme{
	"nocolor": {
		Cursor:     Color(""),
		StatusBar:  Color(""),
		Key:        Color(""),
		String:     Color(""),
		Null:       Color(""),
		Boolean:    Color(""),
		Number:     Color(""),
		LineNumber: Color(""),
	},
	"dark": {
		Cursor:     defaultCursor,
		StatusBar:  defaultStatusBar,
		Key:        defaultKey,
		String:     defaultString,
		Null:       defaultNull,
		Boolean:    defaultBoolean,
		Number:     defaultNumber,
		LineNumber: defaultLineNumber,
		Syntax:     defaultSyntax,
        Error:      defaultError,
	},
	 "light": {
	 	Cursor:     Color("#0066cc"),
	 	StatusBar:  Color("4"),
	 	Key:        Color("#0066cc"),
	 	String:     Color("#22863a"),
	 	Null:       Color("#6f42c1"),
	 	Boolean:    Color("#d73a49"),
	 	Number:     Color("#005cc5"),
	 	LineNumber: Color("#586069"),
        Error:      Color("9"),
	 },
}

func SetCurrentTheme(name string) {
	currentTheme = themes[name]

	lineNumbersCol = lipgloss.NewStyle().
		Align(lipgloss.Right).
		Width(5).
		Foreground(lipgloss.Color(currentTheme.LineNumber))

	cursorStyle = lipgloss.NewStyle().
		Reverse(true)

	keyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(currentTheme.Key))

	stringStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(currentTheme.String))

	nullStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(currentTheme.Null))

	booleanStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(currentTheme.Boolean))

	numberStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(currentTheme.Number))

	syntaxStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(currentTheme.Syntax))

    statusBarStyle = lipgloss.NewStyle().
        Align(lipgloss.Bottom).
        Background(lipgloss.Color(currentTheme.StatusBar))

    blankChar = lipgloss.NewStyle().
        Foreground(lipgloss.Color(currentTheme.LineNumber))

    errorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color(currentTheme.Error))
        
}

func RenderIndent(text string, selected bool) string {
    if selected {
        return lipgloss.NewStyle().
            Background(lipgloss.Color("#414868")).Render(text)
    }

    return text
}

func RenderKey(text string, selected bool) string {
    if selected {
        return keyStyle.Background(lipgloss.Color("#414868")).Render(text)
    }

    return keyStyle.Render(text)
}

func RenderSyntax(text string, hasCursor bool, isSelected bool) string {
    return RenderElement(
        text, hasCursor, isSelected, syntaxStyle)
}


func RenderString(text string, hasCursor bool, isSelected bool) string {
    return RenderElement(
        text, hasCursor, isSelected, stringStyle)
}

func RenderNumber(text string, hasCursor bool, isSelected bool) string {
    return RenderElement(
        text, hasCursor, isSelected, numberStyle)
}

func RenderBoolean(text string, hasCursor bool, isSelected bool) string {
    return RenderElement(
        text, hasCursor, isSelected, booleanStyle)
}

func RenderNull(text string, hasCursor bool, isSelected bool) string {
    return RenderElement(
        text, hasCursor, isSelected, nullStyle)
}

func RenderElement(text string, hasCursor bool, selected bool, style lipgloss.Style) string {
    if hasCursor {
        cursor := cursorStyle.Render(text[:1])

        if selected {
            return cursor +
                style.Background(lipgloss.Color("#414868")).Render(text[1:])
        }
        
        return cursor + style.Render(text[1:])
    }

    if selected {
        return style.Background(lipgloss.Color("#414868")).Render(text)
    }

    return style.Render(text)
}
