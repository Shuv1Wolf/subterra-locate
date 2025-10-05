package operations1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	cref "github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	httpcontr "github.com/pip-services4/pip-services4-go/pip-services4-http-go/controllers"

	clients1 "github.com/Shuv1Wolf/subterra-locate/clients/location-engine/clients/version1"

	"github.com/gorilla/websocket"
)

type LocationEngineOperationsV1 struct {
	*httpcontr.RestOperations
	locationMonitor clients1.ILocationMonitorClientV1
	upgrader        websocket.Upgrader
}

func NewLocationEngineOperationsV1() *LocationEngineOperationsV1 {
	c := LocationEngineOperationsV1{
		RestOperations: httpcontr.NewRestOperations(),
		upgrader: websocket.Upgrader{
			CheckOrigin:     func(r *http.Request) bool { return true },
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
	c.DependencyResolver.Put(context.Background(), "location-monitor", cref.NewDescriptor("location-monitor", "client", "*", "*", "1.0"))
	return &c
}

func (c *LocationEngineOperationsV1) SetReferences(ctx context.Context, references cref.IReferences) {
	c.RestOperations.SetReferences(ctx, references)

	dependency, _ := c.DependencyResolver.GetOneRequired("location-monitor")
	client, ok := dependency.(clients1.ILocationMonitorClientV1)
	if !ok {
		panic("LocationEngineOperationsV1: Cant't resolv dependency 'client' to ILocationMonitorClientV1")
	}
	c.locationMonitor = client
}

func (c *LocationEngineOperationsV1) MonitorDeviceLocationWS(res http.ResponseWriter, req *http.Request) {
	// 1) Upgrade HTTP connection to WebSocket
	conn, err := c.upgrader.Upgrade(res, req, nil)
	if err != nil {
		c.Logger.Warn(context.Background(), "ws upgrade error: %v", err)
		http.Error(res, "websocket upgrade failed", http.StatusBadRequest)
		return
	}
	// Ensure the socket is closed when we return
	defer conn.Close()

	// 2) Parse query params
	q := req.URL.Query()

	orgID := strings.TrimSpace(q.Get("org_id"))
	if orgID == "" {
		_ = writeWSJSON(conn, map[string]any{"error": "missing or empty 'org_id'"})
		return
	}

	mapId := strings.TrimSpace(q.Get("map_id"))

	deviceIDs := make([]string, 0)
	if v := strings.TrimSpace(q.Get("device_ids")); v != "" {
		deviceIDs = splitAndTrim(v) // безопасный сплит
	}

	// 3) Create a cancellable context tied to this connection
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	// 4) Start the backend streaming RPC
	stream, err := c.locationMonitor.MonitorDeviceLocation(ctx, orgID, mapId, deviceIDs)
	if err != nil {
		_ = writeWSJSON(conn, map[string]any{"error": fmt.Sprintf("stream start failed: %v", err)})
		return
	}
	defer func() {
		// If your generated client exposes CloseSend() or Close(), call it here
		if closer, ok := any(stream).(interface{ CloseSend() error }); ok {
			_ = closer.CloseSend()
		}
	}()

	// 5) Wire up pumps
	// - readPump: consumes incoming control frames (close/ping/pong). We don't
	// expect messages from the client, but we must read to detect disconnects.
	// - writePump: writes streaming events to the WebSocket as JSON.

	// Channel to forward events from the stream to the writer
	events := make(chan any, 64)
	// Channel to forward errors/EOF
	errCh := make(chan error, 1)

	var once sync.Once
	shutdown := func(reason error) {
		once.Do(func() {
			cancel()
			close(events)
			if reason != nil {
				c.Logger.Info(context.Background(), "ws closed: %v", reason)
			}
		})
	}

	// readPump: keep reading control frames to detect client disconnect
	go func() {
		defer shutdown(nil)
		_ = conn.SetReadDeadline(time.Now().Add(75 * time.Second))
		conn.SetPongHandler(func(string) error {
			return conn.SetReadDeadline(time.Now().Add(75 * time.Second))
		})
		for {
			// We don't care about payloads from client; just drain.
			if _, _, err := conn.ReadMessage(); err != nil {
				return // triggers shutdown via defer
			}
		}
	}()

	// heartbeat: send ping periodically so intermediaries keep the socket open
	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()

	// backend stream reader
	go func() {
		defer func() {
			// signal writer that stream ended
			errCh <- nil
		}()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			// Replace with your real event type
			msg, err := stream.Recv()
			if err != nil {
				errCh <- err
				return
			}

			// Convert the protobuf/struct into a JSONable value.
			// If msg is already JSON-friendly, you can push it directly.
			// Otherwise, map the fields you care about into a DTO.
			events <- msg
		}
	}()

	// writer loop
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-errCh:
			if err != nil && !isStreamEOF(err) {
				_ = writeWSJSON(conn, map[string]any{"error": err.Error()})
			}
			return
		case ev, ok := <-events:
			if !ok {
				return // channel closed
			}
			if err := writeWSJSON(conn, ev); err != nil {
				return
			}
		case <-pingTicker.C:
			_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *LocationEngineOperationsV1) MonitorBeaconLocationWS(res http.ResponseWriter, req *http.Request) {
	// 1) Upgrade HTTP connection to WebSocket
	conn, err := c.upgrader.Upgrade(res, req, nil)
	if err != nil {
		c.Logger.Warn(context.Background(), "ws upgrade error: %v", err)
		http.Error(res, "websocket upgrade failed", http.StatusBadRequest)
		return
	}
	// Ensure the socket is closed when we return
	defer conn.Close()

	// 2) Parse query params
	q := req.URL.Query()

	orgID := strings.TrimSpace(q.Get("org_id"))
	if orgID == "" {
		_ = writeWSJSON(conn, map[string]any{"error": "missing or empty 'org_id'"})
		return
	}

	mapId := strings.TrimSpace(q.Get("map_id"))

	deviceIDs := make([]string, 0)
	if v := strings.TrimSpace(q.Get("beacon_ids")); v != "" {
		deviceIDs = splitAndTrim(v) // безопасный сплит
	}

	// 3) Create a cancellable context tied to this connection
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	// 4) Start the backend streaming RPC
	stream, err := c.locationMonitor.MonitorBeaconLocation(ctx, orgID, mapId, deviceIDs)
	if err != nil {
		_ = writeWSJSON(conn, map[string]any{"error": fmt.Sprintf("stream start failed: %v", err)})
		return
	}
	defer func() {
		// If your generated client exposes CloseSend() or Close(), call it here
		if closer, ok := any(stream).(interface{ CloseSend() error }); ok {
			_ = closer.CloseSend()
		}
	}()

	// 5) Wire up pumps
	// - readPump: consumes incoming control frames (close/ping/pong). We don't
	// expect messages from the client, but we must read to detect disconnects.
	// - writePump: writes streaming events to the WebSocket as JSON.

	// Channel to forward events from the stream to the writer
	events := make(chan any, 64)
	// Channel to forward errors/EOF
	errCh := make(chan error, 1)

	var once sync.Once
	shutdown := func(reason error) {
		once.Do(func() {
			cancel()
			close(events)
			if reason != nil {
				c.Logger.Info(context.Background(), "ws closed: %v", reason)
			}
		})
	}

	// readPump: keep reading control frames to detect client disconnect
	go func() {
		defer shutdown(nil)
		_ = conn.SetReadDeadline(time.Now().Add(75 * time.Second))
		conn.SetPongHandler(func(string) error {
			return conn.SetReadDeadline(time.Now().Add(75 * time.Second))
		})
		for {
			// We don't care about payloads from client; just drain.
			if _, _, err := conn.ReadMessage(); err != nil {
				return // triggers shutdown via defer
			}
		}
	}()

	// heartbeat: send ping periodically so intermediaries keep the socket open
	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()

	// backend stream reader
	go func() {
		defer func() {
			// signal writer that stream ended
			errCh <- nil
		}()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			// Replace with your real event type
			msg, err := stream.Recv()
			if err != nil {
				errCh <- err
				return
			}

			// Convert the protobuf/struct into a JSONable value.
			// If msg is already JSON-friendly, you can push it directly.
			// Otherwise, map the fields you care about into a DTO.
			events <- msg
		}
	}()

	// writer loop
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-errCh:
			if err != nil && !isStreamEOF(err) {
				_ = writeWSJSON(conn, map[string]any{"error": err.Error()})
			}
			return
		case ev, ok := <-events:
			if !ok {
				return // channel closed
			}
			if err := writeWSJSON(conn, ev); err != nil {
				return
			}
		case <-pingTicker.C:
			_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func writeWSJSON(conn *websocket.Conn, v any) error {
	_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	w, err := conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		_ = w.Close()
		return err
	}
	return w.Close()
}

func splitAndTrim(csv string) []string {
	parts := strings.Split(csv, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}
	return out
}

// isStreamEOF checks whether the streaming client reported a terminal condition.
func isStreamEOF(err error) bool {
	if err == nil {
		return true
	}
	// Adjust this to match your gRPC/transport's EOF semantics.
	return errors.Is(err, context.Canceled) || strings.Contains(strings.ToLower(err.Error()), "eof")
}
