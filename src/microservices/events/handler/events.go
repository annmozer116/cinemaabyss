package handler

import (
	"encoding/json"
	"io"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	"github.com/cinemaabyss/microservices/events/config"
	"github.com/cinemaabyss/microservices/events/kafka"
	"github.com/cinemaabyss/microservices/events/model"
	"github.com/gorilla/mux"
)

type EventHandler struct {
	producer *kafka.Producer
	cfg      *config.Config
}

func NewEventHandler(producer *kafka.Producer, cfg *config.Config) *EventHandler {
	return &EventHandler{
		producer: producer,
		cfg:      cfg,
	}
}

func (h *EventHandler) CreateUserEvent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event := model.Event{
		Type: model.UserEvent,
		Data: string(body),
	}

	if err := h.producer.SendEvent(h.cfg.KafkaUserTopic, event); err != nil {
		http.Error(w, "Failed to send event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	eventData := CreateResponse(event.Type, data)
	json.NewEncoder(w).Encode(eventData)
}
func (h *EventHandler) CreatePaymentEvent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event := model.Event{
		Type: model.PaymentEvent,
		Data: string(body),
	}

	if err := h.producer.SendEvent(h.cfg.KafkaPaymentTopic, event); err != nil {
		http.Error(w, "Failed to send event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	eventData := CreateResponse(event.Type, data)
	json.NewEncoder(w).Encode(eventData)
}

func (h *EventHandler) CreateMovieEvent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event := model.Event{
		Type: model.MovieEvent,
		Data: string(body),
	}

	if err := h.producer.SendEvent(h.cfg.KafkaMovieTopic, event); err != nil {
		http.Error(w, "Failed to send event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	eventData := CreateResponse(event.Type, data)

	json.NewEncoder(w).Encode(eventData)
}

func (h *EventHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": true,
	})
}

func SetupRoutes(r *mux.Router, handler *EventHandler) {
	r.HandleFunc("/api/events/health", handler.HealthCheck).Methods("GET")
	r.HandleFunc("/api/events/user", handler.CreateUserEvent).Methods("POST")
	r.HandleFunc("/api/events/payment", handler.CreatePaymentEvent).Methods("POST")
	r.HandleFunc("/api/events/movie", handler.CreateMovieEvent).Methods("POST")
}

// Так как текущая библиотека не возвращает обратную информацию от kafka о партиции, смещении и идентификаторе события, в качестве MVP используем фиктивные данные
func CreateResponse(eventType model.EventType, payload interface{}) map[string]interface{} {

	eventData := map[string]interface{}{
		"id":        strconv.Itoa(rand.IntN(100)),
		"type":      "",
		"timestamp": time.Now(),
		"payload":   payload,
	}

	switch eventType {
	case model.UserEvent:
		eventData["type"] = "user"
	case model.MovieEvent:
		eventData["type"] = "movie"
	case model.PaymentEvent:
		eventData["type"] = "payment"
	}

	resp := map[string]interface{}{
		"status":    "success",
		"partition": 0,
		"offset":    42,
		"event":     eventData,
	}
	return resp
}
