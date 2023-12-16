package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
)

var site_infos = make(map[string][]byte)

func Startup(dir string) {
	items, _ := os.ReadDir(dir)
	for _, item := range items {
		content, err := os.ReadFile(dir + "/" + item.Name())
		if err != nil {
			log.Fatal(err)
		}
		site_infos[item.Name()] = content
	}
}

func OnPage(link string) []byte {
	var buf bytes.Buffer
	// get HTML data from the website
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}

	if err := res.Write(&buf); err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func main() {

	dir := "C:/Users/Aidan Feess/Documents/Projects/LTSEC/scoring-engine/site_infos"
	Startup(dir)

	actual := OnPage("https://info.cern.ch/")
	bytes.Compare(site_infos["cern_info.html"], actual)

	fmt.Println(site_infos["cern_info.html"])
	fmt.Println(actual)

}
