package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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