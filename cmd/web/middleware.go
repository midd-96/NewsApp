package main

import "net/http"

func (a *application) LoadSession(next http.Handler) http.Handler {
	return a.session.LoadAndSave(next)
}
