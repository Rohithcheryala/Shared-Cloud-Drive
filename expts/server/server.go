package main

import (
	"fmt"
	"log"
	"os"

	"github.com/colinmarc/hdfs"
)

type File struct {
	fileName  string
	owner     string
	shared_to []string
}

type User struct {
	userName                  string
	shared_with_me_file_names []string
}

func (f *File) SharedWith(userName string) bool {
	shared := false
	for _, name := range f.shared_to {
		if name == userName {
			shared = true
			break
		}
	}
	return shared
}

type FileStore struct {
	files []File
	// users []User
}

func FileStoreFromClient(client *hdfs.Client) FileStore {
	userNames := GetUserNames(client)
	userNames = userNames

	ds := FileStore{}
	dirInfo, err := client.ReadDir("/home/")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range dirInfo {
		if f.IsDir() {
			userName := f.Name()
			userFilesInfo, err := client.ReadDir("/home/" + f.Name())
			if err != nil {
				log.Fatal(err)
			}
			// assuming there are no folders in user folder
			for _, ff := range userFilesInfo {
				ds.Store(File{
					fileName:  ff.Name(),
					owner:     userName,
					shared_to: []string{},
				})
			}

		}
	}
	// return ds
	return FileStore{
		files: []File{
			{
				fileName: "f12", owner: "usr1", shared_to: []string{"usr2"},
			},
			{
				fileName: "f23", owner: "usr2", shared_to: []string{"usr3"},
			},
			{
				fileName: "fall", owner: "usr3", shared_to: []string{"usr2", "usr1"},
			},
		},
		// users: []User{},
	}
}

func (s *FileStore) Store(file File) {
	s.files = append(s.files, file)
}

type FileAndOwner struct {
	fileName string
	owner    string
}
type Details struct {
	owned_files []string
	shared      [][2]string
}

func (s *FileStore) GetDetails(userName string) Details {
	details := Details{
		owned_files: []string{},
		shared:      [][2]string{},
	}

	for _, file := range s.files {
		if file.owner == userName {
			details.owned_files = append(details.owned_files, file.fileName)
		}
		if file.SharedWith(userName) {
			details.shared = append(details.shared, [2]string{file.fileName, file.owner})
		}
	}
	return details
}

func main() {
	// createImageHandler := http.HandlerFunc(createImage)
	// http.Handle("/txt", createImageHandler)
	// http.ListenAndServe(":8080", nil)
	client, err := hdfs.New("")
	if err != nil {
		log.Fatal(err)
	}
	// files data structure
	// key - username
	// value - [owned files],[shared_with_me files]

	ds := FileStoreFromClient(client)
	fmt.Println("original")
	fmt.Println(ds)

	fmt.Println("\n\nafter a share operation")
	ds.ShareFileWith("f12", "usr3", "usr1")
	fmt.Println(ds)
	fmt.Println(ds.GetDetails("usr1"))
	fmt.Println(ds.GetDetails("usr2"))
	fmt.Println(ds.GetDetails("usr3"))

	fmt.Println("\n\nafter a delete operation")
	ds.DeleteFile("f12", "usr1")
	fmt.Println(ds)
	fmt.Println(ds.GetDetails("usr1"))
	fmt.Println(ds.GetDetails("usr2"))
	fmt.Println(ds.GetDetails("usr3"))
	// client.Walk("/home/", visit)
}

func GetUserNames(client *hdfs.Client) []string {
	dirInfo, err := client.ReadDir("/home/")
	if err != nil {
		log.Fatal(err)
	}
	userNames := []string{}
	for _, f := range dirInfo {
		if f.IsDir() {
			userNames = append(userNames, f.Name())
		}
	}

	return userNames
}

func (s *FileStore) ShareFileWith(fileName string, userName string, owner string) {
	for index, file := range s.files {
		if file.fileName == fileName && file.owner == owner {
			file.shared_to = append(file.shared_to, userName)
		}
		s.files[index] = file
	}
}

func (s *FileStore) DeleteFile(fileName string, owner string) {
	for index, file := range s.files {
		if file.fileName == fileName && file.owner == owner {
			s.files = append(s.files[:index], s.files[index+1:]...)
			return
		}
	}

}

func visit(path string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", path)
	return nil
}

func deleteFile(filename string) {
	client, err := hdfs.New("")
	if err != nil {
		log.Fatal(err)
	}
	username := "usr1"
	err = client.Remove("/home/" + username + "/" + filename)
	if err != nil {
		log.Fatal(err)
	}
}
