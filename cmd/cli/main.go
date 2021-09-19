// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// You are going to want to copy and paste this line later.
// "github.com/ardanlabs/gotraining/topics/go/exercises/contributors/template/github"

// Call the GitHub API to get a list of repository contributors.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/piotrstrzalka/contributors/internal"
	"github.com/piotrstrzalka/contributors/pkg/github"
)

// Create a type where we can decode contributor json values.
// It needs the fields "login" and "contributions".

func main() {
	repo := flag.String("repo", "golang/go", "Provide name of repo (eg. 'repo/go')")
	csvOut := flag.String("out", "", "Path to save data in csv format")
	flag.Parse()

	tkn, err := internal.GetToken()
	if err != nil {
		fmt.Printf(err.Error())
		tkn = ""
	}

	client, err := github.NewClient("https://api.github.com", tkn)
	if err != nil {
		log.Fatal(err)
	}

	err = printContributors(*repo, client, os.Stdout, *csvOut)
	if err != nil {
		log.Fatal(err)
	}
}

type contributorLister interface {
	ContributorList(string) ([]github.Contributor, error)
}

func printContributors(repo string, c contributorLister, target io.Writer, cPath string) error {
	cr, err := c.ContributorList(repo)
	if err != nil {
		return err
	}
	for i, cl := range cr {
		fmt.Fprintf(target, "%d. %s %d\n", i+1, cl.Login, cl.Contributions)
	}

	if cPath != "" {
		writeToCsv(cr, cPath)
		if err != nil {
			return err
		}
	}
	return nil
}

//todo write test for it
func writeToCsv(data []github.Contributor, path string) error {
	file, err := os.OpenFile(path+".csv", os.O_CREATE|os.O_RDWR, 666)
	if err != nil {
		return err
	}
	defer file.Close()

	cw := csv.NewWriter(file)
	cw.Write([]string{"Login", "Contributions"})
	for _, cl := range data {
		cw.Write([]string{cl.Login, fmt.Sprint(cl.Contributions)})
	}
	cw.Flush()

	return nil
}
