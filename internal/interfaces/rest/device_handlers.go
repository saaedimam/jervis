package rest

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/saaedimam/jervis/internal/services/device"
)

// withDeviceAuth enforces device authentication using constant-time comparison on JERVIS_DEVICE_AUTH_KEY.
func (s *Server) withDeviceAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.deviceAuthKey == "" {
			http.Error(w, "Device endpoints unavailable: JERVIS_DEVICE_AUTH_KEY not configured", http.StatusServiceUnavailable)
			return
		}

		auth := r.Header.Get("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(auth, prefix) {
			http.Error(w, "Unauthorized: missing or invalid Bearer token", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(auth, prefix)
		if subtle.ConstantTimeCompare([]byte(token), []byte(s.deviceAuthKey)) != 1 {
			http.Error(w, "Unauthorized: invalid device credentials", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// withDeviceOrAPIAuth allows either device key or main API key (e.g. for status inspect and operator queueing).
func (s *Server) withDeviceOrAPIAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(auth, prefix) {
			http.Error(w, "Unauthorized: missing Bearer token", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(auth, prefix)

		var valid bool
		if s.deviceAuthKey != "" && subtle.ConstantTimeCompare([]byte(token), []byte(s.deviceAuthKey)) == 1 {
			valid = true
		}
		if s.authKey != "" && subtle.ConstantTimeCompare([]byte(token), []byte(s.authKey)) == 1 {
			valid = true
		}

		if !valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// handleDeviceEvents processes POST /v1/device/events
func (s *Server) handleDeviceEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	var env device.Envelope
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if s.app.Services == nil || s.app.Services.Device == nil {
		http.Error(w, "Device service unavailable", http.StatusServiceUnavailable)
		return
	}

	isDup, err := s.app.Services.Device.IngestEvent(r.Context(), &env)
	if err != nil {
		if errors.Is(err, device.ErrInvalidSchema) ||
			errors.Is(err, device.ErrMissingID) ||
			errors.Is(err, device.ErrMissingDevice) ||
			errors.Is(err, device.ErrInvalidTimestamp) ||
			errors.Is(err, device.ErrMissingType) ||
			errors.Is(err, device.ErrMissingPayload) {
			http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if isDup {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":    "acknowledged",
			"duplicate": true,
			"id":        env.ID,
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":    "persisted",
		"duplicate": false,
		"id":        env.ID,
	})
}

// handleDeviceCommands processes GET /v1/device/commands (long-poll) and POST /v1/device/commands (enqueue)
func (s *Server) handleDeviceCommands(w http.ResponseWriter, r *http.Request) {
	if s.app.Services == nil || s.app.Services.Device == nil {
		http.Error(w, "Device service unavailable", http.StatusServiceUnavailable)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Long-poll: wait up to 25s for pending commands
		deviceID := r.URL.Query().Get("device")
		if strings.TrimSpace(deviceID) == "" {
			deviceID = "alcatel-5022d"
		}

		cmd, err := s.app.Services.Device.WaitForCommand(r.Context(), deviceID, 25*time.Second)
		if err != nil {
			if errors.Is(err, r.Context().Err()) {
				// Client cancelled or disconnected
				return
			}
			http.Error(w, "Internal error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if cmd == nil {
			// Long-poll timeout: 204 No Content
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(cmd)

	case http.MethodPost:
		// Command creation is privileged and strictly restricted to the Operator / Jervis Orchestrator.
		// A phone device token must NEVER be permitted to enqueue commands to itself or other devices.
		auth := r.Header.Get("Authorization")
		token := strings.TrimPrefix(auth, "Bearer ")
		if s.authKey == "" || subtle.ConstantTimeCompare([]byte(token), []byte(s.authKey)) != 1 {
			http.Error(w, "Forbidden: only operator credentials can enqueue commands", http.StatusForbidden)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
		var cmd device.Envelope
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.app.Services.Device.QueueCommand(r.Context(), &cmd); err != nil {
			if errors.Is(err, device.ErrQueueFull) {
				http.Error(w, "Queue full: "+err.Error(), http.StatusTooManyRequests)
				return
			}
			if errors.Is(err, device.ErrInvalidSchema) ||
				errors.Is(err, device.ErrMissingID) ||
				errors.Is(err, device.ErrMissingDevice) ||
				errors.Is(err, device.ErrDisallowedAction) ||
				errors.Is(err, device.ErrMissingCommandID) ||
				errors.Is(err, device.ErrMissingAction) ||
				errors.Is(err, device.ErrUnknownField) {
				http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, "Command rejected: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "queued",
			"id":     cmd.ID,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleDeviceResults processes POST /v1/device/results
func (s *Server) handleDeviceResults(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	var res device.Envelope
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if s.app.Services == nil || s.app.Services.Device == nil {
		http.Error(w, "Device service unavailable", http.StatusServiceUnavailable)
		return
	}

	isDup, err := s.app.Services.Device.IngestResult(r.Context(), &res)
	if err != nil {
		if errors.Is(err, device.ErrInvalidSchema) ||
			errors.Is(err, device.ErrMissingID) ||
			errors.Is(err, device.ErrMissingDevice) ||
			errors.Is(err, device.ErrDisallowedAction) ||
			errors.Is(err, device.ErrMissingCommandID) ||
			errors.Is(err, device.ErrMissingAction) ||
			errors.Is(err, device.ErrUnknownField) {
			http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":    "recorded",
		"duplicate": isDup,
		"id":        res.ID,
	})
}

// handleDeviceStatus processes GET /v1/device/status
func (s *Server) handleDeviceStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.app.Services == nil || s.app.Services.Device == nil {
		http.Error(w, "Device service unavailable", http.StatusServiceUnavailable)
		return
	}

	deviceID := r.URL.Query().Get("device")
	if strings.TrimSpace(deviceID) == "" {
		deviceID = "alcatel-5022d"
	}

	status, err := s.app.Services.Device.GetDeviceStatus(r.Context(), deviceID)
	if err != nil {
		http.Error(w, "Failed to get device status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(status)
}
