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
	"github.com/mauricioabreu/golings/golings/exercises"
	"github.com/mauricioabreu/golings/golings/ui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func WatchCmd(infoFile string) *cobra.Command {
	return &cobra.Command{
		Use:   "watch",
		Short: "Verify exercises when files are edited",
		RunE: func(cmd *cobra.Command, args []string) error {
			fd := int(os.Stdin.Fd())
			oldState, err := term.MakeRaw(fd)
			if err != nil {
				return err
			}
			defer term.Restore(fd, oldState)

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

			inListMode := false
			listCursor := 0
			var listExs []exercises.Exercise

			for {
				select {
				case <-update:
					// File changed: rerun next exercise.
					if !inListMode {
						RunNextExercise(infoFile)
					}

				case <-sigChan:
					// Terminal resized: re-render current state.
					if inListMode {
						ClearScreen()
						_, h := ui.GetTerminalSize()
						ui.PrintInteractiveList(os.Stdout, listExs, listCursor, h)
					} else {
						RefreshUI()
					}

				case ch := <-events:
					if inListMode {
						switch ch {
						case 'j':
							if listCursor < len(listExs)-1 {
								listCursor++
								ClearScreen()
								_, h := ui.GetTerminalSize()
								ui.PrintInteractiveList(os.Stdout, listExs, listCursor, h)
							}
						case 'k':
							if listCursor > 0 {
								listCursor--
								ClearScreen()
								_, h := ui.GetTerminalSize()
								ui.PrintInteractiveList(os.Stdout, listExs, listCursor, h)
							}
						case '\r', '\n':
							inListMode = false
							ex := listExs[listCursor]
							RunExercise(ex, infoFile)
						case 'r':
							ex := listExs[listCursor]
							exercises.UnmarkSolved(ex.Name)
							ClearScreen()
							_, h := ui.GetTerminalSize()
							ui.PrintInteractiveList(os.Stdout, listExs, listCursor, h)
						case 'q':
							inListMode = false
							RefreshUI()
						}
					} else {
						switch ch {
						case 'n':
							MoveToNextAndRun(infoFile)
						case 'h':
							if lastState != nil {
								lastState.ShowHint = !lastState.ShowHint
								RefreshUI()
							}
						case 'l':
							exs, err := exercises.List(infoFile)
							if err == nil {
								inListMode = true
								listExs = exs
								listCursor = 0
								for i, e := range exs {
									if e.State() == exercises.Pending {
										listCursor = i
										break
									}
								}
								ClearScreen()
								_, h := ui.GetTerminalSize()
								ui.PrintInteractiveList(os.Stdout, listExs, listCursor, h)
							}
						case 'r':
							exercises.ResetAll()
							RunNextExercise(infoFile)
						case 'q', 3: // 3 = Ctrl+C
							fmt.Print("\r\n")
							color.Green("Goodbye from golings!\r\n")
							return nil
						}
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
