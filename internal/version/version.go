package version

import "fmt"

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

type Info struct {
	Version string
	Commit  string
	Date    string
}

func Current() Info {
	return Info{
		Version: valueOrDefault(Version, "dev"),
		Commit:  valueOrDefault(Commit, "none"),
		Date:    valueOrDefault(Date, "unknown"),
	}
}

func (i Info) String(name string) string {
	if name == "" {
		name = "aurview"
	}
	return fmt.Sprintf("%s %s\ncommit: %s\ndate: %s\n",
		name,
		valueOrDefault(i.Version, "dev"),
		valueOrDefault(i.Commit, "none"),
		valueOrDefault(i.Date, "unknown"),
	)
}

func valueOrDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
