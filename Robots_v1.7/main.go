package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func main() {
	for i := range 32 {
		for j := range 13 {
			res := md5.Sum([]byte(fmt.Sprintf("rgiskard%02d%02d", i, j)))
			res2 := md5.Sum([]byte(hex.EncodeToString(res[:])))
			if hex.EncodeToString(res2[:]) == "dfb35334bf2a1338fa40e5fbb4ae4753" {
				fmt.Printf("day:%d, month: %d", i, j)
				return
			}
		}
	}
}
