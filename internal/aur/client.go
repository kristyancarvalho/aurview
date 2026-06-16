package aur

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	DefaultBaseURL = "https://aur.archlinux.org/rpc"
	APIVersion     = "5"
	defaultTimeout = 8 * time.Second
)

var (
	ErrEmptyQuery = errors.New("empty AUR search query")
	ErrRateLimit  = errors.New("AUR RPC rate limit exceeded")
)

type SearchBy string

const (
	SearchByNameDesc     SearchBy = "name-desc"
	SearchByName         SearchBy = "name"
	SearchByMaintainer   SearchBy = "maintainer"
	SearchByDepends      SearchBy = "depends"
	SearchByMakeDepends  SearchBy = "makedepends"
	SearchByOptDepends   SearchBy = "optdepends"
	SearchByCheckDepends SearchBy = "checkdepends"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	baseURL string
	http    HTTPClient
	timeout time.Duration

	mu          sync.RWMutex
	searchCache map[string][]Package
	infoCache   map[string]Package
}

func NewClient(httpClient HTTPClient) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		baseURL:     DefaultBaseURL,
		http:        httpClient,
		timeout:     defaultTimeout,
		searchCache: make(map[string][]Package),
		infoCache:   make(map[string]Package),
	}
}

func (c *Client) WithBaseURL(baseURL string) *Client {
	c.baseURL = strings.TrimRight(baseURL, "/")
	return c
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	if timeout > 0 {
		c.timeout = timeout
	}
	return c
}

func (c *Client) Search(ctx context.Context, query string, by SearchBy) ([]Package, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, ErrEmptyQuery
	}
	if by == "" {
		by = SearchByNameDesc
	}

	key := string(by) + "\x00" + strings.ToLower(query)
	c.mu.RLock()
	if cached, ok := c.searchCache[key]; ok {
		c.mu.RUnlock()
		return clonePackages(cached), nil
	}
	c.mu.RUnlock()

	values := url.Values{}
	values.Set("v", APIVersion)
	values.Set("type", "search")
	values.Set("by", string(by))
	values.Set("arg", query)

	var response RPCResponse
	if err := c.get(ctx, values, &response); err != nil {
		return nil, err
	}
	if response.Type == "error" {
		return nil, fmt.Errorf("AUR RPC error: %s", response.Error)
	}

	c.mu.Lock()
	c.searchCache[key] = clonePackages(response.Results)
	c.mu.Unlock()

	return clonePackages(response.Results), nil
}

func (c *Client) Info(ctx context.Context, names ...string) ([]Package, error) {
	values := url.Values{}
	values.Set("v", APIVersion)
	values.Set("type", "info")

	requested := make([]string, 0, len(names))
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		requested = append(requested, name)
	}
	if len(requested) == 0 {
		return nil, ErrEmptyQuery
	}

	cached := make([]Package, 0, len(requested))
	missing := make([]string, 0, len(requested))
	c.mu.RLock()
	for _, name := range requested {
		if pkg, ok := c.infoCache[strings.ToLower(name)]; ok {
			cached = append(cached, pkg.Clone())
			continue
		}
		missing = append(missing, name)
		values.Add("arg[]", name)
	}
	c.mu.RUnlock()

	if len(missing) == 0 {
		return cached, nil
	}

	var response RPCResponse
	if err := c.get(ctx, values, &response); err != nil {
		return nil, err
	}
	if response.Type == "error" {
		return nil, fmt.Errorf("AUR RPC error: %s", response.Error)
	}

	c.mu.Lock()
	for _, pkg := range response.Results {
		c.infoCache[strings.ToLower(pkg.Name)] = pkg.Clone()
	}
	c.mu.Unlock()

	out := append(cached, clonePackages(response.Results)...)
	return out, nil
}

func (c *Client) get(ctx context.Context, values url.Values, target any) error {
	reqCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	endpoint := c.baseURL
	if strings.Contains(endpoint, "?") {
		endpoint += "&" + values.Encode()
	} else {
		endpoint += "?" + values.Encode()
	}

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "aurview/read-only")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return ErrRateLimit
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("AUR RPC returned HTTP %d", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(target); err != nil {
		return fmt.Errorf("decode AUR RPC response: %w", err)
	}
	return nil
}

func clonePackages(pkgs []Package) []Package {
	out := make([]Package, len(pkgs))
	for i, pkg := range pkgs {
		out[i] = pkg.Clone()
	}
	return out
}
