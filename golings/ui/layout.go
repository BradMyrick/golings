package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mauricioabreu/golings/golings/exercises"
	"golang.org/x/term"
)

func GetTerminalSize() (int, int) {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		println(err.Error())

		return 80, 24
	}
	return w, h
}

type UIState struct {
	Exercise       exercises.Exercise
	Result         exercises.Result
	RunError       error
	Progress       float64
	Done           int
	Total          int
	TerminalWidth  int
	TerminalHeight int
}

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("33")).
			MarginBottom(1)

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("42")).
		// Padding(0, 1).
		MarginBottom(1)
		// Border(lipgloss.NormalBorder()).
		// BorderForeground(lipgloss.Color("42"))

	errorHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("196")).
		// Padding(0, 1).
		MarginBottom(1)
		// Border(lipgloss.NormalBorder()).
		// BorderForeground(lipgloss.Color("196"))

	hintStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("220")).
			MarginTop(1)

	menuStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)

	progressStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("63")).
			MarginBottom(1)

	outStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("253")).
			MarginBottom(1)
)

func Render(state UIState) string {
	w := state.TerminalWidth
	if w <= 0 {
		w = 80
	}

	doc := strings.Builder{}

	// Progress (full width)
	progressText := fmt.Sprintf("Progress: %d/%d (%.2f%%)", state.Done, state.Total, state.Progress*100)
	doc.WriteString(progressStyle.Width(w).Render(progressText) + "\n")

	// Current Exercise (full width)
	header := headerStyle.Width(w).Render(fmt.Sprintf("Current exercise: %s", state.Exercise.Path))
	doc.WriteString(header + "\n")

	if state.RunError != nil {
		// Failure box spans width, inner text indented
		msg := fmt.Sprintf("Testing of %s failed! Please try again. Here is the output:", state.Exercise.Path)
		doc.WriteString(
			errorHeaderStyle.
				Width(w).
				Render(msg) + "\n",
		)

		if state.Result.Err != "" {
			doc.WriteString(
				outStyle.
					Width(w).
					Render(state.Result.Err) + "\n",
			)
		}
		if state.Result.Out != "" {
			doc.WriteString(
				outStyle.
					Width(w).
					Render(state.Result.Out) + "\n",
			)
		}

		doc.WriteString(
			hintStyle.
				Width(w).
				Render("If you feel stuck, press 'h' for a hint") + "\n",
		)
	} else {
		// Success
		msg := fmt.Sprintf("Successfully ran %s!", state.Exercise.Path)
		doc.WriteString(
			successStyle.
				Width(w).
				Render(msg) + "\n",
		)

		if state.Result.Out != "" {
			doc.WriteString(
				outStyle.
					Width(w).
					Render(state.Result.Out) + "\n",
			)
		}
		doc.WriteString(
			hintStyle.
				Width(w).
				Render("Exercise done! Press 'n' to move to the next one.") + "\n",
		)
	}

	// Menu bar: full width, inner text right-aligned
	menu := menuStyle.
		Width(w).
		Render(
			lipgloss.NewStyle().
				Width(w). // small inner width so border plus text fits
				Align(lipgloss.Center).
				Render("[n]ext [h]int [l]ist [q]uit"),
		)

	doc.WriteString(menu + "\n")

	return doc.String()
}
