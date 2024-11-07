FROM golang:1.23.1 

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . .

RUN go build -o restaurantservice ./main.go

RUN ls -l /app

RUN go mod tidy

CMD ["./restaurantservice"]