package service

import (
	"encoding/hex"
	"strconv"
	"time"
)

func AsBytes(args []string) [][]byte {
	bytes := make([][]byte, len(args))
	for i, arg := range args {
		bytes[i] = []byte(arg)
	}
	return bytes
}

func NowMillisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetNonce() string {
	return strconv.FormatInt(NowMillisecond(), 10)
}

func BytesToHex(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return hex.EncodeToString(b)
}
