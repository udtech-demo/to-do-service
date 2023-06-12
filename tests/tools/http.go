package tools

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"

	graphql "github.com/hasura/go-graphql-client"
	"github.com/labstack/echo"
)

type localRoundTripper struct {
	handler http.Handler
}

func (l localRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	l.handler.ServeHTTP(w, req)
	return w.Result(), nil
}

func DoQuery(q interface{}, variables map[string]interface{}, access string, e *echo.Echo) error {
	client := graphql.NewClient("/api/v1/query", &http.Client{Transport: localRoundTripper{handler: e}})

	if access != "" {
		client = client.WithRequestModifier(func(r *http.Request) {
			r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
		})
	}

	err := client.Query(context.Background(), q, variables)
	if err != nil {
		return err
	}

	return nil
}

func DoMutate(q interface{}, variables map[string]interface{}, access string, e *echo.Echo) error {
	client := graphql.NewClient("/api/v1/query", &http.Client{Transport: localRoundTripper{handler: e}})

	if access != "" {
		client = client.WithRequestModifier(func(r *http.Request) {
			r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
		})
	}

	err := client.Mutate(context.Background(), q, variables)
	if err != nil {
		return err
	}

	return nil
}

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}
