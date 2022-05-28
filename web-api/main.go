package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type album struct {
	ID     int     `json:"id"` //json:"id" merupakan tag untuk json, kalau tidak pakai tag, defaultnya mengikuti nama variable ID
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: 1, Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: 2, Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: 3, Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func GetAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func PostAlbum(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	if newAlbum.ID <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "album id cannot be 0 or less than 0"})
		return
	}

	for _, a := range albums {
		if a.ID == newAlbum.ID {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "album id already exist"})
			return
		}
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func GetAlbumById(c *gin.Context) {
	id := c.Param("id")
	idc, _ := strconv.Atoi(id)
	for _, a := range albums {
		if a.ID == idc {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})

}

func UpdateAlbum(c *gin.Context) {
	id := c.Param("id")
	idc, _ := strconv.Atoi(id)
	var upAlbum album
	for _, a := range albums {
		if a.ID == idc {
			if err := c.BindJSON(&upAlbum); err != nil {
				return
			}
			//bukan best practice
			albums[a.ID-1].Title = upAlbum.Title
			albums[a.ID-1].Artist = upAlbum.Artist
			albums[a.ID-1].Price = upAlbum.Price

			c.IndentedJSON(http.StatusOK, albums[a.ID-1])
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func DeleteAlbum(c *gin.Context) {
	id := c.Param("id")
	idc, _ := strconv.Atoi(id)
	for _, a := range albums {
		if a.ID == idc {
			albums = append(albums[:a.ID-1], albums[a.ID:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "delete success"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func main() {
	router := gin.Default()
	router.GET("/albums", GetAlbums)
	router.POST("/albums", PostAlbum)
	router.GET("/albums/:id", GetAlbumById)
	router.DELETE("/albums/:id", DeleteAlbum)
	router.PUT("/albums/:id", UpdateAlbum)

	router.Run("localhost:8080")
}
