package util

import (
	"crypto/rand"
	"encoding/hex"
)

// UUID4 returns a UUID version 4.
func UUID4() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	b[6] = (b[6] & 0x0f) | (4 << 4)
	b[8] = (b[8] & 0xbf) | 0x80

	buf := make([]byte, 36)
	dash := byte('-')
	hex.Encode(buf[0:8], b[0:4])
	buf[8] = dash
	hex.Encode(buf[9:13], b[4:6])
	buf[13] = dash
	hex.Encode(buf[14:18], b[6:8])
	buf[18] = dash
	hex.Encode(buf[19:23], b[8:10])
	buf[23] = dash
	hex.Encode(buf[24:], b[10:])

	return string(buf)
}
