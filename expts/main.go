package main

import (
	"fmt"
	"log"

	"github.com/colinmarc/hdfs"
)

func main() {
	client, err := hdfs.New("")
	if err != nil {
		log.Fatal(err)
	}
	err = client.CopyToLocal("/home/hadtmp.txt", "./hadtmp2.txt")
	if err != nil {
		log.Fatal(err)
	}
	file, err := client.Open("/home/hadtmp.txt")
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1000)
	file.ReadAt(buf, 0)

	fmt.Println("start")
	fmt.Println(string(buf))
	fmt.Println("end")
	// => Abominable are the tumblers into which he pours his poison.
}
