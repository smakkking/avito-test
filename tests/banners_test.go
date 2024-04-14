package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

const (
	host = "localhost:8080"
)

func TestGetOldValue(t *testing.T) {
	// сервис должен быть задеплоен к началу теста
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	// create banner
	data := []byte(`{
		"tag_ids": [
		  1, 2, 3
		],
		"feature_id": 2,
		"content": {
		  "title": "some_title",
		  "text": "some_text",
		  "url": "some_url"
		},
		"is_active": true
	  }`)

	e.POST("/banner").
		WithHeader("token", "admin").
		WithBytes(data).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().
		HasValue("banner_id", 1)

	// user gets banner
	e.GET("/user_banner").
		WithHeader("token", "user").
		WithQuery("tag_id", 1).
		WithQuery("feature_id", 2).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		HasValue("title", "some_title").
		HasValue("text", "some_text").
		HasValue("url", "some_url")

	// admin delete banner
	e.DELETE("banner/1").
		WithHeader("token", "admin").
		Expect().
		Status(http.StatusNoContent)

	// user get old banner from cache
	e.GET("/user_banner").
		WithHeader("token", "user").
		WithQuery("tag_id", 1).
		WithQuery("feature_id", 2).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		HasValue("title", "some_title").
		HasValue("text", "some_text").
		HasValue("url", "some_url")

	// user try to get new value but error
	e.GET("/user_banner").
		WithHeader("token", "user").
		WithQuery("tag_id", 1).
		WithQuery("feature_id", 2).
		WithQuery("use_last_revision", true).
		Expect().
		Status(http.StatusNotFound)
}

func TestCreate(t *testing.T) {
	// сервис должен быть задеплоен к началу теста
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	// create banner
	data := []byte(`{
		"tag_ids": [
		  1, 2, 3
		],
		"feature_id": 2,
		"content": {
		  "title": "some_title",
		  "text": "some_text",
		  "url": "some_url"
		},
		"is_active": true
	  }`)
	// неавторизован
	e.POST("/banner").
		WithBytes(data).
		Expect().
		Status(http.StatusUnauthorized)

	// нет доступа
	e.POST("/banner").
		WithHeader("token", "user").
		WithBytes(data).
		Expect().
		Status(http.StatusForbidden)

	// создание успешно
	e.POST("/banner").
		WithHeader("token", "admin").
		WithBytes(data).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().
		HasValue("banner_id", 1)
}
