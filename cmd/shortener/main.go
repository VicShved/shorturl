package main

import (
	"fmt"
	"github.com/VicShved/shorturl/internal/app"
	"io"
	"net/http"
	"strings"
)

var urlmap map[string]string

func HandlePOST(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//if r.Header.Get("Content-Type") != "text/plain" {
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	w.Header().Set("Content-Type", "text/plain")

	urlBytes, _ := io.ReadAll(r.Body)
	fmt.Println("string(urlBytes) = ", string(urlBytes))
	//if err != nil {
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	//if len(urlBytes) == 0 {
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	key := app.Hash(string(urlBytes))
	urlmap[key] = string(urlBytes)
	w.WriteHeader(http.StatusCreated)
	newurl := "http://localhost:8080/" + key
	fmt.Println("newurl = ", newurl)
	w.Write([]byte(newurl))
}

func HandleGET(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	urlstr := strings.TrimPrefix(r.URL.Path, "/")
	fmt.Println("urlstr =", urlstr)
	//if len(urlstr) == 0 {
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	url, exists := urlmap[urlstr]
	fmt.Println("exists = ", exists)

	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("url = ", url)
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	urlmap = make(map[string]string)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", HandlePOST)
	mux.HandleFunc("GET /", HandleGET)
	fmt.Println("Listening on :8080")
	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
