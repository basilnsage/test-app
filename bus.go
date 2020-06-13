package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/basilnsage/test-app/router/protos"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
)

var subbedSvcs = []string{
	"http://localhost:8000/event",
	"http://localhost:8001/event",
}

func extractEventPayload(marshallObj *protos.GenericEvent, ctx *gin.Context) (int, []byte, error) {
	resp := make([]byte, 0)
	event, err := ctx.FormFile("data")
	if err != nil {
		return http.StatusBadRequest, resp, errors.New("unable to extract field data")
	}
	eventFile, err := event.Open()
	if err != nil {
		return http.StatusInternalServerError, resp, err
	}
	eventData, err := ioutil.ReadAll(eventFile)
	if err != nil {
		return http.StatusInternalServerError, resp, err
	}
	err = proto.Unmarshal(eventData, marshallObj)
	if err != nil {
		return http.StatusInternalServerError, resp, err
	}
	resp = eventData
	return http.StatusOK, resp, nil
}

func forwardToSvcs(svcs []string, data *[]byte) (int, error) {
	failedSvcs := make([]string, 0)
	for _, svc := range svcs {
		resp, err := http.Post(svc, "application/octet-stream", bytes.NewReader(*data))
		if err != nil || resp.StatusCode-400 >= 0 {
			failedSvcs = append(failedSvcs, svc)
		}
	}
	if len(failedSvcs) > 0 {
		return http.StatusInternalServerError, fmt.Errorf("unable to forward POST to the following services: %v", failedSvcs)
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
		data, err := ctx.GetRawData()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H {
				"error": "unable to parse binary payload",
			})
			return
		}
		respCode, err := forwardToSvcs(subbedSvcs, &data)
		if err != nil {
			ctx.JSON(respCode, gin.H{"error": err})
		} else {
			ctx.JSON(respCode, gin.H{"message": "all good"})
		}
	})
	router.POST("/verify-event", func(ctx *gin.Context) {
		unmarshalledEvent := &protos.GenericEvent{}
		if statusCode, eventData, err := extractEventPayload(unmarshalledEvent, ctx); err != nil {
			ctx.JSON(statusCode, gin.H {
				"error": fmt.Sprint(err),
			})
		} else {
			log.Print(eventData)
			log.Print(unmarshalledEvent)
			//event := ctx.Param("event")
			//for service := range subbedSvcs {
			//	http.Post(service, "application/octet-stream", )
		}
	})
	if err := router.Run(":8100"); err != nil {
		log.Fatal(err)
	}
}
