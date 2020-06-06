package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	commentJson struct {
		Body string `json:"body" binding:"required"`
		CreatedAt int64 `json:"createdAt"`
	}
)

var commentsByPost = make(map[uuid.UUID][]commentJson)

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:*"},
		AllowWildcard: true,
		AllowMethods: []string{"GET", "POST", "OPTION"},
		AllowHeaders: []string{"Origin", "Content-Type"},
		MaxAge: 12 * time.Hour,
	}))
	router.POST("/posts/:postId/comments", func(ctx *gin.Context) {
		if postId, err := uuid.Parse(ctx.Param("postId")); err == nil {
			comments := commentsByPost[postId]
			comment := commentJson{}
			if ctx.Bind(&comment) == nil {
				comment.CreatedAt = time.Now().Unix()
				comments = append(comments, comment)
			}
			commentsByPost[postId] = comments
			ctx.JSON(http.StatusOK, gin.H {
				"comment": comment,
			})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H {
				"error": "malformed request id",
			})
		}
	})
	router.GET("/posts/:postId/comments", func(ctx *gin.Context) {
		log.Printf("%v", ctx.Param("postId"))
		if postId, err := uuid.Parse(ctx.Param("postId")); err == nil {
			if _, postFound := commentsByPost[postId]; postFound {
				ctx.JSON(http.StatusOK, gin.H{
					"comments": commentsByPost[postId],
				})
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
	router.GET("/posts/:postId/comments/:commentId", func(ctx *gin.Context) {
		postId, postErr := uuid.Parse(ctx.Param("postId"))
		commentId, commentErr := strconv.ParseUint(ctx.Param("commentId"), 10, 0)
		if postErr == nil && commentErr == nil {
			if comments, postFound := commentsByPost[postId]; postFound {
				if commentId < uint64(len(comments)) {
					ctx.JSON(http.StatusOK, gin.H{
						"comment": comments[commentId],
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
