FROM golang:1.23.1
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o orderservice ./main.go

RUN ls -l /app

CMD ["./orderservice"]