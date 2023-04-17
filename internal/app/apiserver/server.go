package apiserver

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kirillofOK/TrueConf/internal/app/model"
	"github.com/kirillofOK/TrueConf/internal/app/store"
	"github.com/sirupsen/logrus"
)

const (
	// for authentication
	ctxKeyUser ctxKey = iota
	ctxKeyRequestID
)

var (
	xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
	xRealIP       = http.CanonicalHeaderKey("X-Real-IP")
)

type ctxKey int16

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func newServer(store store.Store, config *Config) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	s.configureRouter()

	if err := s.configureLogger(config); err != nil {
		s.logger.Fatal(err)
	}

	s.logger.Info("Starting api server")
	return s
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(s.realIPRequest)
	s.router.Use(s.Timeout(60 * time.Second))

	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
	})

	api := s.router.PathPrefix("/api").Subrouter()
	v1 := api.PathPrefix("/v1").Subrouter()

	users := v1.PathPrefix("/users").Subrouter()
	users.HandleFunc("/", s.handleSearchUsers()).Methods("GET")
	users.HandleFunc("/", s.handleCreateUser()).Methods("POST")

	id := users.PathPrefix("/{id}").Subrouter()
	id.HandleFunc("/", s.handleGetUser()).Methods("GET")
	id.HandleFunc("/", s.handleUpdateUser()).Methods("PATCH")
	id.HandleFunc("/", s.handleDeleteUser()).Methods("DELETE")

}

func (s *server) configureLogger(config *Config) error {
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)
	return nil
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)
		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start))
	})
}

func (s *server) Timeout(timeout time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer func() {
				cancel()
				if ctx.Err() == context.DeadlineExceeded {
					w.WriteHeader(http.StatusGatewayTimeout)
				}
			}()

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func (s *server) realIPRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rip := realIP(r); rip != "" {
			r.RemoteAddr = rip
		}
		next.ServeHTTP(w, r)
	})
}

func (s *server) handleSearchUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, s.store.User().SearchUsers())
	}
}

func (s *server) handleCreateUser() http.HandlerFunc {
	type request struct {
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
		Password    string `json:"password,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			CreatedAt:   time.Now(),
			DisplayName: req.DisplayName,
			Email:       req.Email,
			Password:    req.Password,
		}
		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parmas := mux.Vars(r)
		id := parmas["id"]
		u, err := s.store.User().Get(id)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
		}

		s.respond(w, r, http.StatusOK, u)
	}
}

// Need to complete
func (s *server) handleUpdateUser() http.HandlerFunc {
	type request struct {
		DisplayName string `json:"display_name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}
		parmas := mux.Vars(r)
		id := parmas["id"]
		err := s.store.User().Update(id, req.DisplayName)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}

		s.respond(w, r, http.StatusOK, req.DisplayName)
	}
}

// Need to complite
func (s *server) handleDeleteUser() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		parmas := mux.Vars(r)
		id := parmas["id"]
		err := s.store.User().Delete(id)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

//Need to complite

func (s *server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
