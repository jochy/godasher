package main

import (
	"../../godasher"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var baseUrl = ""
var user = ""
var password = ""
var insecure = ""
var client = &http.Client{
	Timeout: time.Second * 10,
}
var configLoaded = false

func Setup(config godasher.Config) {
	configLoaded = false
	jenkinsConfig, ok := config.ExternalConfig["jenkins"]
	if ok {
		baseUrl, ok = jenkinsConfig["baseUrl"]
		if !ok {
			panic("externalconfig.jenkins.baseUrl is required")
		}
		user, ok = jenkinsConfig["user"]
		if !ok {
			panic("externalconfig.jenkins.user is required")
		}
		password, ok = jenkinsConfig["password"]
		if !ok {
			panic("externalconfig.jenkins.password is required")
		}
		insecure, ok = jenkinsConfig["insecure"]
		shouldVerify := ok && strings.ToLower(insecure) == "true"
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: shouldVerify},
		}

		configLoaded = true
	}
}

func Render(component godasher.Component) template.HTML {
	if !configLoaded {
		panic("Missing jenkins configuration. Please add 'externalconfig.jenkins' into your configuration")
	}

	if strings.ToLower(component.View) == "view" {
		name, ok := component.Data["name"]

		if !ok {
			println("No key named 'name' found for component=Jenkins and view='View'. Please add a data key named 'name' which corresponds to the Jenkins's view name.")
		}
		return RenderView("/view/"+name, component)
	} else if strings.ToLower(component.View) == "job" {
		name, ok := component.Data["name"]

		if !ok {
			println("No key named 'name' found for component=Jenkins and view='Job'. Please add a data key named 'name' which corresponds to the Jenkins's job name.")
		}
		return RenderView("/job/"+name, component)
	}

	return template.HTML("<h1 style=\"color:red\">View not found</p>")
}

func RenderView(uri string, component godasher.Component) template.HTML {
	var tpl bytes.Buffer

	tmpl, err := template.New("jenkins_view").Funcs(template.FuncMap{
		"mapValue": func(m map[string]interface{}, key string) interface{} {
			val, ok := m[key]
			if !ok {
				return ""
			}
			return val
		},
		"jobColor": func(job map[string]interface{}) string {
			color := strings.ToLower(fmt.Sprintf("%v", job["color"]))
			class := "alert-secondary"
			if strings.Contains(color, "red") {
				class = "alert-danger"
			} else if strings.Contains(color, "yellow") {
				class = "alert-warning"
			} else if strings.Contains(color, "blue") {
				class = "alert-success"
			} else if strings.Contains(color, "grey") ||
				strings.Contains(color, "disabled") ||
				strings.Contains(color, "ABORTED") ||
				strings.Contains(color, "notbuilt") {
				class = "alert-dark"
			}

			if strings.Contains(color, "_anime") {
				class = class + " " + "jenkins-active"
			}

			return class
		},
		"msToMin": func(n float64) float64 {
			return math.Round(n / 60000.0)
		},
	}).Parse(JenkinsViewTemplate)
	if err != nil {
		log.Printf("Error while loading the Jenkins View template file: %v", err)
	}

	response := make(map[string]interface{})
	jenkinsResponse, err := callJenkins(uri)

	if err != nil {
		return template.HTML(fmt.Sprintf("%v", err))
	}

	if strings.ToLower(component.View) == "job" {
		response["jobs"] = [1]interface{}{jenkinsResponse}
	} else {
		response = jenkinsResponse
	}

	err = tmpl.Execute(&tpl, RenderRequest{
		"jenkins-" + strconv.Itoa(rand.Intn(10000)),
		response,
		strings.ToLower(component.View),
	})
	if err != nil {
		log.Printf("Error while rendering the Jenkins View template: %v", err)
	}

	return template.HTML(tpl.String())
}

func callJenkins(uri string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", baseUrl+uri+"/api/json?pretty=false&depth=2", nil)

	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(user, password)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	bodyText, _ := ioutil.ReadAll(resp.Body)
	var dat map[string]interface{}
	if err := json.Unmarshal(bodyText, &dat); err != nil {
		return nil, err
	}
	return dat, nil
}

type RenderRequest struct {
	RandomId string
	Response map[string]interface{}
	View     string
}
