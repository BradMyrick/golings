package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/mauricioabreu/golings/golings/exercises"
	"github.com/mauricioabreu/golings/golings/ui"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func VerifyCmd(infoFile string) *cobra.Command {
	return &cobra.Command{
		Use:   "verify",
		Short: "Verify all exercises",
		Run: func(cmd *cobra.Command, args []string) {
			allExercises, err := exercises.List(infoFile)
			if err != nil {
				color.Red(err.Error())
				os.Exit(1)
			}

			width, _ := ui.GetTerminalSize()
			barWidth := width / 2
			if barWidth < 10 {
				barWidth = 10
			}
			bar := progressbar.NewOptions(
				len(allExercises),
				progressbar.OptionSetWidth(barWidth),
				progressbar.OptionEnableColorCodes(true),
				progressbar.OptionSetPredictTime(false),
				progressbar.OptionSetElapsedTime(false),
				progressbar.OptionSetDescription("[cyan][reset] Running exercises"),
				progressbar.OptionSetTheme(progressbar.Theme{
					Saucer:        "[yellow]=[reset]",
					SaucerHead:    "[yellow]>[reset]",
					SaucerPadding: " ",
					BarStart:      "[",
					BarEnd:        "]",
				}),
			)
			if err := bar.RenderBlank(); err != nil {
				color.Red(err.Error())
				os.Exit(1)
			}

			for _, exercise := range allExercises {
				result, runErr := exercise.Run()
				if result.Exercise.State() == exercises.Pending || runErr != nil {
					fmt.Print("\n\n")
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
					os.Exit(1)
				}
				bar.Add(1) // nolint
			}

			color.Green("\n\nCongratulations!!!")
			color.Green("You passed all the exercises")
		},
	}
}
