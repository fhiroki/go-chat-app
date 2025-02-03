FROM golang:1.23.3-alpine

WORKDIR /app

COPY . .

RUN apk update && apk add --no-cache gcc musl-dev
RUN CGO_ENABLED=1 GOOS=linux go build -o chat .

EXPOSE 8080

CMD ["./chat"]
