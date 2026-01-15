package main

import (
	"fmt"
	"strings"
)

func (m *model) performSearch() {
	if m.searchBuffer == "" {
		return
	}

	m.searchResults = []SearchMatch{}
	searchTerm := strings.ToLower(m.searchBuffer)

	// Search through all visible nodes
	for virtualLine, realLine := range m.tree.VirtualToRealLines {
		node, exists := m.tree.GetNodeAtLine(realLine)
		if !exists {
			continue
		}

		// Search in key
		if node.Key != "" && strings.Contains(strings.ToLower(node.Key), searchTerm) {
			m.searchResults = append(m.searchResults, SearchMatch{
				VirtualLine: virtualLine,
				Path:        node.Path,
				MatchType:   "key",
				Content:     node.Key,
			})
		}

		// Search in value (only for primitive types)
		if node.Type != ObjectType && node.Type != ArrayType {
			valueStr := nodeValueToString(node)
			if valueStr != "" && strings.Contains(strings.ToLower(valueStr), searchTerm) {
				m.searchResults = append(m.searchResults, SearchMatch{
					VirtualLine: virtualLine,
					Path:        node.Path,
					MatchType:   "value",
					Content:     valueStr,
				})
			}
		}
	}

	// Find the first match at or after current cursor position
	if len(m.searchResults) > 0 {
		firstMatchIndex := m.findFirstMatchFromCursor()

		// If no match found at or after cursor, wrap to first match
		if firstMatchIndex == -1 {
			firstMatchIndex = 0
		}

		m.currentMatchIndex = firstMatchIndex
	} else {
		m.currentMatchIndex = 0
	}

	m.updateSearchStatusBar()
}

func (m *model) navigateToNextMatch() {
	if len(m.searchResults) == 0 {
		return
	}

	// Find first match AFTER current cursor position
	for i, match := range m.searchResults {
		if match.VirtualLine > m.cursorY {
			m.currentMatchIndex = i
			m.navigateToMatch(i)
			return
		}
	}

	// Wrap to first match
	m.currentMatchIndex = 0
	m.navigateToMatch(0)
}

func (m *model) navigateToPreviousMatch() {
	if len(m.searchResults) == 0 {
		return
	}

	// Find last match BEFORE current cursor position
	for i := len(m.searchResults) - 1; i >= 0; i-- {
		if m.searchResults[i].VirtualLine < m.cursorY {
			m.currentMatchIndex = i
			m.navigateToMatch(i)
			return
		}
	}

	// Wrap to last match
	m.currentMatchIndex = len(m.searchResults) - 1
	m.navigateToMatch(len(m.searchResults) - 1)
}

func (m *model) navigateToMatch(index int) {
	if index < 0 || index >= len(m.searchResults) {
		return
	}

	match := m.searchResults[index]
	m.cursorY = match.VirtualLine
	m.updateCurrentPath()
	m.ScrollDown()
	m.ScrollUp()
	m.updateSearchStatusBar()
}

// findFirstMatchFromCursor finds the first match at or after cursor position
// This version includes matches on the current line
func (m *model) findFirstMatchFromCursor() int {
	for i, match := range m.searchResults {
		if match.VirtualLine >= m.cursorY {
			return i
		}
	}
	return -1 // No match found at or after cursor
}

func (m *model) updateSearchStatusBar() {
	if len(m.searchResults) == 0 {
		m.statusBar = "Pattern not found: " + m.searchBuffer
	} else {
		m.statusBar = fmt.Sprintf("/%s [%d/%d]",
			m.searchBuffer, m.currentMatchIndex+1, len(m.searchResults))
	}
}

func nodeValueToString(node *Node) string {
	switch node.Type {
	case StringType:
		if str, ok := node.Value.(string); ok {
			return str
		}

	case NumberType:
		return fmt.Sprintf("%v", node.Value)

	case BoolType:
		if b, ok := node.Value.(bool); ok {
			return fmt.Sprintf("%t", b)
		}

	case NullType:
		return "null"

	case ObjectType, ArrayType:
		// For objects/arrays, we might want to search in their string representation
		// or skip them entirely for basic search
		return ""
	}

	return fmt.Sprintf("%v", node.Value)
}
