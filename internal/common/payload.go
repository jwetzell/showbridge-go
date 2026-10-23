package common

type WrappedPayload struct {
	Payload any
	Source  string
	End     bool
}
