package lib

import (
	"bytes"
)

func TrimBuf(buf []byte) []byte {
	return bytes.Trim(buf, "\x00")
}
