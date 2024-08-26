package tools

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	datas PageData
	data  PageDatauno
	tmp   *template.Template
)

func init() {
	tmp = template.Must(template.ParseGlob("templates/*.html"))
}

func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" || r.Method != "GET" {
		w.WriteHeader(404)
		tmp.ExecuteTemplate(w, "404.html", nil)
		return
	}
	apiURL := "https://groupietrackers.herokuapp.com/api"
	cards, err := FetchArtistData(apiURL)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error fetching artist data: %v", err)
		return
	}
	datas = PageData{Cards: cards}
	if err := tmp.ExecuteTemplate(w, "index.html", datas); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Bad Request: Only POST method is allowed", http.StatusBadRequest)
		return
	}
}

func About(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Bad Request: Only GET method is allowed", http.StatusBadRequest)
	}
	tmp.ExecuteTemplate(w, "aboutus.html", nil)
}

func Bandinfo(w http.ResponseWriter, r *http.Request) {
	// if r.Method != "GET" {
	// 	http.Error(w, "Bad Request: Only GET method is allowed", http.StatusBadRequest)
	// }
	// id := strings.TrimPrefix(r.URL.RawQuery, "=id")
	// if id == "" {
	// 	http.Error(w, "/400", http.StatusBadRequest)
	// 	return
	// }
	// uno, err := strconv.Atoi(id)
	// if uno > len(data.Cards) || uno < 1 {
	// 	w.WriteHeader(404)
	// 	tmp.ExecuteTemplate(w, "404.html", nil)
	// 	return
	// }
	// if err != nil {
	// 	http.Error(w, "/400", http.StatusBadRequest)
	// 	return
	// }
	// casablanca := data.Cards[uno].Locations
	// ArrMaps := GetStaticMapURL(casablanca)
	// maps := MakerMaps(ArrMaps)
	// data = PageData{MapURL: maps}
	// tmp.ExecuteTemplate(w, "bandsinfo.html", data.Cards[uno-1])
	if r.Method != "GET" {
		http.Error(w, "Bad Request: Only GET method is allowed", http.StatusBadRequest)
		return
	}

	id := strings.TrimPrefix(r.URL.RawQuery, "=id")
	if id == "" {
		http.Error(w, "Bad Request: Missing ID", http.StatusBadRequest)
		return
	}

	uno, err := strconv.Atoi(id)
	if err != nil || uno < 1 || uno > len(datas.Cards) {
		// Avoid calling w.WriteHeader explicitly here if http.Error is used later
		tmp.ExecuteTemplate(w, "404.html", nil)
		return
	}

	casablanca := datas.Cards[uno-1].Locations
	ArrMaps := GetStaticMapURL(casablanca)
	maps := MakerMaps(ArrMaps)
	data = PageDatauno{
		Card:  datas.Cards[uno-1],
		MapURL: maps,
	}

	if err := tmp.ExecuteTemplate(w, "bandsinfo.html", data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}
