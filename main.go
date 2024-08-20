package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book
var booksFile = "books.json"

func main() {
	r := gin.Default()

	loadBooks()

	r.GET("/ping", ping)
	r.GET("/books", getBooks)
	r.GET("/books/:id", getBook)
	r.POST("/books", createBook)
	r.PUT("/books/:id", updateBook)
	r.DELETE("/books/:id", deleteBook)
	r.PATCH("/books/:id", patchBook)

	r.Run(":8000")
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func loadBooks() {
	file, err := os.Open(booksFile)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)
	if err := json.Unmarshal(byteValue, &books); err != nil {
		log.Fatalf("Failed to parse JSON: %s", err)
	}
}

func saveBooks() {
	byteValue, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		log.Fatalf("Failed to serialize books: %s", err)
	}

	if err := os.WriteFile(booksFile, byteValue, 0644); err != nil {
		log.Fatalf("Failed to write to file: %s", err)
	}
}

func getBooks(c *gin.Context) {
	c.JSON(http.StatusOK, books)
}

func getBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	for _, book := range books {
		if book.ID == id {
			c.JSON(http.StatusOK, book)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}

func createBook(c *gin.Context) {
	var newBook Book
	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	newBook.ID = len(books) + 1
	books = append(books, newBook)

	saveBooks()

	c.JSON(http.StatusCreated, newBook)
}

func updateBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var updatedBook Book
	if err := c.BindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	for i, book := range books {
		if book.ID == id {
			updatedBook.ID = id
			books[i] = updatedBook
			saveBooks()
			c.JSON(http.StatusOK, updatedBook)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}

func deleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			saveBooks()
			c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}

func patchBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var patchData map[string]interface{}
	if err := c.BindJSON(&patchData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	for i, book := range books {
		if book.ID == id {
			if title, exists := patchData["title"]; exists {
				books[i].Title = title.(string)
			}
			if author, exists := patchData["author"]; exists {
				books[i].Author = author.(string)
			}
			saveBooks()
			c.JSON(http.StatusOK, books[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}
