package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/mauricioabreu/golings/golings/exercises"
)

func PrintList(o io.Writer, exs []exercises.Exercise) {
	color.New(color.Bold, color.FgWhite).Fprintf(o, "%-17s%-46s%-10s\n", "Name", "Path", "State")
	for _, ex := range exs {
		state := ex.State()
		stateStr := state.String()
		if state == exercises.Done {
			stateStr = color.GreenString(stateStr)
		} else {
			stateStr = color.RedString(stateStr)
		}
		fmt.Fprintf(o, "%-17s%-46s%s\n", ex.Name, ex.Path, stateStr)
	}
}

func PrintInteractiveList(o io.Writer, exs []exercises.Exercise, cursor int, height int) {
	maxLines := height - 5
	if maxLines < 5 {
		maxLines = 5
	}

	start := cursor - maxLines/2
	if start < 0 {
		start = 0
	}
	end := start + maxLines
	if end > len(exs) {
		end = len(exs)
		start = end - maxLines
		if start < 0 {
			start = 0
		}
	}

	var sb strings.Builder
	sb.WriteString(color.New(color.Bold, color.FgWhite).Sprintf("  %-17s%-46s%-10s\r\n", "Name", "Path", "State"))

	for i := start; i < end; i++ {
		ex := exs[i]
		state := ex.State()
		stateStr := state.String()
		if state == exercises.Done {
			stateStr = color.GreenString(stateStr)
		} else {
			stateStr = color.RedString(stateStr)
		}

		cursorStr := "  "
		if i == cursor {
			cursorStr = color.CyanString("> ")
		}

		sb.WriteString(fmt.Sprintf("%s%-17s%-46s%s\r\n", cursorStr, ex.Name, ex.Path, stateStr))
	}

	sb.WriteString("\r\nUse 'j'/'k' or up/down arrows to scroll, 'Enter' to select, 'r' to reset, 'q' to quit list.\r\n")
	fmt.Fprint(o, sb.String())
}
