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

package routes

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"whipcode/control"
	"whipcode/podman"
	"whipcode/server"
)

func (l *LanguageID) UnmarshalJSON(b []byte) error {
	var intValue int
	if err := json.Unmarshal(b, &intValue); err == nil {
		l.value = strconv.Itoa(intValue)
		return nil
	}

	var strValue string
	if err := json.Unmarshal(b, &strValue); err == nil {
		l.value = strValue
		return nil
	}

	return json.Unmarshal(b, &l.value)
}

func getLanguageConfig() map[string]map[string]string {
	return map[string]map[string]string{
		"1":  {"cmd": "/usr/bin/python3", "ext": "py", "image": "whipcode-python"},
		"2":  {"cmd": "/usr/bin/node", "ext": "js", "image": "whipcode-nodejs"},
		"3":  {"cmd": "/bin/bash", "ext": "sh", "image": "whipcode-bash"},
		"4":  {"cmd": "/usr/bin/perl", "ext": "pl", "image": "whipcode-perl"},
		"5":  {"cmd": "/usr/bin/lua", "ext": "lua", "image": "whipcode-lua"},
		"6":  {"cmd": "/usr/bin/ruby", "ext": "rb", "image": "whipcode-ruby"},
		"7":  {"cmd": "sh /usr/bin/c-run", "ext": "c", "image": "whipcode-c"},
		"8":  {"cmd": "sh /usr/bin/cpp-run", "ext": "cpp", "image": "whipcode-cpp"},
		"9":  {"cmd": "sh /usr/bin/rs-run", "ext": "rs", "image": "whipcode-rust"},
		"10": {"cmd": "sh /usr/bin/f90-run", "ext": "f90", "image": "whipcode-fortran"},
		"11": {"cmd": "sh /usr/bin/hs-run", "ext": "hs", "image": "whipcode-haskell"},
		"12": {"cmd": "/usr/bin/java", "ext": "java", "image": "whipcode-java"},
		"13": {"cmd": "sh /usr/bin/go-run", "ext": "go", "image": "whipcode-go"},
		"14": {"cmd": "sh /usr/bin/ts-run", "ext": "ts", "image": "whipcode-typescript"},
		"15": {"cmd": "/usr/bin/sbcl --script", "ext": "lisp", "image": "whipcode-lisp"},
		"16": {"cmd": "/usr/bin/racket", "ext": "rkt", "image": "whipcode-racket"},
	}
}

func Run(w http.ResponseWriter, r *http.Request) {
	masterKey := r.Header.Get("X-Master-Key")
	if masterKey == "" {
		server.Send(w, http.StatusUnauthorized, []byte(`{"detail": "unauthorized"}`))
		return
	}

	ks, _ := r.Context().Value(server.KeyStoreContextKey).(*control.KeyStore)
	if !ks.CheckKey(masterKey, r.Context().Value(server.MasterKeyContextKey).([]string)) {
		server.Send(w, http.StatusUnauthorized, []byte(`{"detail": "unauthorized"}`))
		return
	}

	mimeType := r.Header.Get("Content-Type")
	if strings.Split(mimeType, ";")[0] != "application/json" {
		server.Send(w, http.StatusUnsupportedMediaType, []byte(`{"detail": "unsupported media type"}`))
		return
	}

	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		server.Send(w, http.StatusBadRequest, []byte(`{"detail": "invalid request format"}`))
		return
	}

	langConfig, exists := getLanguageConfig()[user.LanguageID.value]
	if !exists {
		server.Send(w, http.StatusBadRequest, []byte(`{"detail": "invalid value for parameter language_id, must be a value from 1 to 13"}`))
		return
	}

	codeBytes, err := base64.StdEncoding.DecodeString(user.Code)
	if err != nil || user.Code == "" {
		server.Send(w, http.StatusBadRequest, []byte(`{"detail": "invalid value for parameter code, must be a base64 encoded string"}`))
		return
	}

	cmd := langConfig["cmd"]
	ext := langConfig["ext"]
	img := langConfig["image"]

	ex, _ := r.Context().Value(server.ExecutorContextKey).(podman.Executor)
	status, result := ex.RunCode(string(codeBytes), cmd, ext, img, r.Context().Value(server.EnableCacheContextKey).(bool))
	resultBytes, _ := json.Marshal(result)

	server.Send(w, status, resultBytes)
}
