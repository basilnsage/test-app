package main

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/basilnsage/test-app/shared"
	"github.com/basilnsage/test-app/shared/protos"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:*"},
		AllowWildcard: true,
		AllowMethods: []string{"GET", "POST", "OPTION"},
		AllowHeaders: []string{"Origin", "Content-Type"},
		MaxAge: 12 * time.Hour,
	}))
	router.POST("/event", func(ctx *gin.Context){
		event := &protos.GenericEvent{}
		data, err := ctx.GetRawData()
		if err != nil {
			ctx.String(http.StatusBadRequest, "could not read data from request")
			return
		}
		err = proto.Unmarshal(data, event)
		if err != nil {
			ctx.String(http.StatusBadRequest, "could not parse data")
		}
		switch eventType := event.Type; eventType {
		case "createComment":
			log.Print("received event createComment")
			if strings.Contains(event.CommentData.Body, "orange") {
				event.CommentData.Status = "rejected"
			} else {
				event.CommentData.Status = "approved"
			}
			event.Type = "moderateComment"
			wireframe, err := proto.Marshal(event)
			if err != nil {
				ctx.String(http.StatusInternalServerError, "unable to marshal response")
				return
			}
			resp, err := http.Post("http://localhost:8100/event", "application/octet-stream", bytes.NewReader(wireframe))
			status, err := shared.RespErrorCheck(resp, err)
			if err != nil {
				ctx.String(status, err.Error())
				return
			}
			ctx.Status(status)
		}
		ctx.Status(http.StatusOK)
	})
	if err := router.Run(":8003"); err != nil {
		log.Fatalf("unable to run Gin server, err: %s", err)
	}
}