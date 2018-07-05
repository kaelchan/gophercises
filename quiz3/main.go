package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type storyArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

func main() {
	var (
		err error
	)
	// read json
	raw, err := ioutil.ReadFile("./gopher.json")
	if err != nil {
		panic(err)
	}
	// unmarshal the data
	stories := make(map[string]storyArc)
	err = json.Unmarshal(raw, &stories)
	if err != nil {
		panic(err)
	}

	// create HTML pages
	/*
	 * Notice that the following usage is incorrect
	 * t, err := template.New(name)
	 * t.ParseFiles(filename)
	 */
	adventureTemplate, err := template.ParseFiles("adventure.html")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/story/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("path:", r.URL.Path, "#")
		splitList := strings.Split(r.URL.Path, "/")
		if len(splitList) != 3 {
			log.Println("Wrong path format:", r.URL.Path)
			return
		}
		arc := splitList[2]
		log.Println("arc:", arc)
		log.Println("title:", stories[arc].Title)
		err = adventureTemplate.Execute(w, stories[arc])
		if err != nil {
			log.Println(err)
		}
	})

	server := &http.Server{
		Addr: "localhost:8080",
	}
	log.Fatal(server.ListenAndServe())
}
