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
	ui.PrintList(os.Stdout, exs)
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
	} else {
		color.Blue("Progress: %d/%d (%.2f%%)\n\n", done, total, progress*100)
	}

	fmt.Printf("Current exercise: %s\n\n", exercise.Path)

	result, err := exercise.Run()
	if err != nil {
		color.Red("Testing of %s failed! Please try again. Here is the output:\n", exercise.Path)
		fmt.Println("")
		color.Red(result.Err)
		color.Red(result.Out)
		fmt.Println("")
		color.Yellow("If you feel stuck, press 'h' for a hint")
	} else {
		color.Green("Successfully ran %s!", exercise.Path)
		fmt.Println("")
		color.Cyan(result.Out)
		color.Yellow("\nExercise done! Press 'n' to move to the next one.")
	}

	fmt.Print("\n[n]ext [h]int [l]ist [q]uit: ")
}

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
