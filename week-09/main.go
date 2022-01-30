package main

import (
	"encoding/binary"
	"fmt"
)

// 协议头
const rawHeaderLen = 16
const packetOffset = 0
const headerOffset = 4
const verOffset = 6
const opOffset = 8
const seqOffset = 12

func main() {
	data := encoder(1, 1, 1, "Hello, World!")
	decoder(data)
}

// decoder data=协议头+数据
func decoder(data []byte) {
	if len(data) <= 16 {
		fmt.Println("data len < 16.")
		return
	}

	packetLen := binary.BigEndian.Uint32(data[packetOffset:headerOffset])
	fmt.Printf("packetLen:%v\n", packetLen)

	headerLen := binary.BigEndian.Uint16(data[headerOffset:verOffset])
	fmt.Printf("headerLen:%v\n", headerLen)

	version := binary.BigEndian.Uint16(data[verOffset:opOffset])
	fmt.Printf("version:%v\n", version)

	operation := binary.BigEndian.Uint32(data[opOffset:seqOffset])
	fmt.Printf("operation:%v\n", operation)

	sequence := binary.BigEndian.Uint32(data[seqOffset:rawHeaderLen])
	fmt.Printf("sequence:%v\n", sequence)

	body := string(data[rawHeaderLen:])
	fmt.Printf("body:%v\n", body)
}

func encoder(version, operation, sequence int, body string) []byte {
	packetLen := len(body) + rawHeaderLen
	ret := make([]byte, packetLen)

	binary.BigEndian.PutUint32(ret[packetOffset:headerOffset], uint32(packetLen))
	binary.BigEndian.PutUint16(ret[headerOffset:verOffset], uint16(rawHeaderLen))
	binary.BigEndian.PutUint16(ret[verOffset:opOffset], uint16(version))
	binary.BigEndian.PutUint32(ret[opOffset:seqOffset], uint32(operation))
	binary.BigEndian.PutUint32(ret[seqOffset:rawHeaderLen], uint32(sequence))

	byteBody := []byte(body)
	copy(ret[rawHeaderLen:], byteBody)
	return ret
}
