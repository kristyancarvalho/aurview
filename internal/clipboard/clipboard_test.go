package clipboard

import (
	"context"
	"errors"
	"testing"
)

type fakeCopier struct {
	err    error
	calls  int
	copied string
}

func (f *fakeCopier) Copy(_ context.Context, text string) error {
	f.calls++
	f.copied = text
	return f.err
}

func TestMultiCopierFallback(t *testing.T) {
	first := &fakeCopier{err: ErrUnavailable}
	second := &fakeCopier{}
	copier := NewMulti(first, second)

	if err := copier.Copy(context.Background(), "paru"); err != nil {
		t.Fatalf("Copy() error = %v", err)
	}
	if first.calls != 1 || second.calls != 1 || second.copied != "paru" {
		t.Fatalf("fallback not used correctly: first=%#v second=%#v", first, second)
	}
}

func TestMultiCopierUnavailable(t *testing.T) {
	copier := NewMulti(&fakeCopier{err: ErrUnavailable})
	err := copier.Copy(context.Background(), "paru")
	if !errors.Is(err, ErrUnavailable) {
		t.Fatalf("err = %v, want ErrUnavailable", err)
	}
}
