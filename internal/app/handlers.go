package app

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func HandlePOST(w http.ResponseWriter, r *http.Request) {
	urlmap := *GetStorage()
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//if r.Header.Get("Content-Type") != "text/plain" {
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	w.Header().Set("Content-Type", "text/plain")
	defer r.Body.Close()
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
	key := Hash(string(urlBytes))
	urlmap[key] = string(urlBytes)
	w.WriteHeader(http.StatusCreated)
	newurl := ServerConfig.BaseURL + "/" + key
	fmt.Println("newurl = ", newurl)
	w.Write([]byte(newurl))
}

func HandleGET(w http.ResponseWriter, r *http.Request) {
	urlmap := *GetStorage()
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
