package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/igolubevJ/bookings/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing: test pass
	default:
		content := fmt.Sprintf("type is not *chi.Mux, type is %T", v)
		t.Error(content)
	}
}
