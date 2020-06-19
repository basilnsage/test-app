FROM golang:alpine3.12

WORKDIR moderationService
COPY ./ ./

RUN go build moderationService.go

CMD ["./moderationService"]
