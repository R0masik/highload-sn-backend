package http

import (
	"fmt"
	"net/http"
	"time"

	"highload-sn-backend/handlers"
	"highload-sn-backend/internal/log"

	"github.com/gorilla/mux"
)

const (
	LoginURL        = "/login"
	RegisterUserURL = "/user/register"
	GetUserURL      = "/user/get/{id}"

	httpPortDefault = "8080"
)

func NewServer() (*http.Server, string) {
	router := mux.NewRouter()
	router.Methods(http.MethodPost).Path(LoginURL).Handler(WithLogger(handlers.LoginHandler))
	router.Methods(http.MethodPost).Path(RegisterUserURL).Handler(WithLogger(handlers.RegisterUserHandler))
	router.Methods(http.MethodGet).Path(GetUserURL).Handler(WithLogger(handlers.GetUserHandler))

	handler := http.NewServeMux()
	handler.Handle("/", router)

	return &http.Server{
		Addr:         ":" + httpPortDefault,
		Handler:      handler,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}, httpPortDefault
}

func WithLogger(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wr := newCodeResponseWriter(w)
		start := time.Now()
		defer func() {
			elapsed := time.Since(start)
			code := fmt.Sprintf("%d", wr.code)
			reqRespInfo := fmt.Sprintf(
				"%s %s %s --> Status %s (%v)",
				r.Proto, r.Method, r.URL.String(), code, elapsed,
			)

			log.Logger().Info(reqRespInfo)
		}()

		h.ServeHTTP(wr, r)
	}
}
