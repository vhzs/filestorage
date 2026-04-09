package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendJSON(t *testing.T) {
	w := httptest.NewRecorder()
	sendJSON(w, http.StatusCreated, map[string]string{"id": "1"})

	if w.Code != 201 {
		t.Errorf("got status %d, want 201", w.Code)
	}

	var result map[string]string
	json.Unmarshal(w.Body.Bytes(), &result)
	if result["id"] != "1" {
		t.Errorf("got id=%s, want 1", result["id"])
	}
}

func TestSendError(t *testing.T) {
	w := httptest.NewRecorder()
	sendError(w, http.StatusBadRequest, "oops")

	if w.Code != 400 {
		t.Errorf("got %d, want 400", w.Code)
	}

	var body errorBody
	json.Unmarshal(w.Body.Bytes(), &body)
	if body.Error != "oops" {
		t.Errorf("error = %q, want oops", body.Error)
	}
}
