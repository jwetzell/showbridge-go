package main

import "syscall/js"

type LogWriter struct {
	Element   js.Value
	Container js.Value
}

func (pw *LogWriter) ScrollToBottom() {
	scrollHeight := pw.Container.Get("scrollHeight").Int()
	clientHeight := pw.Container.Get("clientHeight").Int()
	pw.Container.Set("scrollTop", scrollHeight-clientHeight)
}

func (pw *LogWriter) IsScrolledToBottom() bool {
	scrollHeight := pw.Container.Get("scrollHeight").Int()
	clientHeight := pw.Container.Get("clientHeight").Int()
	scrollTop := pw.Container.Get("scrollTop").Int()
	return scrollHeight-clientHeight <= scrollTop+25
}

func (pw *LogWriter) Write(p []byte) (n int, err error) {
	if !pw.Element.IsUndefined() {
		currentText := pw.Element.Get("textContent").String()
		newText := currentText + string(p)
		pw.Element.Set("textContent", newText)
		if pw.IsScrolledToBottom() {
			pw.ScrollToBottom()
		}
	}
	return len(p), nil
}

func NewLogWriter(id string) *LogWriter {
	element := document.Call("getElementById", id)
	container := element.Get("parentElement")
	return &LogWriter{
		Element:   element,
		Container: container,
	}
}
