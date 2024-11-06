FROM golang:1.23.1 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o restaurantservice ./main.go

RUN ls -l /app

CMD ["./restaurantservice"]