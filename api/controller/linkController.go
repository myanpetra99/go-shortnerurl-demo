package controller

import (
	"database/sql"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func GetOriginalURL(c *gin.Context, db *sql.DB) {
	shortcode := c.Param("short")
	var url string

	err := db.QueryRow("SELECT original FROM links WHERE short = $1 AND expired_at >= $2", shortcode, time.Now().Format("2006-01-02")).Scan(&url)
	if err != nil {
		if err == sql.ErrNoRows {
			c.String(404, "No links found with that short or maybe the link has expired")
			return 
		}
		log.Fatal(err)
		c.String(500, "Internal server error")
		return
	}

	c.Redirect(302, url)
	return
}

func CreateNewShortURL(c *gin.Context, db *sql.DB) {
	originalURL := c.PostForm("url")
	shortCode := c.PostForm("shortcode")
	var short string

	var existingShortcode string
	err := db.QueryRow("SELECT short FROM links WHERE short = $1", shortCode).Scan(&existingShortcode)
	if err == nil {
		c.String(401, "Shortcode already exists")
		return
	} else if err != sql.ErrNoRows {
		log.Fatal(err)
		c.String(500, "Internal server error")
		return
	}

	err = db.QueryRow("INSERT INTO links (original, short, expired_at, created_at, status) VALUES ($1, $2, $3, $4, $5) returning short", originalURL, shortCode, time.Now().Add(7*24*time.Hour).Format("2006-01-02"), time.Now().Format("2006-01-02"), true).Scan(&short)
	if err != nil {
		log.Fatal(err)
		c.String(500, "Internal server error")
		return
	}

	shortCode = "http://localhost:8080/" + short
	c.String(201, shortCode)
	return
}
