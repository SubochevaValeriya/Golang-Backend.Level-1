package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Testing GET request with query parameter (extension)
func TestFiles(t *testing.T) {
	req, err := http.NewRequest("GET", "/files?extension=.yaml", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	files(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Проверяем тело ответа
	expected := []File{{
		Name:      "testfile.yaml",
		Extension: ".yaml",
		Size:      0,
	}}

	var result []File

	json.Unmarshal(rr.Body.Bytes(), &result)

	for i := 0; i < len(expected); i++ {
		if result[i] != expected[i] || len(expected) != len(result) {
			t.Errorf("handler returned unexpected body: got %v want %v",
				result, expected)
		}
	}
}

// Testing uploading of doubles on server
func TestUploading(t *testing.T) {
	file, _ := os.Open("testfile.txt")
	defer file.Close()

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok!")
	}))
	defer ts.Close()
	uploadHandler := &UploadHandler{
		UploadDir: "upload",
		HostAddr:  ts.URL,
	}

	uploadHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `testfile_3`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
