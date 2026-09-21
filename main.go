package main

import (
	"io"
	"net/http"
)

func main() {

	http.HandleFunc("/compare", func(w http.ResponseWriter, r *http.Request) {

	})

	http.HandleFunc("/encode", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		token := string(body)
		digest := token.Hash(token)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(digest.Hex()))
	})

	http.ListenAndServe(":8080", nil)
}
