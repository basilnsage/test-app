package main

import (
	"bytes"
	"fmt"
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

type (
	commentsById map[uuid.UUID]shared.CommentJson
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
		case "moderateComment":
			log.Print("received event: moderateComment")
			log.Print(event)
			log.Print(commentsByPost)
			if event.CommentData.Status != "approved" && event.CommentData.Status != "rejected" {
				ctx.String(http.StatusBadRequest, "approval status not recognized")
				return
			}
			commentId, err := uuid.Parse(event.CommentData.CommentId)
			if err != nil {
				ctx.String(http.StatusBadRequest, fmt.Sprintf("could not parse commentId: %v", event.CommentData.CommentId))
				return
			}
			postId, err := uuid.Parse(event.CommentData.PostId)
			if err != nil {
				ctx.String(http.StatusBadRequest, fmt.Sprintf("could not parse postId: %v", event.CommentData.PostId))
				return
			}
			comments, ok := commentsByPost[postId]
			if !ok {
				ctx.String(http.StatusNotFound, fmt.Sprintf("could not find post for postId: %s", postId))
				return
			}
			comment, ok := comments[commentId]
			if !ok {
				ctx.String(http.StatusNotFound, fmt.Sprintf("could not find comment for commentId: %s", commentId))
				return
			}
			comment.Status = event.CommentData.Status
			event.Type = "updateComment"
			log.Print(event)
			updatedData, err := proto.Marshal(event)
			if err != nil {
				ctx.String(http.StatusInternalServerError, "unable to marshal updateComment event")
				return
			}
			resp, err := http.Post("http://localhost:8100/event", "application/octet-stream", bytes.NewReader(updatedData))
			status, err := shared.RespErrorCheck(resp, err)
			if err != nil {
				ctx.String(status, err.Error())
			}
		}
		ctx.Status(http.StatusOK)
	})
	router.POST("/posts/:postId/comments", func(ctx *gin.Context) {
		postId, err := uuid.Parse(ctx.Param("postId"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "malformed postId"})
		}
		comments, ok := commentsByPost[postId]
		if !ok {
			comments = make(map[uuid.UUID]shared.CommentJson)
			commentsByPost[postId] = comments
		}
		comment := shared.CommentJson{}
		err = ctx.Bind(&comment)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		createdAt := time.Now().Unix()
		id := uuid.New()
		comment.Status = "pending"
		comment.CreatedAt = createdAt
		comment.ID = id.String()
		event := &protos.GenericEvent{}
		event.Type = "createComment"
		event.CommentData = &protos.CommentEvent{
			Body: comment.Body,
			Status: "pending",
			CreatedAt: createdAt,
			CommentId: comment.ID,
			PostId: postId.String(),
		}
		eventBytes, err := proto.Marshal(event)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		comments[id] = comment
		commentsByPost[postId] = comments
		resp, err := http.Post("http://localhost:8100/event", "application-octet-stream", bytes.NewReader(eventBytes))
		_, err = shared.RespErrorCheck(resp, err)
		if err != nil {
			log.Print("[ERROR] - a downstream issue happened during comment creation, error: %s", err.Error())
		}
		ctx.JSON(http.StatusOK, gin.H{"comment": comment})
	})
	router.GET("/posts/:postId/comments", func(ctx *gin.Context) {
		if postId, err := uuid.Parse(ctx.Param("postId")); err == nil {
			if _, postFound := commentsByPost[postId]; postFound {
				commentsArray := make([]shared.CommentJson, 0)
				for _, value := range commentsByPost[postId] {
					commentsArray = append(commentsArray, value)
				}
				ctx.JSON(http.StatusOK, gin.H{
					"comments": commentsArray,
				})
			} else {
				ctx.JSON(http.StatusOK, gin.H{
					"comments": make([]shared.CommentJson, 0),
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
