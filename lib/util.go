package lib

import (
	"log"
	"bytes"
	"encoding/json"
	"net"
)

func TrimBuf(buf []byte) []byte {
	return bytes.Trim(buf, "\x00")
}

func RespError(e error) ([]byte, error) {
	type resp struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	str, err := json.Marshal(resp{"error", e.Error()})
	if err != nil {
		panic(err)
	}
	return str, nil
}

func RespOk(msg string) ([]byte, error) {
	type resp struct {
		Status  string `json: "status"`
		Message string `json: "message"`
	}
	str, err := json.Marshal(resp{"ok", msg})
	if err != nil {
		panic(err)
	}
	return str, nil
}

// A wrapper around *netTCPConn.Write but ensures all bytes from buffer
// are sent
func SendAll(conn *net.TCPConn, buf []byte) error {
	bufSize := len(buf)
	wrote, err := conn.Write(buf)
	if err != nil {
		return err
	}
	if wrote == bufSize {
		return nil
	}
	log.Printf("socket sent %d of %d retrying left\n", wrote, bufSize)
	return SendAll(conn, buf[wrote:])
}
