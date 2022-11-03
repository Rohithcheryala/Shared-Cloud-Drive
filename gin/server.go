package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/colinmarc/hdfs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type File struct {
	fileName  string
	owner     string
	shared_to []string
}

func (f *File) AlreadySharedWithUser(userName string) bool {
	shared := false
	for _, name := range f.shared_to {
		if name == userName {
			shared = true
			break
		}

	}
	return shared
}

// type User struct {
// 	userName                  string
// 	shared_with_me_file_names []string
// }

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
	// userNames := GetUserNames(client)
	// userNames = userNames

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
	return ds
	// return FileStore{
	// 	files: []File{{fileName: "f12", owner: "usr1", shared_to: []string{"usr2"}}, {fileName: "f23", owner: "usr2", shared_to: []string{"usr3"}}, {fileName: "fall", owner: "usr3", shared_to: []string{"usr2", "usr1"}}},
	// }
}

func (s *FileStore) Store(file File) {
	s.files = append(s.files, file)
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
	// fmt.Printf("a%vb", userName)
	if userName == "Dropdown" {
		return details
	}
	for _, file := range s.files {
		fmt.Println(file.fileName)
		if file.owner == userName {
			fmt.Println(file.fileName, "is owned by", userName)
			details.owned_files = append(details.owned_files, file.fileName)
		}
		if file.SharedWith(userName) && file.owner != userName {
			fmt.Println(file.fileName, "is shared to", userName)
			details.shared = append(details.shared, [2]string{file.fileName, file.owner})
		}
	}
	return details
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
			if file.AlreadySharedWithUser(userName) {
				return
			}
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

func (s *FileStore) UserHasFile(userName string, fileName string) bool {
	// fmt.Println(s)
	hasIt := false
	for _, file := range s.files {
		if file.fileName == fileName && file.owner == userName {
			hasIt = true
			break
		}
	}
	return hasIt
}
func uploadMetrics(client *hdfs.Client) {
	for replication := 1; replication < 4; replication += 1 {
		f200, _ := os.Open("./1_200")
		f200size, _ := os.Stat("./1_200")
		buffer := make([]byte, f200size.Size())
		f200.Read(buffer)
		filepath := "/home/usr1/" + "1_200"
		time_start := time.Now()
		hdptr, err := client.CreateFile(filepath, replication, 1048576, 0777)
		if err != nil {
			client.Remove(filepath)
			log.Fatal("errhd ", err)
		}

		_, err = hdptr.Write(buffer)
		if err != nil {
			log.Fatal("errwr", err)
		}
		hdptr.Flush()
		hdptr.Close()
		f200.Close()
		elapsed := time.Since(time_start)
		fmt.Printf("replication: %v time: %v ms\n", replication, elapsed)
		err = client.Remove(filepath)
		if err != nil {
			log.Fatal("errrm", err)
		}
	}
}

func downloadMetrics(client *hdfs.Client) {
	filepath := "/home/usr1/1_200"
	for replication := 1; replication < 4; replication += 1 {

		client.Remove("/home/usr1/1_200")
		hdptr, err := client.CreateFile(filepath, replication, 1048576, 0777)
		if err != nil {
			log.Fatal(err)
		}

		//
		f200, _ := os.Open("./1_200")
		f200size, _ := os.Stat("./1_200")
		buffer := make([]byte, f200size.Size())
		f200.Read(buffer)

		hdptr.Write(buffer)
		hdptr.Flush()
		hdptr.Close()

		ch1 := make(chan time.Duration)
		// ch2 := make(chan time.Duration)

		go func() {
			time_start := time.Now()
			client.CopyToLocal(filepath, "./usr1/")
			elapsed := time.Since(time_start)
			ch1 <- elapsed
		}()

		// go func() {
		// 	time_start := time.Now()
		// 	client.CopyToLocal(filepath, "./usr2/")
		// 	elapsed := time.Since(time_start)
		// 	ch2 <- elapsed
		// }()

		el1 := <-ch1
		// el2 := <-ch2

		// fmt.Printf("replication: %v time1: %v ms time2: %v ms avg:%v\n", replication, el1, el2, (el1+el2)/2)
		fmt.Printf("replication: %v time1: %v \n", replication, el1)
		client.Remove(filepath)
	}

}
func main() {
	client, err := hdfs.New("")
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.Use(cors.Default())
	// files data structure
	// key - username
	// value - [owned files],[shared_with_me files]

	ds := FileStoreFromClient(client)

	r.POST("/upload", func(c *gin.Context) {
		fmt.Println("-1")
		fmt.Println(math.MaxInt64 - (math.MaxInt64 / 10))
		fmt.Println(651 * 1024 * 1024)
		// parse form for 32 mb
		// err := c.Request.ParseMultipartForm(32 << 20)
		err := c.Request.ParseMultipartForm(math.MaxInt64 - (math.MaxInt64 / 10))
		if err != nil {
			fmt.Println("err56:", err)
		}
		fmt.Println("-2")
		// get file
		file, h, err := c.Request.FormFile("file")
		userName := c.Request.FormValue("userName")
		fmt.Printf("username:%v<-\n", userName)
		replication, _ := strconv.Atoi(c.Request.FormValue("replication"))
		fmt.Println("repl:", replication)
		if replication == 0 {
			replication = 1
		}

		if err != nil {
			fmt.Println("err0: ", err)
		}
		fmt.Println("making buf")
		buffer := make([]byte, h.Size)

		fmt.Println("read buf")
		file.Read(buffer)

		filepath := "/home/" + userName + "/" + h.Filename
		if ds.UserHasFile(userName, h.Filename) {
			fmt.Println("file already exists", filepath)
			err = client.Remove(filepath)
			if err != nil {
				fmt.Println("err45:", err)
			}
			fmt.Println("file deleted")
			ds.DeleteFile(h.Filename, userName)
		}

		// time start
		time_start := time.Now()
		hdptr, err := client.CreateFile(filepath, replication, 1048576, 0777)

		if err != nil {
			fmt.Println(userName)
			fmt.Println(h.Filename)
			fmt.Println(filepath)
			fmt.Println("err-1: ", err)
		}
		hdptr.Write(buffer)
		hdptr.Flush()
		hdptr.Close()
		// time end
		elapsed := time.Since(time_start).Milliseconds()
		if err != nil {
			fmt.Println("err2: ", err)
			return
		}
		ds.Store(File{
			fileName:  h.Filename,
			owner:     userName,
			shared_to: []string{},
		})

		c.Header("Access-Control-Allow-Origin", "*")
		// c.Header("Content-Disposition", "attachment;filename="+h.Filename)
		// c.Data(http.StatusOK, "application/octet-stream", buff)
		c.JSON(http.StatusOK, gin.H{
			"file":      "uploaded",
			"fileName":  h.Filename,
			"time":      elapsed,
			"time_unit": "millisec",
		})
	})

	r.POST("/download", func(c *gin.Context) {
		// parse form for 32 mb
		c.Request.ParseMultipartForm(32 << 20)
		// get file
		fileName := c.Request.FormValue("fileName")
		ownerName := c.Request.FormValue("ownerName")
		fmt.Println(fileName, ownerName)
		if fileName == "" {
			fmt.Println("err0: ", err)
		}
		filepath := "/home/" + ownerName + "/" + fileName

		time_start := time.Now()
		data, _ := client.ReadFile(filepath)
		err := client.CopyToLocal(filepath, "./"+fileName)
		if err != nil {
			fmt.Println("err89: ", err)
		}
		elapsed := time.Since(time_start).Milliseconds()
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("time", strconv.FormatInt(elapsed, 10))
		c.Header("time_units", "millisec")
		c.Header("Content-Disposition", "attachment;filename="+fileName)
		c.Data(http.StatusOK, "application/octet-stream", data)
		// c.JSON(http.StatusOK, gin.H{
		// 	"file":      base64.StdEncoding.EncodeToString(data),
		// 	"fileName":  fileName,
		// 	"time":      elapsed,
		// 	"time_unit": "millisec",
		// })
		// c.FileAttachment("./"+fileName, fileName)
		fmt.Println("tile elapsed:", elapsed, "ms")
	})

	r.POST("/share", func(c *gin.Context) {
		// parse form for 32 mb
		c.Request.ParseMultipartForm(32 << 20)
		// get file
		fileName := c.Request.FormValue("fileName")
		shareName := c.Request.FormValue("shareName")
		ownerName := c.Request.FormValue("ownerName")

		ds.ShareFileWith(fileName, shareName, ownerName)
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusOK, gin.H{
			"shared": "success",
		})
		fmt.Println(ds)
	})

	r.POST("/delete", func(c *gin.Context) {
		fmt.Println("-", ds)

		// parse form for 32 mb
		c.Request.ParseMultipartForm(32 << 20)
		// get file
		fileName := c.Request.FormValue("fileName")
		ownerName := c.Request.FormValue("ownerName")

		// hdfs logic
		err = client.Remove("/home/" + ownerName + "/" + fileName)
		if err != nil {
			fmt.Println(err)
		}

		ds.DeleteFile(fileName, ownerName)
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusOK, gin.H{
			// "file":     buff,
			// "fileName": fileName,
			"deleted": "success",
		})
		fmt.Println("=", ds)
	})
	r.POST("/getData", func(c *gin.Context) {
		c.Request.ParseMultipartForm(32 << 20)

		c.Header("Access-Control-Allow-Origin", "*")
		userName := c.Request.FormValue("userName")
		details := ds.GetDetails(userName)
		c.JSON(http.StatusOK, gin.H{
			"message":      "getpong",
			"owner_files":  details.owned_files,
			"shared_files": details.shared,
		})
		// fmt.Print("\n\n\n")
		// fmt.Println(ds)
		// fmt.Print("\n\n\n")
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
