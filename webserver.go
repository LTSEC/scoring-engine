package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func OnPage(link string) string {
	// get HTML data from the website
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	// read that data
	content, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	// Sanitize the HTML
	p := bluemonday.StripTagsPolicy()
	html := p.Sanitize(string(content))

	return strings.Replace(string(html), " ", "", -1)
}

func GetTemplate(directory string) string {
	dat, err := os.ReadFile(directory)
	if err != nil {
		log.Fatal(err)
	}

	// Sanitize the HTML
	p := bluemonday.StripTagsPolicy()
	html := p.Sanitize(string(dat))

	return strings.Replace(string(html), " ", "", -1)
}

func main() {

	dir := "C:/Users/Aidan Feess/Documents/Projects/LTSEC/scoring-engine/cern_info.html"

	actual := OnPage("https://info.cern.ch/")
	expected := GetTemplate(dir)

	fmt.Println(expected == actual)
}
