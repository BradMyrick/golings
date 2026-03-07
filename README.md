# golings

[![build and test](https://github.com/bradmyrick/golings/actions/workflows/test.yml/badge.svg)](https://github.com/bradmyrick/golings/actions/workflows/test.yml)

![gopher](misc/gopher-dance.gif)

> rustlings but for golang this time

This is a heavily modified fork of the original [golings](https://github.com/mauricioabreu/golings) by [BradMyrick](https://github.com/BradMyrick).

## What makes this fork different?

This version has been enhanced to more accurately replicate the **Rustlings** experience, including:

- **Enhanced CLI UI**: A completely overhauled terminal interface using `lipgloss` for a premium, responsive feel.
- **Interactive Watch Mode**: No need to hit Enter! Real-time single-key commands ('n' for next, 'h' for hint, 'l' for list).
- **Interactive Exercise List**: A scrollable, searchable list within the watch mode that lets you jump to any exercise.
- **State Persistence**: Uses a local `.golings-state` file to track progress accurately, removing the need for `// I AM NOT DONE` markers.
- **Improved Hint System**: Hints can be toggled on and off dynamically during the watch session.

## Installing

First, you need to have `go` installed. You can install it by visiting the [Go downloads page](https://go.dev/dl/)

### Option 1: GO install

```sh
go install github.com/bradmyrick/golings/golings@v0.0.1
```

Add `go/bin` to your PATH if you want to run golings anywhere in your terminal.

### Option 2: DevContainer

1. Clone the repository and open it in VSCode.
2. You will be prompted to reopen the code in a devcontainer.
3. Open a new terminal and run `golings watch`.

## Doing exercises

All the exercises can be found in the directory `exercises/<topic>`.

Clone the repository:

```sh
git clone https://github.com/bradmyrick/golings.git
```

To run the exercises in the recommended order while taking advantage of fast feedback loop, use the _watch_ command:

```sh
golings watch
```

### Key Commands in Watch Mode:
- `n`: Move to the next pending exercise.
- `h`: Toggle hint for the current exercise.
- `l`: Open the interactive list view to scroll and select exercises.
- `q`: Quit.

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md)

## Other 'lings

* [rustlings](https://github.com/rust-lang/rustlings)
* [ziglings](https://github.com/ratfactor/ziglings)

