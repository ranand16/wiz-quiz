package wizard

type Question struct {
	Question string                  `json:"question"`
	Callback func(args ...any) error `json:"callback"` // We are adding callback function so that it can be used for validation or any other side effects
}
