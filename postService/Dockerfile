FROM golang:alpine3.12

WORKDIR postService
COPY ./ ./

RUN go build postService.go

CMD ["./postService"]
