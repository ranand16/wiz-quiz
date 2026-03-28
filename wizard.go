// Package wizard is a small library for running interactive terminal wizards.
// You give it a list of questions, it handles the TUI, and hands back the answers.
// Built on top of Bubble Tea, so it follows the Elm architecture (Init/Update/View).
package wizard

import (
	"fmt"

	ti "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// model holds everything the wizard needs to know at any point in time.
// It's unexported — callers never touch this directly, only through RunQuestions.
type model struct {
	questions []Question
	inputs    []ti.Model // one input field per question, all initialised upfront
	index     int        // the question the user is currently on
	done      bool
	err       string // last validation error — shown inline below the input
}

// Init kicks off the cursor blink. Without this the text input just sits there looking dead.
func (m model) Init() tea.Cmd {
	return ti.Blink
}

// Update handles incoming events. Most of the interesting logic is in the
// keypress handling — everything else just gets forwarded to the active input.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			// Run the callback before moving on. If it returns an error, we
			// stay on the current question and show the message — the user
			// has to fix their answer before they can continue.
			if cb := m.questions[m.index].Callback; cb != nil {
				if err := cb(m.inputs[m.index].Value()); err != nil {
					m.err = err.Error()
					return m, nil
				}
			}
			m.err = ""
			if m.index == len(m.questions)-1 {
				// That was the last question, we're done.
				m.done = true
				return m, tea.Quit
			}
			m.index++
			return m, m.inputs[m.index].Focus()
		}
	}
	// Any other keypress goes straight to the active text input.
	var cmd tea.Cmd
	m.inputs[m.index], cmd = m.inputs[m.index].Update(msg)
	return m, cmd
}

// View is called by Bubble Tea after every Update to produce the string that
// gets rendered to the terminal. Keep it fast — no side effects here.
func (m model) View() string {
	if m.done {
		return fmt.Sprintf("\n Done! %s\n", m.inputs[0].Value())
	}

	errMsg := ""
	if m.err != "" {
		// Render the validation error right below the input so it's hard to miss.
		errMsg = fmt.Sprintf("\n⚠ Error: %s", m.err)
	}

	return fmt.Sprintf(
		"Setup Wizard [%d/%d]\n\n%s\n\n%s%s",
		m.index+1,
		len(m.questions),
		m.questions[m.index].Question,
		m.inputs[m.index].View(),
		errMsg,
	) + "\n\n(press enter to continue)\n"
}

// RunQuestions runs the wizard and blocks until the user answers everything or
// quits with ctrl+c / esc. Answers come back in the same order as the questions.
func RunQuestions(q []Question) (answers []string, err error) {
	updatedModel, err := tea.NewProgram(initializeModel(q)).Run()
	if err != nil {
		fmt.Println("There was an error", err)
	}
	m := updatedModel.(model)
	answers = make([]string, len(m.questions))
	for i := range m.questions {
		answers[i] = m.inputs[i].Value()
	}
	return answers, err
}

// initializeModel wires up the starting state. We pre-create all the input
// fields here rather than on demand — this way any answers already typed are
// preserved if we ever add back/forward navigation in the future.
func initializeModel(q []Question) model {
	questions := make([]string, len(q))
	for i, question := range q {
		questions[i] = question.Question
	}

	inputs := make([]ti.Model, len(questions))
	for i := range questions {
		inputs[i] = ti.New()
		inputs[i].Placeholder = "Type your answer here"
		inputs[i].Focus()
	}

	return model{
		questions: q,
		inputs:    inputs,
		index:     0,
		done:      false,
	}
}
