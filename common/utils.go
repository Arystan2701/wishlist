package common

import "time"

func Timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func VerifyPhone(phone string) bool {
	if len(phone) != 12 {
		return false
	}
	return true
}
