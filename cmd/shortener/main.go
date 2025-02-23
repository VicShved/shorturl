package main

import (
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"strconv"
	"strings"
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

	good1 := (r.Method != http.MethodPost)
	good2 := (r.Method != http.MethodGet)
	good := good1 && good2
	if good {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	if r.Method == http.MethodPost {
		urlBytes, err := io.ReadAll(r.Body)
		fmt.Println("string(urlBytes) = ", string(urlBytes))
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
		newUrl := "http://localhost:8080/" + hash(string(urlBytes))
		fmt.Println("newUrl = ", newUrl)
		w.Write([]byte(newUrl))
		return
	}

	if r.Method == http.MethodGet {
		urlstr := strings.TrimPrefix(r.URL.Path, "/")
		fmt.Println("urlstr =", urlstr)
		if len(urlstr) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		url, exists := urlMap[urlstr]

		if !exists {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Println("url = ", url)
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)

	}
}

func main() {
	urlMap = make(map[string]string)
	mux := http.NewServeMux()
	mux.HandleFunc("/", handle)
	fmt.Println("Listening on :8080")
	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
