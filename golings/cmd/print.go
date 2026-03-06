package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
	"github.com/mauricioabreu/golings/golings/exercises"
	"github.com/mauricioabreu/golings/golings/ui"
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
	ClearScreen()
	exs, err := exercises.List(infoFile)
	if err != nil {
		color.Red("Failed to list exercises")
	}
	width, _ := ui.GetTerminalSize()
	ui.PrintList(os.Stdout, exs, width)
}

func RunNextExercise(infoFile string) {
	ClearScreen()

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
	fmt.Print(ui.Render(*lastState))
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

	ClearScreen()
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
