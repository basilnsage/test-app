FROM golang:alpine3.12

WORKDIR eventBus
COPY ./ ./

RUN go build eventBus.go

CMD ["./eventBus"]
