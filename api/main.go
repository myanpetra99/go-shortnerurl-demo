package handler

import (
	"fmt"
	"net/http"
)

/*
func Handler() {

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "This is home page")
	})

	r.GET("/:short", func(c *gin.Context) {
		db, err := config.ConnectToDB()
		if err != nil {
			log.Fatalf("Unable to connect to database: %v", err)
		}
		controller.GetOriginalURL(c, db)
	})

	r.POST("/api/shorten", func(c *gin.Context) {
		db, err := config.ConnectToDB()
		if err != nil {
			log.Fatalf("Unable to connect to database: %v", err)
		}
		controller.CreateNewShortURL(c, db)
	})

	r.Run(":8080")
}
*/

func VercelHandler(w http.ResponseWriter, r http.Request) {
	fmt.Fprintf(w, "Hello, Vercel!")
}
