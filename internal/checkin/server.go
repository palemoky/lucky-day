package checkin

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	qrcode "github.com/skip2/go-qrcode"
	"golang.org/x/time/rate"

	"github.com/palemoky/lucky-day/internal/i18n"
	"github.com/palemoky/lucky-day/internal/model"
)

//go:embed templates/*.html
var templatesFS embed.FS

// Server represents the check-in server
type Server struct {
	port           int
	participants   []model.Participant
	mu             sync.RWMutex
	nextID         int
	server         *http.Server
	translator     *i18n.Translator
	newParticipant chan model.Participant
	limiter        *rate.Limiter // Rate limiter for check-in requests
}

// NewServer creates a new check-in server
func NewServer(port int, translator *i18n.Translator) *Server {
	return &Server{
		port:           port,
		participants:   make([]model.Participant, 0),
		nextID:         1,
		translator:     translator,
		newParticipant: make(chan model.Participant, 100),
		limiter:        rate.NewLimiter(rate.Limit(10), 20), // 10 requests/sec, burst of 20
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Serve check-in page
	mux.HandleFunc("/", s.handleCheckInPage)

	// Handle check-in submission
	mux.HandleFunc("/checkin", s.handleCheckIn)

	// Get participant count
	mux.HandleFunc("/count", s.handleCount)

	// Serve QR code image
	mux.HandleFunc("/qr", s.handleQRCode)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	go func() {
		// Start server silently to avoid screen flicker
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Only log errors, not normal startup
			log.Printf("Server error: %v\n", err)
		}
	}()

	return nil
}

// Stop stops the HTTP server
func (s *Server) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

// handleCheckInPage serves the check-in HTML page
func (s *Server) handleCheckInPage(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "zh"
	}

	tmpl, err := template.ParseFS(templatesFS, "templates/checkin.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"Lang":            lang,
		"Title":           s.getTranslation(lang, "qr.title"),
		"NamePlaceholder": s.getTranslation(lang, "qr.name_placeholder"),
		"DeptPlaceholder": s.getTranslation(lang, "qr.dept_placeholder"),
		"Submit":          s.getTranslation(lang, "qr.submit"),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}

// handleCheckIn handles check-in form submission
func (s *Server) handleCheckIn(w http.ResponseWriter, r *http.Request) {
	// Rate limiting
	if !s.limiter.Allow() {
		http.Error(w, "Too many requests, please try again later", http.StatusTooManyRequests)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Input validation
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Validate name length (1-100 characters)
	if len(name) > 100 {
		http.Error(w, "Name is too long (max 100 characters)", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	participant := model.Participant{
		ID:             s.nextID,
		Name:           name,
		WinningHistory: []model.WinningRecord{},
	}
	s.participants = append(s.participants, participant)
	s.nextID++
	s.mu.Unlock()

	// Notify via channel
	select {
	case s.newParticipant <- participant:
	default:
		// Channel full, skip notification
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Check-in successful",
		"id":      participant.ID,
	}) // Ignore encoding error
}

// handleCount returns the current participant count
func (s *Server) handleCount(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	count := len(s.participants)
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]int{
		"count": count,
	}) // Ignore encoding error
}

// handleQRCode serves the QR code image
func (s *Server) handleQRCode(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("http://localhost:%d/?lang=%s", s.port, s.translator.GetLanguage())

	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	_, _ = w.Write(png) // Ignore write error
}

// GetParticipants returns the current list of participants
func (s *Server) GetParticipants() []model.Participant {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy
	participants := make([]model.Participant, len(s.participants))
	copy(participants, s.participants)
	return participants
}

// GetParticipantCount returns the number of checked-in participants
func (s *Server) GetParticipantCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.participants)
}

// GetNewParticipantChannel returns the channel for new participant notifications
func (s *Server) GetNewParticipantChannel() <-chan model.Participant {
	return s.newParticipant
}

// GenerateQRCode generates a QR code image file
func GenerateQRCode(url, outputPath string) error {
	return qrcode.WriteFile(url, qrcode.Medium, 256, outputPath)
}

// getTranslation is a helper to get translations
func (s *Server) getTranslation(lang, key string) string {
	if lang == "en" {
		s.translator.SetLanguage(i18n.English)
	} else {
		s.translator.SetLanguage(i18n.Chinese)
	}
	return s.translator.T(key)
}

// GetURL returns the check-in URL
func (s *Server) GetURL() string {
	return fmt.Sprintf("http://localhost:%d/?lang=%s", s.port, s.translator.GetLanguage())
}

// SaveToExcel saves checked-in participants to Excel file
func (s *Server) SaveToExcel(filePath string) error {
	// This will be implemented using the datasource.SaveParticipantsToExcel function
	// For now, just return nil
	return nil
}

// WaitForParticipants waits for a certain duration to collect participants
func (s *Server) WaitForParticipants(duration time.Duration) {
	time.Sleep(duration)
}
