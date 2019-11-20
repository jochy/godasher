package main

import (
	"../../godasher"
	"encoding/json"
	"fmt"
	"github.com/yalp/jsonpath"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: time.Second * 10,
}

var refresh = make(map[string]int64)
var lastRefresh = make(map[string]string)

func Setup(config godasher.Config) {
}

func Render(component godasher.Component) template.HTML {
	val, err := retrieveMetric(component)

	if err != nil {
		return template.HTML(fmt.Sprintf("%v", err))
	}

	color := component.Data[val]
	size := component.Data["size"]

	return template.HTML(fmt.Sprintf(`
<div style="background-color: %v;height: 100%%;text-align: center; display: flex">
	<div style="margin:auto;font-size:%v">
		<b>%v</b>
	</div>
</div>`, color, size, val))
}

func retrieveMetric(component godasher.Component) (string, error) {
	previousRefresh, exists := refresh[component.Data["url"]]
	if exists  {
		refreshInterval, nil := strconv.Atoi(component.Data["refreshIntervalInSeconds"])
		if time.Now().Unix() - previousRefresh < int64(refreshInterval) {
			return lastRefresh[component.Data["url"]], nil
		}
	}

	req, err := http.NewRequest("GET", component.Data["url"], nil)
	if err != nil {
		return "nil", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "nil", err
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "nil", err
	}

	format := strings.ToLower(component.Data["format"])
	path := component.Data["path"]

	if format == "json" && path != "" {
		var data interface{}
		err = json.Unmarshal(bodyText, &data)

		if err != nil {
			return "nil", err
		}

		val, err := jsonpath.Read(data, path)
		if err != nil {
			return "nil", err
		}

		lastRefresh[component.Data["url"]] = val.(string)
		refresh[component.Data["url"]] = time.Now().Unix()
		return val.(string), nil
	}

	return string(bodyText), nil
}
