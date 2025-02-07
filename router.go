package resty

import (
	"net/http"

	"github.com/go-zoo/bone"
)

type Router interface {
	// Get add a new route to the Router with the Get method
	Get(path string, action ActionFunc) Router

	// Post add a new route to the Router with the Post method
	Post(path string, action ActionFunc) Router

	// Put add a new route to the Router with the Put method
	Put(path string, action ActionFunc) Router

	// Delete add a new route to the Router with the Delete method
	Delete(path string, action ActionFunc) Router

	// Head add a new route to the Router with the Head method
	Head(path string, action ActionFunc) Router

	// Patch add a new route to the Router with the Patch method
	Patch(path string, action ActionFunc) Router

	// Options add a new route to the Router with the Options method
	Options(path string, action ActionFunc) Router

	ServeHTTP(rw http.ResponseWriter, request *http.Request)
}

type router struct {
	mux *bone.Mux
}

func (r *router) Get(path string, action ActionFunc) Router {
	r.mux.Get(path, HandleAction(action))
	return r
}

func (r *router) Post(path string, action ActionFunc) Router {
	r.mux.Post(path, HandleAction(action))
	return r
}

func (r *router) Put(path string, action ActionFunc) Router {
	r.mux.Put(path, HandleAction(action))
	return r
}

func (r *router) Delete(path string, action ActionFunc) Router {
	r.mux.Delete(path, HandleAction(action))
	return r
}

func (r *router) Head(path string, action ActionFunc) Router {
	r.mux.Head(path, HandleAction(action))
	return r
}

func (r *router) Patch(path string, action ActionFunc) Router {
	r.mux.Patch(path, HandleAction(action))
	return r
}

func (r *router) Options(path string, action ActionFunc) Router {
	r.mux.Options(path, HandleAction(action))
	return r
}

func (r *router) ServeHTTP(rw http.ResponseWriter, request *http.Request) {
	r.mux.ServeHTTP(rw, request)
}

func NewRouter() Router {
	return &router{mux: bone.New()}
}
