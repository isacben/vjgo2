package main

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	lineNumbersCol lipgloss.Style
	cursorStyle    lipgloss.Style
	keyStyle       lipgloss.Style
	stringStyle    lipgloss.Style
	nullStyle      lipgloss.Style
	booleanStyle   lipgloss.Style
	numberStyle    lipgloss.Style
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
}

var (
	currentTheme = themes["nocolor"]

	defaultCursor     = Color("#bb9af7")
	defaultStatusBar  = Color("5")
	defaultKey        = Color("#7dcfff")
	defaultString     = Color("#9ece6a")
	defaultNull       = Color("#565f89")
	defaultBoolean    = Color("#ff9e64")
	defaultNumber     = Color("#ff9e64")
	defaultLineNumber = Color("#565f89")
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
		// Background(lipgloss.Color("#414868"))

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
}
