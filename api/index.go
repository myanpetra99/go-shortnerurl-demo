package handler

import (
	"log"
	"net/http"
	"database/sql"
	"github.com/gin-gonic/gin"
	"time"
	"fmt"
)

var app *gin.Engine

func Routes(r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "This is home page")
	})
	r.GET("/:short", func(c *gin.Context) {
		db, err := ConnectToDB()
		if err != nil {
			log.Fatalf("Unable to connect to database: %v", err)
		}
		GetOriginalURL(c, db)
	})
	r.POST("/api/shorten", func(c *gin.Context) {
		db, err := ConnectToDB()
		if err != nil {
			log.Fatalf("Unable to connect to database: %v", err)
		}
		CreateNewShortURL(c, db)
	})
}

func init() {
	app = gin.New()
	r := app.Group("/api")
	Routes(r)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}


type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}


func ConnectToDB() (*sql.DB, error) {
	var cfg DBConfig
	cfg.Host = "localhost"
	cfg.Port = 5432
	cfg.User = "postgres"
	cfg.Password = "1234"
	cfg.DBName = "pendek.in"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// Try to connect
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to database")
	return db, nil
}



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