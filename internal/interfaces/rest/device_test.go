package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/saaedimam/jervis/internal/app"
	"github.com/saaedimam/jervis/internal/interfaces/rest"
	"github.com/saaedimam/jervis/internal/services/device"
)

func setupTestServer(t *testing.T, apiAuth, devAuth string) (*rest.Server, func()) {
	t.Helper()
	cfg := app.DefaultConfig()
	cfg.DatabasePath = ":memory:"
	a, err := app.New(context.Background(), cfg)
	if err != nil {
		t.Fatalf("app.New: %v", err)
	}

	srv := rest.NewServerWithConfig(a, rest.ServerConfig{
		Host:          "127.0.0.1",
		Port:          0,
		APIAuthKey:    apiAuth,
		DeviceAuthKey: devAuth,
	})

	cleanup := func() {
		_ = a.Close()
	}
	return srv, cleanup
}

func TestDeviceAuth(t *testing.T) {
	const devSecret = "test-device-key-123"
	srv, cleanup := setupTestServer(t, "api-key", devSecret)
	defer cleanup()

	doPost := func(authHeader string, payload string) int {
		req := httptest.NewRequest(http.MethodPost, "/v1/device/events", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		if authHeader != "" {
			req.Header.Set("Authorization", authHeader)
		}
		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, req)
		return rec.Code
	}

	// 1. No Authorization header -> 401
	if code := doPost("", `{"schema":"phone-agent.v1"}`); code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", code)
	}

	// 2. Wrong token -> 401
	if code := doPost("Bearer wrong-key", `{"schema":"phone-agent.v1"}`); code != http.StatusUnauthorized {
		t.Fatalf("expected 401 with wrong key, got %d", code)
	}

	// 3. Valid auth but missing fields -> 400 (not 401)
	if code := doPost("Bearer "+devSecret, `{"schema":"phone-agent.v1"}`); code != http.StatusBadRequest {
		t.Fatalf("expected 400 on invalid payload, got %d", code)
	}
}

func TestDeviceEventsEndpoint(t *testing.T) {
	const devSecret = "test-device-key"
	srv, cleanup := setupTestServer(t, "api-key", devSecret)
	defer cleanup()

	validEvent := device.Envelope{
		Schema:    device.SchemaVersion,
		ID:        "evt-http-001",
		Device:    "alcatel-5022d",
		Timestamp: time.Now().UnixMilli(),
		Type:      device.EventSMSReceived,
		Payload:   json.RawMessage(`{"sender":"+8801700000000","body":"Hello Mac","sms_timestamp":1726880000000}`),
	}
	body, _ := json.Marshal(validEvent)

	// Ingest event -> 201 Created
	req := httptest.NewRequest(http.MethodPost, "/v1/device/events", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+devSecret)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}

	// Duplicate event -> 200 OK acknowledged
	reqDup := httptest.NewRequest(http.MethodPost, "/v1/device/events", bytes.NewReader(body))
	reqDup.Header.Set("Authorization", "Bearer "+devSecret)
	reqDup.Header.Set("Content-Type", "application/json")
	recDup := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recDup, reqDup)

	if recDup.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for duplicate event, got %d: %s", recDup.Code, recDup.Body.String())
	}

	var resp map[string]any
	_ = json.Unmarshal(recDup.Body.Bytes(), &resp)
	if resp["duplicate"] != true {
		t.Fatalf("expected duplicate: true in response, got %v", resp)
	}
}

func TestDeviceCommandsEndpoint(t *testing.T) {
	const devSecret = "test-device-key"
	const apiSecret = "test-api-key"
	srv, cleanup := setupTestServer(t, apiSecret, devSecret)
	defer cleanup()

	// 0. Device credentials cannot enqueue commands -> 403 Forbidden
	reqDevForbidden := httptest.NewRequest(http.MethodPost, "/v1/device/commands", bytes.NewReader([]byte(`{}`)))
	reqDevForbidden.Header.Set("Authorization", "Bearer "+devSecret)
	recDevForbidden := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recDevForbidden, reqDevForbidden)
	if recDevForbidden.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden when device tries to enqueue command, got %d: %s",
			recDevForbidden.Code, recDevForbidden.Body.String())
	}

	// 1. Post disallowed action -> 400 Bad Request
	badCmd := device.Envelope{
		Schema:    device.SchemaVersion,
		ID:        "cmd-bad-01",
		Device:    "alcatel-5022d",
		Timestamp: time.Now().UnixMilli(),
		Type:      device.TypeCommand,
		Payload:   json.RawMessage(`{"command_id":"bad-1","action":"arbitrary_exec"}`),
	}
	badBody, _ := json.Marshal(badCmd)
	reqBad := httptest.NewRequest(http.MethodPost, "/v1/device/commands", bytes.NewReader(badBody))
	reqBad.Header.Set("Authorization", "Bearer "+apiSecret)
	recBad := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recBad, reqBad)

	if recBad.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for disallowed action, got %d: %s", recBad.Code, recBad.Body.String())
	}

	// 2. Post allowed command -> 202 Accepted
	cmdPayload, _ := json.Marshal(device.CommandPayload{
		CommandID: "cmd-sms-01",
		Action:    device.ActionSendSMS,
		Recipient: "+8801700000000",
		Text:      "Acknowledged.",
	})
	goodCmd := device.Envelope{
		Schema:    device.SchemaVersion,
		ID:        "env-cmd-sms-01",
		Device:    "alcatel-5022d",
		Timestamp: time.Now().UnixMilli(),
		Type:      device.TypeCommand,
		Payload:   cmdPayload,
	}
	goodBody, _ := json.Marshal(goodCmd)
	reqGood := httptest.NewRequest(http.MethodPost, "/v1/device/commands", bytes.NewReader(goodBody))
	reqGood.Header.Set("Authorization", "Bearer "+apiSecret)
	recGood := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recGood, reqGood)

	if recGood.Code != http.StatusAccepted {
		t.Fatalf("expected 202 Accepted for valid command, got %d: %s", recGood.Code, recGood.Body.String())
	}

	// 3. Long-poll GET /v1/device/commands retrieves the queued command
	reqPoll := httptest.NewRequest(http.MethodGet, "/v1/device/commands?device=alcatel-5022d", nil)
	reqPoll.Header.Set("Authorization", "Bearer "+devSecret)
	recPoll := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recPoll, reqPoll)

	if recPoll.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from long-poll, got %d: %s", recPoll.Code, recPoll.Body.String())
	}

	var fetchedCmd device.Envelope
	if err := json.Unmarshal(recPoll.Body.Bytes(), &fetchedCmd); err != nil {
		t.Fatalf("failed decoding fetched command: %v", err)
	}
	if fetchedCmd.ID != goodCmd.ID {
		t.Fatalf("expected command ID %s, got %s", goodCmd.ID, fetchedCmd.ID)
	}
}

func TestDeviceResultsAndStatusEndpoint(t *testing.T) {
	const devSecret = "test-device-key"
	const apiSecret = "test-api-key"
	srv, cleanup := setupTestServer(t, apiSecret, devSecret)
	defer cleanup()

	// Post execution result
	resPayload, _ := json.Marshal(device.ResultPayload{
		CommandID: "cmd-sms-01",
		Action:    device.ActionSendSMS,
		Status:    device.StatusSuccess,
		Data:      json.RawMessage(`{"sent":true,"parts":1}`),
	})
	resEnv := device.Envelope{
		Schema:    device.SchemaVersion,
		ID:        "res-env-01",
		Device:    "alcatel-5022d",
		Timestamp: time.Now().UnixMilli(),
		Type:      device.TypeResult,
		Payload:   resPayload,
	}
	body, _ := json.Marshal(resEnv)

	req := httptest.NewRequest(http.MethodPost, "/v1/device/results", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+devSecret)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for result, got %d: %s", rec.Code, rec.Body.String())
	}

	// Query status
	reqStatus := httptest.NewRequest(http.MethodGet, "/v1/device/status?device=alcatel-5022d", nil)
	reqStatus.Header.Set("Authorization", "Bearer "+apiSecret)
	recStatus := httptest.NewRecorder()
	srv.Handler().ServeHTTP(recStatus, reqStatus)

	if recStatus.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for status, got %d: %s", recStatus.Code, recStatus.Body.String())
	}

	var status device.DeviceStatus
	if err := json.Unmarshal(recStatus.Body.Bytes(), &status); err != nil {
		t.Fatalf("failed decoding status response: %v", err)
	}
	if status.DeviceID != "alcatel-5022d" {
		t.Fatalf("expected device_id 'alcatel-5022d', got %s", status.DeviceID)
	}
	if !status.Connected {
		t.Fatal("expected device to be marked connected after result activity")
	}
}
