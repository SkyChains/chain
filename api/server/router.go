// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package server

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	"github.com/skychains/chain/utils/set"
)

var (
	errUnknownBaseURL  = errors.New("unknown base url")
	errUnknownEndpoint = errors.New("unknown endpoint")
	errAlreadyReserved = errors.New("route is either already aliased or already maps to a handle")
)

type router struct {
	lock   sync.RWMutex
	router *mux.Router

	routeLock      sync.Mutex
	reservedRoutes set.Set[string]                    // Reserves routes so that there can't be alias that conflict
	aliases        map[string][]string                // Maps a route to a set of reserved routes
	routes         map[string]map[string]http.Handler // Maps routes to a handler
}

func newRouter() *router {
	return &router{
		router:         mux.NewRouter(),
		reservedRoutes: set.Set[string]{},
		aliases:        make(map[string][]string),
		routes:         make(map[string]map[string]http.Handler),
	}
}

func (r *router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	r.router.ServeHTTP(writer, request)
}

func (r *router) GetHandler(base, endpoint string) (http.Handler, error) {
	r.routeLock.Lock()
	defer r.routeLock.Unlock()

	urlBase, exists := r.routes[base]
	if !exists {
		return nil, errUnknownBaseURL
	}
	handler, exists := urlBase[endpoint]
	if !exists {
		return nil, errUnknownEndpoint
	}
	return handler, nil
}

func (r *router) AddRouter(base, endpoint string, handler http.Handler) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.routeLock.Lock()
	defer r.routeLock.Unlock()

	return r.addRouter(base, endpoint, handler)
}

func (r *router) addRouter(base, endpoint string, handler http.Handler) error {
	if r.reservedRoutes.Contains(base) {
		return fmt.Errorf("%w: %s", errAlreadyReserved, base)
	}

	return r.forceAddRouter(base, endpoint, handler)
}

func (r *router) forceAddRouter(base, endpoint string, handler http.Handler) error {
	endpoints := r.routes[base]
	if endpoints == nil {
		endpoints = make(map[string]http.Handler)
	}
	url := base + endpoint
	if _, exists := endpoints[endpoint]; exists {
		return fmt.Errorf("failed to create endpoint as %s already exists", url)
	}

	endpoints[endpoint] = handler
	r.routes[base] = endpoints

	// Name routes based on their URL for easy retrieval in the future
	route := r.router.Handle(url, handler)
	if route == nil {
		return fmt.Errorf("failed to create new route for %s", url)
	}
	route.Name(url)

	var err error
	if aliases, exists := r.aliases[base]; exists {
		for _, alias := range aliases {
			if innerErr := r.forceAddRouter(alias, endpoint, handler); err == nil {
				err = innerErr
			}
		}
	}
	return err
}

func (r *router) AddAlias(base string, aliases ...string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.routeLock.Lock()
	defer r.routeLock.Unlock()

	for _, alias := range aliases {
		if r.reservedRoutes.Contains(alias) {
			return fmt.Errorf("%w: %s", errAlreadyReserved, alias)
		}
	}

	for _, alias := range aliases {
		r.reservedRoutes.Add(alias)
	}

	r.aliases[base] = append(r.aliases[base], aliases...)

	var err error
	if endpoints, exists := r.routes[base]; exists {
		for endpoint, handler := range endpoints {
			for _, alias := range aliases {
				if innerErr := r.forceAddRouter(alias, endpoint, handler); err == nil {
					err = innerErr
				}
			}
		}
	}
	return err
}
