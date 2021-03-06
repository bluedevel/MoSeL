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
package handler

import (
	"net/http"
	"encoding/json"

	"github.com/bluedevel/mosel/api"
	"github.com/bluedevel/mosel/moseld/server/context"
	"github.com/bluedevel/mosel/moselserver"
)

// Handler for providing a login and credential validation functionality.
type loginHandler struct {
	moselserver.SecureHandler

	ctx *context.MoseldServerContext
}

func NewLoginHandler(ctx *context.MoseldServerContext) loginHandler {
	return loginHandler{
		ctx: ctx,
	}
}

func (handler loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := api.NewLoginResponse()

	user, passwd, enabled := r.BasicAuth()

	if !enabled || !handler.ctx.Auth.Authenticate(user, passwd) {
		resp.Successful = false
	} else {
		key, validTo := handler.ctx.Sessions.NewSession()
		resp.Key = key
		resp.ValidTo = validTo
		resp.Successful = true
	}

	json.NewEncoder(w).Encode(resp)
}

func (handler loginHandler) GetPath() string {
	return "/login"
}

func (handler loginHandler) Secure() bool {
	return true
}
