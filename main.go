package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var hostname string

func indexPage(c *gin.Context) {
	c.HTML(200, "index.tmpl", nil)
}

func showFile(file_folder string) []string {
	var pathDirOut []string
	var pathOut []string
	var allOut []string

	files, _ := ioutil.ReadDir(file_folder)
	for _, file := range files {
		if file.IsDir() {
			path := showFile(file_folder + "/" + file.Name())
			pathDirOut = append(pathDirOut, path...)
		} else {
			path := (file_folder + "/" + file.Name())
			pathOut = append(pathOut, path)
		}
	}
	for _, out := range pathDirOut {
		allOut = append(allOut, strings.ReplaceAll(out, "upload_file/", hostname+"/static/"))
	}
	for _, out := range pathOut {
		allOut = append(allOut, strings.ReplaceAll(out, "upload_file/", hostname+"/file/"))
	}
	return allOut
}

func uploadPage(c *gin.Context) {
	c.HTML(200, "upload.tmpl", nil)
}

func listPage(c *gin.Context) {
	hostname = c.Request.Host

	var showFileOut []string
	var outStr string
	showFileOut = showFile("upload_file")
	for _, out := range showFileOut {
		filename := strings.ReplaceAll(out, hostname, "")
		outStr += fmt.Sprintf("<a href='//%s' class='item'><i class='file icon'></i> %s</a>", out, filename)
		if out != "" {
			outStr = fmt.Sprintf("%s", outStr) + "\n"
		}
	}
	if outStr == "" {
		outStr = "唉呀！目前沒有任何檔案 嗚嗚"
	}
	c.HTML(200, "list.tmpl", gin.H{
		"showFile":  template.HTML(outStr),
		"IPAddress": c.ClientIP(),
	})
}

func pageNotAvailable(c *gin.Context) {
	c.HTML(404, "404.tmpl", nil)
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

	if _, err := os.Stat("upload_file"); os.IsNotExist(err) {
		err = os.Mkdir("upload_file", 0755)
		if err != nil {
			fmt.Println(err)
		}
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

	fmt.Print("\n-------------------")
	fmt.Print("\nFile and Folder System")
	fmt.Print("\nPort listing at 80")
	fmt.Print("\nRepo: https://github.com/steveyiyo/file_folder")
	fmt.Print("\nAuthor: SteveYi")
	fmt.Print("\n-------------------\n\n")

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.LoadHTMLGlob("static/*")

	router.GET("/", indexPage)
	router.GET("/list", listPage)
	router.GET("/upload", uploadPage)
	router.POST("/upload", uploadFile)
	router.StaticFS("file/", http.Dir("./upload_file"))
	router.NoRoute(pageNotAvailable)

	router.Run(":80")
	// router.RunTLS(":443", "certificate/ssl.pem", "certificate/ssl.key")
}
