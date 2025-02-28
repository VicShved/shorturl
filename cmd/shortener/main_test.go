package main

import (
	"context"
	"fmt"
	"github.com/VicShved/shorturl/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPost(t *testing.T) {
	type want struct {
		status int
		//response    string
		contentType string
	}
	tests := []struct {
		name   string
		method string
		url    string
		want   want
	}{
		{
			name:   "test1",
			method: http.MethodPost,
			url:    "https://google.com",
			want: want{
				status: 201,
				//response:    "http://localhost:8080/",
				contentType: "text/plain",
			},
		},
		{
			name:   "test2",
			method: http.MethodPost,
			url:    "https://rbc.ru/",

			want: want{
				status: 201,
				//response:    "http://localhost:8080/",
				contentType: "text/plain",
			},
		},
	}

	app.ServerConfig.BaseURL = "http://localhost:8080"
	baseurl := app.ServerConfig.BaseURL
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/", strings.NewReader(test.url))
			w := httptest.NewRecorder()
			app.HandlePOST(w, request)
			res := w.Result()
			err := res.Body.Close()
			assert.Nil(t, err)
			assert.Equal(t, test.want.status, res.StatusCode)
			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, baseurl+"/"+app.Hash(test.url), string(body))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestGet(t *testing.T) {
	type want struct {
		status         int
		locationheader string
		contentType    string
	}
	tests := []struct {
		name    string
		method  string
		suffics string
		want    want
	}{
		{
			name:    "test1",
			method:  http.MethodGet,
			suffics: "1",
			want: want{
				status:         307,
				locationheader: "http://ya.ru/",
			},
		},
		{
			name:    "test2",
			method:  http.MethodGet,
			suffics: "2",
			want: want{
				status:         307,
				locationheader: "https://google.com/",
			},
		},
		{
			name:    "null in key",
			method:  http.MethodGet,
			suffics: "",
			want: want{
				status:         400,
				locationheader: "https://google.com/",
			},
		},
	}

	urlmap := *app.GetStorage()
	for _, test := range tests {
		if test.suffics != "" {
			urlmap[test.suffics] = test.want.locationheader
		}
		t.Run(test.name, func(t *testing.T) {
			target := "/{key}"
			fmt.Println("target", target)
			request := httptest.NewRequest(test.method, target, nil)
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("key", test.suffics)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			app.HandleGET(w, request)
			res := w.Result()
			err := res.Body.Close()
			assert.Nil(t, err)
			assert.Equal(t, test.want.status, res.StatusCode)
			if test.suffics != "" {
				assert.Equal(t, test.want.locationheader, res.Header.Get("Location"))
			}
		})
	}
}
