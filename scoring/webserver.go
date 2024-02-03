package scoring

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var site_infos = make(map[int][]byte)

// Runs on startup of the web server checker to read the html of each site into memory.
func web_startup(dir string) error {
	// Get the directory of the real HTML
	items, _ := os.ReadDir(dir)
	iter := 0

	// Add each real HTML file to a list as byte arrays
	for _, item := range items {
		content, err := os.ReadFile(dir + "/" + item.Name())
		if err != nil {
			return err
		}
		site_infos[iter] = []byte(strings.ReplaceAll(string(content), "\n", "")) // Converting from byte to string to byte to remove stray eol
		iter += 1
	}

	return nil
}

// Returns the HTML data on the given website, takes a link as an input and returns a byte array.
func onPage(link string) ([]byte, error) {
	// Get HTML data from the website
	res, err := http.Get(link)

	if err != nil {
		return make([]byte, 0), err
	}

	defer res.Body.Close()

	// Read it into memory as bytes
	res_body, err := io.ReadAll(res.Body)

	if err != nil {
		return make([]byte, 0), err
	}

	return res_body, nil
}

// Iterates through the websites provided and returns a list of booleans indicating which websites are up and which are down.
func CheckWeb(dir string, site_ips []string) ([]bool, error) {

	var return_sites []bool
	err := web_startup(dir)
	if err != nil {
		return make([]bool, 0), err
	}

	for i := 0; i < len(site_ips); i++ {

		boolArray, err := onPage(site_ips[i])
		if err != nil {
			return make([]bool, 0), err
		}

		pagehtml := bytes.TrimSuffix(boolArray, []byte{10})                  // Trim byte 10 (eof) from end of file
		site_info := bytes.ReplaceAll(site_infos[i], []byte{13}, []byte{10}) // Exchange byte 13 for byte 10 (im not sure why eof is at the end of every line)

		fmt.Println(bytes.TrimSpace(pagehtml))
		fmt.Println(bytes.TrimSpace(site_info))

		webserv_up := bytes.Equal(bytes.TrimSpace(site_info), bytes.TrimSpace(pagehtml))
		return_sites = append(return_sites, webserv_up)
	}

	return return_sites, nil
}
