package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/mauricioabreu/golings/golings/exercises"
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

			bar := progressbar.NewOptions(
				len(allExercises),
				progressbar.OptionSetWidth(50),
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
				result, _ := exercise.Run()
				if result.Exercise.State() == exercises.Pending || result.Err != "" {
					fmt.Print("\n\n")
					color.Cyan("Failed to compile the exercise %s\n\n", exercise.Path)
					color.White("Check the output below: \n\n")
					if result.Err != "" {
						color.Red(result.Err)
					}
					color.Red(result.Out)
					color.Yellow("If you feel stuck, ask a hint by executing `golings hint %s`", exercise.Name)
					os.Exit(1)
				}
				bar.Add(1) // nolint
			}

			color.Green("\n\nCongratulations!!!")
			color.Green("You passed all the exercises")
		},
	}
}
