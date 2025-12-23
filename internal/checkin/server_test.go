package checkin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/palemoky/lucky-day/internal/i18n"
	"github.com/palemoky/lucky-day/internal/model"
)

func TestNewServer(t *testing.T) {
	translator := i18n.NewTranslator(i18n.Chinese)
	server := NewServer(8888, translator)

	if server == nil {
		t.Fatal("Expected server to be created, got nil")
	}

	if server.port != 8888 {
		t.Errorf("Expected port 8888, got %d", server.port)
	}

	if server.nextID != 1 {
		t.Errorf("Expected nextID to be 1, got %d", server.nextID)
	}

	if len(server.participants) != 0 {
		t.Errorf("Expected empty participants list, got %d participants", len(server.participants))
	}

	if server.limiter == nil {
		t.Error("Expected rate limiter to be initialized")
	}
}

func TestHandleCheckIn_Success(t *testing.T) {
	translator := i18n.NewTranslator(i18n.Chinese)
	server := NewServer(8888, translator)

	// Create form data
	form := url.Values{}
	form.Add("name", "张三")

	req := httptest.NewRequest(http.MethodPost, "/checkin", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	server.handleCheckIn(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Parse JSON response
	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || !success {
		t.Error("Expected success to be true")
	}

	// Verify participant was added
	participants := server.GetParticipants()
	if len(participants) != 1 {
		t.Errorf("Expected 1 participant, got %d", len(participants))
	}

	if participants[0].Name != "张三" {
		t.Errorf("Expected participant name '张三', got '%s'", participants[0].Name)
	}
}

func TestHandleCheckIn_EmptyName(t *testing.T) {
	translator := i18n.NewTranslator(i18n.Chinese)
	server := NewServer(8888, translator)

	form := url.Values{}
	form.Add("name", "")

	req := httptest.NewRequest(http.MethodPost, "/checkin", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	server.handleCheckIn(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandleCheckIn_WhitespaceOnly(t *testing.T) {
	translator := i18n.NewTranslator(i18n.Chinese)
	server := NewServer(8888, translator)

	form := url.Values{}
	form.Add("name", "   ")

	req := httptest.NewRequest(http.MethodPost, "/checkin", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	server.handleCheckIn(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandleCheckIn_NameTooLong(t *testing.T) {
	translator := i18n.NewTranslator(i18n.Chinese)
	server := NewServer(8888, translator)

	// Create a name with 101 characters
	longName := strings.Repeat("a", 101)
	form := url.Values{}
	form.Add("name", longName)

	req := httptest.NewRequest(http.MethodPost, "/checkin", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	server.handleCheckIn(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHandleCheckIn_MethodNotAllowed(t *testing.T) {
	translator := i18n.NewTranslator(i18n.Chinese)
	server := NewServer(8888, translator)

	req := httptest.NewRequest(http.MethodGet, "/checkin", nil)
	w := httptest.NewRecorder()

	server.handleCheckIn(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestHandleCheckIn_RateLimiting(t *testing.T) {
	translator := i18n.NewTranslator(i18n.Chinese)
	server := NewServer(8888, translator)

	// Make requests rapidly to trigger rate limiting
	form := url.Values{}
	form.Add("name", "测试用户")

	successCount := 0
	rateLimitedCount := 0

	// Try 30 requests (burst is 20, so some should be rate limited)
	for i := 0; i < 30; i++ {
		req := httptest.NewRequest(http.MethodPost, "/checkin", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		server.handleCheckIn(w, req)

		switch w.Code {
		case http.StatusOK:
			successCount++
		case http.StatusTooManyRequests:
			rateLimitedCount++
		}
	}

	// Should have some successful requests and some rate limited
	if successCount == 0 {
		t.Error("Expected some successful requests")
	}

	if rateLimitedCount == 0 {
		t.Error("Expected some requests to be rate limited")
	}

	t.Logf("Success: %d, Rate limited: %d", successCount, rateLimitedCount)
}

func TestHandleCount(t *testing.T) {
	translator := i18n.NewTranslator(i18n.Chinese)
	server := NewServer(8888, translator)

	// Add some participants
	server.participants = []model.Participant{
		{ID: 1, Name: "User1"},
		{ID: 2, Name: "User2"},
		{ID: 3, Name: "User3"},
	}

	req := httptest.NewRequest(http.MethodGet, "/count", nil)
	w := httptest.NewRecorder()

	server.handleCount(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]int
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if count, ok := response["count"]; !ok || count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}
}

func TestGetParticipants(t *testing.T) {
	translator := i18n.NewTranslator(i18n.Chinese)
	server := NewServer(8888, translator)

	// Add participants
	server.participants = []model.Participant{
		{ID: 1, Name: "User1"},
		{ID: 2, Name: "User2"},
	}

	participants := server.GetParticipants()

	if len(participants) != 2 {
		t.Errorf("Expected 2 participants, got %d", len(participants))
	}

	// Verify it's a copy (modifying returned slice shouldn't affect server)
	participants[0].Name = "Modified"
	if server.participants[0].Name == "Modified" {
		t.Error("GetParticipants should return a copy, not the original slice")
	}
}

func TestGetParticipantCount(t *testing.T) {
	translator := i18n.NewTranslator(i18n.Chinese)
	server := NewServer(8888, translator)

	if count := server.GetParticipantCount(); count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}

	server.participants = []model.Participant{
		{ID: 1, Name: "User1"},
		{ID: 2, Name: "User2"},
		{ID: 3, Name: "User3"},
	}

	if count := server.GetParticipantCount(); count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}
}

func TestGetURL(t *testing.T) {
	translator := i18n.NewTranslator(i18n.Chinese)
	server := NewServer(8888, translator)

	url := server.GetURL()
	expected := "http://localhost:8888/?lang=zh"

	if url != expected {
		t.Errorf("Expected URL '%s', got '%s'", expected, url)
	}

	// Test with English
	translator.SetLanguage(i18n.English)
	url = server.GetURL()
	expected = "http://localhost:8888/?lang=en"

	if url != expected {
		t.Errorf("Expected URL '%s', got '%s'", expected, url)
	}
}

func TestGetNewParticipantChannel(t *testing.T) {
	translator := i18n.NewTranslator(i18n.Chinese)
	server := NewServer(8888, translator)

	ch := server.GetNewParticipantChannel()
	if ch == nil {
		t.Error("Expected non-nil channel")
	}

	// Test that channel receives participants
	go func() {
		form := url.Values{}
		form.Add("name", "测试")
		req := httptest.NewRequest(http.MethodPost, "/checkin", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		server.handleCheckIn(w, req)
	}()

	select {
	case participant := <-ch:
		if participant.Name != "测试" {
			t.Errorf("Expected participant name '测试', got '%s'", participant.Name)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for participant notification")
	}
}

func TestGenerateQRCode(t *testing.T) {
	// Test QR code generation
	url := "http://localhost:8888"
	outputPath := "test_qr.png"

	err := GenerateQRCode(url, outputPath)
	if err != nil {
		t.Errorf("Failed to generate QR code: %v", err)
	}

	// Clean up
	// Note: In a real test, you'd want to check if file exists and delete it
}

func TestConcurrentCheckIns(t *testing.T) {
	translator := i18n.NewTranslator(i18n.Chinese)
	server := NewServer(8888, translator)

	// Simulate concurrent check-ins
	const numGoroutines = 10
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			form := url.Values{}
			form.Add("name", "User"+string(rune('A'+id)))

			req := httptest.NewRequest(http.MethodPost, "/checkin", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			// Add small delay to spread out requests
			time.Sleep(100 * time.Millisecond)
			server.handleCheckIn(w, req)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Check that all participants were added (accounting for rate limiting)
	count := server.GetParticipantCount()
	if count == 0 {
		t.Error("Expected some participants to be added")
	}

	t.Logf("Successfully added %d participants concurrently", count)
}

func TestWaitForParticipants(t *testing.T) {
	translator := i18n.NewTranslator(i18n.Chinese)
	server := NewServer(8888, translator)

	start := time.Now()
	server.WaitForParticipants(100 * time.Millisecond)
	elapsed := time.Since(start)

	if elapsed < 100*time.Millisecond {
		t.Errorf("Expected to wait at least 100ms, waited %v", elapsed)
	}
}
