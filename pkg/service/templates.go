package service

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

var views = make(map[string]*template.Template)

func init() {
	pwd, _ := os.Getwd()
	loadTemplate("layout", pwd+"/pkg/service/views/basic-layout.html")
	loadTemplate("results", pwd+"/pkg/service/views/results.html")
	loadTemplate("search", pwd+"/pkg/service/views/search.html")
}

func loadTemplate(name, path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	temp, err := template.New(name).Parse(string(data))
	if err != nil {
		log.Fatal(err)
	}

	views[name] = temp
}

func executeTemplate(name string, vars map[string]interface{}) []byte {
	buf := &bytes.Buffer{}
	if err := views[name].Execute(buf, vars); err != nil {
		log.Println(err)
		return []byte("Error while processing template")
	}
	return buf.Bytes()
}
