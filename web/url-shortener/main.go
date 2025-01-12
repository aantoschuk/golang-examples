package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type URLStruct struct {
	URL  string
}

type FormatedURL struct {
	Formated string
}

const baseURL = "http://localhost:80"

const base62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
var m = make(map[string]string)


// Generate a random Base62 ID
func generateBase62ID  ()(string, error) {
	var result string

	bytes := make([]byte, 6)

	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	for _, b := range bytes {
		result += string(base62[b%62])
	}

	return result, nil
}

func SimpleHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the template
	tmpl := template.Must(template.ParseFiles("./template/index.html"))

	// For GET requests, show the form with a "Hello" message
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	// For POST requests, capture the URL from the form
	data := URLStruct{
		URL: r.FormValue("url"),
	}

	if data.URL == "" {
		http.Error(w, "Please fill URL", http.StatusBadRequest)
	}


	id, err := generateBase62ID()
	if err != nil {
		log.Fatal("Error generating base62 ID")
	}

	m[id] = data.URL 

	result := "http://localhost:80/short/" + id
	// Send the response with the data
	err = tmpl.Execute(w, FormatedURL{Formated: result})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Redirect if URL is existed in the database 
func RedirectHandler (w http.ResponseWriter, r *http.Request) {
	fullURL := r.URL.Path

	pathArray := strings.Split(fullURL, "/")
	id := pathArray[len(pathArray)-1]

	value, exists := m[id]

	if !exists  {
		fmt.Println("ID not found in map")
		return
	}

	http.Redirect(w, r, value, http.StatusFound)
}

func main () {

    fs := http.FileServer(http.Dir("static/"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", SimpleHandler)
	http.HandleFunc("/short/", RedirectHandler)

	// Log any occured errors
	log.Fatal(http.ListenAndServe(":80", nil))
}
