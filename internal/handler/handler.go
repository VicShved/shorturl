package handler

import (
	"encoding/json"
	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/middware"
	"github.com/VicShved/shorturl/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

type Handler struct {
	serv    *service.ShortenService
	baseurl string
}

func GetHandler(serv *service.ShortenService, baseurl string) *Handler {
	return &Handler{
		serv:    serv,
		baseurl: baseurl,
	}
}

func (h Handler) InitRouter(mdwr []func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()
	for _, mw := range mdwr {
		router.Use(mw)
	}
	router.Post("/", h.HandlePOST)
	router.Post("/api/shorten", h.HandlePostJSON)
	router.Get("/{key}", h.HandleGET)
	return router
}

func (h Handler) HandlePostJSON(w http.ResponseWriter, r *http.Request) {
	type inJSON struct {
		URL string `json:"url"`
	}
	var indata inJSON
	type outJSON struct {
		Result string `json:"result"`
	}
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
	urlbytes, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(urlbytes, &indata)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	key := app.Hash(indata.URL)
	h.serv.Save(key, indata.URL)
	w.WriteHeader(http.StatusCreated)
	newurl := h.baseurl + "/" + key
	//fmt.Println("newurl = ", newurl)
	var outdata outJSON
	outdata.Result = newurl
	resp, err := json.Marshal(outdata)
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
	middware.Log.Info("", zap.String("url", indata.URL), zap.String("response", string(resp)))
}

func (h Handler) HandlePOST(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain")
	defer r.Body.Close()
	urlBytes, _ := io.ReadAll(r.Body)
	url := string(urlBytes)
	//fmt.Println("string(urlBytes) = ", url)
	key := app.Hash(url)
	h.serv.Save(key, url)
	w.WriteHeader(http.StatusCreated)
	newurl := h.baseurl + "/" + key
	//fmt.Println("newurl = ", newurl)
	w.Write([]byte(newurl))
}

func (h Handler) HandleGET(w http.ResponseWriter, r *http.Request) {

	urlstr := chi.URLParam(r, "key")
	//fmt.Println("urlstr =", urlstr)

	url, exists := h.serv.Read(urlstr)
	//fmt.Println("exists = ", exists)

	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//fmt.Println("url = ", url)
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
