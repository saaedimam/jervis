package working

import (
	"testing"

	"github.com/ioriimasu/jervis/internal/memory/working/model"
)

func TestWorkingMemory_SlidingWindow(t *testing.T) {
	wm := New(3) // Capacity of 3

	e1 := model.NewEntry("1", "content1", nil)
	e2 := model.NewEntry("2", "content2", nil)
	e3 := model.NewEntry("3", "content3", nil)
	e4 := model.NewEntry("4", "content4", nil)

	_ = wm.Add(e1)
	_ = wm.Add(e2)
	_ = wm.Add(e3)

	if len(wm.All()) != 3 {
		t.Errorf("expected 3 entries, got %d", len(wm.All()))
	}

	// Adding 4th entry should prune 1st
	_ = wm.Add(e4)

	if len(wm.All()) != 3 {
		t.Errorf("expected 3 entries after pruning, got %d", len(wm.All()))
	}

	if _, ok := wm.Get("1"); ok {
		t.Error("entry 1 should have been pruned")
	}

	if _, ok := wm.Get("4"); !ok {
		t.Error("entry 4 should exist")
	}

	all := wm.All()
	if all[0].ID() != "2" || all[1].ID() != "3" || all[2].ID() != "4" {
		t.Errorf("unexpected order: %s, %s, %s", all[0].ID(), all[1].ID(), all[2].ID())
	}
}

func TestWorkingMemory_UpdateEntry(t *testing.T) {
	wm := New(3)
	e1 := model.NewEntry("1", "old", nil)
	e2 := model.NewEntry("2", "val", nil)
	_ = wm.Add(e1)
	_ = wm.Add(e2)

	e1Updated := model.NewEntry("1", "new", nil)
	_ = wm.Add(e1Updated)

	if len(wm.All()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(wm.All()))
	}

	ent, _ := wm.Get("1")
	if ent.Content() != "new" {
		t.Error("entry was not updated")
	}

	// e1 should now be at the end of the FIFO queue if we re-added it
	all := wm.All()
	if all[0].ID() != "2" || all[1].ID() != "1" {
		t.Error("update should move entry to the end of FIFO")
	}
}

func TestWorkingMemory_Clear(t *testing.T) {
	wm := New(5)
	_ = wm.Add(model.NewEntry("1", "a", nil))
	wm.Clear()
	if len(wm.All()) != 0 {
		t.Error("Clear failed")
	}
	if wm.Capacity() != 5 {
		t.Error("Capacity should be preserved after clear")
	}
}
