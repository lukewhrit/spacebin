/*
 * Copyright 2020-2024 Luke Whritenour

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 *     http://www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"net/http"

	"github.com/lukewhrit/spacebin/internal/util"
)

func (s *Server) GetConfig(w http.ResponseWriter, r *http.Request) {
	if err := util.WriteJSON(w, http.StatusOK, s.Config); err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
