package config

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func SetupServer(serverAddress string) (*http.Server, *chi.Mux) {
	router := chi.NewRouter()
	server := &http.Server{
		Addr:    serverAddress,
		Handler: router,
	}

	return server, router
}
