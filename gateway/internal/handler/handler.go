// Package handler implements the HTTP transport for the message service.
package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/ReshetovItsMe/secret-message/gateway/internal/service"
)

// Handler is the HTTP surface of the gateway. It stays thin: validation and
// response shaping only — all logic lives in service.Service.
type Handler struct {
	svc  *service.Service
	log  *slog.Logger
	val  *validator.Validate
}

// New creates a Handler.
func New(svc *service.Service, log *slog.Logger) *Handler {
	return &Handler{
		svc: svc,
		log: log,
		val: validator.New(validator.WithRequiredStructEnabled()),
	}
}

type createMessageRequest struct {
	Message string `json:"message" validate:"required,max=100000"`
}

type createMessageResponse struct {
	MessageID string `json:"messageId"`
}

type getMessageResponse struct {
	Message string `json:"message"`
}

type errorBody struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Create handles POST /message.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createMessageRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	if err := h.val.Struct(req); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "message is required and must be a string")
		return
	}

	id, err := h.svc.CreateMessage(r.Context(), req.Message)
	if err != nil {
		h.log.Error("create message failed", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL", "failed to create message")
		return
	}

	writeJSON(w, http.StatusCreated, createMessageResponse{MessageID: id})
}

// Get handles GET /message?messageId=...
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("messageId")
	if id == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "messageId query parameter is required")
		return
	}

	message, err := h.svc.GetMessage(r.Context(), id)
	if errors.Is(err, service.ErrMessageNotFound) {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "message does not exist or has already been read")
		return
	}
	if err != nil {
		h.log.Error("get message failed", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL", "failed to read message")
		return
	}

	writeJSON(w, http.StatusOK, getMessageResponse{Message: message})
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorBody{Error: errorDetail{Code: code, Message: message}})
}
