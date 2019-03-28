package models

import "github.com/edouardparis/lntop/app"

type Models struct {
	App *app.App
}

func New(app *app.App) *Models {
	return &Models{App: app}
}
