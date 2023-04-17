package apiserver

import (
	"net/http"

	"github.com/kirillofOK/TrueConf/internal/app/store/jsonstore"
)

func Start(config *Config) error {
	st := config.StoreURL

	store := jsonstore.New(st)
	srv := newServer(store, config)
	return http.ListenAndServe(config.BindAddr, srv.router)
}
