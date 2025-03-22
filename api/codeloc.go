package api

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const CodetabsBaseurl = "https://api.codetabs.com/v1/loc/"

type CodeLocResponseEntry struct {
	Language    string `json:"language"`
	Lines       int    `json:"lines"`
	Files       int    `json:"files"`
	LinesOfCode int    `json:"linesOfCode"`
}

// ShieldsEndpointResponse supplies response declared at https://shields.io/badges/endpoint-badge
type ShieldsEndpointResponse struct {
	// SchemaVersion is required and always the number 1
	SchemaVersion int `json:"schemaVersion"`
	// Label is required
	Label string `json:"label"`
	// Message is required
	Message string `json:"message"`
	Color   string `json:"color,omitempty"`
}

//goland:noinspection GoUnusedExportedFunction called by Vercel
func CodeLoc(w http.ResponseWriter, r *http.Request) {
	githubRepo := r.URL.Query().Get("github")
	if githubRepo == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Query parameter 'github' is required"))
		return
	}

	language := r.URL.Query().Get("language")
	if language == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Query parameter 'language' is required"))
		return
	}

	matched, err := regexp.MatchString("[a-zA-Z0-9-_.]+/[a-zA-Z0-9-_.]+", githubRepo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error matching github repo: " + err.Error()))
		return
	}

	if !matched {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid github repo: " + githubRepo))
		return
	}

	resp, err := http.Get(CodetabsBaseurl +
		"?github=" + githubRepo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error retrieving from code loc: " + err.Error()))
		return
	}

	var entries []CodeLocResponseEntry

	err = json.NewDecoder(resp.Body).Decode(&entries)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error decoding response: " + err.Error()))
		return
	}

	var linesOfCode = 0
	for _, entry := range entries {
		if strings.EqualFold(entry.Language, language) {
			linesOfCode = entry.LinesOfCode
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")

	shieldsResponse := ShieldsEndpointResponse{
		SchemaVersion: 1,
		Label:         "LOC (" + language + ")",
		Message:       strconv.Itoa(linesOfCode),
	}

	err = json.NewEncoder(w).Encode(shieldsResponse)
	if err != nil {
		log.Fatal(err)
	}
}
