package main

import (
	"../../godasher"
	"html/template"
)

func Setup(config godasher.Config) {
	println("Store config variables here (api key for example)")
}

func Render(component godasher.Component) template.HTML {
	return template.HTML("<p>My awesome rendering plugin</p>")
}
