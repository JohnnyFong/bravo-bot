package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

func ApiQuery(query string, secert string) string {
	t := time.Now().UnixNano() / int64(time.Millisecond)
	q := query + "&timestamp=" + strconv.FormatInt(t, 10)
	key := []byte(secert)

	sig := hmac.New(sha256.New, key)
	sig.Write([]byte(q))

	return q + "&signature=" + fmt.Sprintf("%x", (sig.Sum(nil)))
}
