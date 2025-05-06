package api

import (
	"Mxx/api/media"
	"Mxx/api/models"
	"Mxx/api/session"
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGetSessionRoute(t *testing.T) {
	router := GetApiRouter()

	// Simulate a GET request to /session
	req, _ := http.NewRequest("GET", "/session", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	var response models.SessionResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.SessionId == "" {
		t.Errorf("Expected session_id in response, got empty")
	}
	t.Logf("Response : %+v\n", response)
}

func TestUploadRoute(t *testing.T) {
	router := GetApiRouter()

	// Simulate a POST request to /upload without a session ID
	req, _ := http.NewRequest("POST", "/medias", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the response
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
	if !strings.Contains(w.Body.String(), "Session ID is required") {
		t.Fatalf("Response body does not contain the expected error message")
	}

	// simulate with expired session
	sessionId := session.GenerateSessionId()
	session.Update(sessionId, time.Now().Add(-time.Hour))
	req, _ = http.NewRequest("POST", "/medias", nil)
	req.Header.Set("X-Session-Id", sessionId)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	// Simulate a POST request to /upload with not contains file
	sessionId = session.GenerateSessionId()
	session.Update(sessionId, time.Now())
	req, _ = http.NewRequest("POST", "/medias", nil)
	req.Header.Set("X-Session-Id", sessionId)
	// remove dir after test
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Failed to remove media directory: %s", err)
		}
	}(sessionId)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the response
	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	// test upload file , file is using []byte
	fileContent := []byte("test file content")
	fileName := "test.txt"
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		t.Fatal(err)
	}
	_, err = part.Write(fileContent)
	if err != nil {
		t.Fatal(err)
	}

	_ = writer.Close()

	req, _ = http.NewRequest("POST", "/medias", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Session-Id", sessionId)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the response
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check if the file was saved correctly
	targetPath := filepath.Join("data/media", sessionId, fileName)
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		t.Fatalf("Uploaded file does not exist: %v", err)
	}
}

func TestGetSubtitlesRoute(t *testing.T) {
	router := GetApiRouter()

	// Simulate a POST request to /subtitles with media not uploaded
	sessionId := session.GenerateSessionId()
	session.Update(sessionId, time.Now())
	req, _ := http.NewRequest("POST", "/medias/subtitles", nil)
	req.Header.Set("X-Session-Id", sessionId)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the response
	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected status code 404, got %d", w.Code)
	}

	// Add a media path to the manager
	manager := media.GetMediaManager()
	manager.SetMediaPath(sessionId, "../TestSrc/test_ffmpeg.mp4")
	req, _ = http.NewRequest("POST", "/medias/subtitles", nil)
	req.Header.Set("X-Session-Id", sessionId)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", w.Code)
	}

	for {
		req, _ = http.NewRequest("GET", "/medias/task", nil)
		req.Header.Set("X-Session-Id", sessionId)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", w.Code)
		}
		var response models.TaskStateResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if response.TaskState != "Running" {
			t.Logf("Task complett or failed : %s", response.TaskState)
			break
		}
		time.Sleep(5 * time.Second)
	}
	req, _ = http.NewRequest("GET", "/medias/subtitles", nil)
	req.Header.Set("X-Session-Id", sessionId)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", w.Code)
	}
	var response models.ValueResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response.Value == nil {
		t.Fatalf("Expected result in response, got nil")
	}
	t.Log(response.Value.(string))
}
