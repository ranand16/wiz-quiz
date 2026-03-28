# wiz-quiz

A tiny Go library for running interactive terminal wizards. You give it a list of questions, it handles the TUI, and hands back the answers in order. Supports inline validation so you can reject a bad answer before letting the user move on.

## Install

```bash
go get github.com/ranand16/wiz-quiz@v0.1.2
```

## Usage

```go
package main

import (
    "errors"
    "fmt"
    "strings"

    wizard "github.com/ranand16/wiz-quiz"
)

func main() {
    questions := []wizard.Question{
        {
            Question: "What is your project name?",
            // No Callback means no validation — anything goes.
        },
        {
            Question: "What is your email?",
            Callback: func(val string) error {
                if !strings.Contains(val, "@") {
                    return errors.New("that doesn't look like a valid email")
                }
                return nil
            },
        },
        {
            Question: "Pick a license (MIT, Apache, GPL):",
        },
    }

    answers, err := wizard.RunQuestions(questions)
    if err != nil {
        fmt.Println("something went wrong:", err)
        return
    }

    fmt.Println("Project:", answers[0])
    fmt.Println("Email:", answers[1])
    fmt.Println("License:", answers[2])
}
```

## API

### `RunQuestions(q []Question) ([]string, error)`

Starts the wizard and blocks until the user finishes all questions or quits. Returns answers in the same order as the input slice. If the user quits early (`ctrl+c` / `esc`), answers typed so far are still returned for whatever questions were completed.

### `Question`

| Field | Type | Description |
|-------|------|-------------|
| `Question` | `string` | The prompt shown to the user |
| `Callback` | `func(string) error` | Optional. Return an error to block moving forward |

Set `Callback` to `nil` if you don't need validation on a particular step.

## Controls

| Key | Action |
|-----|--------|
| `enter` | Submit current answer and advance |
| `ctrl+c` / `esc` | Quit the wizard |

## How it works

wiz-quiz is a thin wrapper around [Bubble Tea](https://github.com/charmbracelet/bubbletea). It follows the Elm architecture (Init / Update / View) under the hood — the internal `model` type is unexported, so you never have to think about it. `RunQuestions` is the only entry point.
