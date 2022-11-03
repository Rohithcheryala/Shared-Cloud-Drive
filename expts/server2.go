package main

import (
	"fmt"
	"log"

	"github.com/colinmarc/hdfs"
)

type Block struct {
	userName   string
	ownedfiles map[string][]string
}


type FileStore struct {
	data map[string]Block
}

func newFileStore() FileStore {
	return FileStore{}
}

func (s *FileStore) Store(fileName string, owner string, shared_to []string) {

}

type FileAndOwner struct {
	fileName string
	owner    string
}
type Details struct {
	owned_files []string
	shared      []FileAndOwner
}

func (s *FileStore) GetDetails(userName string) Details {
	details := Details{
		owned_files: []string{},
		shared:      []FileAndOwner{},
	}

	return details
}

func main() {
	client, err := hdfs.New("")
	if err != nil {
		log.Fatal(err)
	}
	// files data structure
	// key - username
	// value - [[owned file, shared_with_usernames]]

	ds := newFileStore()

	dirInfo, err := client.ReadDir("/home/")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range dirInfo {
		// fmt.Println(f, f.IsDir())
		if f.IsDir() {
			userName := f.Name()
			userFilesInfo, err := client.ReadDir("/home/" + f.Name())
			if err != nil {
				log.Fatal(err)
			}
			for _, ff := range userFilesInfo {
				ds.Store(ff.Name(), userName, []string{})
			}
		}
	}

	// fmt.Println(ds)
	fmt.Println(ds.GetDetails("usr1"))

	// client.Walk("/home/", visit)
}
