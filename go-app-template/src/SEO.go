package main

import (
	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

// =================
// This is the HTML render component for SEO optimization.
// This is currently stubbed to return basically nothing - but can be filled in
// and/or replaced with the src-wasm code for full SEO integration.
// =================

type SEO struct {
	app.Compo
}

func (S *SEO) Render() app.UI {
	return app.Div()
}
