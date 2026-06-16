package history

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestStoreAddDedupLimitAndNavigate(t *testing.T) {
	store := New(3)
	store.Add("paru")
	store.Add("yay")
	store.Add("trizen")
	store.Add("paru")
	store.Add("aurutils")

	wantItems := []string{"trizen", "paru", "aurutils"}
	if got := store.Items(); !reflect.DeepEqual(got, wantItems) {
		t.Fatalf("Items() = %#v, want %#v", got, wantItems)
	}

	if got, ok := store.Prev(); !ok || got != "aurutils" {
		t.Fatalf("Prev() = %q,%v", got, ok)
	}
	if got, ok := store.Prev(); !ok || got != "paru" {
		t.Fatalf("Prev() = %q,%v", got, ok)
	}
	if got, ok := store.Next(); !ok || got != "aurutils" {
		t.Fatalf("Next() = %q,%v", got, ok)
	}
	if got, ok := store.Next(); ok || got != "" {
		t.Fatalf("Next() = %q,%v, want empty end", got, ok)
	}
}

func TestStoreLoadSave(t *testing.T) {
	path := filepath.Join(t.TempDir(), "aurview", "history")
	store := New(10)
	store.Add("paru")
	store.Add("yay")
	if err := store.Save(path); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded := New(10)
	if err := loaded.Load(path); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got, want := loaded.Items(), store.Items(); !reflect.DeepEqual(got, want) {
		t.Fatalf("loaded = %#v, want %#v", got, want)
	}
}
