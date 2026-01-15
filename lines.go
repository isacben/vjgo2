package main

import (
	"strings"
)

type line struct {
	num     int
	content string
}

type LineType string

const (
	ContentLine      LineType = "content"
	ContentWithBrace LineType = "content_with_brace"
	OpenBracket      LineType = "open_bracket"
	CloseBracket     LineType = "close_bracket"
)

type LineMetadata struct {
	LineNumber     int
	LineType       LineType
	Content        string
	NodePath       string
	NodeType       NodeType
	Key            string
	Value          interface{}
	Indent         int
	IsCollapsed    bool
	HasChildren    bool
	BracketChar    string // "{", "}", "[", "]"
	IsArrayElement bool
	IsLastChild    bool // for comma handling
}

type VisibleLines struct {
	firstLine     int
	total         int
	content       []line
	linesOnScreen []line
}

type VisibleLines2 struct {
	firstLine     int
	total         int
	content       []LineMetadata
	linesOnScreen []LineMetadata
}

func NewVisibleLines(firstLine int, total int, content string) *VisibleLines {
	vl := &VisibleLines{}
	vl.UpdateContent(content)
	vl.UpdateVisibleLines(firstLine, total)
	return vl
}

func NewVisibleLines2(firstLine int, total int, content []LineMetadata) *VisibleLines2 {
	vl := &VisibleLines2{}
	vl.UpdateContent2(content)
	vl.UpdateVisibleLines2(firstLine, total)
	return vl
}

func (vl *VisibleLines) UpdateContent(content string) {
	lines := strings.Split(content, "\n")
	// clear slice
	vl.content = vl.content[:0]
	for i, l := range lines {
		ln := line{i, l}
		vl.content = append(vl.content, ln)
	}
}

func (vl *VisibleLines2) UpdateContent2(content []LineMetadata) {
	// clear slice
	vl.content = vl.content[:0]

	// equivalent to:
	// for _,line := range content {
	// 	vl.content = append(vl.content, line)
	// }
	vl.content = append(vl.content, content...)
}

func (vl *VisibleLines) UpdateVisibleLines(firstLine int, total int) {
	vl.firstLine = firstLine
	vl.total = total

	// clear slice
	vl.linesOnScreen = vl.linesOnScreen[:0]

	// copy the lines that should be on the screen
	for _, line := range vl.content {
		if line.num >= vl.firstLine && len(vl.linesOnScreen) < vl.total {
			vl.linesOnScreen = append(vl.linesOnScreen, line)
		}
	}
}

func (vl *VisibleLines2) UpdateVisibleLines2(firstLine int, total int) {
	vl.firstLine = firstLine
	vl.total = total

	// clear slice
	vl.linesOnScreen = vl.linesOnScreen[:0]

	// copy the lines that should be on the screen
	for _, line := range vl.content {
		if line.LineNumber >= vl.firstLine && len(vl.linesOnScreen) < vl.total {
			vl.linesOnScreen = append(vl.linesOnScreen, line)
		}
	}
}
