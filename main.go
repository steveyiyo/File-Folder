package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var hostname = "localhost:80"

func indexPage(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"showFile": showFile("upload_file"),
	})
}

func showFile(file_folder string) string {
	var pathOut string
	var allOut []string

	files, _ := ioutil.ReadDir(file_folder)
	for _, file := range files {
		if file.IsDir() {
			path := showFile(file_folder + "/" + file.Name())
			pathOut = strings.ReplaceAll(path, "upload_file/", hostname+"/static/")
			allOut = append(allOut, pathOut)
			fmt.Println(allOut)
		} else {
			path := (file_folder + "/" + file.Name())
			pathOut = (strings.ReplaceAll(path, "upload_file/", hostname+"/file/"))
			allOut = append(allOut, pathOut)
			fmt.Println(allOut)
		}
	}
	if pathOut == "" {
		pathOut = "唉呀！目前沒有任何檔案 嗚嗚"
	}
	return pathOut
}

func uploadPage(c *gin.Context) {
	c.HTML(200, "upload.html", nil)
}

func pageNotAvailable(c *gin.Context) {
	c.HTML(404, "404.html", nil)
}

func uploadFile(c *gin.Context) {

	type Result struct {
		Success   bool
		Message   string
		File_Name string
	}

	var r Result

	file, header, err := c.Request.FormFile("upload_file")
	if err != nil {
		r = Result{false, "Bad Request!", ""}
		c.JSON(400, r)
		return
	}

	nsec := time.Now().UnixNano()
	filename := strconv.FormatInt(nsec, 10) + "_" + header.Filename

	out, err := os.Create("upload_file/" + filename)
	if err != nil {
		r = Result{false, "Error!", ""}
		c.JSON(400, r)
	}

	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		r = Result{false, "Error!", ""}
		c.JSON(400, r)
	} else {
		r = Result{true, "Upload Success!", filename}
		c.JSON(201, r)
	}
}

func main() {

	fmt.Print("\n-------------------\n")
	fmt.Print("\nFile and Folder System")
	fmt.Print("\nPort listing at 80")
	fmt.Print("\nRepo: https://github.com/steveyiyo/file_folder")
	fmt.Print("\nAuthor: SteveYi")
	fmt.Print("\n-------------------\n\n")

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.LoadHTMLGlob("static/*")

	router.GET("/", indexPage)
	router.GET("/upload", uploadPage)
	router.POST("/upload", uploadFile)
	router.StaticFS("file/", http.Dir("./upload_file"))
	router.NoRoute(pageNotAvailable)

	router.Run(":80")
	// router.RunTLS(":443", "certificate/ssl.pem", "certificate/ssl.key")
}
