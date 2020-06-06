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
	postJson struct {
		Title string `json:"title" binding:"required"`
		Body string `json:"body" binding:"required"`
		Author string `json:"author" binding:"required"`
		CreatedAt int64 `json:"createdAt"`
		ID uuid.UUID `json:id`
	}
)

var posts = make([]postJson, 0)

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:*"},
		AllowWildcard: true,
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Origin"},
		MaxAge: 12 * time.Hour,
	}))
	router.POST("/posts", func(ctx *gin.Context) {
		json := postJson{}
		if ctx.Bind(&json) == nil {
			json.CreatedAt = time.Now().Unix()
			json.ID = uuid.New()
			posts = append(posts, json)
			ctx.Status(http.StatusCreated)
		}
	})
	router.GET("/posts/:postId", func(ctx *gin.Context) {
		postId, err := strconv.ParseUint(ctx.Param("postId"), 10, 0)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H {
				"error": "malformed postId",
			})
		} else {
			if l := uint64(len(posts)); postId < l {
				post := posts[postId]
				ctx.JSON(http.StatusOK, gin.H {
					"post": post,
				})
			} else {
				ctx.JSON(http.StatusNotFound, gin.H {
					"error": "post not found",
				})
			}
		}
	})
	router.GET("/posts", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H {
			"posts": posts,
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
