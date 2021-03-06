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
package commons

import "net/http"

func HttpNoContent(w http.ResponseWriter) {
	HttpError(w, http.StatusNoContent)
}

func HttpUnauthorized(w http.ResponseWriter) {
	HttpError(w, http.StatusUnauthorized)
}

func HttpBadRequest(w http.ResponseWriter) {
	HttpError(w, http.StatusBadRequest)
}

func HttpError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func HttpCheckError(err error, status int, w http.ResponseWriter) {
	if err != nil {
		HttpError(w, status)
	}
}
