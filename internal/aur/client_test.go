package aur

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRPCResponseParsing(t *testing.T) {
	payload := `{
		"version": 5,
		"type": "search",
		"resultcount": 1,
		"results": [{
			"ID": 1,
			"Name": "paru",
			"PackageBaseID": 2,
			"PackageBase": "paru",
			"Version": "2.1.0-1",
			"Description": "Feature packed AUR helper",
			"URL": "https://github.com/Morganamilo/paru",
			"NumVotes": 500,
			"Popularity": 12.5,
			"OutOfDate": null,
			"Maintainer": "maint",
			"FirstSubmitted": 1600000000,
			"LastModified": 1700000000,
			"URLPath": "/cgit/aur.git/snapshot/paru.tar.gz",
			"License": ["GPL"],
			"Depends": ["pacman"],
			"MakeDepends": ["cargo"],
			"CheckDepends": ["git"],
			"OptDepends": ["bat: preview"],
			"Conflicts": ["paru-bin"],
			"Provides": ["aur-helper"],
			"Keywords": ["aur", "pacman"]
		}]
	}`

	var resp RPCResponse
	if err := json.Unmarshal([]byte(payload), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.Version != 5 || resp.Type != "search" || resp.ResultCount != 1 {
		t.Fatalf("unexpected metadata: %#v", resp)
	}
	pkg := resp.Results[0]
	if pkg.Name != "paru" || pkg.MaintainerName() != "maint" || pkg.IsOutOfDate() {
		t.Fatalf("unexpected package: %#v", pkg)
	}
	if got := pkg.AURURL(); got != "https://aur.archlinux.org/packages/paru" {
		t.Fatalf("AURURL() = %q", got)
	}
}

func TestClientSearchBuildsReadOnlyRPCRequestAndCaches(t *testing.T) {
	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		query := r.URL.Query()
		if query.Get("v") != APIVersion || query.Get("type") != "search" || query.Get("by") != string(SearchByNameDesc) || query.Get("arg") != "paru" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"version":5,"type":"search","resultcount":1,"results":[{"ID":1,"Name":"paru","PackageBase":"paru","Version":"1","Description":"desc","NumVotes":1,"Popularity":1,"FirstSubmitted":1,"LastModified":1}]}`))
	}))
	defer server.Close()

	client := NewClient(server.Client()).WithBaseURL(server.URL)
	for i := 0; i < 2; i++ {
		pkgs, err := client.Search(context.Background(), "paru", SearchByNameDesc)
		if err != nil {
			t.Fatalf("Search() error = %v", err)
		}
		if len(pkgs) != 1 || pkgs[0].Name != "paru" {
			t.Fatalf("unexpected packages: %#v", pkgs)
		}
	}
	if calls != 1 {
		t.Fatalf("calls = %d, want cached single request", calls)
	}
}

func TestClientSearchErrors(t *testing.T) {
	t.Run("empty query", func(t *testing.T) {
		_, err := NewClient(nil).Search(context.Background(), " ", SearchByNameDesc)
		if !errors.Is(err, ErrEmptyQuery) {
			t.Fatalf("err = %v, want ErrEmptyQuery", err)
		}
	})

	t.Run("rate limit", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
		}))
		defer server.Close()

		_, err := NewClient(server.Client()).WithBaseURL(server.URL).Search(context.Background(), "paru", SearchByNameDesc)
		if !errors.Is(err, ErrRateLimit) {
			t.Fatalf("err = %v, want ErrRateLimit", err)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{nope`))
		}))
		defer server.Close()

		_, err := NewClient(server.Client()).WithBaseURL(server.URL).Search(context.Background(), "paru", SearchByNameDesc)
		if err == nil || !strings.Contains(err.Error(), "decode AUR RPC response") {
			t.Fatalf("err = %v, want decode error", err)
		}
	})
}
