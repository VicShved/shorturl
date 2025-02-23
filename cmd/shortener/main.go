package main

import (
	"hash/fnv"
	"io"
	"net/http"
	"strconv"
)

var urlMap map[string]string

func hash(s string) string {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return ""
	}
	return strconv.Itoa(int(h.Sum32()))
}

func handle(w http.ResponseWriter, r *http.Request) {
	good := (r.Method != http.MethodPost)
	good = good || (r.Method != http.MethodGet)
	good = good || (r.Header.Get("Content-Type") != "text/plain")
	if good {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	if r.Method == http.MethodPost {
		urlBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if len(urlBytes) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		urlMap[hash(string(urlBytes))] = string(urlBytes)
		w.WriteHeader(http.StatusCreated)
		return
	}
	if r.Method == http.MethodGet {
		urlstr := r.URL.Path
		if len(urlstr) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		url, exists := urlMap[urlstr]
		if !exists {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)

	}
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", handle)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
