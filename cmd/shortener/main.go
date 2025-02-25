package main

import (
	"fmt"
	"github.com/VicShved/shorturl/internal/app"
	"github.com/go-chi/chi/v5"
	"net/http"
)

//var urlmap map[string]string

//func HandlePOST(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodPost {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	//if r.Header.Get("Content-Type") != "text/plain" {
//	//	w.WriteHeader(http.StatusBadRequest)
//	//	return
//	//}
//
//	w.Header().Set("Content-Type", "text/plain")
//	defer r.Body.Close()
//	urlBytes, _ := io.ReadAll(r.Body)
//	fmt.Println("string(urlBytes) = ", string(urlBytes))
//	//if err != nil {
//	//	w.WriteHeader(http.StatusBadRequest)
//	//	return
//	//}
//	//if len(urlBytes) == 0 {
//	//	w.WriteHeader(http.StatusBadRequest)
//	//	return
//	//}
//	key := app.Hash(string(urlBytes))
//	urlmap[key] = string(urlBytes)
//	w.WriteHeader(http.StatusCreated)
//	newurl := "http://localhost:8080/" + key
//	fmt.Println("newurl = ", newurl)
//	w.Write([]byte(newurl))
//}
//
//func HandleGET(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodGet {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	urlstr := strings.TrimPrefix(r.URL.Path, "/")
//	fmt.Println("urlstr =", urlstr)
//	//if len(urlstr) == 0 {
//	//	w.WriteHeader(http.StatusBadRequest)
//	//	return
//	//}
//
//	url, exists := urlmap[urlstr]
//	fmt.Println("exists = ", exists)
//
//	if !exists {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	fmt.Println("url = ", url)
//	w.Header().Set("Location", url)
//	w.WriteHeader(http.StatusTemporaryRedirect)
//}

func main() {
	router := chi.NewRouter()
	router.Post("/", app.HandlePOST)
	router.Get("/{key}", app.HandleGET)
	//mux := http.NewServeMux()
	//mux.HandleFunc("POST /", app.HandlePOST)
	//mux.HandleFunc("GET /", app.HandleGET)
	var config = app.InitServerConfig()
	fmt.Println("Start URL=", config.StartBaseURL)
	fmt.Println("Result URL=", config.ResultBaseURL)

	//err := http.ListenAndServe(`localhost:8080`, mux)
	err := http.ListenAndServe(config.StartBaseURL, router)
	if err != nil {
		panic(err)
	}
}
