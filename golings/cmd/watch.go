package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	// commented out until bug below is resolved
	//"golang.org/x/term"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	//"golang.org/x/term"
)

func WatchCmd(infoFile string) *cobra.Command {
	return &cobra.Command{
		Use:   "watch",
		Short: "Verify exercises when files are edited",
		RunE: func(cmd *cobra.Command, args []string) error {
		/* 
			THIS causes a bug where the first letter of a new line
			is printed at the end of the line obove the rest of the message
			removing the code for oldstate and fd below fixes the issue, 
			but you must hit enter after making your menu selection.
		*/
			
			// **bug begin**
			/*
			fd := int(os.Stdin.Fd())
			oldState, err := term.MakeRaw(fd)
			if err != nil {
				return err
			}
			defer term.Restore(fd, oldState)
			*/
			// **bug end**

			// Initial run.
			RunNextExercise(infoFile)

			update := make(chan string)
			go WatchEvents(update)

			// Key events.
			events := make(chan byte)
			go func() {
				b := make([]byte, 1)
				for {
					_, err := os.Stdin.Read(b)
					if err != nil {
						return
					}
					events <- b[0]
				}
			}()

			// Window resize.
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGWINCH)

			for {
				select {
				case <-update:
					// File changed: rerun next exercise.
					RunNextExercise(infoFile)

				case <-sigChan:
					// Terminal resized: re-render current state.
					RefreshUI()

				case ch := <-events:
					switch ch {
					case 'n':
						MoveToNextAndRun(infoFile)
					case 'h':
						RefreshUI()
						PrintHint(infoFile)
						color.Green("Hint: 'n' to move to next exercise, 'l' for list of exercises")
						// redraw whole UI after extra output
					case 'l':
						RefreshUI()
						PrintList(infoFile)
						color.Green("List of exercises: 'l' to list, 'n' for next")

					case 'q', 3: // 3 = Ctrl+C
						fmt.Println()
						color.Green("Bye from kdor/golings")
						return nil
					case '\r', '\n':
						// TODO: ignore Enter in raw mode if i need it
					default:
						//TODO: ignore other keys or optionally show a small message
					}
				}
			}
		},
	}
}

func WatchEvents(updateF chan<- string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	path, _ := os.Getwd()
	directories := fmt.Sprintf("%s/exercises", path)

	err = filepath.WalkDir(directories, func(path_dir string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
			return err
		}
		if d.IsDir() {
			err = watcher.Add(path_dir)
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal("Error in file path:", err.Error())
	}

	for event := range watcher.Events {
		if event.Has(fsnotify.Write) || event.Has(fsnotify.Rename) {
			updateF <- event.Name
		}
	}
}
