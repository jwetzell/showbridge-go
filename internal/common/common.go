package common

type contextKey string

const RouterContextKey contextKey = contextKey("router")
const SourceContextKey contextKey = contextKey("source")
const ModulesContextKey contextKey = contextKey("modules")
const SenderContextKey contextKey = contextKey("sender")
