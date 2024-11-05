
FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV PORT=8083

EXPOSE 8083

RUN go build -o user_service main.go

CMD ["./user_service"]
