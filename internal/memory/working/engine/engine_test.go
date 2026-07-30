package engine

import (
	"github.com/ioriimasu/jervis/internal/memory/working/model"
	"testing"
)

func TestEngine(t *testing.T) {
	e := New(2)

	if e.Capacity() != 2 {
		t.Errorf("expected capacity 2, got %d", e.Capacity())
	}

	e1 := model.NewEntry("1", "a", nil)
	_ = e.Add(e1)

	if ent, ok := e.Get("1"); !ok || ent != e1 {
		t.Error("Get failed")
	}

	if _, ok := e.Get("missing"); ok {
		t.Error("Get should fail for missing ID")
	}

	// Default capacity
	eDefault := New(0)
	if eDefault.Capacity() != 50 {
		t.Error("Default capacity should be 50")
	}
}

func TestEngine_AddSliding(t *testing.T) {
	e := New(1)
	e1 := model.NewEntry("1", "a", nil)
	e2 := model.NewEntry("2", "b", nil)
	_ = e.Add(e1)
	_ = e.Add(e2)

	if _, ok := e.Get("1"); ok {
		t.Error("1 should be pruned")
	}
	if len(e.All()) != 1 {
		t.Error("should have 1 entry")
	}
}

func TestEngine_Clear(t *testing.T) {
	e := New(10)
	_ = e.Add(model.NewEntry("1", "a", nil))
	e.Clear()
	if len(e.All()) != 0 {
		t.Error("Clear failed")
	}
}

func TestEngine_AddUpdate(t *testing.T) {
	e := New(5)
	e1 := model.NewEntry("1", "a", nil)
	_ = e.Add(e1)

	e1b := model.NewEntry("1", "b", nil)
	_ = e.Add(e1b)

	if len(e.All()) != 1 {
		t.Errorf("expected 1, got %d", len(e.All()))
	}
	ent, _ := e.Get("1")
	if ent.Content() != "b" {
		t.Error("update failed")
	}
}
