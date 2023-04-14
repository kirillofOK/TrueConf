package apiserver

import "net/http"

func Start(config *Config) error {
	srv := newServer(config)
	return http.ListenAndServe(config.BindAddr, srv.router)
}
