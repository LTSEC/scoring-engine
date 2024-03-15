package scoring

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
)

var site_info []byte

// Runs on startup of the web server checker to read the html of the site into memory.
// Takes a string which represents the directory that will be read from
func web_startup(dir string) error {

	content, err := os.ReadFile(dir)
	if err != nil {
		return err
	}
	site_info = []byte(strings.ReplaceAll(string(content), "\n", "")) // Converting from byte to string to byte to remove stray eol

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
// Takes in a directory and list of strings which are ip addresses
func CheckWeb(dir string, ip string) (bool, error) {

	err := web_startup(dir)
	if err != nil {
		return false, err
	}

	pagehtml, err := onPage("http://" + ip)
	if err != nil {
		return false, err
	}

	pagehtml = bytes.TrimSuffix(pagehtml, []byte{10})                // Trim byte 10 (eof) from end of file
	site_info := bytes.ReplaceAll(site_info, []byte{13}, []byte{10}) // Exchange byte 13 for byte 10 (im not sure why eof is at the end of every line)

	webserv_up := bytes.Equal(bytes.TrimSpace(site_info), bytes.TrimSpace(pagehtml))

	return webserv_up, nil
}
