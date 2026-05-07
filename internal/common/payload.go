package common

type WrappedPayload struct {
	Payload any
	Router  RouteIO
	Modules map[string]Module
	Source  string
	End     bool
}
