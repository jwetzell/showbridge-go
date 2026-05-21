package common

type WrappedPayload struct {
	Payload      any
	InputHandler InputHandler
	Modules      map[string]Module
	Source       string
	End          bool
}
