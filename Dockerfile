FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN ls -la /app

COPY cmd/deathfirearsenal .

RUN go build -o main

EXPOSE 8080:8080

CMD ["./main"]
