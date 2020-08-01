// Convert JSON choose-your-own-adventure into html pages
//
// implementation of Gophercise Exercise #3
//   https://github.com/gophercises/cyoa/
//
// primary features:
//	use html/template package to create HTML pages
//	custom http.Handler
//	handle JSON
//
// flow:
//      load all segments
//      start with "Intro" segment
//      show the title, story, and options
//      user selects a segment
//	if there are no options, the story is completed
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// types  -------------------------------------------
type Segment struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

type Segments map[string]Segment

// loading data  -------------------------------------------
func loadJsonfile(jsonFilename string) []byte {
	jsonFile, err := os.Open(jsonFilename)
	if err != nil {
		fmt.Println("ERROR: loadJsonfile")
		log.Fatal(err)
	}
	data, _ := ioutil.ReadAll(jsonFile)
	defer jsonFile.Close()
	return data
}

func loadSegments(data []byte) Segments {
	var segments Segments
	err := json.Unmarshal([]byte(data), &segments)
	if err != nil {
		fmt.Println("ERROR: loadSegments")
		log.Fatal(err)
	}
	return segments
}

// http version -------------------------------------------
func httpVersion(segments Segments) {

	baseTemplate, err := template.ParseFiles("layout.html")
	if err != nil {
		fmt.Println("ERROR: serve(): template.ParseFiles")
		log.Fatal(err)
	}
	tmpl := template.Must(baseTemplate, err)

	for name, segment := range segments {

		fmt.Println("Registering http handler for:", name)
		closureSegment := segment

		http.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
			err := tmpl.Execute(w, closureSegment)
			if err != nil {
				fmt.Println("ERROR: tmpl.Execute()")
				log.Fatal(err)
			}
		})
	}

	// register one more special redirect, to start at the "/intro" page
	fmt.Println("Registering special http handler for: /")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/intro", http.StatusMovedPermanently)
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("ERROR: http.ListenAndServe()")
		log.Fatal(err)
	}
}

// console version -------------------------------------------
func showSegment(segment Segment) string {
	fmt.Printf(" \n %s\n %s\n \n", segment.Title, segment.Story)
	for _, option := range segment.Options {
		fmt.Println(option.Arc, ": ", option.Text)
	}
	var chosen string
	fmt.Scanln(&chosen)
	return chosen
}

func consoleVersion(segments Segments) {
	const STARTING_SEGMENT = "intro"
	chosen := showSegment(segments[STARTING_SEGMENT])
	for chosen != "" {
		chosen = showSegment(segments[chosen])
	}
}

func main() {

	jsonFilename := "gopher.json"
	json := loadJsonfile(jsonFilename)
	segments := loadSegments(json)

	httpVersion(segments)
	//consoleVersion(segments)
}
