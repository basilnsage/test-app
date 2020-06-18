package main

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/basilnsage/test-app/shared"
	"github.com/basilnsage/test-app/shared/protos"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
)

var (
	posts = make(map[uuid.UUID]shared.PostJson)
	eventBus = "http://localhost:8100/event"
)

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:*"},
		AllowWildcard: true,
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Origin", "Content-Type"},
		MaxAge: 12 * time.Hour,
	}))
	router.POST("/event", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	router.POST("/posts", func(ctx *gin.Context) {
		post := shared.PostJson{}
		err := ctx.Bind(&post)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		createdAt := time.Now().Unix()
		id := uuid.New()
		post.CreatedAt = createdAt
		post.ID = id
		event := &protos.GenericEvent{}
		event.Type = "createPost"
		event.PostData = &protos.PostEvent{
			Title: post.Title,
			Body: post.Body,
			CreatedAt: createdAt,
			Uuid: id.String(),
		}
		eventBytes, err := proto.Marshal(event)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		resp, err := http.Post(eventBus, "application-octet-stream", bytes.NewReader(eventBytes))
		_, err = shared.RespErrorCheck(resp, err)
		if err != nil {
			log.Print("[ERROR] - a downstream issue happened during post creation, error: %s", err.Error())
			return
		}
		posts[post.ID] = post
		ctx.Status(http.StatusCreated)
	})
	router.GET("/posts/:postId", func(ctx *gin.Context) {
		if postId, err := uuid.Parse(ctx.Param("postId")); err == nil {
			if post, postFound := posts[postId]; postFound {
				ctx.JSON(http.StatusOK, gin.H {
					"post": post,
				})
			} else {
				ctx.JSON(http.StatusNotFound, gin.H {
					"error": "post not found",
				})
			}
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H {
				"error": "malformed postId",
			})
		}
	})
	router.GET("/posts", func(ctx *gin.Context) {
		postsArray := make([]shared.PostJson, 0)
		for _, value := range posts {
			postsArray = append(postsArray, value)
		}
		ctx.JSON(200, gin.H {
			"posts": postsArray,
		})
	})
	router.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H {
			"hello": "world",
		})
	})
	if err := router.Run(":8000"); err != nil {
		log.Fatalf("unable to run Gin server, err: %s", err)
	}
}
