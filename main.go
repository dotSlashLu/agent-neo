package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strings"
	"unsafe"

	"github.com/dotSlashLu/agent-neo/lib"
)

const magic = 0x53b

var config = &lib.Config{}

type protoHeader struct {
	Magic int32
}

type protoBody struct {
	FnName    string
	ParamsLen uint32
	Params    []byte
}

func main() {
	flags := parseFlags()
	if err := lib.ParseConfig(flags.configFile, config); err != nil {
		panic(fmt.Sprintf("Error parsing config file: %s", err.Error()))
	}
	log.Printf("read config: %+v\n", config)
	log.Printf("registered modules: %+v\n", registeredModules)
	sock, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		panic(err.Error())
	}
	defer sock.Close()
	log.Println("listening on port ", config.Port)
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
	if err = binary.Read(r, config.Endianness_, &h); err != nil {
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

/*
	struct {
		FnName    [32]byte
		ParamsLen uint32
		Params    []byte
	}
	python struct fmt: 32si{x}s
*/
func parseBody(conn net.Conn) *protoBody {
	buf := make([]byte, 32)
	if _, err := conn.Read(buf); err != nil {
		panic(err.Error())
	}
	fnName := string(bytes.Trim(buf, "\x00"))
	fmt.Printf("parsed fnName: %v\n", fnName)

	buf = make([]byte, 4)
	if _, err := conn.Read(buf); err != nil {
		panic(err.Error())
	}
	paramLen := config.Endianness_.Uint32(buf)
	fmt.Printf("parsed paramLen: %v\n", paramLen)

	paramBuf := make([]byte, paramLen)
	if _, err := conn.Read(paramBuf); err != nil {
		panic(err.Error())
	}
	fmt.Printf("params buff: %v\n", paramBuf)
	return &protoBody{fnName, paramLen, paramBuf}
}

func call(fnFull string, params []byte) ([]byte, error) {
	fnSlice := strings.Split(fnFull, ".")
	modName := fnSlice[0]
	fnName := fnSlice[1]
	fmt.Println("mod", modName, "fn", fnName, "params", params)
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
	ret, _ := call(body.FnName, body.Params)
	fmt.Printf("ret %s\n", ret)
	conn.Write(ret)
	fmt.Println("handle over, byebye")
}
