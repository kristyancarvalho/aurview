package tui

import (
	"context"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/kristyancarvalho/aurview/internal/aur"
	"github.com/kristyancarvalho/aurview/internal/ranking"
)

type mouseTestClient struct{}

func (mouseTestClient) Search(context.Context, string) ([]aur.Package, error) {
	return nil, nil
}

func (mouseTestClient) Info(context.Context, string, string) ([]aur.Package, error) {
	return nil, nil
}

func TestMouseClickSearchFocusesInput(t *testing.T) {
	model := mouseTestModel()
	model.focus = focusList

	updated, _ := model.Update(tea.MouseMsg{
		X:      4,
		Y:      1,
		Action: tea.MouseActionPress,
		Button: tea.MouseButtonLeft,
	})
	got := updated.(Model)
	if got.focus != focusSearch {
		t.Fatalf("focus = %v, want search", got.focus)
	}
}

func TestMouseClickPackageRowSelectsResult(t *testing.T) {
	model := mouseTestModel()

	updated, _ := model.Update(tea.MouseMsg{
		X:      2,
		Y:      5,
		Action: tea.MouseActionPress,
		Button: tea.MouseButtonLeft,
	})
	got := updated.(Model)
	if got.focus != focusList {
		t.Fatalf("focus = %v, want list", got.focus)
	}
	if got.selected != 1 {
		t.Fatalf("selected = %d, want 1", got.selected)
	}
}

func TestMouseWheelScrollsListSelection(t *testing.T) {
	model := mouseTestModel()

	updated, _ := model.Update(tea.MouseMsg{
		X:      2,
		Y:      4,
		Action: tea.MouseActionPress,
		Button: tea.MouseButtonWheelDown,
	})
	got := updated.(Model)
	if got.selected != 3 {
		t.Fatalf("selected = %d, want 3", got.selected)
	}
	if got.focus != focusList {
		t.Fatalf("focus = %v, want list", got.focus)
	}
}

func TestMouseWheelScrollsDetail(t *testing.T) {
	model := mouseTestModel()

	updated, _ := model.Update(tea.MouseMsg{
		X:      90,
		Y:      5,
		Action: tea.MouseActionPress,
		Button: tea.MouseButtonWheelDown,
	})
	got := updated.(Model)
	if got.detailScroll != 3 {
		t.Fatalf("detailScroll = %d, want 3", got.detailScroll)
	}
	if got.focus != focusDetail {
		t.Fatalf("focus = %v, want detail", got.focus)
	}
}

func TestMouseHitAreaNarrowLayout(t *testing.T) {
	model := mouseTestModel()
	model.width = 80
	model.height = 24

	if got := model.hitArea(2, 1).kind; got != hitSearch {
		t.Fatalf("search hit = %v", got)
	}
	if got := model.hitArea(2, 5); got.kind != hitListRow || got.index != 1 {
		t.Fatalf("list row hit = %#v", got)
	}
	if got := model.hitArea(2, 20).kind; got != hitDetail {
		t.Fatalf("detail hit = %v", got)
	}
}

func mouseTestModel() Model {
	maint := "alice"
	model := New(Options{Client: mouseTestClient{}})
	model.width = 120
	model.height = 24
	model.focus = focusList
	for i, name := range []string{"paru", "yay", "aurutils", "trizen", "pacseek"} {
		model.results = append(model.results, ranking.RankedPackage{
			Package: aur.Package{
				Name:         name,
				PackageBase:  name,
				Version:      "1.0.0-1",
				Description:  "test package",
				LastModified: int64(1_700_000_000 + i),
				Maintainer:   &maint,
			},
			Score: float64(100 - i),
		})
	}
	return model
}
