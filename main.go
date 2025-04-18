package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type Artwork struct {
	ObjectID          int    `json:"objectID"`
	Title             string `json:"title"`
	ArtistDisplayName string `json:"artistDisplayName"`
	PrimaryImage      string `json:"primaryImage"`
}

func main() {
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/artwork", artworkHandler)
	http.HandleFunc("/artwork-template", artworkTemplateHandler)
	fmt.Println("server is running: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func getArtwork(objectID string) (Artwork, error) {
	apiURL := fmt.Sprintf("https://collectionapi.metmuseum.org/public/collection/v1/objects/%s", objectID)
	resp, err := http.Get(apiURL)
	if err != nil {
		return Artwork{}, err
	}
	defer resp.Body.Close()

	var artwork Artwork
	if err := json.NewDecoder(resp.Body).Decode(&artwork); err != nil {
		return Artwork{}, err
	}
	return artwork, nil
}

func artworkHandler(w http.ResponseWriter, r *http.Request) {
	artwork, err := getArtwork(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "failed to get artwork", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artwork)
}

func artworkTemplateHandler(w http.ResponseWriter, r *http.Request) {
	artwork, err := getArtwork(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "failed to get artwork", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/artwork.html")
	if err != nil {
		http.Error(w, "failed to load template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, artwork)
}
