package engine

import (
	"context"
	"testing"
	"time"
	"github.com/ioriimasu/jervis/internal/memory/contracts"
	"github.com/ioriimasu/jervis/internal/memory/store/sqlite"
	"github.com/ioriimasu/jervis/internal/runtime/eventbus/events"
	"github.com/ioriimasu/jervis/internal/runtime/types"
)

func setupTestEngine(t *testing.T) (*Engine, func()) {
	driver, err := sqlite.New(":memory:")
	if err != nil {
		t.Fatalf("failed to create sqlite driver: %v", err)
	}
	
	if err := driver.Initialize(context.Background()); err != nil {
		t.Fatalf("failed to initialize schema: %v", err)
	}

	return New(driver), func() { driver.Close() }
}

func TestEngine(t *testing.T) {
	e, cleanup := setupTestEngine(t)
	defer cleanup()
	
	ctx := context.Background()
	
	id1, _ := types.NewEventID("1")
	ev1, err := events.NewBuilder().SetID(id1).SetType("a.b.c").SetSource("s").SetPayload("p").Build()
	if err != nil {
		t.Fatal(err)
	}
	
	if err := e.Append(ctx, ev1); err != nil {
		t.Fatalf("Append failed: %v", err)
	}
	
	res, err := e.Query(ctx, contracts.Filter{Type: "a.b.c"})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if len(res) != 1 {
		t.Errorf("expected 1 event, got %d", len(res))
	}
}

func TestEngine_QueryFilters(t *testing.T) {
	e, cleanup := setupTestEngine(t)
	defer cleanup()
	
	ctx := context.Background()
	now := types.Now()
	
	id1, _ := types.NewEventID("1")
	ev1, err := events.NewBuilder().SetID(id1).SetType("test.event.occurred").SetSource("s").SetPayload("p").SetTimestamp(now).Build()
	if err != nil {
		t.Fatal(err)
	}
	_ = e.Append(ctx, ev1)
	
	// Filter by past
	res, _ := e.Query(ctx, contracts.Filter{To: now.Time().Add(-1 * time.Second)})
	if len(res) != 0 {
		t.Error("should be 0 for 'To' in past")
	}
	
	// Filter by future
	res, _ = e.Query(ctx, contracts.Filter{From: now.Time().Add(1 * time.Second)})
	if len(res) != 0 {
		t.Error("should be 0 for 'From' in future")
	}
	
	// Limit
	res, _ = e.Query(ctx, contracts.Filter{Limit: 1})
	if len(res) != 1 {
		t.Error("Limit failed")
	}
}

func TestEngine_QueryTypeFilter(t *testing.T) {
	e, cleanup := setupTestEngine(t)
	defer cleanup()
	
	ctx := context.Background()
	id1, _ := types.NewEventID("1")
	ev1, _ := events.NewBuilder().SetID(id1).SetType("test.event.a").SetSource("s").SetPayload("p").Build()
	_ = e.Append(ctx, ev1)
	
	res, _ := e.Query(ctx, contracts.Filter{Type: "other.event"})
	if len(res) != 0 {
		t.Error("Type filter failed")
	}
}
