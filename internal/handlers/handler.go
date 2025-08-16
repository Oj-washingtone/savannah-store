package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
} // no space after json: for the struct tags

// get albums respond with json
var albums = []Album{
	{ID: "1", Title: "This is an amazing album", Artist: "jalang'O", Price: 100.943},
	{ID: "2", Title: "This is the next big thing", Artist: "An amazing artist", Price: 23424.23},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func GetAlbums(c *gin.Context) {
	// gin context is the most important part of gin
	/*
		it carries request details
		validates and serializes json
		and more
	*/
	c.JSON(http.StatusOK, albums)

	/*
		Context.IndentJSON -> To serialize the struct into and add it to the response
		you can replace it with Context.JSON  to send a compact json
		in practice, the indented form is easier to work with when debugging
	*/
}

// handler to add new item
func AddToAlbum(c *gin.Context) {
	// logic to add to the album

	var newAlbum Album

	// call bind json to bind the received json
	// to newAlbum
	err := c.BindJSON(&newAlbum)

	if err != nil {
		return
	}

	// add the new album to the slice
	albums = append(albums, newAlbum)
	c.JSON(201, newAlbum)
}

func GetAlbumById(c *gin.Context) {
	id := c.Param("id")
	// Context.Param("name") gets the params from request mapping

	// loop through albums to find by id

	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(200, a)
			return
		}
	}

	c.IndentedJSON(404, gin.H{"message": "Item not found"})

}
