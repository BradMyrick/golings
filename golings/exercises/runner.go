package exercises

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

)

type Result struct {
	Exercise Exercise
	Out      string
	Err      string
}

func (e Exercise) Run() (Result, error) {
	if e.Mode == "compile" {
		return runCompile(e)
	}
	return runTest(e)
}

func runCompile(e Exercise) (Result, error) {
	cmd := exec.Command("go", "run", "-tags=golings", fmt.Sprintf("./%s", e.Path))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return Result{Exercise: e, Out: stdout.String(), Err: stderr.String()}, err
}

func runTest(e Exercise) (Result, error) {
	src, err := os.ReadFile(e.Path)
	if err != nil {
		return Result{Exercise: e, Err: err.Error()}, err
	}

	tmp, err := os.MkdirTemp("", "golings-*")
	if err != nil {
		return Result{Exercise: e, Err: err.Error()}, err
	}
	defer os.RemoveAll(tmp)
	if err := os.WriteFile(filepath.Join(tmp, "main.go"),
		[]byte("package main\nfunc main(){}\n"), 0644); err != nil {
		return Result{Exercise: e, Err: err.Error()}, err
	}
	version := strings.Split(runtime.Version()[2:], ".")
	modContent := fmt.Sprintf("module golings_exercise\n\ngo %s.%s\n", version[0], version[1])


	if err := os.WriteFile(filepath.Join(tmp, "go.mod"),
		[]byte(modContent), 0600); err != nil {
			return Result{Exercise: e, Err: err.Error()}, err
		}

	dest := filepath.Join(tmp, "main_test.go")
	if err := os.WriteFile(dest, src, 0644); err != nil {
		return Result{Exercise: e, Err: err.Error()}, err
	}
	cmd := exec.Command("go", "test", "-v", "-race", "-tags=golings")
	cmd.Dir = tmp
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	return Result{Exercise: e, Out: stdout.String(), Err: stderr.String()}, err

}

func BuildArgs(e Exercise) []string {
	if e.Mode == "compile" {
		return []string{"run", "-tags=golings", fmt.Sprintf("./%s", e.Path)}
	}
	return []string{"test", "-v", "-race", "-tags=golings", fmt.Sprintf("./%s", e.Path)}
}

