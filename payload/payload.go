package payload

type WrappedPayload struct {
	Payload any
	Source  string
	End     bool
}
