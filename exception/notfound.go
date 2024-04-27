package exception

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type notFound func(w http.ResponseWriter, r *http.Request, params httprouter.Params)
type handler func(w http.ResponseWriter, r *http.Request)

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r)
}

func NotFound(n notFound) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.Params{}
		n(w,r,params)
	}
}


func NotFoundHandler() notFound {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		w.Header().Add("Content-Type", "application/json")
		WriteJson(w, http.StatusNotFound, "not found path", "path yang kamu minta tidak ditemukan")
	}
}
