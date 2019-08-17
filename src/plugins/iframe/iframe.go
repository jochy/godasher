package main

import (
	"../../godasher"
	"fmt"
	"html/template"
	"strings"
)

func Setup(config godasher.Config) {

}

func Render(component godasher.Component) template.HTML {
	attribute := "src"
	value := component.Data["url"]

	if val, ok := component.Data["src"]; ok {
		attribute = "srcdoc"
		value = val
	}

	return template.HTML(fmt.Sprintf("<iframe %v=\"%v\" style=\"width: 100%%; height: 100%%; border:none;\"></iframe>",
		attribute,
		strings.ReplaceAll(value, "\"", "'")))
}
