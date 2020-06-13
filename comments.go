package main

import (
	"bytes"

	"log"
	"net/http"
	"time"

	"github.com/basilnsage/test-app/router/protos"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
)

type (
	commentJson struct {
		Body string `json:"body" binding:"required"`
		CreatedAt int64 `json:"createdAt"`
		ID uuid.UUID `json:"id"`
	}
	commentsById map[uuid.UUID]commentJson
)

var (
	commentsByPost = make(map[uuid.UUID]commentsById)
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
	router.POST("/event", func(ctx *gin.Context) {
		data, err := ctx.GetRawData()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		event := &protos.GenericEvent{}
		if err = proto.Unmarshal(data, event); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "data did not match expected format",
			})
		}
		ctx.JSON(http.StatusOK, gin.H{"message": event})
	})
	router.POST("/posts/:postId/comments", func(ctx *gin.Context) {
		postId, err := uuid.Parse(ctx.Param("postId"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "malformed postId"})
		}
		comments := commentsByPost[postId]
		if comments == nil {
			comments = make(map[uuid.UUID]commentJson)
			commentsByPost[postId] = comments
		}
		comment := commentJson{}
		err = ctx.Bind(&comment)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		createdAt := time.Now().Unix()
		id := uuid.New()
		comment.CreatedAt = createdAt
		comment.ID = id
		event := &protos.GenericEvent{}
		event.Type = "createComment"
		event.CommentData = &protos.CommentEvent{
			Body: comment.Body,
			CreatedAt: createdAt,
			Uuid: id.String(),
		}
		eventBytes, err := proto.Marshal(event)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		resp, err := http.Post("http://localhost:8100/event", "application-octet-stream", bytes.NewReader(eventBytes))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		} else if resp.StatusCode - 400 >= 0 {
			ctx.JSON(resp.StatusCode, gin.H{"error": resp.Status})
			return
		}
		comments[comment.ID] = comment
		commentsByPost[postId] = comments
		ctx.JSON(http.StatusOK, gin.H{"comment": comment})
	})
	router.GET("/posts/:postId/comments", func(ctx *gin.Context) {
		if postId, err := uuid.Parse(ctx.Param("postId")); err == nil {
			if _, postFound := commentsByPost[postId]; postFound {
				commentsArray := make([]commentJson, 0)
				for _, value := range commentsByPost[postId] {
					commentsArray = append(commentsArray, value)
				}
				ctx.JSON(http.StatusOK, gin.H{
					"comments": commentsArray,
				})
			} else {
				ctx.JSON(http.StatusOK, gin.H{
					"comments": make([]commentJson, 0),
				})
			}
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H {
				"error": "malformed request id",
			})
		}
	})
	router.GET("/posts/:postId/comments/:commentId", func(ctx *gin.Context) {
		postId, postErr := uuid.Parse(ctx.Param("postId"))
		commentId, commentErr := uuid.Parse(ctx.Param("commentId"))
		if postErr == nil && commentErr == nil {
			if comments, postFound := commentsByPost[postId]; postFound {
				if comment, commentFound := comments[commentId]; commentFound {
					ctx.JSON(http.StatusOK, gin.H{
						"comment": comment,
					})
				} else {
					ctx.JSON(http.StatusNotFound, gin.H{
						"error": "comment does not exist",
					})
				}
			} else {
				ctx.JSON(http.StatusNotFound, gin.H{
					"error": "post does not exist",
				})
			}
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H {
				"error": "malformed request id",
			})
		}
	})
	router.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H {
			"hello": "world",
		})
	})
	if err := router.Run(":8001"); err != nil {
		log.Fatalf("unable to run Gin server, err: %s", err)
	}
}
