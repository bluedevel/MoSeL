/*
 * Copyright 2016 Robin Engel
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package moselserver

import (
	"net/http"
	"log"
	"fmt"
	"github.com/gorilla/mux"

	"github.com/bluedevel/mosel/config"
	"strconv"
	"reflect"
)

// The abstract http-server type underlying the mosel servers.
type MoselServer struct {
	Config  moselconfig.MoselServerConfig
	Context *MoselServerContext

	Filters   []Filter
	Handlers  []http.Handler
	InitFuncs []func() error
}

// Boot up the server
// 1: Run init functions
// 2: Init request handlers and wrap them with gorilla/mux
// 3: Handle the gorilla/mux router
func (server *MoselServer) Run() error {

	//initializing server context
	server.InitFuncs = append([]func() error{
		server.initAuth,
		server.initDataSources,
		server.initSessionCache,
	}, server.InitFuncs...)

	err := server.initContext()

	if ! server.Context.IsInitialized {

		if err != nil {
			return err
		}

		return fmt.Errorf("Mosel Server - Run: Context wasn't initialized correctly")
	}

	//init router and handlers
	r := mux.NewRouter()
	server.initHandler(r)
	http.Handle("/", r)

	addr := server.Config.Http.BindAddress
	log.Printf("Binding http server to %s", addr)

	//do async jobs after initialization here
	errors := make(chan error)

	go func() {
		errors <- http.ListenAndServe(addr, nil)
	}()

	return <-errors
}

/*
 * Initialize Context
 */

// Initialize the server context.
// The configured init functions will be called and on success server.Context.IsInitialized will be set to true.
func (server *MoselServer) initContext() error {
	server.Context = &MoselServerContext{}

	for _, fn := range server.InitFuncs {
		err := fn()

		if (err != nil) {
			return err
		}
	}

	server.Context.IsInitialized = true
	return nil
}

// Initialize the configured authentication method
func (server *MoselServer) initAuth() error {
	config := server.Config

	var enabledCount int = 0

	if config.AuthStatic.Enabled {
		enabledCount++
		server.Context.Auth = &AuthStatic{
			Users: config.Users,
		}
	}

	if config.AuthSys.Enabled {
		enabledCount++
		server.Context.Auth = &AuthSys{
			AllowedUsers: config.AuthSys.AllowedUsers,
		}
	}

	if config.AuthMySQL.Enabled {
		enabledCount++
	}

	if config.AuthTrue.Enabled {
		enabledCount++
		log.Println("Using AuthTrue! This is for debug purposes only, make sure you don't deploy this in production")
		server.Context.Auth = &AuthTrue{}
	}

	if enabledCount > 1 {
		return fmt.Errorf("More than one auth service enabled")
	} else if enabledCount == 0 {
		return fmt.Errorf("No auth service configured")
	}

	return nil
}

// Initialize the configured data sources
func (server *MoselServer) initDataSources() error {
	server.Context.DataSources = make(map[string]dataSource)

	for name, config := range server.Config.DataSources {
		var ds dataSource
		var err error

		log.Printf("Initializing data source %s of type %s", name, config.Type)

		if config.Type == "mysql" {
			ds, err = server.initMySql(config.Type, config.Connection)
		} else if config.Type == "mongo" {
			ds, err = server.initMongo(config.Type, config.Connection)
		} else {
			log.Printf("Data source type '%s' not supported; %s",
				config.Type, err)
			continue
		}

		if err != nil {
			log.Printf("Failed to register data source %s of type %s; %s",
				name, config.Type, err)
			continue
		}

		log.Printf("Register data source %s of type %s", name, ds.GetType())
		server.Context.DataSources[name] = ds
	}

	return nil
}

// Initialize the session cache
// This gets executed even if sessions are disabled by the config as the decision on weather to use session is
// taken in the authInterceptor
func (server *MoselServer) initSessionCache() error {
	c := NewSessionCache()
	server.Context.Sessions = *c
	return nil
}

/*
 * Initialize Handler
 */

type WrapHandler struct {
	f func(http.ResponseWriter, *http.Request)
}

func (h WrapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.f(w, r)
}

// Initialize the configured http handlers and wrap them into a gorilla/mux router
func (server *MoselServer) initHandler(r *mux.Router) {
	authFilter := newAuthFilter(server)

	for handlerIndex, _ := range server.Handlers {

		h := server.Handlers[handlerIndex]

		f := func(w http.ResponseWriter, r *http.Request) {
			//h.ServeHTTPContext(server.Context, w, r)
			h.ServeHTTP(w, r)
		}

		secure := false
		if sh, ok := h.(SecureHandler); ok {
			f = chainFilter(f, authFilter)
			secure = sh.Secure()
		}

		for filterIndex, _ := range server.Filters {
			f = chainFilter(f, server.Filters[filterIndex])
		}

		var path string
		if crh, ok := h.(CustomRouteHandler); ok {
			crh.ConfigureRoute(r, WrapHandler{f: f})
			path = "<custom>"
		} else {
			path = getHandlerPath(h)
			r.HandleFunc(path, f)
		}

		log.Printf("Handling %s - secure=%s", path, strconv.FormatBool(secure))
	}
}

func getHandlerPath(h http.Handler) string {
	if pathInfo, ok := h.(PathInfo); ok {
		return pathInfo.GetPath()
	}

	t := reflect.TypeOf(h)
	return "/" + t.Name()
}

func chainFilter(f http.HandlerFunc, filter Filter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter.Apply(w, r, func() {
			f(w, r)
		})
	}
}
