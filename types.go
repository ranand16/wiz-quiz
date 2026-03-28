package wizard

// Question is the building block of a wizard run. Fill in Question with
// whatever prompt you want to show, and optionally provide a Callback to
// validate the answer before the wizard lets the user move on.
//
// If Callback returns an error, the message is shown inline and the user
// stays on that question until they get it right. Set it to nil if you
// don't need any validation for a particular step.
type Question struct {
	Question string                  `json:"question"`
	Callback func(args ...any) error `json:"callback"`
}
