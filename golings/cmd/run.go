package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/bradmyrick/golings/golings/exercises"
	"github.com/bradmyrick/golings/golings/ui"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func RunCmd(infoFile string) *cobra.Command {
	return &cobra.Command{
		Use:   "run next | <exercise name>",
		Short: "Run a single exercise",
		Long: `example next pending exercise : golings run next
example specific exercise : golings run variables1`,
		Args:          cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var exercise exercises.Exercise
			var err error
			if args[0] == "next" {
				exercise, err = exercises.NextPending(infoFile)
			} else {
				exercise, err = exercises.Find(args[0], infoFile)
			}

			spinner := RunSpinner(exercise.Name)

			if errors.Is(err, exercises.ErrExerciseNotFound) {
				color.White("No exercise found for '%s'", args[0])
				return err
			}

			result, runErr := exercise.Run()

			spinner.Close()

			progress, done, total, _ := exercises.Progress(infoFile)
			w, h := ui.GetTerminalSize()
			state := ui.UIState{
				Exercise:       exercise,
				Result:         result,
				RunError:       runErr,
				Progress:       float64(progress),
				Done:           done,
				Total:          total,
				TerminalWidth:  w,
				TerminalHeight: h,
			}
			fmt.Print(ui.Render(state))

			return runErr
		},
	}
}

func RunSpinner(exercise string) *progressbar.ProgressBar {
	spinner := progressbar.NewOptions(
		-1, // a negative number makes turns the progress bar into a spinner
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription(color.WhiteString("Running exercise: %s", exercise)),
		progressbar.OptionOnCompletion(func() {
			color.White("\nRunning complete!\n\n")
		}),
	)
	go func() {
		for x := 0; x < 100; x++ {
			spinner.Add(1) // nolint
			time.Sleep(250 * time.Millisecond)
		}
	}()

	return spinner
}
