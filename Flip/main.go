package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strings"
)

func startConn() (*net.Conn, []byte) {
	conn, err := net.Dial("tcp", "10.10.18.183:1337")
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 2048)
	if _, err := conn.Read(buf); err != nil {
		log.Fatal(err)
	}
	clear(buf)

	if _, err := conn.Read(buf); err != nil {
		log.Fatal(err)
	}

	if _, err := conn.Write([]byte("admil")); err != nil {
		log.Fatal(err)
	}

	if _, err := conn.Read(buf); err != nil {
		log.Fatal(err)
	}
	clear(buf)

	if _, err := conn.Write([]byte("sUp3rPaSs1")); err != nil {
		log.Fatal(err)
	}

	if _, err := conn.Read(buf); err != nil {
		log.Fatal(err)
	}

	cipher := []byte(strings.TrimSpace(strings.Split(string(buf), ":")[1]))

	dec := make([]byte, hex.DecodedLen(len(cipher)))
	hex.Decode(dec, cipher)
	dec = bytes.TrimRight(dec, "\x00")
	return &conn, dec
}

func main() {
	conn, err := net.Dial("tcp", "10.10.18.183:1337")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for i := range 8 {
		conn, cipher := startConn()
		cipher[4] ^= (1 << i)

		dst := make([]byte, hex.EncodedLen(len(cipher)))
		hex.Encode(dst, cipher)

		if _, err := (*conn).Write(dst); err != nil {
			(*conn).Close()
			log.Fatal(err)
		}

		res := make([]byte, 2048)
		if _, err := (*conn).Read(res); err != nil {
			(*conn).Close()
			log.Fatal(err)
		}
		if strings.Contains(string(res), "No way! You got it!\nA nice flag for you:") {
			fmt.Println(string(res))
			(*conn).Close()
			return
		}
		clear(res)
		(*conn).Close()
	}
}
