package main

import (
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	str := []byte("hello world")
	res := base64.StdEncoding.EncodeToString(str)
	str1, err := base64.StdEncoding.DecodeString(res)
	if err != nil {
		log.Fatal(err)
	}
	if string(str) == string(str1) {
		fmt.Println("here")
	}
	fmt.Println(res)
}
