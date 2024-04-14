package tests

import (
	"net/url"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

const (
	host = "localhost:8080"
)

func TestCreateGetUpdateGetnew(t *testing.T) {
	// сервис должен быть задеплоен к началу теста
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	// create banner

	e.POST
}
