package main

import (
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var urlmap map[string]string

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

	//if r.Header.Get("Content-Type") != "text/plain" {
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	w.Header().Set("Content-Type", "text/plain")

	if r.Method == http.MethodPost {
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
		key := hash(string(urlBytes))
		urlmap[key] = string(urlBytes)
		w.WriteHeader(http.StatusCreated)
		newurl := "http://localhost:8080/" + key
		fmt.Println("newurl = ", newurl)
		w.Write([]byte(newurl))
		return
	}

	if r.Method == http.MethodGet {
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
}

func main() {
	urlmap = make(map[string]string)
	mux := http.NewServeMux()
	mux.HandleFunc("/", handle)
	fmt.Println("Listening on :8080")
	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
