package controllers

import "github.com/linewin/linend"

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}
