package main

import (
	"strings"
)

type line struct {
	num     int
	content string
}

type VisibleLines struct {
	firstLine     int
	total         int
	content       []line
	linesOnScreen []line
}

func NewVisibleLines(firstLine int, total int, content string) *VisibleLines {
	vl := &VisibleLines{}
	vl.UpdateContent(content)
	vl.UpdateVisibleLines(firstLine, total)
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
