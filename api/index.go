package api

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var app *gin.Engine

func Routes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "https://nextjs-urlshortner-kecilin.vercel.app")
		return
	})
	r.GET("/:short", func(c *gin.Context) {
		db, err := ConnectToDB()
		if err != nil {
			log.Fatalf("Unable to connect to database: %v", err)
		}
		GetOriginalURL(c, db)
	})
	r.POST("/shorten", func(c *gin.Context) {
		db, err := ConnectToDB()
		if err != nil {
			log.Fatalf("Unable to connect to database: %v", err)
		}
		CreateNewShortURL(c, db)
	})
}

func init() {
	app = gin.New()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}

	app.Use(cors.New(config))
	Routes(app)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func ConnectToDB() (*sql.DB, error) {
	var cfg DBConfig
	cfg.Host = os.Getenv("DB_HOST")
	cfg.Port = os.Getenv("DB_PORT")
	cfg.User = os.Getenv("DB_USER")
	cfg.Password = os.Getenv("DB_PASSWORD")
	cfg.DBName = os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
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
			c.Redirect(302, "https://nextjs-urlshortner-kecilin.vercel.app/not-found")
			return
		}
		log.Fatal(err)
		c.String(500, "Internal server error")
		return
	}

	c.Redirect(302, url)
	return
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}

type ShortURLRequest struct {
	Url       string `json:"url"`
	Shortcode string `json:"shortcode"`
}

func CreateNewShortURL(c *gin.Context, db *sql.DB) {
	var requestData ShortURLRequest

	// Bind JSON data from request to requestData
	if err := c.BindJSON(&requestData); err != nil {
		c.String(400, "Invalid data format")
		return
	}

	originalURL := requestData.Url
	shortCode := requestData.Shortcode

	if originalURL == "" {
		c.String(400, "URL cannot be empty")
		return
	}

	var short string

	// If shortCode is empty, generate a random string
	if shortCode == "" {
		shortCode = randomString(6)
	}

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

	shortCode = "https://kecilin.vercel.app/" + short
	c.String(201, shortCode)
	return
}
