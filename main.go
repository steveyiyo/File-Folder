package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func indexPage(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func pageNotAvailable(c *gin.Context) {
	c.HTML(404, "404.html", nil)
}

func main() {

	fmt.Print("\n-------------------\n\n")
	fmt.Print("\nFile and Folder System")
	fmt.Print("\nPort listing at 13275")
	fmt.Print("\nRepo: https://github.com/steveyiyo/file_folder")
	fmt.Print("\nAuthor: SteveYi")
	fmt.Print("\n-------------------\n\n")

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.LoadHTMLGlob("static/*")

	router.GET("/", indexPage)
	router.StaticFS("static/", http.Dir("./static"))
	router.NoRoute(pageNotAvailable)

	router.Run(":13275")
}
