package exercises

import (
	"errors"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

var ErrExerciseNotFound = errors.New("exercise not found")
var ErrNoPendingExercises = errors.New("no pending exercises")

type Info struct {
	Exercises []Exercise
}

func List(infoFile string) ([]Exercise, error) {
	var info Info

	data, err := os.ReadFile(infoFile)
	if err != nil {
		return info.Exercises, err
	}

	if err := toml.Unmarshal(data, &info); err != nil {
		return info.Exercises, err
	}

	return info.Exercises, nil
}

func NextPending(infoFile string) (Exercise, error) {
	allExercises, err := List(infoFile)
	if err != nil {
		return Exercise{}, err
	}

	solved := GetSolved()

	for _, exercise := range allExercises {
		if !solved[exercise.Name] {
			return exercise, nil
		}
	}

	return Exercise{}, ErrNoPendingExercises
}

func Find(exercise string, infoFile string) (Exercise, error) {
	exs, err := List(infoFile)
	if err != nil {
		return Exercise{}, err
	}

	for _, ex := range exs {
		if ex.Name == exercise {
			return ex, nil
		}
	}

	return Exercise{}, ErrExerciseNotFound
}

func GetSolved() map[string]bool {
	solved := make(map[string]bool)
	data, err := os.ReadFile(".golings-state")
	if err != nil {
		return solved
	}
	for _, line := range strings.Split(string(data), "\n") {
		if line != "" {
			solved[line] = true
		}
	}
	return solved
}

func MarkSolved(name string) {
	solved := GetSolved()
	if solved[name] {
		return
	}
	f, err := os.OpenFile(".golings-state", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(name + "\n")
}

func Progress(infoFile string) (float32, int, int, error) {
	allExercises, err := List(infoFile)
	if err != nil {
		return 0.0, 0, 0, err
	}
	solved := GetSolved()
	totalDone := 0
	for _, exercise := range allExercises {
		if solved[exercise.Name] {
			totalDone++
		}
	}

	total := len(allExercises)
	return float32(totalDone) / float32(total), totalDone, total, nil
}
