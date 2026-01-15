package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO (isaac): delete this test
func TestNewVisibleLines(t *testing.T) {
	tests := []struct {
		name      string
		firstLine int
		total     int
		content   string
		expected  []line
	}{
		{
			"all lines visible",
			0,
			10,
			"line0\nline1\nline2",
			[]line{{0, "line0"}, {1, "line1"}, {2, "line2"}},
		},
		{
			"two lines visible",
			1,
			2,
			"line0\nline1\nline2",
			[]line{{1, "line1"}, {2, "line2"}},
		},
		{
			"four lines visible",
			2,
			4,
			"line0\nline1\nline2\nline3\n" +
				"line4\nline5\nline6\nline7\nline7",
			[]line{
				{2, "line2"}, {3, "line3"}, {4, "line4"},
				{5, "line5"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vl := NewVisibleLines(tt.firstLine, tt.total, tt.content)
			assert.Equal(t, tt.expected, vl.linesOnScreen)
		})
	}
}

func TestNewVisibleLines2(t *testing.T) {
	tests := []struct {
		name      string
		firstLine int
		total     int
		content   []LineMetadata
		expected  []string
	}{
		{
			"all lines visible",
			0,
			10,
			[]LineMetadata{
				{LineNumber: 0, Content: "line0"},
				{LineNumber: 1, Content: "line1"},
				{LineNumber: 2, Content: "line2"}},
			[]string{"line0", "line1", "line2"},
		},
		{
			"two lines visible",
			1,
			2,
			[]LineMetadata{
				{LineNumber: 0, Content: "line0"},
				{LineNumber: 1, Content: "line1"},
				{LineNumber: 2, Content: "line2"}},
			[]string{"line1", "line2"},
		},
		{
			"four lines visible",
			2,
			4,
			[]LineMetadata{
				{LineNumber: 0, Content: "line0"},
				{LineNumber: 1, Content: "line1"},
				{LineNumber: 2, Content: "line2"},
				{LineNumber: 3, Content: "line3"},
				{LineNumber: 4, Content: "line4"},
				{LineNumber: 5, Content: "line5"},
				{LineNumber: 6, Content: "line6"},
				{LineNumber: 7, Content: "line7"},
				{LineNumber: 8, Content: "line7"}},
			[]string{"line2", "line3", "line4", "line5"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vl := NewVisibleLines2(tt.firstLine, tt.total, tt.content)

			// extract the Content of each LineMetadata struct
			var actual []string
			for _, line := range vl.linesOnScreen {
				actual = append(actual, line.Content)
			}

			assert.Equal(t, tt.expected, actual)
		})
	}
}
