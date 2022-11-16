package main

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/controllers"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type People struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Member       []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Location     string   `json:"location"`
	Relation     Relation
}

type Relation struct {
	Id    int                 `json:"id"`
	Dates map[string][]string `json:"datesLocations"`
}

func Unmarshal(w http.ResponseWriter, r *http.Request) []People {
	url := "https://groupietrackers.herokuapp.com/api/artists"
	res, err := http.Get(url)
	if err != nil {
		controllers.HandlerErrors(w, 500)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		controllers.HandlerErrors(w, 500)
	}

	var artists []People

	jsonErr := json.Unmarshal(body, &artists)

	if jsonErr != nil {
		controllers.HandlerErrors(w, 500)
	}

	return artists
}

func UnmarshalArtist(w http.ResponseWriter, r *http.Request, id int) People {
	artists := Unmarshal(w, r)
	return artists[id]
}

func UnmarshalRelation(w http.ResponseWriter, r *http.Request, id int) Relation {
	url_dates := "https://groupietrackers.herokuapp.com/api/relation/" + strconv.Itoa(id+1)
	res, err := http.Get(url_dates)
	if err != nil {
		controllers.HandlerErrors(w, 500)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		controllers.HandlerErrors(w, 500)
	}

	var relations Relation

	jsonErr := json.Unmarshal(body, &relations)

	if jsonErr != nil {
		controllers.HandlerErrors(w, 500)
	}
	return relations
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		controllers.HandlerErrors(w, 404)
		return
	}
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		controllers.HandlerErrors(w, 500)
		return
	}
	artists := Unmarshal(w, r)
	err = tmpl.Execute(w, artists)
	if err != nil {
		controllers.HandlerErrors(w, 500)
		return
	}
}

func Description(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("Description.html")
	if err != nil {
		controllers.HandlerErrors(w, 500)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		controllers.HandlerErrors(w, 500)
		return
	}
	artist := UnmarshalArtist(w, r, id-1)
	relation := UnmarshalRelation(w, r, id-1)
	artist.Relation = relation
	err = tmpl.Execute(w, artist)
	if err != nil {
		controllers.HandlerErrors(w, 500)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/Description", Description)
	fileServer := http.FileServer(http.Dir("./resources/"))
	mux.Handle("/resources/", http.StripPrefix("/resources/", fileServer))
	log.Println("Запуск веб-сервера на http://localhost:8080/ ")
	fmt.Println("Server is listening...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
