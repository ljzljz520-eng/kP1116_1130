package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"showroom/internal/display"
	"showroom/internal/gesture"
	"showroom/internal/model"
	"showroom/internal/workflow"
)

type Handler struct {
	flow        *workflow.Orchestrator
	galleryName string
}

func NewHandler(flow *workflow.Orchestrator, galleryName string) *Handler {
	return &Handler{flow: flow, galleryName: galleryName}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready", "gallery": h.galleryName})
}

func (h *Handler) Gesture(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var request struct {
		SessionID string `json:"session_id"`
		Signal    string `json:"signal"`
		Frame     int    `json:"frame"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid gesture request")
		return
	}
	signal, err := gesture.Parse(request.Signal, request.Frame)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	decision, state, err := h.flow.HandleSignal(r.Context(), request.SessionID, signal)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"decision": decision, "state": state})
}

func (h *Handler) CustomPhrase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var request struct {
		SessionID string          `json:"session_id"`
		Mode      model.SceneMode `json:"mode"`
		ID        string          `json:"id"`
		Text      string          `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid phrase request")
		return
	}
	state, err := h.flow.ConfirmCustom(r.Context(), request.SessionID, request.Mode, request.ID, request.Text)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, state)
}

func (h *Handler) End(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	sessionID := strings.TrimSpace(r.URL.Query().Get("session_id"))
	if sessionID == "" {
		writeError(w, http.StatusBadRequest, "session_id is required")
		return
	}
	state, err := h.flow.EndDisplay(r.Context(), sessionID)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, state)
}

func (h *Handler) Snapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "GET required")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"gallery": h.galleryName, "hint": "use gesture endpoint", "panel": h.flow.Panel()})
}

func (h *Handler) Panel(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		writeJSON(w, http.StatusOK, h.flow.Panel())
		return
	}
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "GET or POST required")
		return
	}
	var request struct {
		SessionID string              `json:"session_id"`
		Action    display.PanelAction `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid panel request")
		return
	}
	snapshot, err := h.flow.ApplyPanelAction(r.Context(), request.SessionID, request.Action)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, snapshot)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		return
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": fmt.Sprintf("%s", message)})
}
