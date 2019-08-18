package main

import (
	"./godasher"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"plugin"
	"strconv"
	"strings"
)

var config = godasher.Config{}

func home(w http.ResponseWriter, r *http.Request) {
	readConfigFile()
	initPlugins()

	renderPage(w, config.Dashboards[0], 0)
}

func next(w http.ResponseWriter, r *http.Request) {
	readConfigFile()
	initPlugins()

	dashboardNum, _ := strconv.Atoi(r.URL.Query()["dashboard"][0])

	if dashboardNum >= len(config.Dashboards) {
		dashboardNum = 0
	}

	renderPage(w, config.Dashboards[dashboardNum], dashboardNum)
}
func renderPage(w http.ResponseWriter, dashboard godasher.Dashboard, dashboardNumber int) {
	t, _ := template.New("index.html").Funcs(template.FuncMap{
		"loop": func(n int) []struct{} {
			return make([]struct{}, n)
		},
		"render": func(component godasher.Component) template.HTML {
			p, err := plugin.Open("plugins/" + component.Type + ".so")

			if err != nil {
				log.Fatalf("Unable to load the plugin: %v", err)
			}

			symbol, err := p.Lookup("Render")
			if err != nil {
				log.Fatalf("Render method not found: %v", err)
			}

			renderFunc, ok := symbol.(func(godasher.Component) template.HTML)
			if !ok {
				panic("Plugin has no 'Render(godasher.Component) template.HTML' function")
			}

			return renderFunc(component)
		},
		"sum": func(a int, b int) int {
			return a + b
		},
		"maxColumns": func(DashboardNum int) int {
			max := 0
			for _, component := range config.Dashboards[DashboardNum].Components {
				if component.Column+component.Width > max {
					max = component.Column + component.Width
				}
			}
			return max
		},
	}).ParseFiles("index.html")
	err := t.Execute(w, godasher.RequestContext{Config: config, Dashboard: dashboard, DashboardNumber: dashboardNumber})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func main() {
	readConfigFile()
	initPlugins()

	http.HandleFunc("/", home)
	http.HandleFunc("/next", next)
	err := http.ListenAndServe(":"+strconv.Itoa(config.Port), nil)

	if err != nil {
		log.Fatalf("error: %v", err)
	} else {
		log.Println("Started")
	}
}

func readConfigFile() {
	data, readError := ioutil.ReadFile("config.yml")
	if readError != nil {
		log.Fatalf("error: %v", readError)
	}
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func initPlugins() {
	fileInfo, err := ioutil.ReadDir("plugins")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	for _, file := range fileInfo {
		if strings.Contains(file.Name(), ".so") {
			p, err := plugin.Open("plugins/" + file.Name())

			if err != nil {
				log.Fatalf("Unable to load the plugin: %v", err)
			}

			symbol, err := p.Lookup("Setup")
			if err != nil {
				log.Fatalf("Setup method not found: %v", err)
			}

			setup, ok := symbol.(func(godasher.Config))
			if !ok {
				panic("Plugin has no 'Setup(godasher.Config)' function")
			}

			setup(config)
		}
	}
}
