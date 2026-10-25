package routing

type WrappedPayload struct {
	Payload any
	Source  string
	End     bool
}
