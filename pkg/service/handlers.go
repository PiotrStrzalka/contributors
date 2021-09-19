package service

import (
	"html/template"
	"log"
	"net/http"

	"github.com/piotrstrzalka/contributors/internal"
	"github.com/piotrstrzalka/contributors/pkg/github"
)

var client *github.Client

func init() {
	token, err := internal.GetToken()
	if err != nil {
		log.Fatal(err)
	}
	client, err = github.NewClient("https://api.github.com", token)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Println("Empty token, data fetch will be limited")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	repo := r.FormValue("repo")
	log.Println("Request for:", repo)

	var cons []github.Contributor
	if repo != "" {
		var err error
		cons, err = client.ContributorList(repo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	fv := map[string]interface{}{"repo": repo}
	site := render(fv, cons)
	// w.Write([]byte("Hello"))
	w.Write(site)
}

func render(fv map[string]interface{}, results []github.Contributor) []byte {
	if results != nil {
		vars := map[string]interface{}{"Items": results}
		markup := executeTemplate("results", vars)
		fv["Results"] = template.HTML(string(markup))
	}

	markup := executeTemplate("search", fv)

	vars := map[string]interface{}{"LayoutContent": template.HTML(string(markup))}
	return executeTemplate("layout", vars)
}
