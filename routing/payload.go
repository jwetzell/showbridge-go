package routing

type WrappedPayload struct {
	Payload any
	Source  string
	// TODO(jwetzell): this can probably be removed with the new middleware-esque layout
	End bool
}
