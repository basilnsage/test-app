FROM golang:alpine3.12

WORKDIR queryService
COPY ./ ./

RUN go build queryService.go

CMD ["./queryService"]
