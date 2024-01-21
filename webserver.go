package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var site_infos = make(map[int][]byte)
var site_ips = []string{
	"https://info.cern.ch/",
}

func Startup(dir string) { // Run on startup to read the site htmls to memory for comparison later
	// Get the directory of the real HTML
	items, _ := os.ReadDir(dir)
	iter := 0

	// Add each real HTML file to a list as byte arrays
	for _, item := range items {
		content, err := os.ReadFile(dir + "/" + item.Name())
		if err != nil {
			log.Fatal(err)
		}
		site_infos[iter] = []byte(strings.ReplaceAll(string(content), "\n", "")) // Converting from byte to string to byte to remove stray eol
		iter += 1
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

	/// TODO: get directory base on user parameters in console when calling the file
	/// files will need to be named according to each team
	dir := "E:/Projects/scoring-engine/site_infos"
	Startup(dir)

	/// TODO: get each website somehow
	// 	actual := OnPage("https://info.cern.ch/")

	/// TODO: do this for every site in site_infos
	// webserv_up := bytes.Equal(site_infos["cern_info.html"], actual)
	// fmt.Println(webserv_up)

	for i := 0; i < len(site_ips); i++ {

		pagehtml := bytes.TrimSuffix(OnPage(site_ips[i]), []byte{10})        // Trim byte 10 (eof) from end of file
		site_info := bytes.ReplaceAll(site_infos[i], []byte{13}, []byte{10}) // Exchange byte 13 for byte 10 (im not sure why eof is at the end of every line)

		fmt.Println(bytes.TrimSpace(pagehtml))
		fmt.Println(bytes.TrimSpace(site_info))

		webserv_up := bytes.Equal(bytes.TrimSpace(site_info), bytes.TrimSpace(pagehtml))
		fmt.Println(webserv_up)

	}

}
