package timeline

import (
	"context"
	"github.com/ioriimasu/jervis/internal/memory/contracts"
	"github.com/ioriimasu/jervis/internal/memory/store/sqlite"
	"github.com/ioriimasu/jervis/internal/runtime/eventbus/events"
	"github.com/ioriimasu/jervis/internal/runtime/types"
	"testing"
	"time"
)

func setupTestTimeline(t *testing.T) (*Timeline, func()) {
	driver, err := sqlite.New(":memory:")
	if err != nil {
		t.Fatalf("failed to create sqlite driver: %v", err)
	}

	if err := driver.Initialize(context.Background()); err != nil {
		t.Fatalf("failed to initialize schema: %v", err)
	}

	return New(driver), func() { driver.Close() }
}

func TestTimeline_AppendAndQuery(t *testing.T) {
	tl, cleanup := setupTestTimeline(t)
	defer cleanup()

	ctx := context.Background()

	// Create a test event
	id1, _ := types.NewEventID("ev1")
	ev1, err := events.NewBuilder().
		SetID(id1).
		SetType("runtime.test.occurred").
		SetSource("test").
		SetPayload(map[string]any{"key": "val1"}).
		Build()
	if err != nil {
		t.Fatal(err)
	}

	id2, _ := types.NewEventID("ev2")
	ev2, _ := events.NewBuilder().
		SetID(id2).
		SetType("runtime.other.occurred").
		SetSource("test").
		SetPayload(map[string]any{"key": "val2"}).
		Build()

	_ = tl.Append(ctx, ev1)
	_ = tl.Append(ctx, ev2)

	t.Run("Query All", func(t *testing.T) {
		results, err := tl.Query(ctx, contracts.Filter{})
		if err != nil {
			t.Fatal(err)
		}
		if len(results) != 2 {
			t.Errorf("expected 2 events, got %d", len(results))
		}
	})

	t.Run("Filter by Type", func(t *testing.T) {
		results, _ := tl.Query(ctx, contracts.Filter{Type: "runtime.test.occurred"})
		if len(results) != 1 {
			t.Errorf("expected 1 event, got %d", len(results))
		}
		if results[0].ID().String() != "ev1" {
			t.Error("wrong event returned")
		}
	})

	t.Run("Filter by Time Range", func(t *testing.T) {
		now := time.Now().UTC()
		// ev1 and ev2 were just created, so their timestamp should be around now.
		results, _ := tl.Query(ctx, contracts.Filter{
			From: now.Add(-1 * time.Minute),
			To:   now.Add(1 * time.Minute),
		})
		if len(results) != 2 {
			t.Errorf("expected 2 events in time range, got %d", len(results))
		}

		results, _ = tl.Query(ctx, contracts.Filter{
			From: now.Add(1 * time.Minute),
		})
		if len(results) != 0 {
			t.Error("expected 0 events in future time range")
		}
	})

	t.Run("Limit Results", func(t *testing.T) {
		results, _ := tl.Query(ctx, contracts.Filter{Limit: 1})
		if len(results) != 1 {
			t.Errorf("expected 1 event due to limit, got %d", len(results))
		}
	})
}
