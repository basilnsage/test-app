package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type posts map[uuid.UUID]shared.PostJson
type busEvents struct {
	Events [][]byte `json:"events"`
}

var savedPosts posts = make(map[uuid.UUID]shared.PostJson)

func (p posts) post(event *protos.PostEvent) (int, string) {
	id, _ := uuid.Parse(event.Uuid)
	if _, ok := p[id]; !ok {
		post := shared.PostJson{
			Title:     event.Title,
			Body:      event.Body,
			CreatedAt: event.CreatedAt,
			ID:        id,
			Comments:  make(map[string]shared.CommentJson),
		}
		p[id] = post
		return http.StatusCreated, "post created"
	} else {
		return http.StatusOK, "post already exists"
	}
}

// add a comment to an existing post
func (p posts) addComment(event *protos.CommentEvent) (int, string) {
	postId, err := uuid.Parse(event.PostId)
	if err != nil {
		return http.StatusBadRequest, "could not process postId"
	}
	post, ok := p[postId]
	if !ok {
		return http.StatusNotFound, "could not find specified post"
	}
	commentId := event.CommentId
	_, ok = post.Comments[commentId]
	if ok {
		return http.StatusBadRequest, "comment already exists"
	}
	newComment := shared.CommentJson{
		Body: event.Body,
		Status: event.Status,
		CreatedAt: event.CreatedAt,
		ID: event.CommentId,
	}
	post.Comments[commentId] = newComment
	return http.StatusCreated, "comment created"
}

func (p posts) updateComment(event *protos.CommentEvent) (int, string) {
	postId, err := uuid.Parse(event.PostId)
	if err != nil {
		return http.StatusBadRequest, "could not parse postId"
	}
	commentId := event.CommentId
	post, ok := p[postId]
	if !ok {
		log.Printf("[ERROR] - post does not exist, upstream existenance checks may have failed")
		return http.StatusNotFound, fmt.Sprintf("could not find post: %v", postId)
	}
	comment, ok := post.Comments[commentId]
	if !ok {
		log.Printf("[ERROR] - comment does not exist, upstream existenance checks may have failed")
		return http.StatusNotFound, fmt.Sprintf("could not find comment: %v", commentId)
	}
	comment.Status = event.Status
	comment.Body = event.Body
	comment.ID = event.CommentId
	comment.CreatedAt = event.CreatedAt
	post.Comments[commentId] = comment
	p[postId] = post
	return http.StatusOK, "comment updated"
}
func handleEvent(data *[]byte) (int, string) {
	event := &protos.GenericEvent{}
	err := proto.Unmarshal(*data, event)
	if err != nil {
		return http.StatusBadRequest, "could not parse data"
	}
	switch eventType := event.Type; eventType {
	case "createPost":
		log.Print("received event createPost")
		respCode, respStatus := savedPosts.post(event.PostData)
		return respCode, respStatus
	case "createComment":
		log.Print("received event createComment")
		respCode, respStatus := savedPosts.addComment(event.CommentData)
		return respCode, respStatus
	case "updateComment":
		log.Print("received event updateComment")
		respCode, respStatus := savedPosts.updateComment(event.CommentData)
		return respCode, respStatus
	}
	return http.StatusOK, ""
}

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
		data, err := ctx.GetRawData()
		if err != nil {
			ctx.String(http.StatusBadRequest, "could not read data from request")
			return
		}
		respCode, status := handleEvent(&data)
		ctx.String(respCode, status)
	})
	router.GET("/posts", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"posts": savedPosts,
		})
	})

	log.Print("checking for events")
	resp, err := http.Get("http://localhost:8100/events")
	if err != nil {
		log.Printf("unable to sync events: %s", err.Error())
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[ERROR] - unable to process event, error: %s", err.Error())
		}
		defer func() {
			_ = resp.Body.Close()
		}()
		events := &busEvents{}
		err = json.Unmarshal(body, events)
		for _, event := range events.Events {
			log.Print("processing cached event")
			resp, status := handleEvent(&event)
			log.Printf("response: %d, status %s", resp, status)
		}
	}

	if err := router.Run(":8002"); err != nil {
		log.Fatalf("gin router failed to start: %v", err)
	}
}
