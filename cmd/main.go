// Main body of the program lives here
package main

import (
	"log"
	"net/http"

	"github.com/LTSEC/scoring-engine/web"
)

var debug int = 1

func main() {

	http.HandleFunc("/", web.Root)

	if debug == 1 {
		log.Fatal(http.ListenAndServe("127.0.0.1:80", nil))
	} else {
		log.Fatal(http.ListenAndServe("0.0.0.0:80", nil))
	}
}
