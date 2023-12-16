package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var site_infos = make(map[string][]byte)

func Startup(dir string) {
	// Get the directory of the real HTML
	items, _ := os.ReadDir(dir)

	// Add each real HTML file to a list as byte arrays
	for _, item := range items {
		content, err := os.ReadFile(dir + "/" + item.Name())
		if err != nil {
			log.Fatal(err)
		}
		site_infos[item.Name()] = content
	}
}

func OnPage(link string) []byte {
	// Get HTML data from the website
	res, err := http.Get(link)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	// Read it into memory as bytes
	res_body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	return res_body
}

func main() {

	// TODO: get directory base on user parameters in console when calling the file
	dir := "C:/Users/Aidan Feess/Documents/Projects/LTSEC/scoring-engine/site_infos"
	Startup(dir)

	/// TODO: get each website somehow
	actual := OnPage("https://info.cern.ch/")

	// TODO: do this for every site in site_infos
	webserv_up := bytes.Equal(site_infos["cern_info.html"], actual)
	fmt.Println(webserv_up)

}
