package main

import (
	"Gee/gee"
	"log"
	"net/http"
)

func main() {
	r := gee.New()

	r.GET("/hello/:key", func(c *gee.Context) {
		// expect /hello/geektutu
		c.String(http.StatusOK, "Hello %s, you're at %s\n", c.Param("key"), c.Path)
	})
	r.GET("/hello/key/user", func(c *gee.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "%s\n", "/hello/key")
	})
	r.GET("/hello/*filepath", func(c *gee.Context) {
		// expect /assets/file1.txt
		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	log.Fatal(r.Run(":9999"))
}
