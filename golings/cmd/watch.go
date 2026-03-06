package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func WatchCmd(infoFile string) *cobra.Command {
	return &cobra.Command{
		Use:   "watch",
		Short: "Verify exercises when files are edited",
		RunE: func(cmd *cobra.Command, args []string) error {
			oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}
			defer term.Restore(int(os.Stdin.Fd()), oldState)

			RunNextExercise(infoFile)

			update := make(chan string)
			go WatchEvents(update)

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

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGWINCH)

			for {
				select {
				case <-update:
					RunNextExercise(infoFile)
				case <-sigChan:
					RefreshUI()
				case char := <-events:
					switch char {
					case 'n':
						// For now, just run next exercise.
						// We'll improve this to only move if current is done.
						MoveToNextAndRun(infoFile)
					case 'h':
						PrintHint(infoFile)
					case 'l':
						PrintList(infoFile)
					case 'q', 3: // 3 is Ctrl+C in raw mode
						color.Green("\nBye by golings o/")
						return nil
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
