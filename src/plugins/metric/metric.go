package main

import (
	"../../godasher"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yalp/jsonpath"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var db *sql.DB
var client = &http.Client{
	Timeout: time.Second * 10,
}

func init() {
	/*
		var err error
		db, err = sql.Open("sqlite3", "./data/db.sqlite")
		checkErr(err)
		executeQuery(`CREATE TABLE IF NOT EXISTS METRIC(
								id INTEGER PRIMARY KEY AUTOINCREMENT,
								name VARCHAR(255) NULL,
								lastRefresh DATE NULL
							);`)
		executeQuery(`CREATE TABLE IF NOT EXISTS METRIC_VALUE(
								id INTEGER PRIMARY KEY AUTOINCREMENT,
								metric INTEGER,
								date DATE,
								value VARCHAR(255)
							);`)
	*/
}

func Setup(config godasher.Config) {
	// Nothing
}

func Render(component godasher.Component) template.HTML {
	val, err := retrieveMetric(component)

	if err != nil {
		return template.HTML(fmt.Sprintf("%v", err))
	}

	var tpl bytes.Buffer

	tmpl, err := template.New("metric.html").Funcs(template.FuncMap{}).ParseFiles("plugins/metric.html")
	if err != nil {
		return template.HTML(fmt.Sprintf("%v", err))
	}

	err = tmpl.Execute(&tpl, RenderMetric{RandomId: "truc", Value: fmt.Sprintf("%v", val)})
	if err != nil {
		return template.HTML(fmt.Sprintf("%v", err))
	}

	return template.HTML(tpl.String())
}

func retrieveMetric(component godasher.Component) (interface{}, error) {
	req, err := http.NewRequest("GET", component.Data["url"], nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	format := strings.ToLower(component.Data["format"])
	path := component.Data["path"]

	if format == "json" && path != "" {
		var data interface{}
		err = json.Unmarshal(bodyText, &data)

		if err != nil {
			return nil, err
		}

		val, err := jsonpath.Read(data, path)
		if err != nil {
			return nil, err
		}

		return val, nil
	}

	return string(bodyText), nil
}

func executeQuery(query string, args ...interface{}) sql.Result {
	stmt, err := db.Prepare(query)
	checkErr(err)

	res, err := stmt.Exec(args...)
	checkErr(err)
	return res
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type RenderMetric struct {
	RandomId string
	Value    string
}
