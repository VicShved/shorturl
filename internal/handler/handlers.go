package handler

import (
	"encoding/json"
	"fmt"
	"github.com/VicShved/shorturl/internal/app"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
)

func HandlePOST(w http.ResponseWriter, r *http.Request) {
	urlmap := *app.GetStorage()

	w.Header().Set("Content-Type", "text/plain")
	defer r.Body.Close()
	urlBytes, _ := io.ReadAll(r.Body)
	fmt.Println("string(urlBytes) = ", string(urlBytes))
	key := app.Hash(string(urlBytes))
	urlmap[key] = string(urlBytes)
	w.WriteHeader(http.StatusCreated)
	newurl := app.ServerConfig.BaseURL + "/" + key
	fmt.Println("newurl = ", newurl)
	w.Write([]byte(newurl))
}

func HandleGET(w http.ResponseWriter, r *http.Request) {
	urlmap := *app.GetStorage()

	urlstr := chi.URLParam(r, "key")
	fmt.Println("urlstr =", urlstr)

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

func HandlePostJSON(w http.ResponseWriter, r *http.Request) {
	type inJSON struct {
		URL string `json:"url"`
	}
	var indata inJSON
	type outJSON struct {
		Result string `json:"result"`
	}
	urlmap := *app.GetStorage()
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
	urlbytes, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(urlbytes, &indata)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	key := app.Hash(indata.URL)
	urlmap[key] = string(indata.URL)
	w.WriteHeader(http.StatusCreated)
	newurl := app.ServerConfig.BaseURL + "/" + key
	fmt.Println("newurl = ", newurl)
	var outdata outJSON
	outdata.Result = newurl
	resp, err := json.Marshal(outdata)
	fmt.Println("resp = ", string(resp))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	lenth, err := w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(lenth))
}
