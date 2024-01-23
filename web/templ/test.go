package main

import (
	"net/http"
	//"time"
	"fmt"
	"github.com/a-h/templ"
)


func main() {
	services := []string{"ftp", "apache", "inginx"}
	teams := []string{"team a", "team b", "team c"}
	//component := hello("Danison")
	table := table(teams, services, true)
	
	//http.Handle("/", templ.Handler(component))
	http.Handle("/", templ.Handler(table))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}


/*
func main() {
	http.Handle("/", templ.Handler(timeComponent(time.Now())))
	http.Handle("/404", templ.Handler(notFoundComponent(), templ.WithStatus(http.StatusNotFound)))

	http.ListenAndServe(":8080", nil)
}
*/

/*

func NewNowHandler(now func() time.Time) NowHandler { // constructs a new NowHandler struct,
	return NowHandler{Now: now}						  // which returns the current time, 
}													  // which was passed in as a parameter

type NowHandler struct { // structs have fields
	Now func() time.Time // Now field with type func() time.Time
}

func (nh NowHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	timeComponent(nh.Now()).Render(r.Context(), w)
}

func main() {
	http.Handle("/", NewNowHandler(time.Now))

	http.ListenAndServe(":8080", nil)
}
*/