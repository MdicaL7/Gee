package main

import (
	"Gee/gee"
	"log"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	log.Fatal(r.Run(":9999"))
}
