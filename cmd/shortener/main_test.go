package main

import (
	"github.com/VicShved/shorturl/internal/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortener(t *testing.T) {
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
	localurl := app.InitServerConfig().ResultBaseURL
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/", strings.NewReader(test.url))
			w := httptest.NewRecorder()
			app.HandlePOST(w, request)
			res := w.Result()
			//defer func(Body io.ReadCloser) {
			//	err := Body.Close()
			//	assert.Nil(t, err)
			//}(res.Body)
			err := res.Body.Close()
			assert.Nil(t, err)
			assert.Equal(t, test.want.status, res.StatusCode)
			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, localurl+"/"+app.Hash(test.url), string(body))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
