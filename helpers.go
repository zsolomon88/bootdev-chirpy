package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

func getBadWords() []string {
	return []string{"kerfuffle", "sharbert", "fornax"}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "applicatoin/json")
	w.WriteHeader(code)
	type errReturnVals struct {
		ErrorVal string `json:"error"`
	}
	errBody := errReturnVals{
		ErrorVal: msg,
	}
	dat, err := json.Marshal(errBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.Write(nil)
		return
	}
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "applicatoin/json")
	w.WriteHeader(code)

	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		w.Write(nil)
		return
	}

	w.Write(dat)
}

func replaceWord(subject string, search string, replace string) string {
	searchRegex := regexp.MustCompile("(?i)" + search)
	return searchRegex.ReplaceAllString(subject, replace)
}
