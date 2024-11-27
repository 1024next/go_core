# bullfight-service/Dockerfile
FROM golang:1.19-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o bullfight-service .

EXPOSE 8081

CMD ["./bullfight-service"]
