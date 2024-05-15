package router

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ghostship-dev/authservice/core/datatypes"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

type Router struct {
	mux        *http.ServeMux
	middleware []func(http.Handler) http.Handler
}

func New() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var handler http.Handler = r.mux
	for _, mw := range r.middleware {
		handler = mw(handler)
	}
	handler.ServeHTTP(w, req)
}

func (r *Router) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, r)
}

func (r *Router) Use(mw func(http.Handler) http.Handler) {
	r.middleware = append(r.middleware, mw)
}

func (r *Router) Get(pattern string, handler HandlerFunc) {
	r.mux.HandleFunc("GET "+pattern, func(w http.ResponseWriter, req *http.Request) {
		handleError(w, handler(w, req))
	})
}

func (r *Router) Post(pattern string, handler HandlerFunc) {
	r.mux.HandleFunc("POST "+pattern, func(w http.ResponseWriter, req *http.Request) {
		handleError(w, handler(w, req))
	})
}

func (r *Router) Put(pattern string, handler HandlerFunc) {
	r.mux.HandleFunc("PUT "+pattern, func(w http.ResponseWriter, req *http.Request) {
		handleError(w, handler(w, req))
	})
}

func (r *Router) Patch(pattern string, handler HandlerFunc) {
	r.mux.HandleFunc("PATCH "+pattern, func(w http.ResponseWriter, req *http.Request) {
		handleError(w, handler(w, req))
	})
}

func (r *Router) Delete(pattern string, handler HandlerFunc) {
	r.mux.HandleFunc("DELETE "+pattern, func(w http.ResponseWriter, req *http.Request) {
		handleError(w, handler(w, req))
	})
}

func handleError(w http.ResponseWriter, err error) {
	if err != nil {
		var requestError datatypes.RequestErrorInterface
		if errors.As(err, &requestError) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(requestError.StatusCode())
			_, err = w.Write([]byte(requestError.Error()))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *Router) Group(prefix string) *Router {
	group := &Router{
		mux:        http.NewServeMux(),
		middleware: r.middleware,
	}
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if strings.HasPrefix(req.URL.Path, prefix) {
				req.URL.Path = strings.TrimPrefix(req.URL.Path, prefix)
				group.ServeHTTP(w, req)
			} else {
				next.ServeHTTP(w, req)
			}
		})
	})
	return group
}

func (r *Router) IsolatedGroup(prefix string) *Router {
	group := &Router{
		mux: http.NewServeMux(),
	}
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if strings.HasPrefix(req.URL.Path, prefix) {
				req.URL.Path = strings.TrimPrefix(req.URL.Path, prefix)
				group.ServeHTTP(w, req)
			} else {
				next.ServeHTTP(w, req)
			}
		})
	})
	return group
}
