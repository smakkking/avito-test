FROM golang:alpine3.19

COPY . .

RUN go build -mod vendor -o main ./cmd/service/service.go

CMD ./main