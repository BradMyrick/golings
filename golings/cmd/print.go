package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/bradmyrick/golings/golings/exercises"
	"github.com/bradmyrick/golings/golings/ui"
)

func PrintHint(infoFile string) {
	exercise, err := exercises.NextPending(infoFile)
	if err != nil {
		color.Red("Failed to find next exercise")
		return
	}
	fmt.Printf("\nHint for %s:\n", exercise.Name)
	color.Yellow(exercise.Hint)
}

func PrintSpecificHint(name string, infoFile string) {
	exercise, err := exercises.Find(name, infoFile)
	if err != nil {
		color.Red("Failed to find exercise: %s", name)
		return
	}
	fmt.Printf("\nHint for %s:\n", exercise.Name)
	color.Yellow(exercise.Hint)
}

func PrintList(infoFile string) {
	exs, err := exercises.List(infoFile)
	if err != nil {
		color.Red("Failed to list exercises")
	}
	ui.PrintList(os.Stdout, exs)
}

func RunNextExercise(infoFile string) {

	exercise, err := exercises.NextPending(infoFile)
	if err != nil {
		color.Green("You've completed all exercises! Great job!")
		return
	}

	RunExercise(exercise, infoFile)
}

func RunExercise(exercise exercises.Exercise, infoFile string) {
	progress, done, total, err := exercises.Progress(infoFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	result, runErr := exercise.Run()

	w, h := ui.GetTerminalSize()
	lastState = &ui.UIState{
		Exercise:       exercise,
		Result:         result,
		RunError:       runErr,
		Progress:       float64(progress),
		Done:           done,
		Total:          total,
		TerminalWidth:  w,
		TerminalHeight: h,
		ShowHint:       false,
	}

	RefreshUI()
}

func RefreshUI() {
	if lastState == nil {
		return
	}
	w, h := ui.GetTerminalSize()
	lastState.TerminalWidth = w
	lastState.TerminalHeight = h
	ClearScreen()

	out := ui.Render(*lastState)
	out = strings.ReplaceAll(out, "\r\n", "\n")
	out = strings.ReplaceAll(out, "\n", "\r\n")

	fmt.Print(out)
}

var lastState *ui.UIState

func MoveToNextAndRun(infoFile string) {
	current, err := exercises.NextPending(infoFile)
	if err == nil {
		exercises.MarkSolved(current.Name)
	}

	next, err := exercises.NextPending(infoFile)
	if err != nil {
		color.Green("\nYou've reached the end! No more exercises.")
		return
	}

	RunExercise(next, infoFile)
}

func ClearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			color.Red("Clear terminal command error")
		}
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			color.Red("Clear terminal command error")
		}
	}
}
