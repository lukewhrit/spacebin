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
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lukewhrit/spacebin/internal/util"
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

func TestHandleCreateBodyJSON(t *testing.T) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(map[string]interface{}{
		"content": "Hello, world!",
	})

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Type", "application/json")
	body, err := util.HandleCreateBody(400000, req)

	require.NoError(t, err)
	require.Equal(t, "Hello, world!", body.Content)
}

func TestHandleCreateBodyMultipart(t *testing.T) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	fw, _ := writer.CreateFormField("content")
	io.Copy(fw, strings.NewReader("Hello, world!"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	body, err := util.HandleCreateBody(400000, req)

	require.NoError(t, err)
	require.Equal(t, "Hello, world!", body.Content)
}

func TestHandleCreateBodyNoContent(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", &bytes.Buffer{})
	body, err := util.HandleCreateBody(400000, req)

	require.Error(t, err)
	require.Equal(t, "", body.Content)
}

// TestHandleCreateBodyInvalidJSON tests HandleCreateBody with invalid JSON
func TestHandleCreateBodyInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	_, err := util.HandleCreateBody(400000, req)

	require.Error(t, err)
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

// TestValidateBodySigninRequest tests ValidateBody with SigninRequest
func TestValidateBodySigninRequest(t *testing.T) {
	// Valid signin request
	validReq := util.SigninRequest{
		Username: "testuser",
		Password: "validpassword123",
	}
	require.NoError(t, util.ValidateBody(100, validReq))

	// Missing username
	invalidReq := util.SigninRequest{
		Username: "",
		Password: "validpassword123",
	}
	require.Error(t, util.ValidateBody(100, invalidReq))

	// Password too short
	invalidReq2 := util.SigninRequest{
		Username: "testuser",
		Password: "short",
	}
	require.Error(t, util.ValidateBody(100, invalidReq2))

	// Password too long (>128 chars)
	longPassword := strings.Repeat("a", 129)
	invalidReq3 := util.SigninRequest{
		Username: "testuser",
		Password: longPassword,
	}
	require.Error(t, util.ValidateBody(100, invalidReq3))
}

// TestValidateBodySignupRequest tests ValidateBody with SignupRequest
func TestValidateBodySignupRequest(t *testing.T) {
	// Valid signup request
	validReq := util.SignupRequest{
		Username: "testuser",
		Password: "validpassword123",
	}
	require.NoError(t, util.ValidateBody(100, validReq))

	// Missing username
	invalidReq := util.SignupRequest{
		Username: "",
		Password: "validpassword123",
	}
	require.Error(t, util.ValidateBody(100, invalidReq))

	// Password too short
	invalidReq2 := util.SignupRequest{
		Username: "testuser",
		Password: "short",
	}
	require.Error(t, util.ValidateBody(100, invalidReq2))

	// Password too long (>128 chars)
	longPassword := strings.Repeat("a", 129)
	invalidReq3 := util.SignupRequest{
		Username: "testuser",
		Password: longPassword,
	}
	require.Error(t, util.ValidateBody(100, invalidReq3))

	// Missing password
	invalidReq4 := util.SignupRequest{
		Username: "testuser",
		Password: "",
	}
	require.Error(t, util.ValidateBody(100, invalidReq4))
}

// TestHandleSignupBodyJSON tests HandleSignupBody with JSON
func TestHandleSignupBodyJSON(t *testing.T) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(map[string]interface{}{
		"username": "testuser",
		"password": "testpassword1234",
	})

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Type", "application/json")
	body, err := util.HandleSignupBody(400000, req)

	require.NoError(t, err)
	require.Equal(t, "testuser", body.Username)
	require.Equal(t, "testpassword1234", body.Password)
}

// TestHandleSignupBodyMultipart tests HandleSignupBody with multipart
func TestHandleSignupBodyMultipart(t *testing.T) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	writer.WriteField("username", "testuser")
	writer.WriteField("password", "testpassword1234")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	body, err := util.HandleSignupBody(400000, req)

	require.NoError(t, err)
	require.Equal(t, "testuser", body.Username)
	require.Equal(t, "testpassword1234", body.Password)
}

// TestHandleSignupBodyNoContent tests HandleSignupBody with no content type
func TestHandleSignupBodyNoContent(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", &bytes.Buffer{})
	body, err := util.HandleSignupBody(400000, req)

	require.NoError(t, err)
	require.Equal(t, "", body.Username)
	require.Equal(t, "", body.Password)
}

// TestHandleSignupBodyInvalidJSON tests HandleSignupBody with invalid JSON
func TestHandleSignupBodyInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	_, err := util.HandleSignupBody(400000, req)

	require.Error(t, err)
}

// TestHandleSigninBodyJSON tests HandleSigninBody with JSON
func TestHandleSigninBodyJSON(t *testing.T) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(map[string]interface{}{
		"username": "testuser",
		"password": "testpassword1234",
	})

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Type", "application/json")
	body, err := util.HandleSigninBody(400000, req)

	require.NoError(t, err)
	require.Equal(t, "testuser", body.Username)
	require.Equal(t, "testpassword1234", body.Password)
}

// TestHandleSigninBodyMultipart tests HandleSigninBody with multipart
func TestHandleSigninBodyMultipart(t *testing.T) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	writer.WriteField("username", "testuser")
	writer.WriteField("password", "testpassword1234")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	body, err := util.HandleSigninBody(400000, req)

	require.NoError(t, err)
	require.Equal(t, "testuser", body.Username)
	require.Equal(t, "testpassword1234", body.Password)
}

// TestHandleSigninBodyNoContent tests HandleSigninBody with no content type
func TestHandleSigninBodyNoContent(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", &bytes.Buffer{})
	body, err := util.HandleSigninBody(400000, req)

	require.NoError(t, err)
	require.Equal(t, "", body.Username)
	require.Equal(t, "", body.Password)
}

// TestHandleSigninBodyInvalidJSON tests HandleSigninBody with invalid JSON
func TestHandleSigninBodyInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	_, err := util.HandleSigninBody(400000, req)

	require.Error(t, err)
}

// TestHandleSignupBodyMultipartError tests HandleSignupBody with multipart parse error
func TestHandleSignupBodyMultipartError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("invalid multipart"))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=----test")
	_, err := util.HandleSignupBody(1, req) // Small maxSize to trigger error

	require.Error(t, err)
}

// TestHandleSigninBodyMultipartError tests HandleSigninBody with multipart parse error
func TestHandleSigninBodyMultipartError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("invalid multipart"))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=----test")
	_, err := util.HandleSigninBody(1, req) // Small maxSize to trigger error

	require.Error(t, err)
}

// TestHandleCreateBodyMultipartError tests HandleCreateBody with multipart parse error
func TestHandleCreateBodyMultipartError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("invalid multipart data"))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=----boundary")
	_, err := util.HandleCreateBody(1, req) // Very small maxSize to trigger error

	require.Error(t, err)
}

// TestRenderError is tested indirectly through server.StaticDocument error paths
// Testing it directly in the util package would require complex embed.FS setup
// that mirrors the server package structure. The function is fully covered
// by server integration tests.
