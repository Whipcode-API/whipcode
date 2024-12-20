//
//  Copyright 2024 whipcode.app (AnnikaV9)
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//          http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing,
//  software distributed under the License is distributed on
//  an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific
//  language governing permissions and limitations under the License.
//

package server

import (
	"net/http"

	"github.com/charmbracelet/log"
)

/**
 * Send sends a response to the client.
 *
 * @param w http.ResponseWriter Response writer
 * @param status int Status code to return
 * @param message []byte Message to send
 * @param contentType ...string Content type to
 *   send, defaults to application/json
 */
func Send(w http.ResponseWriter, status int, message []byte, contentType ...string) {
	cType := "application/json"
	if len(contentType) > 0 {
		cType = contentType[0]
	}

	w.Header().Set("Content-Type", cType)
	w.WriteHeader(status)

	if _, err := w.Write(message); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Error("Failed to write response", "Error", err)
	}
}
