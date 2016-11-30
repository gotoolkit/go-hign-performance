package main

import (
	"log"
)

type Address struct {
	Addr string
}

type User struct {
	Address
	Name string
}

func main() {
	strs := make([]*string, 10)
	for key, _ := range strs {
		v := "sdfse"
		strs[key] = &v
	}
	for _, value := range strs {
		log.Println(*value)
	}

	change(strs)

	for _, value := range strs {
		log.Println(*value)
	}


}

func change(strs []*string)  {
	v := "123"
	strs[2] = &v
}

func constantDemo() {
	var a int

	for a := 0; a < 10; a++ {
		log.Printf("a : %d\n", a)
	}

	log.Printf("out a : %d", a)
}

func sliceDemo() {
	var arr [100]int = [100]int{1, 2, 3, 4}
	sli := make([]int, 9, 20)
	sli = arr[0:9]
	//log.Println(sli[:cap(sli)])
	log.Println("length: ", len(sli))
	log.Println("cap: ", cap(sli))
	log.Println(sli[:10])
}

