package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"time"
	"log"
)

func Byte2Gzip(data []byte) []byte {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(data)
	gz.Close()
	return buf.Bytes()
}

func main() {
	Test()
	time.Sleep(time.Second * 10)
	log.Println("done")
}

func Test() {
	for i := 0; i < 10000; i++ {
		go Test100000(10) //drop go then ok
		//time.Sleep(time.Nanosecond * 1)//add sleep then ok
	}
}

func Test100000(t int) {
	for i := 0; i < t; i++ {
		println(fmt.Sprintf("%d", len(Byte2Gzip([]byte("dc483e80a7a0bd9ef71d8cf973673924")))))
	}
}