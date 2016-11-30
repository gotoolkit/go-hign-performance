package main

import (
	"encoding/hex"
	"log"
)

func main() {
	h := []byte("zzz")
	log.Println("zzz -> []byte: ", h)
	log.Println("zzz -> []byte -> hex string: ",  hex.EncodeToString(h))
	res := make([]byte, len(h) * 2)
	hex.Encode(res, h)
	log.Println("res ([]byte): ", res)
	log.Println("string(res) : ", string(res))
}