package main

import (
	"fmt"

	"github.com/Thigamore/MD5/hash"
	"github.com/Thigamore/MD5/io"
)

func main() {
	toHash, err := io.GetFile()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", toHash)
	str, _ := hash.Hash(toHash)
	fmt.Println(fmt.Sprintf("toHash: %s", toHash))
	fmt.Printf("Str: %s\n", str)
}
