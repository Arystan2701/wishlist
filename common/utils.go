package common

import "time"

func Timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func VerifySearchQuery(query string) string {
	if len(query) == 0 {
		return query
	}

	if query[0] == '@' {
		query = query[1:]
	}
	return query
}
