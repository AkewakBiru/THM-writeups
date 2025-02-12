package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "10.10.86.123:1337")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 2048)
	if _, err := conn.Read(buf); err != nil {
		log.Fatal(err)
	}

	hexenc := strings.TrimSpace(strings.Split(string(buf), ":")[1])
	enc := make([]byte, hex.DecodedLen(len(hexenc)))
	hex.Decode(enc, []byte(hexenc))
	enc = bytes.TrimRight(enc, "\x00")
	clear(buf)

	var key []byte = make([]byte, 5)
	key[0] = enc[0] ^ byte('T')
	key[1] = enc[1] ^ byte('H')
	key[2] = enc[2] ^ byte('M')
	key[3] = enc[3] ^ byte('{')
	key[4] = enc[len(enc)-1] ^ byte('}')

	plain := make([]byte, len(enc))

	for i := 0; i < len(enc); i++ {
		plain[i] = enc[i] ^ key[i%len(key)]
	}

	fmt.Println("First flag:", string(plain))
	if _, err := conn.Read(buf); err != nil {
		log.Fatal(err)
	}

	clear(buf)
	if _, err := conn.Write(key); err != nil {
		log.Fatal(err)
	}

	if _, err := conn.Read(buf); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Second flag:", string(buf))
}

// con ^ var = b -> con => var ^ b = var2 ^ +
// con ^ var2 = +
