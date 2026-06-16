package aur

import (
	"fmt"
	"strings"
)

type RPCResponse struct {
	Version     int       `json:"version"`
	Type        string    `json:"type"`
	ResultCount int       `json:"resultcount"`
	Results     []Package `json:"results"`
	Error       string    `json:"error,omitempty"`
}

type Package struct {
	Source         string   `json:"-"`
	SourceType     string   `json:"-"`
	SourceURL      string   `json:"-"`
	ID             int      `json:"ID"`
	Name           string   `json:"Name"`
	PackageBaseID  int      `json:"PackageBaseID"`
	PackageBase    string   `json:"PackageBase"`
	Version        string   `json:"Version"`
	Description    string   `json:"Description"`
	URL            string   `json:"URL"`
	NumVotes       int      `json:"NumVotes"`
	Popularity     float64  `json:"Popularity"`
	OutOfDate      *int64   `json:"OutOfDate"`
	Maintainer     *string  `json:"Maintainer"`
	FirstSubmitted int64    `json:"FirstSubmitted"`
	LastModified   int64    `json:"LastModified"`
	URLPath        string   `json:"URLPath"`
	License        []string `json:"License"`
	Depends        []string `json:"Depends"`
	MakeDepends    []string `json:"MakeDepends"`
	CheckDepends   []string `json:"CheckDepends"`
	OptDepends     []string `json:"OptDepends"`
	Conflicts      []string `json:"Conflicts"`
	Provides       []string `json:"Provides"`
	Keywords       []string `json:"Keywords"`
}

func (p Package) DisplaySource() string {
	if strings.TrimSpace(p.Source) == "" || strings.EqualFold(p.Source, "aur") {
		return "AUR"
	}
	return p.Source
}

func (p Package) AURURL() string {
	if p.PackageBase != "" {
		return fmt.Sprintf("https://aur.archlinux.org/packages/%s", p.PackageBase)
	}
	return fmt.Sprintf("https://aur.archlinux.org/packages/%s", p.Name)
}

func (p Package) MaintainerName() string {
	if p.Maintainer == nil || *p.Maintainer == "" {
		return "orphan"
	}
	return *p.Maintainer
}

func (p Package) IsOrphan() bool {
	return p.Maintainer == nil || *p.Maintainer == ""
}

func (p Package) IsOutOfDate() bool {
	return p.OutOfDate != nil && *p.OutOfDate > 0
}

func (p Package) Clone() Package {
	p.License = append([]string(nil), p.License...)
	p.Depends = append([]string(nil), p.Depends...)
	p.MakeDepends = append([]string(nil), p.MakeDepends...)
	p.CheckDepends = append([]string(nil), p.CheckDepends...)
	p.OptDepends = append([]string(nil), p.OptDepends...)
	p.Conflicts = append([]string(nil), p.Conflicts...)
	p.Provides = append([]string(nil), p.Provides...)
	p.Keywords = append([]string(nil), p.Keywords...)
	if p.OutOfDate != nil {
		v := *p.OutOfDate
		p.OutOfDate = &v
	}
	if p.Maintainer != nil {
		v := *p.Maintainer
		p.Maintainer = &v
	}
	return p
}
