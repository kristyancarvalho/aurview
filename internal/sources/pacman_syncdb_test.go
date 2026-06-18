package sources

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/kristyancarvalho/aurview/internal/aur"
)

func TestPacmanSyncDBSourceSearchAndInfo(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "core.db")
	writeSyncDB(t, dbPath, map[string]string{
		"pacman-7.1.0-1/desc": `%NAME%
pacman

%BASE%
pacman

%VERSION%
7.1.0-1

%DESC%
A library-based package manager

%URL%
https://archlinux.org/pacman/

%LICENSE%
GPL-2.0-or-later

%BUILDDATE%
1700000000

%PACKAGER%
Arch Linux

%DEPENDS%
bash
glibc

%PROVIDES%
libalpm.so
`,
		"linux-7.0.1-1/desc": `%NAME%
linux

%VERSION%
7.0.1-1

%DESC%
The Linux kernel and modules
`,
	})
	source := NewPacmanSyncDBSource("core", "core", dbPath)

	results, err := source.Search(context.Background(), "package manager")
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Search() returned %d packages, want 1: %#v", len(results), results)
	}
	pkg := results[0]
	if pkg.Name != "pacman" || pkg.DisplaySource() != "core" || pkg.SourceType != "pacman-syncdb" {
		t.Fatalf("package source metadata = %#v, want labeled core pacman-syncdb", pkg)
	}
	if got, want := pkg.Depends, []string{"bash", "glibc"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Depends = %#v, want %#v", got, want)
	}

	info, err := source.Info(context.Background(), "pacman")
	if err != nil {
		t.Fatalf("Info() error = %v", err)
	}
	if info.URL != "https://archlinux.org/pacman/" || info.MaintainerName() != "Arch Linux" {
		t.Fatalf("Info() = %#v, want parsed URL and packager", info)
	}
}

func TestPacmanSyncDBSourceMissingDatabaseReturnsNoResults(t *testing.T) {
	source := NewPacmanSyncDBSource("missing", "missing", filepath.Join(t.TempDir(), "missing.db"))

	results, err := source.Search(context.Background(), "pacman")
	if err != nil {
		t.Fatalf("Search() error = %v, want nil for missing database", err)
	}
	if len(results) != 0 {
		t.Fatalf("Search() returned %#v, want no results", results)
	}
}

func TestPacmanSyncDBSourceSearchesPackagerQueries(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "core.db")
	writeSyncDB(t, dbPath, map[string]string{
		"pacman-7.1.0-1/desc": `%NAME%
pacman

%VERSION%
7.1.0-1

%DESC%
Arch package manager

%PACKAGER%
Arch Linux
`,
		"other-1.0.0-1/desc": `%NAME%
other

%VERSION%
1.0.0-1

%DESC%
Other package

%PACKAGER%
Someone Else
`,
	})
	source := NewPacmanSyncDBSource("core", "core", dbPath)

	results, err := source.Search(context.Background(), "developer:arch")
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if len(results) != 1 || results[0].Name != "pacman" || results[0].MaintainerName() != "Arch Linux" {
		t.Fatalf("Search() = %#v, want only Arch Linux packager match", results)
	}
}

func TestMultiClientMixesAURAndLocalRepositories(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "core.db")
	writeSyncDB(t, dbPath, map[string]string{
		"pacman-7.1.0-1/desc": `%NAME%
pacman

%VERSION%
7.1.0-1

%DESC%
Arch package manager
`,
	})
	client := NewMultiClient([]Source{
		fakeSource{name: "aur", pkgs: []aur.Package{{Name: "paru", Description: "AUR helper"}}},
		NewPacmanSyncDBSource("core", "core", dbPath),
	})

	results, err := client.Search(context.Background(), "pa")
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	got := []string{results[0].DisplaySource(), results[1].DisplaySource()}
	if want := []string{"AUR", "core"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("DisplaySource() values = %#v, want %#v", got, want)
	}
}

func writeSyncDB(t *testing.T, dbPath string, entries map[string]string) {
	t.Helper()

	file, err := os.Create(dbPath)
	if err != nil {
		t.Fatalf("create sync db: %v", err)
	}
	defer file.Close()

	gz := gzip.NewWriter(file)
	defer gz.Close()
	tw := tar.NewWriter(gz)
	defer tw.Close()

	for name, desc := range entries {
		data := []byte(desc)
		if err := tw.WriteHeader(&tar.Header{Name: name, Mode: 0o644, Size: int64(len(data))}); err != nil {
			t.Fatalf("write tar header: %v", err)
		}
		if _, err := tw.Write(data); err != nil {
			t.Fatalf("write tar entry: %v", err)
		}
	}
}
