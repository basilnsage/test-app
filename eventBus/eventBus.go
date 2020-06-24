package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/basilnsage/test-app/shared"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	subbedSvcs = []string{
		"http://posts-svc:8000/event",
		"http://comments-svc:8001/event",
		"http://query-svc:8002/event",
		"http://moderation-svc:8003/event",
	}
	events = make([][]byte, 0)
)

func forwardToSvcs(svcs []string, data *[]byte) (int, error) {
	failed := false
	for _, svc := range svcs {
		resp, err := http.Post(svc, "application/octet-stream", bytes.NewReader(*data))
		status, err := shared.RespErrorCheck(resp, err)
		if err != nil {
			failed = true
			log.Printf("[ERROR] - %s failed with status code: %d, error: %s", svc, status, err.Error())
		}
	}
	if failed {
		return http.StatusInternalServerError, fmt.Errorf("a downstream service encountered an issue")
	} else {
		return http.StatusOK, nil
	}
}

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:*"},
		AllowWildcard: true,
		AllowMethods: []string{"GET", "POST", "OPTION"},
		AllowHeaders: []string{"Origin", "Content-Type"},
		MaxAge: 12 * time.Hour,
	}))

	// define routes
	router.POST("/event", func(ctx *gin.Context) {
		log.Print("received event")
		data, err := ctx.GetRawData()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H {
				"error": "unable to parse binary payload",
			})
			return
		}
		events = append(events, data)
		respCode, err := forwardToSvcs(subbedSvcs, &data)
		if err != nil {
			ctx.String(respCode, err.Error())
		} else {
			ctx.Status(respCode)
		}
	})
	router.GET("/events", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H {
			"events": events,
		})
	})
	if err := router.Run(":8100"); err != nil {
		log.Fatal(err)
	}
}
