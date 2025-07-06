package service

import (
	"reflect"
	"slices"
	"testing"

	"github.com/VicShved/shorturl/internal/app"
	"github.com/VicShved/shorturl/internal/repository"
)

func TestSave(t *testing.T) {
	user, _ := app.GetNewUUID()
	tests := []struct {
		short  string
		long   string
		userID string
		want   error
	}{
		{
			short:  "1",
			long:   "10",
			userID: user,
			want:   nil,
		},
		{
			short:  "1",
			long:   "10",
			userID: user,
			want:   nil,
		},
		{
			short:  "2",
			long:   "20",
			userID: user,
			want:   nil,
		},
	}

	app.ServerConfig.BaseURL = "http://localhost:8080"
	baseurl := app.ServerConfig.BaseURL
	repo := repository.GetFileRepository(app.ServerConfig.FileStoragePath)
	serv := GetService(repo, baseurl)

	for _, test := range tests {
		t.Run(test.short, func(t *testing.T) {
			err := serv.Save(test.short, test.long, test.userID)
			if err != test.want {
				t.Errorf("result Save expected %s got %s", test.want, err)
			}
		})
	}
}

func TestRead(t *testing.T) {
	type testStruct struct {
		short  string
		long   string
		userID string
		want   string
	}
	user, _ := app.GetNewUUID()
	tests := []testStruct{
		{
			short:  "1",
			long:   "10",
			userID: user,
			want:   "10",
		},
		{
			short:  "1",
			long:   "10",
			userID: "",
			want:   "10",
		},
		{
			short:  "2",
			long:   "20",
			userID: user,
			want:   "20",
		},
	}

	app.ServerConfig.BaseURL = "http://localhost:8080"
	baseurl := app.ServerConfig.BaseURL
	repo := repository.GetFileRepository(app.ServerConfig.FileStoragePath)
	serv := GetService(repo, baseurl)

	for _, test := range tests {
		t.Run(test.short, func(t *testing.T) {
			serv.Save(test.short, test.long, test.userID)
		})
	}
	newTestElement := testStruct{
		short:  "01",
		long:   "01",
		userID: user,
		want:   "",
	}
	tests = append(tests, newTestElement)
	for _, test := range tests {
		t.Run(test.short, func(t *testing.T) {
			long, _, _ := serv.Read(test.short, test.userID)
			if long != test.want {
				t.Errorf("result Read (%s) expected %s got %s", test.short, test.want, long)
			}
		})
	}
}

func TestGetUserURLs(t *testing.T) {
	type testStruct struct {
		short  string
		long   string
		userID string
	}
	user, _ := app.GetNewUUID()
	user2, _ := app.GetNewUUID()
	user3, _ := app.GetNewUUID()
	tests := []testStruct{
		{
			short:  "1",
			long:   "10",
			userID: user,
		},
		{
			short:  "1",
			long:   "10",
			userID: user2,
		},
		{
			short:  "2",
			long:   "20",
			userID: user3,
		},
	}

	app.ServerConfig.BaseURL = "http://localhost:8080"
	baseurl := app.ServerConfig.BaseURL
	repo := repository.GetFileRepository(app.ServerConfig.FileStoragePath)
	serv := GetService(repo, baseurl)

	wants := map[string][]UserURLRespJSON{}
	for _, test := range tests {
		t.Run(test.short, func(t *testing.T) {
			serv.Save(test.short, test.long, test.userID)
			wants[test.userID] = append(wants[test.userID], UserURLRespJSON{app.ServerConfig.BaseURL + "/" + test.short, test.long})
		})
	}
	wants["nobody"] = []UserURLRespJSON{}
	for userID, urls := range wants {
		t.Run(userID, func(t *testing.T) {
			userURLs, err := serv.GetUserURLs(userID)
			if err != nil {
				t.Errorf("serv.GetUserURLs(userID) return error %s", err)
			}

			if userURLs != nil {
				if *userURLs != nil {
					if !reflect.DeepEqual(*userURLs, urls) {
						t.Errorf("result GetUserURLs expected %s got %s", urls, *userURLs)
					}
				} else {
					if !slices.Equal(*userURLs, urls) {
						t.Errorf("result GetUserURLs expected %s got %s", urls, *userURLs)
					}
				}
			} else {
				t.Errorf("result GetUserURLs expected %s got nil", urls)
			}

		})
	}
}
