FROM golang:1.23.3-alpine

WORKDIR /app

COPY . .

RUN go build -o chat .

EXPOSE 8080

CMD ["./chat"]
