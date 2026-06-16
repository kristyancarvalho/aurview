package platform

import "time"

func UnixDate(unix int64) string {
	if unix <= 0 {
		return "-"
	}
	return time.Unix(unix, 0).Local().Format("2006-01-02")
}

func UnixDateTime(unix int64) string {
	if unix <= 0 {
		return "-"
	}
	return time.Unix(unix, 0).Local().Format("2006-01-02 15:04")
}

func OptionalUnixDate(unix *int64) string {
	if unix == nil || *unix <= 0 {
		return "-"
	}
	return UnixDate(*unix)
}
