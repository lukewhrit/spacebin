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

package util_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/orca-group/spirit/internal/util"
	"github.com/stretchr/testify/require"
)

func TestValidateBody(t *testing.T) {
	require.NoError(t, util.ValidateBody(100, util.CreateRequest{
		Content: "Test",
	}))

	require.Error(t, util.ValidateBody(2, util.CreateRequest{
		Content: "Test",
	}))

	require.Error(t, util.ValidateBody(2, util.CreateRequest{
		Content: "",
	}))
}

func TestCountLines(t *testing.T) {
	content := "Line 1\nLine 2"

	lines := util.CountLines(content)

	require.Equal(t, lines, template.HTML("<div>1</div><div>2</div>"))
}

func TestHandleBodyJSON(t *testing.T) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(map[string]interface{}{
		"content": "Hello, world!",
	})

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Type", "application/json")
	body, err := util.HandleBody(400000, req)

	require.NoError(t, err)
	require.Equal(t, "Hello, world!", body.Content)
}

func TestHandleBodyMultipart(t *testing.T) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	fw, _ := writer.CreateFormField("content")
	io.Copy(fw, strings.NewReader("Hello, world!"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	body, err := util.HandleBody(400000, req)

	require.NoError(t, err)
	require.Equal(t, "Hello, world!", body.Content)
}

func TestHandleBodyNoContent(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", &bytes.Buffer{})
	body, err := util.HandleBody(400000, req)

	require.NoError(t, err)
	require.Equal(t, "", body.Content)
}

func TestWriteJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := util.WriteJSON[map[string]interface{}](w, 200, map[string]interface{}{
			"test": "test",
		})

		require.NoError(t, err)
	}))
	defer server.Close()

	res, err := http.Get(server.URL)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)

	x, _ := io.ReadAll(res.Body)
	var body map[string]interface{}
	json.Unmarshal(x, &body)

	require.Equal(t, body, map[string]interface{}{
		"payload": map[string]interface{}{
			"test": "test",
		},
		"error": "",
	})
}

func TestWriteError(t *testing.T) {
	e := errors.New("some error")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := util.WriteError(w, http.StatusInternalServerError, e)
		require.NoError(t, err)
	}))
	defer server.Close()

	res, err := http.Get(server.URL)
	require.NoError(t, err)

	require.Equal(t, http.StatusInternalServerError, res.StatusCode)

	x, _ := io.ReadAll(res.Body)
	var body map[string]interface{}
	json.Unmarshal(x, &body)

	require.Equal(t, body, map[string]interface{}{
		"payload": map[string]interface{}{},
		"error":   e.Error(),
	})
}
