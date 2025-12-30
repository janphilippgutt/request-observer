package httpapi

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type ObserveRequest struct {
	Source string `json:"source"`
	Event  string `json:"event"`
	UserID string `json:"user_id"`
}

func ObserveHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	// implementation follows
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	var observeRequest ObserveRequest
	err := json.NewDecoder(r.Body).Decode(&observeRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if rand.Float64() < 0.2 {
		http.Error(w, "simulated processing failure", http.StatusInternalServerError)
		return
	}

	// Simulate work (50–500ms)
	delay := time.Duration(50+rand.Intn(450)) * time.Millisecond
	time.Sleep(delay)

	w.WriteHeader(http.StatusAccepted)
}
