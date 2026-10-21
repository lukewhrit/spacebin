/*
* Copyright 2020-2026 Luke Whritenour

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
	"testing/fstest"

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

func TestHandleCreateBodyEncryptedJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"content":"{}","encrypted":true}`))
	req.Header.Set("Content-Type", "application/json")
	body, err := util.HandleCreateBody(400000, req)

	require.NoError(t, err)
	require.Equal(t, "{}", body.Content)
	require.True(t, body.Encrypted)
}

func TestHandleCreateBodyEncryptedMultipart(t *testing.T) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	writer.WriteField("encrypted", "true")
	fw, _ := writer.CreateFormFile("content", "envelope.json")
	io.Copy(fw, strings.NewReader("{}"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	body, err := util.HandleCreateBody(400000, req)

	require.NoError(t, err)
	require.Equal(t, "{}", body.Content)
	require.True(t, body.Encrypted)
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
		util.WriteJSON[map[string]interface{}](w, 200, map[string]interface{}{
			"test": "test",
		})
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
		util.WriteError(w, http.StatusInternalServerError, e)
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

// TestHandleCreateBodyMultipartFileUpload tests HandleCreateBody with a file upload field
func TestHandleCreateBodyMultipartFileUpload(t *testing.T) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// CreateFormFile writes a file part, not a plain form field
	fw, err := writer.CreateFormFile("content", "test.txt")
	require.NoError(t, err)
	io.Copy(fw, strings.NewReader("Hello from file!"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	body, err := util.HandleCreateBody(400000, req)

	require.NoError(t, err)
	require.Equal(t, "Hello from file!", body.Content)
}

// TestHandleCreateBodyMultipartNoContentField tests HandleCreateBody when neither
// a plain text field nor a file field named "content" exists in the form
func TestHandleCreateBodyMultipartNoContentField(t *testing.T) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	writer.WriteField("other", "value")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	_, err := util.HandleCreateBody(400000, req)

	require.Error(t, err)
}

// TestHandleCreateBodyMultipartFileReadError tests HandleCreateBody when reading an uploaded file fails
func TestHandleCreateBodyMultipartFileReadError(t *testing.T) {
	util.SetReadAll(t, func(io.Reader) ([]byte, error) {
		return nil, errors.New("read error")
	})

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	fw, err := writer.CreateFormFile("content", "test.txt")
	require.NoError(t, err)
	io.Copy(fw, strings.NewReader("Hello from file!"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	_, err = util.HandleCreateBody(400000, req)

	require.ErrorContains(t, err, "read error")
}

var errorTemplateFS = fstest.MapFS{
	"web/error.html": {Data: []byte("{{.Status}}: {{.Error}}")},
}

// TestRenderError tests that RenderError writes the error page as HTML
func TestRenderError(t *testing.T) {
	rr := httptest.NewRecorder()

	err := util.RenderError(errorTemplateFS, rr, http.StatusNotFound, errors.New("<missing>"))

	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, rr.Code)
	require.Equal(t, "text/html", rr.Header().Get("Content-Type"))
	require.Equal(t, "404 Not Found: &lt;missing&gt;", rr.Body.String())
}

type failingWriter struct{ *httptest.ResponseRecorder }

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("write error") }

// TestRenderErrorWriteFailure tests that RenderError returns errors from writing the page
func TestRenderErrorWriteFailure(t *testing.T) {
	err := util.RenderError(errorTemplateFS, failingWriter{httptest.NewRecorder()}, http.StatusNotFound, errors.New("missing"))

	require.ErrorContains(t, err, "write error")
}

func TestHandleCollectionBodyJSON(t *testing.T) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(map[string]interface{}{
		"title":     "My notes",
		"documents": []string{"AbCd1234", "https://spaceb.in/EfGh5678.go"},
	})

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Type", "application/json")
	body, err := util.HandleCollectionBody(400000, req)

	require.NoError(t, err)
	require.Equal(t, "My notes", body.Title)
	require.Equal(t, []string{"AbCd1234", "EfGh5678"}, body.Documents)
}

func TestHandleCollectionBodyMultipart(t *testing.T) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	writer.WriteField("title", "My notes")
	writer.WriteField("documents", "AbCd1234, EfGh5678")
	writer.WriteField("documents", "IjKl9012")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	body, err := util.HandleCollectionBody(400000, req)

	require.NoError(t, err)
	require.Equal(t, "My notes", body.Title)
	require.Equal(t, []string{"AbCd1234", "EfGh5678", "IjKl9012"}, body.Documents)
}

func TestHandleCollectionBodyInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	_, err := util.HandleCollectionBody(400000, req)

	require.Error(t, err)
}

func TestHandleCollectionBodyMultipartError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("invalid multipart data"))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=----boundary")
	_, err := util.HandleCollectionBody(1, req)

	require.Error(t, err)
}

func TestHandleCollectionBodyUnsupportedContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("title=x"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	_, err := util.HandleCollectionBody(400000, req)

	require.ErrorContains(t, err, "unsupported Content-Type")
}

func TestParseDocumentIDs(t *testing.T) {
	tests := map[string][]string{
		"":                                       {},
		"AbCd1234":                               {"AbCd1234"},
		"AbCd1234,EfGh5678":                      {"AbCd1234", "EfGh5678"},
		" AbCd1234 ,, EfGh5678\nIjKl9012\t":      {"AbCd1234", "EfGh5678", "IjKl9012"},
		"https://spaceb.in/AbCd1234":             {"AbCd1234"},
		"https://spaceb.in/AbCd1234/":            {"AbCd1234"},
		"https://spaceb.in/AbCd1234.go":          {"AbCd1234"},
		"https://spaceb.in/AbCd1234/raw":         {"AbCd1234"},
		"https://spaceb.in/AbCd1234?reader=true": {"AbCd1234"},
		"/AbCd1234#L12":                          {"AbCd1234"},
		"https://spaceb.in/":                     {},
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			require.Equal(t, expected, util.ParseDocumentIDs(input))
		})
	}
}
