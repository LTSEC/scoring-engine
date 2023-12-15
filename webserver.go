package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func OnPage(link string) string {
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	content, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func main() {
	expected := OnPage("https://info.cern.ch/")
	actual := OnPage("https://info.cern.ch/")

	fmt.Println(expected == actual)
}
