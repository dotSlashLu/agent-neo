package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"unsafe"
)

const magic = 0x53b

type protoHeader struct {
	Magic   int32
	BodyLen int32
}

type protoBody struct {
	FnNameLen int
	FnName    string
	Params    []byte
}

func main() {
	sock, err := net.Listen("tcp", ":18103")
	defer sock.Close()
	if err != nil {
		panic(err.Error())
	}
	log.Println("listening on port 18103")
	for {
		conn, err := sock.Accept()
		if err != nil {
			panic(fmt.Sprintf("Error accepting: %v", err.Error()))
		}
		go handleConn(conn)
	}
}

func parseHeader(conn net.Conn) *protoHeader {
	size := unsafe.Sizeof(protoHeader{})
	buf := make([]byte, size, size)
	n, err := conn.Read(buf)
	fmt.Println("read ", n, buf)
	if err != nil {
		panic(fmt.Sprintf("Error reading from conn: %v\n", err.Error()))
	}
	r := bytes.NewReader(buf)
	h := protoHeader{}
	if err = binary.Read(r, binary.LittleEndian, &h); err != nil {
		panic(fmt.Sprintf("Error parsing header: ", err.Error()))
	} else {
		fmt.Printf("parse successful %v\n", h)
	}
	if h.Magic != magic {
		panic(fmt.Sprintf("Bad format, magic not right, received %#x",
			h.Magic))
	}
	return &h
}

func handleConn(conn net.Conn) {
	// python header proto test: sock = socket.socket(); sock.connect(("localhost", 18103)); d = struct.pack("<ii", 0x53b, 2); sock.send(d)
	// problems of a single conn should not affect the whole agent
	defer func() {
		if reason := recover(); reason != nil {
			log.Printf("recovered from connection handling error: %s\n",
				reason)
		}
	}()
	defer conn.Close()
	header := parseHeader(conn)
	fmt.Printf("parsed header %+v\n", header)
}
