package app

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func HandlePOST(w http.ResponseWriter, r *http.Request) {
	urlmap := *GetStorage()

	w.Header().Set("Content-Type", "text/plain")
	defer r.Body.Close()
	urlBytes, _ := io.ReadAll(r.Body)
	fmt.Println("string(urlBytes) = ", string(urlBytes))
	key := Hash(string(urlBytes))
	urlmap[key] = string(urlBytes)
	w.WriteHeader(http.StatusCreated)
	newurl := ServerConfig.BaseURL + "/" + key
	fmt.Println("newurl = ", newurl)
	w.Write([]byte(newurl))
}

func HandleGET(w http.ResponseWriter, r *http.Request) {
	urlmap := *GetStorage()

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
