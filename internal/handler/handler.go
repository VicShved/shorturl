package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/logger"
	"github.com/VicShved/shorturl/internal/repository"
	"github.com/VicShved/shorturl/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type reqJSON struct {
	URL string `json:"url"`
}

type respJSON struct {
	Result string `json:"result"`
}

type Handler struct {
	serv    *service.ShortenService
	baseurl string
}

func GetHandler(serv *service.ShortenService) *Handler {
	return &Handler{
		serv: serv,
	}
}

func (h Handler) InitRouter(mdwr []func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()
	for _, mw := range mdwr {
		router.Use(mw)
	}
	router.Post("/", h.HandlePOST)
	router.Post("/api/shorten", h.HandlePostJSON)
	router.Post("/api/shorten/batch", h.HandleBatchPOST)
	router.Get("/{key}", h.HandleGET)
	router.Get("/ping", h.PingDB)
	router.Get("/api/user/urls", h.GetUserURLs)
	return router
}

func (h Handler) HandlePostJSON(w http.ResponseWriter, r *http.Request) {
	var indata reqJSON
	userID := r.Context().Value(app.ContextUser).(string)
	logger.Log.Debug("Context User ", zap.Any("ID", userID))

	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
	urlbytes, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(urlbytes, &indata)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newurl, key := h.serv.GetShortURLFromLong(&indata.URL)

	err = h.serv.Save(*key, indata.URL, string(userID))

	if err != nil && errors.Is(err, repository.ErrPKConflict) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	var outdata respJSON
	outdata.Result = *newurl
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
	logger.Log.Info("", zap.String("url", indata.URL), zap.String("response", string(resp)))
}

func (h Handler) HandlePOST(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(app.ContextUser).(string)
	logger.Log.Debug("Context User ", zap.Any("ID", userID))

	w.Header().Set("Content-Type", "text/plain")
	defer r.Body.Close()
	urlBytes, _ := io.ReadAll(r.Body)
	url := string(urlBytes)
	newurl, key := h.serv.GetShortURLFromLong(&url)
	err := h.serv.Save(*key, url, string(userID))

	if err != nil && errors.Is(err, repository.ErrPKConflict) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	//fmt.Println("newurl = ", newurl)
	w.Write([]byte(*newurl))
}

func (h Handler) HandleGET(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(app.ContextUser).(string)
	logger.Log.Debug("Context User ", zap.Any("ID", userID))

	urlstr := chi.URLParam(r, "key")
	//fmt.Println("urlstr =", urlstr)

	url, exists := h.serv.Read(urlstr, string(userID))
	//fmt.Println("exists = ", exists)

	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//fmt.Println("url = ", url)
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h Handler) PingDB(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(app.ContextUser).(string)
	logger.Log.Debug("Context User ", zap.Any("ID", userID))

	err := h.serv.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h Handler) HandleBatchPOST(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(app.ContextUser).(string)
	logger.Log.Debug("Context User ", zap.Any("ID", userID))

	var indata []service.BatchReqJSON
	w.Header().Set("Content-Type", "application/json")
	urlbytes, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(urlbytes, &indata)
	// logger.Log.Info("indata", zap.Any("len", indata))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	results, err := h.serv.Batch(&indata, string(userID))
	if err != nil && errors.Is(err, repository.ErrPKConflict) {
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)

	resp, err := json.Marshal(results)
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
	logger.Log.Info("Batch handled", zap.String("response", string(resp)))
}

func (h Handler) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(app.ContextUser).(string)
	logger.Log.Debug("Context User ", zap.Any("ID", userID))

	w.Header().Set("Content-Type", "application/json")
	outdata, err := h.serv.GetUserURLs(string(userID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Length", strconv.Itoa(lenth))
}
