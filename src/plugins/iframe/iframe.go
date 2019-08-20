package main

import (
	"../../godasher"
	"fmt"
	"html/template"
	"strings"
)

func Setup(config godasher.Config) {
	// Nothing
}

func Render(component godasher.Component) template.HTML {
	attribute := "src"
	value := component.Data["url"]

	if val, ok := component.Data["src"]; ok {
		attribute = "srcdoc"
		value = val
	} else if shouldProxy := component.Data["proxy-port"]; shouldProxy != "" {
		port := component.Data["proxy-port"]
		hostname := component.Data["proxy-hostname"]
		protocol := component.Data["proxy-protocol"]
		base := component.Data["base"]
		uri := component.Data["url"]
		insecure := component.Data["insecure"]

		if port == "" {
			panic("The 'proxy-port' must be set when proxy is enabled")
		}
		if hostname == "" {
			panic("The 'proxy-hostname' must be set when proxy is enabled")
		}
		if protocol == "" {
			panic("The 'proxy-protocol' must be set when proxy is enabled")
		}
		if base == "" {
			panic("The 'base' must be set when proxy is enabled")
		}

		proxy, err := godasher.NewProxy(base, strings.ToLower(insecure) == "true")
		if err != nil {
			return template.HTML(fmt.Sprintf("%v", err))
		}

		proxy.StartServer(port)
		value = strings.ReplaceAll(uri, base, protocol+"://"+hostname+":"+port)
	}

	return template.HTML(fmt.Sprintf("<iframe %v=\"%v\" style=\"width: 100%%; height: 100%%; border:none;\"></iframe>",
		attribute,
		strings.ReplaceAll(value, "\"", "'")))
}
