FROM golang:alpine3.12

WORKDIR commentService
COPY ./ ./

RUN go build commentService.go

CMD ["./commentService"]
