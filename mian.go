package wizard

import (
	"fmt"

	ti "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	questions []Question
	inputs    []ti.Model
	index     int
	done      bool
	err       string
}

// ======================================
// Init function of the ELM architecture
func (m model) Init() tea.Cmd {
	return ti.Blink
}

// =====================================
// Updte function of the ELM architecture
// return m and command to be executed
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) { // this is a syntax for type switch since we dont know the type of msg that will be entered by the user
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc": // if the user presses ctrl + c or esc, we will quit the program
			return m, tea.Quit
		case "enter": // if the user presses enter
			// Run callback for validation
			if cb := m.questions[m.index].Callback; cb != nil {
				if err := cb(m.inputs[m.index].Value()); err != nil {
					m.err = err.Error()
					return m, nil
				}
			}
			m.err = ""
			if m.index == len(m.questions)-1 { // if we are at the last question, then we are done and we can quit the program
				m.done = true
				return m, tea.Quit
			}
			// otherwise go to next question and focus on the input
			m.index++
			return m, m.inputs[m.index].Focus()
		}
	}
	// otherwise, we will update the input box with the user's input
	var cmd tea.Cmd
	m.inputs[m.index], cmd = m.inputs[m.index].Update(msg)
	return m, cmd
}

// =====================================
// View function
func (m model) View() string {
	// If the wizard is finished, show a summary message
	if m.done {
		return fmt.Sprintf("\n Done! Created package: %s\n", m.inputs[0].Value())
	}

	// Render the progress header, the current question title, and the input box
	errMsg := ""
	if m.err != "" {
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

// =====================================
// RunQuestions is the function that will be called from the main.go file to run the wizard and get the answers from the user
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

// ======================================
// 1. Inintialie the model with the quesiton I want to ask about the github codebase
func initializeModel(q []Question) model {
	// 1. The questions
	questions := make([]string, len(q))
	for i, question := range q {
		questions[i] = question.Question
	}
	// 2. Get the inputs
	inputs := make([]ti.Model, len(questions))
	for i := range questions {
		inputs[i] = ti.New()
		inputs[i].Placeholder = "Type your answer here"
		inputs[i].Focus()
	}

	//  We will initialize the Model
	model := model{
		questions: q,
		inputs:    inputs,
		index:     0,
		done:      false,
	}

	// 3. return and start asking the questions
	return model

}
