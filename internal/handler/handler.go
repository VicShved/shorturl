// Package for handler http request
package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/VicShved/shorturl/internal/logger"
	"github.com/VicShved/shorturl/internal/middware"
	"github.com/VicShved/shorturl/internal/repository"
	"github.com/VicShved/shorturl/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// type reqJSON
type reqJSON struct {
	URL string `json:"url"`
}

// type respJSON
type respJSON struct {
	Result string `json:"result"`
}

// type Handler
type Handler struct {
	serv *service.ShortenService
	cidr string
}

// func GetHandler
func GetHandler(serv *service.ShortenService, cidr string) *Handler {
	return &Handler{
		serv: serv,
		cidr: cidr,
	}
}

// func (h Handler) InitRouter
func (h Handler) InitRouter(mdwr []func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()
	for _, mw := range mdwr {
		router.Use(mw)
	}

	router.Mount("/debug", middleware.Profiler())
	router.Post("/", h.HandlePOST)
	router.Post("/api/shorten", h.HandlePostJSON)
	router.Post("/api/shorten/batch", h.HandleBatchPOST)
	router.Get("/{key}", h.HandleGET)
	router.Get("/ping", h.PingDB)
	router.Get("/api/user/urls", h.GetUserURLs)
	router.Delete("/api/user/urls", h.DelUserURLs)
	router.Get("/api/internal/stats", h.GetStats)
	return router
}

// HandlePostJSON godoc
// @Tags api
// @Summary Сохранение url
// @Description Запрос на сохранение длинного url
// @ID HandlePostJSON
// @Accept  json
// @Produce json
// @Param url body reqJSON true "Исходный url"
// @Success 200 {object} respJSON
// @Failure 400 {string} string "Неверный запрос"
// @Failure 403 {string} string "Ошибка авторизации"
// @Failure 409 {string} string "Уже есть такой url"
// @Failure 500 {string} string "Внутренняя ошибка"
// @Security AuthorizationCook
// @Router /api/shorten [post]
func (h Handler) HandlePostJSON(w http.ResponseWriter, r *http.Request) {
	var indata reqJSON
	// Вытаскиваю userID из контекста
	userID := r.Context().Value(middware.ContextUser).(string)
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

	err = h.serv.Save(*key, indata.URL, userID)

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
	logger.Log.Debug("", zap.String("url", indata.URL), zap.String("response", string(resp)))
}

// HandlePost godoc
// @Tags
// @Summary Сохранение url
// @Description Запрос на сохранение длинного url
// @ID HandlePost
// @Accept  plain/text
// @Produce plain/text
// @Param url body string true "Исходный url"
// @Success 200 {string} string
// @Failure 400 {string} string "Неверный запрос"
// @Failure 403 {string} string "Ошибка авторизации"
// @Failure 409 {string} string "Уже есть такой url"
// @Failure 500 {string} string "Внутренняя ошибка"
// @Security AuthorizationCook
// @Router / [post]
func (h Handler) HandlePOST(w http.ResponseWriter, r *http.Request) {
	// Вытаскиваю userID из контекста
	userID := r.Context().Value(middware.ContextUser).(string)
	logger.Log.Debug("Context User ", zap.Any("ID", userID))

	w.Header().Set("Content-Type", "text/plain")
	defer r.Body.Close()
	urlBytes, _ := io.ReadAll(r.Body)
	url := string(urlBytes)
	newurl, key := h.serv.GetShortURLFromLong(&url)
	err := h.serv.Save(*key, url, userID)

	if err != nil && errors.Is(err, repository.ErrPKConflict) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	w.Write([]byte(*newurl))
}

// HandleGET godoc
// @Tags
// @Summary Получение url
// @Description Запрос на получение длинного url из короткого в header Location
// @ID HandleGet
// @Accept  plain/text
// @Produce plain/text
// @Param key path string true "Короткий url"
// @Success 307
// @Failure 400 {string} string "Неверный запрос"
// @Failure 403 {string} string "Ошибка авторизации"
// @Failure 410 {string} string "Удален"
// @Failure 500 {string} string "Внутренняя ошибка"
// @Security AuthorizationCook
// @Router / [get]
func (h Handler) HandleGET(w http.ResponseWriter, r *http.Request) {
	// Вытаскиваю userID из контекста
	userID := r.Context().Value(middware.ContextUser).(string)
	logger.Log.Debug("Context User ", zap.Any("ID", userID))

	urlstr := chi.URLParam(r, "key")
	//fmt.Println("urlstr =", urlstr)

	url, exists, isDeleted := h.serv.Read(urlstr, userID)
	//fmt.Println("exists = ", exists)

	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if isDeleted {
		w.WriteHeader(http.StatusGone)
		return
	}

	//fmt.Println("url = ", url)
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// PingDB godoc
// @Tags
// @Summary Проверка БД
// @Description Запрос на проверку работы БД
// @ID PingDB
// @Success 200
// @Failure 500
// @Security AuthorizationCook
// @Router /ping [get]
func (h Handler) PingDB(w http.ResponseWriter, r *http.Request) {
	// Вытаскиваю userID из контекста
	userID := r.Context().Value(middware.ContextUser).(string)
	logger.Log.Debug("Context User ", zap.Any("ID", userID))

	err := h.serv.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// HandleBatchPOST godoc
// @Tags api
// @Summary Сохранение группы url
// @Description Запрос на сохранение группы длинных url
// @ID HandleBatchPOST
// @Accept  json
// @Produce json
// @Param urls body service.BatchReqJSON true "Исходные urls"
// @Success 201 {object} service.BatchRespJSON
// @Failure 400
// @Failure 403
// @Failure 409
// @Failure 500
// @Security AuthorizationCook
// @Router /api/shorten/batch [post]
func (h Handler) HandleBatchPOST(w http.ResponseWriter, r *http.Request) {
	// Вытаскиваю userID из контекста
	userID := r.Context().Value(middware.ContextUser).(string)
	logger.Log.Debug("Context User ", zap.Any("ID", userID))

	var indata []service.BatchReqJSON
	w.Header().Set("Content-Type", "application/json")
	urlbytes, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(urlbytes, &indata)
	// logger.Log.Debug("indata", zap.Any("len", indata))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	results, err := h.serv.Batch(&indata, userID)
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
	logger.Log.Debug("Batch handled", zap.String("response", string(resp)))
}

// GetUserURLs godoc
// @Tags api
// @Summary Все urls пользователя
// @Description Запрос на все urls пользователя
// @ID GetUserURLs
// @Accept  json
// @Produce json
// @Success 200 {object} service.UserURLRespJSON
// @Success 204
// @Failure 500
// @Security AuthorizationCook
// @Router /api/user/urls [get]
func (h Handler) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	// Вытаскиваю userID из контекста
	userID := r.Context().Value(middware.ContextUser).(string)
	logger.Log.Debug("Context User ", zap.Any("ID", userID))

	w.Header().Set("Content-Type", "application/json")
	outdata, err := h.serv.GetUserURLs(userID)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(*outdata) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp, err := json.Marshal(outdata)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lenth, err := w.Write(resp)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(lenth))
}

// DelUserURLs godoc
// @Tags api
// @Summary Удаление url пользователя
// @Description Запрос на удаление urls пользователя
// @ID DelUserURLs
// @Accept  json
// @Produce json
// @Param urls body []string true "urls на удаление"
// @Success 202
// @Failure 400
// @Failure 500
// @Security AuthorizationCook
// @Router /api/user/urls [delete]
func (h Handler) DelUserURLs(w http.ResponseWriter, r *http.Request) {
	// Вытаскиваю userID из контекста
	userID := r.Context().Value(middware.ContextUser).(string)
	logger.Log.Debug("Context User ", zap.Any("ID", userID))

	var indata []string
	urlbytes, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(urlbytes, &indata)
	logger.Log.Debug("indata", zap.Any("len", indata))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.serv.DelUserURLs(&indata, userID)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// GetStats godoc
// @Tags api/internal
// @Summary Статистика
// @Description Запрос возвратит количество сокращенных URL и пользователей сервиса
// @ID GetStats
// @Produce json
// @Success 200 {object} service.Stats
// @Failure 403
// @Failure 500
// @Security TRUSTED_SUBNET
// @Router /api/internal/stats [get]
func (h Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	address := r.Header.Get("X-Real-IP")
	logger.Log.Debug("GetStats", zap.String("Address", address))
	isInSubNet := isInSubNet(address, h.cidr)
	if !isInSubNet {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	outdata, err := h.serv.GetStats()
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp, err := json.Marshal(outdata)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lenth, err := w.Write(resp)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(lenth))

}
