package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/fulldump/apitest"
	"github.com/fulldump/biff"

	"github.com/SeniorGo/seniorgoacademy/auth"
	"github.com/SeniorGo/seniorgoacademy/persistence"
)

func TestNewApi(t *testing.T) {

	// cursePersistence := persistence.NewInMemory[Curse]()
	cursePersistence, err := persistence.NewInDisk[Curse](t.TempDir())
	biff.AssertNil(err)

	user := auth.User{
		ID:      "user-test",
		Nick:    "user-nick",
		Picture: "user-picture",
		Email:   "user@email.com",
	}

	authHeaderBytes, _ := json.Marshal(map[string]any{"user": user})
	authHeader := string(authHeaderBytes)

	h := NewApi("testversion", "", cursePersistence)
	a := apitest.NewWithHandler(h)

	t.Run("Request /version", func(t *testing.T) {
		r := a.Request("GET", "/version").
			WithHeader(auth.XGlueAuthentication, authHeader).
			Do()

		biff.AssertEqual(r.StatusCode, http.StatusOK)
		biff.AssertEqual(r.BodyString(), "testversion")
	})

	t.Run("List curses (empty)", func(t *testing.T) {
		r := a.Request("GET", "/v0/curses").
			WithHeader(auth.XGlueAuthentication, authHeader).
			Do()

		biff.AssertEqual(r.StatusCode, http.StatusOK)

		expectedBody := []map[string]interface{}{}
		biff.AssertEqualJson(r.BodyJson(), expectedBody)
	})

	t.Run("Create Curse", func(t *testing.T) {
		r := a.Request("POST", "/v0/curses").
			WithHeader(auth.XGlueAuthentication, authHeader).
			WithBodyJson(map[string]string{
				"title":       "My Curse",
				"description": "This is my description",
			}).
			Do()

		biff.AssertEqual(r.StatusCode, http.StatusCreated)

		body := r.BodyJsonMap()
		expectedBody := map[string]interface{}{
			"id":                body["id"],
			"author":            user,
			"title":             "My Curse",
			"description":       "This is my description",
			"creation_time":     body["creation_time"],
			"modification_time": body["modification_time"],
		}
		biff.AssertEqualJson(r.BodyJsonMap(), expectedBody)
	})

	t.Run("List curses (1)", func(t *testing.T) {
		r := a.Request("GET", "/v0/curses").
			WithHeader(auth.XGlueAuthentication, authHeader).
			Do()

		biff.AssertEqual(r.StatusCode, http.StatusOK)
		biff.AssertEqual(len(r.BodyJson().([]any)), 1)
	})

	t.Run("Create Curse - Error title validation", func(t *testing.T) {
		r := a.Request("POST", "/v0/curses").
			WithHeader(auth.XGlueAuthentication, authHeader).
			WithBodyJson(map[string]string{
				"title": strings.Repeat("a", 1025),
			}).
			Do()

		biff.AssertEqual(r.StatusCode, http.StatusBadRequest)

		expectedBody := map[string]interface{}{
			"error": map[string]interface{}{
				"title":       "Bad Request",
				"description": "title is too long",
			},
		}
		biff.AssertEqual(r.BodyJsonMap(), expectedBody)
	})

}
