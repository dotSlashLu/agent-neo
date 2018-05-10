package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strings"
	"unsafe"
	// "reflect"
)

const magic = 0x53b

var registeredModules modules

type protoHeader struct {
	Magic int32
}

type protoBody struct {
	FnNameLen uint32
	ParamsLen uint32
	FnName    string
	Params    []byte
}

func main() {
	registeredModules = registerModules()
	log.Printf("registered modules: %+v\n", registeredModules)
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
	// python header proto test:
	// import socket; import struct; sock = socket.socket(); sock.connect(("localhost", 18103)); d = struct.pack("<i", 0x53b); sock.send(d)
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

func parseBody(conn net.Conn) *protoBody {
	buf := make([]byte, 8)
	if _, err := conn.Read(buf); err != nil {
		panic(err.Error())
	}
	fnNameLen := binary.LittleEndian.Uint32(buf[:4])
	paramLen := binary.LittleEndian.Uint32(buf[4:])
	fmt.Printf("parsed fnNameLen: %v, paramLen: %v\n", fnNameLen, paramLen)

	buf = make([]byte, fnNameLen, fnNameLen)
	if _, err := conn.Read(buf); err != nil {
		panic(err.Error())
	}
	fnName := string(buf)
	fmt.Printf("parsed fnName: %v\n", fnName)

	paramBuf := make([]byte, paramLen)
	if _, err := conn.Read(paramBuf); err != nil {
		panic(err.Error())
	}
	fmt.Printf("params buff: %v\n", paramBuf)
	return &protoBody{fnNameLen, paramLen, fnName, paramBuf}
}

func call(fnFull string, params []byte) (string, error) {
	fnSlice := strings.Split(fnFull, ".")
	modName := fnSlice[0]
	fnName := fnSlice[1]
	fmt.Println("mod", modName, "fn", fnName)
	mod := registeredModules[modName]
	fmt.Printf("mod %T %+v\n", mod, mod)
	return mod.Call(fnName, params)
}

func handleConn(conn net.Conn) {
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
	body := parseBody(conn)
	fmt.Printf("parsed body %+v\n", body)
	ret, _ := call(body.FnName, body.Params)
	fmt.Printf("ret %v\n", ret)
	fmt.Println("handle over, byebye")
}
